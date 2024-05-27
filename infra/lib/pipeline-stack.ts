import * as cdk from 'aws-cdk-lib';
import { Construct } from 'constructs';
import { CodePipeline, CodePipelineSource, ShellStep } from 'aws-cdk-lib/pipelines';
import * as codestarconnections from 'aws-cdk-lib/aws-codestarconnections';
import { Artifact } from 'aws-cdk-lib/aws-codepipeline';
import { PipelineProject } from 'aws-cdk-lib/aws-codebuild';
import { BUILDSPEC } from './constants/buildspec';

const REPO = "payloadops/plato";

export class PipelineStack extends cdk.Stack {
    constructor(scope: Construct, id: string, stages: cdk.Stage[], props?: cdk.StackProps) {
      super(scope, id, props);

      const connection = new codestarconnections.CfnConnection(this, 'MyConnection', {
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

      const project = new PipelineProject(this, 'Project', {
        buildSpec: BUILDSPEC
      })
      const output = new Artifact('Output');
      
      stages.forEach(stage => pipeline.addStage(stage));
    }
  }
