
import { Template } from 'aws-cdk-lib/assertions';
import * as cdk from 'aws-cdk-lib';
import { VpcStack } from '../lib/vpc-stack';
import { EcsStack } from '../lib/ecs-stack';


test('ECS Cluster Created', () => {
  const app = new cdk.App();
  const vpcStack = new VpcStack(app, 'TestVpcStack');
  const stack = new EcsStack(app, 'TestEcsStack', vpcStack);
  const template = Template.fromStack(stack);
  
  template.hasResource('AWS::ECS::Cluster', {});
});
