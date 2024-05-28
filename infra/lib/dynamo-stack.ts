import * as cdk from 'aws-cdk-lib';
import * as dynamodb from 'aws-cdk-lib/aws-dynamodb';
import { Construct } from 'constructs';
import Stages from './constants/stages';

const REPLICATIONS_REGIONS: string[] = [];

export class DynamoStack extends cdk.Stack {
  constructor(scope: Construct, id: string, stage: string, props?: cdk.StackProps) {
    super(scope, id, props);
    new dynamodb.Table(this, 'ProjectsTable', {
        tableName: "projects",
        partitionKey: { name: 'pk', type: dynamodb.AttributeType.STRING},
        sortKey: { name: 'sk', type: dynamodb.AttributeType.STRING},
        replicationRegions: stage === Stages.PROD ? REPLICATIONS_REGIONS : undefined,
        billingMode: dynamodb.BillingMode.PAY_PER_REQUEST,
        tableClass: dynamodb.TableClass.STANDARD
      })

    new dynamodb.Table(this, 'PromptsTable', {
        tableName: "prompts",
        partitionKey: { name: 'pk', type: dynamodb.AttributeType.STRING},
        sortKey: { name: 'sk', type: dynamodb.AttributeType.STRING},
        replicationRegions: REPLICATIONS_REGIONS,
        billingMode: dynamodb.BillingMode.PAY_PER_REQUEST,
        tableClass: dynamodb.TableClass.STANDARD
    })

    new dynamodb.Table(this, 'BranchesTable', {
        tableName: "branches",
        partitionKey: { name: 'pk', type: dynamodb.AttributeType.STRING},
        sortKey: { name: 'sk', type: dynamodb.AttributeType.STRING},
        replicationRegions: REPLICATIONS_REGIONS,
        billingMode: dynamodb.BillingMode.PAY_PER_REQUEST,
        tableClass: dynamodb.TableClass.STANDARD
      })

    new dynamodb.Table(this, 'CommitsTable', {
        tableName: "commits",
        partitionKey: { name: 'pk', type: dynamodb.AttributeType.STRING},
        sortKey: { name: 'sk', type: dynamodb.AttributeType.STRING},
        replicationRegions: REPLICATIONS_REGIONS,
        billingMode: dynamodb.BillingMode.PAY_PER_REQUEST,
        tableClass: dynamodb.TableClass.STANDARD
    })
  }
}
