import * as cdk from 'aws-cdk-lib';
import { Construct } from "constructs";
import { VpcStack } from '../lib/vpc-stack';
import { DynamoStack } from '../lib/dynamo-stack';
import { EcsStack } from '../lib/ecs-stack';
import { ElastiCacheStack } from '../lib/elasticache-stack';
import Regions from '../lib/constants/regions';
import { disambiguator } from '../lib/util/disambiguator';

export class Stage extends cdk.Stage {
    public readonly ecsStack: EcsStack
    constructor(scope: Construct, id: string, stage: string, props?: cdk.StageProps) {
      super(scope, id, props);
      const region = props?.env?.region!;
      const vpcStack = new VpcStack(this, disambiguator('PlatoVpcStack', stage, region), stage, {
        env: { account: props?.env?.account, region: props?.env?.region },
      });
      
    //   new ElastiCacheStack(scope, disambiguator('ElasticacheStack'), vpcStack, {
    //     env: { account: props?.env?.account, region: props?.env?.region },
    //   });
      
      let ecsStack = new EcsStack(this, disambiguator('EcsStack', stage, region), vpcStack, stage, {
        env: { account: props?.env?.account, region: props?.env?.region },
      });

      this.ecsStack = ecsStack
      if (props?.env?.region === Regions.US_EAST_1) {
        new DynamoStack(this, disambiguator('DynamoStack', stage, region), stage, {
            env: { account: props?.env?.account, region: props?.env?.region },
          })
      }
    }
}