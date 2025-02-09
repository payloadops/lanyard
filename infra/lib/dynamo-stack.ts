import * as cdk from 'aws-cdk-lib';
import * as dynamodb from 'aws-cdk-lib/aws-dynamodb';
import { Construct } from 'constructs';
import Stages from './constants/stages';
import Regions from './constants/regions';

interface DynamoStackProps extends cdk.StackProps {
  stage: string;
}

const REPLICATIONS_REGIONS: string[] = [
  // Regions.US_WEST_2,
  // Regions.EU_WEST_2,
  // Regions.EU_CENTRAL_1,
];

export class DynamoStack extends cdk.Stack {
  constructor(scope: Construct, id: string, props?: DynamoStackProps) {
    super(scope, id, props);
    new dynamodb.Table(this, 'ServicesTable', {
        tableName: "Services",
        partitionKey: { name: 'pk', type: dynamodb.AttributeType.STRING},
        sortKey: { name: 'sk', type: dynamodb.AttributeType.STRING},
        replicationRegions: props?.stage === Stages.PROD ? REPLICATIONS_REGIONS : undefined,
        billingMode: dynamodb.BillingMode.PAY_PER_REQUEST,
        tableClass: dynamodb.TableClass.STANDARD,
        // removalPolicy: cdk.RemovalPolicy.RETAIN
      })
      
    const apiKeysTable = new dynamodb.Table(this, 'APIKeysTable', {
      tableName: "APIKeys",
      partitionKey: { name: 'pk', type: dynamodb.AttributeType.STRING},
      replicationRegions: props?.stage === Stages.PROD ? REPLICATIONS_REGIONS : undefined,
      billingMode: dynamodb.BillingMode.PAY_PER_REQUEST,
      tableClass: dynamodb.TableClass.STANDARD,
      // removalPolicy: cdk.RemovalPolicy.RETAIN
    })
  }
}
