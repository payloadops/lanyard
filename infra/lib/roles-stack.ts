import * as cdk from 'aws-cdk-lib';
import { Construct } from 'constructs';
import { StackProps } from 'aws-cdk-lib';
import * as iam from 'aws-cdk-lib/aws-iam';

interface RolesStackProps extends StackProps {
  stage: string;
}

export class RolesStack extends cdk.Stack {
    readonly ecsTaskRoleArn: string;
  constructor(scope: Construct, id: string, props: RolesStackProps) {
    super(scope, id, props);

    const region = props?.env?.region!

    const ecsExecutionRole = new iam.Role(this, 'ecsExecutionRole', {
      assumedBy: new iam.ServicePrincipal('ecs-tasks.amazonaws.com'),
      description: 'Role for ECS tasks to interact with ECR and other AWS services',
    });
    
    // Add ECR related permissions to the role
    ecsExecutionRole.addToPolicy(new iam.PolicyStatement({
      actions: [
        'ecr:GetAuthorizationToken',
        'ecr:BatchCheckLayerAvailability',
        'ecr:GetDownloadUrlForLayer',
        'ecr:BatchGetImage'
      ],
      resources: ['*'],
    }));
    
    // If you are using specific ECR repositories, replace '*' with specific ARN(s)
    ecsExecutionRole.addToPolicy(new iam.PolicyStatement({
      actions: [
        'ecr:GetDownloadUrlForLayer',
        'ecr:BatchGetImage'
      ],
      resources: ['*'],
    }));

    const ecsTaskRole = new iam.Role(this, 'ecsTaskRole', {
      assumedBy: new iam.ServicePrincipal('ecs-tasks.amazonaws.com'),
      description: 'Role for ECS tasks to interact with ECR and other AWS services',
    });
    
    ecsTaskRole.addToPolicy(new iam.PolicyStatement({
      actions: [
        'dynamodb:Query',
        'dynamodb:GetItem',
        'dynamodb:Scan',
        'dynamodb:PutItem',
        'dynamodb:UpdateItem',
        'dynamodb:DeleteItem'
      ],
      resources: ['*'],
    }));

    ecsTaskRole.addToPolicy(new iam.PolicyStatement({
      actions: [
        's3:ListBucket',
        's3:GetBucketLocation',
        's3:GetObject',
        's3:PutObject',
        's3:DeleteObject',
        's3:ListBucketMultipartUploads',
        's3:AbortMultipartUpload',
        's3:ListMultipartUploadParts'
      ],
      resources: [
        `arn:aws:s3:::${props.stage}-${region}-s3stack-*`, // Bucket-level actions
        `arn:aws:s3:::${props.stage}-${region}-s3stack-*/*` // Object-level actions
      ],
    }));

    this.ecsTaskRoleArn = ecsTaskRole.roleArn;

    new cdk.CfnOutput(this, 'ecsTaskRoleArn', {
        value: ecsTaskRole.roleArn,
        exportName: `ecsTaskRole-${this.region}`,
      });
    new cdk.CfnOutput(this, 'ecsExecutionRoleArn', {
    value: ecsExecutionRole.roleArn,
    exportName: `ecsExecutionRole-${this.region}`,
    });
  }
}
