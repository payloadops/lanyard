import * as cdk from 'aws-cdk-lib';
import { Construct } from 'constructs';
import * as ec2 from "aws-cdk-lib/aws-ec2";
import * as ecs from "aws-cdk-lib/aws-ecs";
import * as ecr from "aws-cdk-lib/aws-ecr";
import * as ecs_patterns from "aws-cdk-lib/aws-ecs-patterns";
import * as route53 from 'aws-cdk-lib/aws-route53';
import * as route53Targets from 'aws-cdk-lib/aws-route53-targets';
import { aws_logs } from 'aws-cdk-lib';
import { VpcStack } from './vpc-stack';


export class EcsStack extends cdk.Stack {
  constructor(scope: Construct, id: string, vpcStack: VpcStack, props?: cdk.StackProps) {
    super(scope, id, props);

    const vpc = vpcStack.vpc;

    const cluster = new ecs.Cluster(this, "PlatoCluster", {
      vpc: vpc
    });

    const repository = new ecr.Repository(this, 'Repository');

    const securityGroup = new ec2.SecurityGroup(this, 'ServiceSecurityGroup', { vpc });
    securityGroup.addIngressRule(ec2.Peer.anyIpv4(), ec2.Port.tcp(80), 'Allow HTTP traffic');
    securityGroup.addIngressRule(ec2.Peer.anyIpv4(), ec2.Port.tcp(443), 'Allow HTTPS traffic');

    // Create a load-balanced Fargate service and make it public
    const fargateService = new ecs_patterns.ApplicationLoadBalancedFargateService(this, "PlatoFargateService", {
      cluster: cluster, // Required
      cpu: 256, // Default is 256
      desiredCount: 1, // Default is 1
      taskImageOptions: { 
        image: ecs.ContainerImage.fromEcrRepository(repository),
        containerPort: 8080,
        // logDriver: ecs.LogDrivers.awsLogs({
        //   streamPrefix: "ecs",
        //   logGroup: new aws_logs.LogGroup(this, "LogGroup", {
        //     logGroupName: "/ecs/PlatoCluster",
        //     removalPolicy: cdk.RemovalPolicy.DESTROY,
        //   }),
        // }),
       },
      memoryLimitMiB: 512, // Default is 512
      publicLoadBalancer: true, // Default is true,
      securityGroups: [securityGroup]
    });

    const scaling = fargateService.service.autoScaleTaskCount({ maxCapacity: 10 });
    scaling.scaleOnCpuUtilization('CpuScaling', {
      targetUtilizationPercent: 70,
      scaleInCooldown: cdk.Duration.minutes(10),
      scaleOutCooldown: cdk.Duration.minutes(10),
    });

    const zone = new route53.HostedZone(this, 'PlatoZone', {
      zoneName: "payloadops.com"
    });

    // Create a subdomain A record for the API pointing to the ALB
    new route53.ARecord(this, 'ApiAliasRecord', {
      zone: zone,
      recordName: 'api',  // Subdomain
      target: route53.RecordTarget.fromAlias(new route53Targets.LoadBalancerTarget(fargateService.loadBalancer)),
    });
  }
}
