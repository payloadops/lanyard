import * as cdk from 'aws-cdk-lib';
import { Construct } from "constructs";
import { VpcStack } from '../lib/vpc-stack';
import { DynamoStack } from '../lib/dynamo-stack';
import { EcsStack } from '../lib/ecs-stack';
import Regions from '../lib/constants/regions';
import { disambiguator } from '../lib/util/disambiguator';
import { S3Stack } from './s3-stack';


export class Stage extends cdk.Stage {
    constructor(scope: Construct, id: string, stage: string, props?: cdk.StageProps) {
      super(scope, id, props);
      const region = props?.env?.region!;
      let bucketName: string | undefined = undefined;

      const vpcStack = new VpcStack(this, disambiguator('PlatoVpcStack', stage, region), stage, {
        env: { account: props?.env?.account, region: props?.env?.region },
      });

      if (props?.env?.region === Regions.US_EAST_1) {
        new DynamoStack(this, disambiguator('DynamoStack', stage, region), stage, {
            env: { account: props?.env?.account, region: props?.env?.region },
          })
        const s3Stack = new S3Stack(this, disambiguator('S3Stack', stage, region), stage, {
            env: { account: props?.env?.account, region: props?.env?.region },
        });
        bucketName = s3Stack.bucketName
      }

      if (!bucketName) {
        bucketName = cdk.Fn.importValue('BucketName-us-east-1');
      }
      
    //   new ElastiCacheStack(scope, disambiguator('ElasticacheStack'), vpcStack, {
    //     env: { account: props?.env?.account, region: props?.env?.region },
    //   });
      
      new EcsStack(this, disambiguator('EcsStack', stage, region), vpcStack, stage, bucketName, {
        env: { account: props?.env?.account, region: props?.env?.region },
      });
    }
}