import * as cdk from 'aws-cdk-lib';
import { Construct } from 'constructs';
import * as ec2 from "aws-cdk-lib/aws-ec2";
import * as ecs from "aws-cdk-lib/aws-ecs";
import * as ecr from "aws-cdk-lib/aws-ecr";
import * as ecs_patterns from "aws-cdk-lib/aws-ecs-patterns";
import * as route53 from 'aws-cdk-lib/aws-route53';
import * as certificatemanager from 'aws-cdk-lib/aws-certificatemanager'
import { aws_logs, StackProps } from 'aws-cdk-lib';
import { VpcStack } from './vpc-stack';
import { disambiguator } from './util/disambiguator';
import Stages from './constants/stages';
import Regions from './constants/regions';
import Accounts from './constants/accounts';
import * as iam from 'aws-cdk-lib/aws-iam';
import * as ssm from 'aws-cdk-lib/aws-secretsmanager';
import { DOMAIN } from './constants/domain';
import Subdomains from './constants/subdomains';
import { ApplicationProtocol } from 'aws-cdk-lib/aws-elasticloadbalancingv2';

interface ECSStackProps extends StackProps {
  imageTag: string;
  stage: string;
  vpcStack: VpcStack;
  bucketName: string;
}

export class EcsStack extends cdk.Stack {
  constructor(scope: Construct, id: string, props: ECSStackProps) {
    super(scope, id, props);

    const MIN_CAPACITY = props.stage === Stages.PROD ? 2 : 0;
    const MAX_CAPACITY = props.stage === Stages.PROD ? 10 : 10;

    const region = props?.env?.region!
    const vpc = props.vpcStack.vpc;

    const ecsExecutionRole = iam.Role.fromRoleArn(this, disambiguator('ecsExecutionRole', props.stage, region), cdk.Fn.importValue(`ecsExecutionRole-${region}`));

    const ecsTaskRole = iam.Role.fromRoleArn(this, disambiguator('ecsTaskRole', props.stage, region), cdk.Fn.importValue(`ecsTaskRole-${this.region}`));
    
    const cluster = new ecs.Cluster(this, disambiguator('Cluster', props.stage, region), {
      vpc: vpc
    });

    const securityGroup = new ec2.SecurityGroup(this, disambiguator('ServiceSecurityGroup', props.stage, region), { vpc });
    securityGroup.addIngressRule(ec2.Peer.anyIpv4(), ec2.Port.tcp(443), 'Allow HTTPS traffic');

    const ecrRepository = ecr.Repository.fromRepositoryArn(this, disambiguator('ServiceRepository', props.stage, region), `arn:aws:ecr:${Regions.US_EAST_1}:${Accounts.DEV}:repository/app`)

    // Create a new Secrets Manager secret to store JWT_SECRET
    const ecsSecret = new ssm.Secret(this, disambiguator('JwtSecret', props.stage, region), {
      secretName: 'plato-secret',
      generateSecretString: {
        secretStringTemplate: JSON.stringify({
          JWT_SECRET: 'CHANGE_ME',
        }),
        generateStringKey: 'unused',
      },
    });

    // Allow the ecs task role to read JWT_SECRET
    ecsSecret.grantRead(ecsTaskRole)

    const subdomain = Accounts.DEV === this.account ? Subdomains.DEV : Subdomains.PROD;
    const domain = `${subdomain}.${DOMAIN}`;

    const zone = new route53.HostedZone(this,  disambiguator('PlatoZone', props.stage, region), {
      zoneName: DOMAIN
    });

    const certificate = new certificatemanager.Certificate(this, 'ServiceCertificate', {
      domainName: domain, // Adjust based on your subdomain and domain logic
      validation: certificatemanager.CertificateValidation.fromDns(zone), // This will handle DNS validation automatically
    });

    // Create a load-balanced Fargate service and make it public
    const fargateService = new ecs_patterns.ApplicationLoadBalancedFargateService(this, disambiguator('PlatoFargateService', props.stage, region), {
      cluster: cluster, // Required
      cpu: 256, // Default is 256
      desiredCount: 1, // Default is 1
      healthCheck: {
         command: [ "CMD-SHELL", "curl -f http://localhost:8080/v1/health || exit 1" ],
         interval: cdk.Duration.seconds(30),
         retries: 5,
         startPeriod: cdk.Duration.seconds(60),
         timeout: cdk.Duration.seconds(5),
      },
      taskImageOptions: { 
        environment: {
          "REGION": region,
          "STAGE": props.stage,
          "PROMPT_BUCKET": props.bucketName,
        },
        secrets: {
          "JWT_SECRET": ecs.Secret.fromSecretsManager(ecsSecret, "JWT_SECRET"),
        },
        taskRole: ecsTaskRole,
        executionRole: ecsExecutionRole,
        image: ecs.ContainerImage.fromRegistry(ecrRepository.repositoryUriForTag(props.imageTag)),
        containerPort: 8080,
        logDriver: ecs.LogDrivers.awsLogs({
          streamPrefix: "ecs",
          logGroup: new aws_logs.LogGroup(this, "LogGroup", {
            logGroupName: "/ecs/PlatoCluster",
            removalPolicy: cdk.RemovalPolicy.DESTROY,
          }),
        }),
       },
      memoryLimitMiB: 512, // Default is 512
      publicLoadBalancer: true, // Default is true,
      securityGroups: [securityGroup],
      domainName: domain,
      domainZone: zone,
      assignPublicIp: true,
      protocol: ApplicationProtocol.HTTPS,
      listenerPort: 443,
      certificate: certificate,
      redirectHTTP: true,
    });

    fargateService.targetGroup.configureHealthCheck({
      path: "/v1/health",
    });

    const scaling = fargateService.service.autoScaleTaskCount({ minCapacity: MIN_CAPACITY, maxCapacity: MAX_CAPACITY});
    scaling.scaleOnCpuUtilization('CpuScaling', {
      targetUtilizationPercent: 70,
      scaleInCooldown: cdk.Duration.minutes(5),
      scaleOutCooldown: cdk.Duration.minutes(5),
    });
  }
}
