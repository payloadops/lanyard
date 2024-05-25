import * as cdk from 'aws-cdk-lib';
import * as elasticache from 'aws-cdk-lib/aws-elasticache';
import { Construct } from 'constructs';
import * as ec2 from 'aws-cdk-lib/aws-ec2';
import { VpcStack } from './vpc-stack';

export class ElastiCacheStack extends cdk.Stack {
    constructor(scope: Construct, id: string, vpcStack: VpcStack, props?: cdk.StackProps) {
      super(scope, id, props);

      const vpc = vpcStack.vpc;
  
      // Create a security group for the Redis cluster
      const sg = new ec2.SecurityGroup(this, 'RedisSecurityGroup', {
        vpc,
        description: 'Allow Redis connection',
        allowAllOutbound: true
      });
      sg.addIngressRule(ec2.Peer.anyIpv4(), ec2.Port.tcp(6379));
  
      // Create a subnet group
      const subnetGroup = new elasticache.CfnSubnetGroup(this, 'RedisSubnetGroup', {
        description: 'Subnet group for Redis',
        subnetIds: vpc.publicSubnets.map(subnet => subnet.subnetId),
        cacheSubnetGroupName: 'plato-redis-subnet-group'
      });
  
      // Create the Redis cluster
      const redisCluster = new elasticache.CfnCacheCluster(this, 'RedisCluster', {
        cacheNodeType: 'cache.t3.micro',
        engine: 'redis',
        numCacheNodes: 1,
        cacheSubnetGroupName: subnetGroup.cacheSubnetGroupName,
        vpcSecurityGroupIds: [sg.securityGroupId]
      });
    }
  }