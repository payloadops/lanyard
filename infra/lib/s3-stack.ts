import * as cdk from 'aws-cdk-lib';
import * as s3 from 'aws-cdk-lib/aws-s3';
import { Construct } from 'constructs';
import * as iam from 'aws-cdk-lib/aws-iam';

export class S3Stack extends cdk.Stack {
    readonly bucketName: string; 
    constructor(scope: Construct, id: string, props?: cdk.StackProps) {
      super(scope, id, props);
      const bucket = new s3.Bucket(this, "PromptsBucket", {
        versioned: true,
        publicReadAccess: false,
        // removalPolicy: cdk.RemovalPolicy.RETAIN
      })

      const bucketPolicy = new s3.BucketPolicy(this, 'BucketPolicy', { bucket });
      bucketPolicy.document.addStatements(new iam.PolicyStatement({
        actions: ['s3:GetObject'],
        resources: [bucket.arnForObjects('*')],
        principals: [new iam.AccountRootPrincipal()],
      }));

      new cdk.CfnOutput(this, 'BucketNameOutput', {
        value: bucket.bucketName,
        exportName: `BucketName-${this.region}`,
      });

      this.bucketName = bucket.bucketName
    }
}
