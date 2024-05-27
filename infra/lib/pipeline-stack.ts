import * as cdk from 'aws-cdk-lib';
import { Construct } from 'constructs';
import { CodePipeline, CodePipelineSource, ShellStep } from 'aws-cdk-lib/pipelines';
import * as codestarconnections from 'aws-cdk-lib/aws-codestarconnections';
import * as ecr from 'aws-cdk-lib/aws-ecr';

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

      const ecrRepository = new ecr.Repository(this, 'Repository');
      const dockerBuildStep = new ShellStep('BuildAndPushDockerImage', {
        commands: [
            'cd app',
            'docker build -t $ECR_URI/app:$CODEBUILD_RESOLVED_SOURCE_VERSION .',
            'aws ecr get-login-password --region $AWS_DEFAULT_REGION | docker login --username AWS --password-stdin $ECR_URI',
            'docker push $ECR_URI/app:$CODEBUILD_RESOLVED_SOURCE_VERSION',
            `echo $CODEBUILD_RESOLVED_SOURCE_VERSION > image_tag.txt`
        ],
        env: {
          'ECR_URI': ecrRepository.repositoryUri
        },
      });

      pipeline.addWave('BuildAndPushImage', {
        post: [
          dockerBuildStep
        ]
      });
      
      stages.forEach(stage => pipeline.addStage(stage));
    }
  }
