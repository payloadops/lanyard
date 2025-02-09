#!/usr/bin/env node
import 'source-map-support/register';
import * as cdk from 'aws-cdk-lib';
import Stages from '../lib/constants/stages';
import Accounts from '../lib/constants/accounts';
import Regions from '../lib/constants/regions';
import { Stage } from '../lib/stage';
import { PipelineStack } from '../lib/pipeline-stack';


const imageTag = process.env.CODEBUILD_RESOLVED_SOURCE_VERSION!;

const app = new cdk.App();

const stages = [
  new Stage(app, `${Stages.DEV}-${Regions.US_EAST_1}`, Stages.DEV, {
    env: {account: Accounts.DEV, region: Regions.US_EAST_1},
    imageTag: imageTag
  }), 
  new Stage(app, `${Stages.PROD}-${Regions.US_EAST_1}`, Stages.PROD, {
    env: {account: Accounts.PROD, region: Regions.US_EAST_1},
    imageTag: imageTag
  }), 
  // new Stage(app, `${Stages.PROD}-${Regions.US_WEST_2}`, Stages.PROD, {
  //   env: {account: Accounts.PROD, region: Regions.US_WEST_2},
  //   imageTag: imageTag
  // }),
  // new Stage(app, `${Stages.PROD}-${Regions.EU_WEST_2}`, Stages.PROD, {
  //   env: {account: Accounts.PROD, region: Regions.EU_WEST_2},
  //   imageTag: imageTag
  // }),  
  // new Stage(app, `${Stages.PROD}-${Regions.EU_CENTRAL_1}`, Stages.PROD, {
  //   env: {account: Accounts.PROD, region: Regions.EU_CENTRAL_1},
  //   imageTag: imageTag
  // }),  
]

new PipelineStack(app, 'PipelineStack', stages, {
  env: {
    region: Regions.US_EAST_1,
    account: Accounts.DEV
  }
})
