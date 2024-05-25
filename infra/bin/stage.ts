import * as cdk from 'aws-cdk-lib';
import { Construct } from "constructs";
import { VpcStack } from '../lib/vpc-stack';
import { DynamoStack } from '../lib/dynamo-stack';
import { EcsStack } from '../lib/ecs-stack';
import { ElastiCacheStack } from '../lib/elasticache-stack';

export class Stage extends cdk.Stage {

    constructor(scope: Construct, id: string, props?: cdk.StageProps) {
      super(scope, id, props);

      const vpcStack = new VpcStack(scope, 'PlatoVpcStack', {
        env: { account: props?.env?.account, region: props?.env?.region },
      });
      
    //   new ElastiCacheStack(scope, 'ElasticacheStack', vpcStack, {
    //     env: { account: props?.env?.account, region: props?.env?.region },
    //   });
      
      new EcsStack(scope, 'EcsStack', vpcStack, {
        env: { account: props?.env?.account, region: props?.env?.region },
      });
      
      new DynamoStack(scope, 'DynamoStack', {
        env: { account: props?.env?.account, region: props?.env?.region },
      })
    }
}