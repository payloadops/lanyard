
import { Template } from 'aws-cdk-lib/assertions';
import * as cdk from 'aws-cdk-lib';
import { VpcStack } from '../lib/vpc-stack';
import { ElastiCacheStack } from '../lib/elasticache-stack';

test('ElastiCache Cluster Created', () => {
  const app = new cdk.App();
  const vpcStack = new VpcStack(app, 'TestVpcStack');
  const stack = new ElastiCacheStack(app, 'TestElastiCacheStack', vpcStack);
  const template = Template.fromStack(stack);
  
  template.hasResource('AWS::ElastiCache::CacheCluster', {});
});
