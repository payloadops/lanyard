
import { Template } from 'aws-cdk-lib/assertions';
import * as cdk from 'aws-cdk-lib';
import { DynamoStack } from '../lib/dynamo-stack';
import Stages from '../lib/constants/stages';


test('DynamoDB Table Created', () => {
  const app = new cdk.App();
  const stack = new DynamoStack(app, 'TestDynamoStack', Stages.DEV);
  const template = Template.fromStack(stack);
  
  template.hasResourceProperties('AWS::DynamoDB::Table', {
    BillingMode: 'PAY_PER_REQUEST'
  });
});
