import * as cdk from 'aws-cdk-lib';
import { Construct } from 'constructs';
import { CodePipeline, CodePipelineSource, ShellStep } from 'aws-cdk-lib/pipelines';
import * as codestarconnections from 'aws-cdk-lib/aws-codestarconnections';

const REPO = "payloadops/plato";

export class PipelineStack extends cdk.Stack {
    constructor(scope: Construct, id: string, stages: cdk.Stage[], props?: cdk.StackProps) {
      super(scope, id, props);

      const connection = new codestarconnections.CfnConnection(this, 'MyConnection', {
        connectionName: 'GitHubConnection',
        providerType: 'GitHub', // or 'Bitbucket', 'GitHubEnterpriseServer'
      });

      const pipeline = new CodePipeline(this, 'Pipeline', {
        pipelineName: 'Pipeline',
        selfMutation: true,
        synth: new ShellStep('Synth', {
          input: CodePipelineSource.connection(REPO, 'main', {
            connectionArn: connection.attrConnectionArn
          }),
          commands: [
              'cd infra',
              'npm ci',
              'npm run build',
              'npx cdk synth'
            ],
        })
      });
      
      stages.forEach(stage => pipeline.addStage(stage));
    }
  }
