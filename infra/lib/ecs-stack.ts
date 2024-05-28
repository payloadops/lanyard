import * as cdk from 'aws-cdk-lib';
import { Construct } from 'constructs';
import * as ec2 from "aws-cdk-lib/aws-ec2";
import * as ecs from "aws-cdk-lib/aws-ecs";
import * as ecr from "aws-cdk-lib/aws-ecr";
import * as ecs_patterns from "aws-cdk-lib/aws-ecs-patterns";
import * as route53 from 'aws-cdk-lib/aws-route53';
import * as certificatemanager from 'aws-cdk-lib/aws-certificatemanager'
import { aws_logs } from 'aws-cdk-lib';
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


export class EcsStack extends cdk.Stack {
  constructor(scope: Construct, id: string, vpcStack: VpcStack, stage: string, props?: cdk.StackProps) {
    super(scope, id, props);

    const region = props?.env?.region!
    const vpc = vpcStack.vpc;

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
    
    // Add ECR related permissions to the role
    ecsTaskRole.addToPolicy(new iam.PolicyStatement({
      actions: [
        's3:*',
        'dynamo:*'
      ],
      resources: ['*'],
    }));
    
    const cluster = new ecs.Cluster(this, disambiguator('Cluster', stage, region), {
      vpc: vpc
    });

    const securityGroup = new ec2.SecurityGroup(this, disambiguator('ServiceSecurityGroup', stage, region), { vpc });
    securityGroup.addIngressRule(ec2.Peer.anyIpv4(), ec2.Port.tcp(443), 'Allow HTTPS traffic');

    const ecrRepository = ecr.Repository.fromRepositoryArn(this, disambiguator('ServiceRepository', stage, region), `arn:aws:ecr:${Regions.US_EAST_1}:${Accounts.DEV}:repository/app`)

    // Create a new Secrets Manager secret to store JWT_SECRET
    const ecsSecret = new ssm.Secret(this, disambiguator('JwtSecret', stage, region), {
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

    const zone = new route53.HostedZone(this,  disambiguator('PlatoZone', stage, region), {
      zoneName: DOMAIN
    });

    const certificate = new certificatemanager.Certificate(this, 'ServiceCertificate', {
      domainName: domain, // Adjust based on your subdomain and domain logic
      validation: certificatemanager.CertificateValidation.fromDns(zone), // This will handle DNS validation automatically
    });

    // Create a load-balanced Fargate service and make it public
    const fargateService = new ecs_patterns.ApplicationLoadBalancedFargateService(this, disambiguator('PlatoFargateService', stage, region), {
      cluster: cluster, // Required
      cpu: 256, // Default is 256
      desiredCount: 2, // Default is 1
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
          "STAGE": stage,
        },
        secrets: {
          "JWT_SECRET": ecs.Secret.fromSecretsManager(ecsSecret, "JWT_SECRET"),
        },
        taskRole: ecsTaskRole,
        executionRole: ecsExecutionRole,
        image: ecs.ContainerImage.fromRegistry(ecrRepository.repositoryUriForTag()),
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
      protocol: ApplicationProtocol.HTTPS,
      listenerPort: 443,
      certificate: certificate,
      redirectHTTP: true,
    });

    fargateService.targetGroup.configureHealthCheck({
      path: "/v1/health",
    });

    const scaling = fargateService.service.autoScaleTaskCount({ minCapacity: 2, maxCapacity: 10 });
    scaling.scaleOnCpuUtilization('CpuScaling', {
      targetUtilizationPercent: 70,
      scaleInCooldown: cdk.Duration.minutes(10),
      scaleOutCooldown: cdk.Duration.minutes(10),
    });
  }
}
