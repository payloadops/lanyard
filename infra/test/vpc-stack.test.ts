
import { Template } from 'aws-cdk-lib/assertions';
import * as cdk from 'aws-cdk-lib';
import { VpcStack } from '../lib/vpc-stack';

test('VPC Created with Correct CIDR', () => {
  const app = new cdk.App();
  const stack = new VpcStack(app, 'TestVpcStack');
  const template = Template.fromStack(stack);
  
  template.hasResourceProperties('AWS::EC2::VPC', {
    CidrBlock: '10.0.0.0/16'
  });
});
