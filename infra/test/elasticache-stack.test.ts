
import { Template } from 'aws-cdk-lib/assertions';
import * as cdk from 'aws-cdk-lib';
import { VpcStack } from '../lib/vpc-stack';
import { ElastiCacheStack } from '../lib/elasticache-stack';
import Stages from '../lib/constants/stages';

test('ElastiCache Cluster Created', () => {
  const app = new cdk.App();
  const vpcStack = new VpcStack(app, 'TestVpcStack', Stages.DEV);
  const stack = new ElastiCacheStack(app, 'TestElastiCacheStack', {
    vpcStack: vpcStack,
    stage: ""
  });
  const template = Template.fromStack(stack);
  
  template.hasResource('AWS::ElastiCache::CacheCluster', {});
});
