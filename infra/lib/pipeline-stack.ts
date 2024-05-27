import * as cdk from 'aws-cdk-lib';
import { Construct } from 'constructs';
import { CodePipeline, CodePipelineSource, ShellStep } from 'aws-cdk-lib/pipelines';
import { Repository } from 'aws-cdk-lib/aws-ecr';
import Regions from './constants/regions';
import Stages from './constants/stages';
import { disambiguator } from './util/disambiguator';

const REPO = "payloadops/plato";

export class PipelineStack extends cdk.Stack {
    constructor(scope: Construct, id: string, stages: cdk.Stage[], props?: cdk.StackProps) {
      super(scope, id, props);

      const pipeline = new CodePipeline(this, 'Pipeline', {
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