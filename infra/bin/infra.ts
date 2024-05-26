#!/usr/bin/env node
import 'source-map-support/register';
import * as cdk from 'aws-cdk-lib';
import Stages from '../lib/constants/stages';
import Accounts from '../lib/constants/accounts';
import Regions from '../lib/constants/regions';
import { Stage } from '../lib/stage';
import { PipelineStack } from '../lib/pipeline-stack';


const app = new cdk.App();

const stages = [
  new Stage(app, `${Stages.STAGING}-${Regions.US_EAST_1}`, Stages.STAGING, {
    env: {account: Accounts.STAGING, region: Regions.US_EAST_1}
  })
]

new PipelineStack(app, 'PipelineStack', stages, {
  env: {
    region: Regions.US_WEST_2,
    account: Accounts.STAGING
  }
})
