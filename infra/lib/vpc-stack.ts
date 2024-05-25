import { Stack, StackProps } from 'aws-cdk-lib';
import { Vpc } from 'aws-cdk-lib/aws-ec2';
import { Construct } from 'constructs';
import * as ec2 from 'aws-cdk-lib/aws-ec2';
import { disambiguator } from './util/disambiguator';

export class VpcStack extends Stack {
  public readonly vpc: Vpc;

  constructor(scope: Construct, id: string, stage: string, props?: StackProps) {
    super(scope, id, props);
    const region = props?.env?.region!

    this.vpc = new ec2.Vpc(this, disambiguator("PlatoVpc", stage, region), {
        maxAzs: 2 // Default is all AZs in region
    });
  }
}
