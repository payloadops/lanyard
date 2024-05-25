#!/usr/bin/env node
import 'source-map-support/register';
import * as cdk from 'aws-cdk-lib';
import { CodePipeline, CodePipelineSource, ShellStep } from 'aws-cdk-lib/pipelines';

import Stages from '../lib/constants/stages';
import Accounts from '../lib/constants/accounts';
import Regions from '../lib/constants/regions';
import { Stage } from '../lib/stage';

const REPO = "payload/plato"
const app = new cdk.App();

const stages = [
  new Stage(app, `${Stages.STAGING}-${Regions.US_EAST_1}`, Stages.STAGING, {
    env: {account: Accounts.STAGING, region: Regions.US_EAST_1}
  })
]

const pipeline = new CodePipeline(app, 'Pipeline', {
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