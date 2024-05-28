import * as cdk from 'aws-cdk-lib';
import { Construct } from 'constructs';
import {
    CodeBuildStep,
    CodePipeline,
    CodePipelineSource,
    ShellStep,
    ManualApprovalStep,
} from 'aws-cdk-lib/pipelines';
import * as codestarconnections from 'aws-cdk-lib/aws-codestarconnections';
import * as ecr from 'aws-cdk-lib/aws-ecr';
import * as iam from 'aws-cdk-lib/aws-iam';
import { LinuxBuildImage } from 'aws-cdk-lib/aws-codebuild';
import Accounts from './constants/accounts';

const REPO = "payloadops/plato";

export class PipelineStack extends cdk.Stack {
    constructor(scope: Construct, id: string, stages: cdk.Stage[], props?: cdk.StackProps) {
      super(scope, id, props);

      const connection = new codestarconnections.CfnConnection(this, 'Connection', {
        connectionName: 'GitHubConnection',
        providerType: 'GitHub',
      });

      const pipeline = new CodePipeline(this, 'Pipeline', {
        pipelineName: 'Pipeline',
        selfMutation: true,
        synth: new ShellStep('Synth', {
          input: CodePipelineSource.connection(REPO, 'main', {
            connectionArn: connection.attrConnectionArn,
          }),
          commands: [
              'cd infra',
              'npm ci',
              'npm run build',
              'npx cdk synth',
            ],
          primaryOutputDirectory: 'infra/cdk.out',
        })
      });

      const ecrRepository = new ecr.Repository(this, 'Repository', {repositoryName: "app"});

      const codeBuildRole = new iam.Role(this, 'CodeBuildRole', {
        assumedBy: new iam.ServicePrincipal('codebuild.amazonaws.com'),
        description: 'Role for CodeBuild to access ECR and other resources',
      });
  
      codeBuildRole.addToPolicy(new iam.PolicyStatement({
        actions: [
          'ecr:GetAuthorizationToken',
          'ecr:GetDownloadUrlForLayer',
          'ecr:BatchGetImage',
          'ecr:GetLoginPassword',
          'ecr:BatchCheckLayerAvailability',
          'ecr:InitiateLayerUpload',
          'ecr:UploadLayerPart',
          'ecr:CompleteLayerUpload',
          'ecr:PutImage',
        ],
        resources: [`*`],
      }));

      const dockerBuildStep = new CodeBuildStep('BuildAndPushDockerImage', {
        commands: [
            'cd app',
            'docker build -t $ECR_URI:latest .',
            'aws ecr get-login-password --region $AWS_DEFAULT_REGION | docker login --username AWS --password-stdin $ECR_URI',
            'docker push $ECR_URI:latest',
        ],
        env: {
            'ECR_URI': ecrRepository.repositoryUri,
            'AWS_DEFAULT_REGION': this.region
        },
        buildEnvironment: {
            buildImage: LinuxBuildImage.STANDARD_5_0,
            privileged: true, // necessary for Docker operations
        },
        role: codeBuildRole // Explicitly specify the IAM role
    });

      pipeline.addWave('BuildAndPushImage', {
        post: [
          dockerBuildStep
        ]
      });
      
      stages.forEach(stage => {
        if (stage.account !== Accounts.PROD) {
          pipeline.addStage(stage, {
              pre: [
                  // ENDPOINT needs the alb uri or address of the ecs task at each stage
                  new CodeBuildStep('RunE2ETests', {
                      commands: [
                          'export ENDPOINT=""',
                          'cd app',
                          'go mod download',
                          'go test -v ./e2e --tags=e2e'
                      ],
                      buildEnvironment: {
                          buildImage: LinuxBuildImage.STANDARD_5_0
                      },
                      role: codeBuildRole, // Ensure the role has the necessary permissions
                  }),
                  new ManualApprovalStep('OverrideE2ETests')
              ],
              post: [
                  // I'm nervous that this might incur substantial costs.
                  new ShellStep('BakeTime', {
                      commands: ['sleep 3600'] // Simulate 1-hour bake time
                  }),
                  new ManualApprovalStep('OverrideBakeTime')
              ]
          })
        } else {
          pipeline.addStage(stage)
        }
      });
    }
  }
