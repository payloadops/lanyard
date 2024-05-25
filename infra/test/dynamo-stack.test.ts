
import { Template } from 'aws-cdk-lib/assertions';
import * as cdk from 'aws-cdk-lib';
import { DynamoStack } from '../lib/dynamo-stack';


test('DynamoDB Table Created', () => {
  const app = new cdk.App();
  const stack = new DynamoStack(app, 'TestDynamoStack');
  const template = Template.fromStack(stack);
  
  template.hasResourceProperties('AWS::DynamoDB::Table', {
    BillingMode: 'PAY_PER_REQUEST'
  });
});
