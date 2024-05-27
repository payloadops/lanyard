
import { Template } from 'aws-cdk-lib/assertions';
import * as cdk from 'aws-cdk-lib';
import { VpcStack } from '../lib/vpc-stack';
import { EcsStack } from '../lib/ecs-stack';
import Stages from '../lib/constants/stages';


test('ECS Cluster Created', () => {
  const app = new cdk.App();
  const vpcStack = new VpcStack(app, 'TestVpcStack', Stages.DEV);
  const stack = new EcsStack(app, 'TestEcsStack', vpcStack, Stages.DEV);
  const template = Template.fromStack(stack);
  
  template.hasResource('AWS::ECS::Cluster', {});
});
