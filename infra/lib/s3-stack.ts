import * as cdk from 'aws-cdk-lib';
import * as s3 from 'aws-cdk-lib/aws-s3';
import { Construct } from 'constructs';

export class S3Stack extends cdk.Stack {
    constructor(scope: Construct, id: string, stage: string, props?: cdk.StackProps) {
      super(scope, id, props);
      new s3.Bucket(this, "PromptsBucket", {
        bucketName: "prompts-bucket",
        versioned: true,
        publicReadAccess: false,
        // removalPolicy: cdk.RemovalPolicy.RETAIN
      })
    }
}
