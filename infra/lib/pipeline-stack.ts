import * as cdk from 'aws-cdk-lib';
import { Construct } from 'constructs';
import { CodePipeline, CodePipelineSource, ShellStep } from 'aws-cdk-lib/pipelines';

const REPO = "payload/plato";

export class PipelineStack extends cdk.Stack {
    constructor(scope: Construct, id: string, stages: cdk.Stage[], props?: cdk.StackProps) {
      super(scope, id, props);
  
      const pipeline = new CodePipeline(scope, 'Pipeline', {
        pipelineName: 'Pipeline',
        selfMutation: true,
        synth: new ShellStep('Synth', {
          input: CodePipelineSource.gitHub(REPO, 'main'),
          commands: [
              'cd infra',
              'npm ci',
              'npm run build',
              'npx cdk synth'
            ],
          primaryOutputDirectory: 'infra/cdk.out',
        })
      });
      
      stages.forEach(stage => pipeline.addStage(stage))
    }
  }