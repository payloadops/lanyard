#!/usr/bin/env node
import 'source-map-support/register';
import * as cdk from 'aws-cdk-lib';
import { CodePipeline, CodePipelineSource, ShellStep } from 'aws-cdk-lib/pipelines';
import { Stage } from './stage';

// import { ElasticacheStack } from '../lib/elasticache-stack';
const REPO = "payload/plato"
const app = new cdk.App();

const stages = [
  new Stage(app, 'dev', {
    env: {account: "", region: ""}
  })
]

const pipeline = new CodePipeline(app, 'Pipeline', {
  pipelineName: 'Pipeline',
  synth: new ShellStep('Synth', {
    input: CodePipelineSource.gitHub(REPO, 'main'),
    commands: ['npm ci', 'npm run build', 'npx cdk synth']
  })
});

stages.forEach(stage => pipeline.addStage(stage))