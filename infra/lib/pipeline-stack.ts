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
import {Stage} from "./stage";
import Subdomains from './constants/subdomains';
import { DOMAIN } from './constants/domain';

const REPO = "payloadops/lanyard";

export class PipelineStack extends cdk.Stack {
    constructor(scope: Construct, id: string, stages: Stage[], props?: cdk.StackProps) {
      super(scope, id, props);

      const connection = new codestarconnections.CfnConnection(this, 'Connection', {
        connectionName: 'GitHubConnection',
        providerType: 'GitHub',
      });

      const pipeline = new CodePipeline(this, 'Pipeline', {
        crossAccountKeys: true,
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

      const policyStatement = new iam.PolicyStatement({
        actions: [
          'ecr:GetDownloadUrlForLayer',
          'ecr:BatchGetImage',
          'ecr:BatchCheckLayerAvailability',
          'ecr:PutImage',
          // Add more ECR actions as necessary
        ],
        principals: [new iam.AccountPrincipal(Accounts.PROD)]
      });
      
      ecrRepository.addToResourcePolicy(policyStatement);

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
            'docker build -t $ECR_URI:$CODEBUILD_RESOLVED_SOURCE_VERSION .',
            'aws ecr get-login-password --region $AWS_DEFAULT_REGION | docker login --username AWS --password-stdin $ECR_URI',
            'docker push $ECR_URI:$CODEBUILD_RESOLVED_SOURCE_VERSION',
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


      const waves = new Map();

      stages.forEach(stage => {
        const wave = waves.has(stage.stageName) ? waves.get(stage.stageName) : pipeline.addWave(stage.stageName);
        if (stage.account === Accounts.DEV) {
          wave.addStage(stage, {
            post: [
              new CodeBuildStep('RunE2ETests', {
                commands: [
                  'cd app',
                  'go mod download',
                  'go test -v ./e2e --tags=e2e'
                ],
                buildEnvironment: {
                  buildImage: LinuxBuildImage.fromDockerRegistry(
                    'public.ecr.aws/docker/library/golang:1.22-alpine',
                  )
                },
                role: codeBuildRole, // Ensure the role has the necessary permissions
                env: {
                  BASE_URL: `http://${Subdomains.DEV}.${DOMAIN}`
                }
              }),
              new ManualApprovalStep('Manual Approval'),
            ]
          });
        } 
        else {
          wave.addStage(stage);
        }
      });
    }
  }
