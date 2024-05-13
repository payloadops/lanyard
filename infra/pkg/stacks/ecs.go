package stacks

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awscertificatemanager"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsecr"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsecs"
	"github.com/aws/aws-cdk-go/awscdk/v2/awselasticloadbalancingv2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsroute53"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsroute53targets"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	constants "github.com/payloadops/plato/infra/pkg/const"
)

func CreateEcsFargateAPI(scope constructs.Construct, vpc awsec2.IVpc, stage string, region string) awsecs.IFargateService {
	stack := awscdk.NewStack(scope, jsii.String("ECSApiGatewayStack"), nil)

	taskRole := awsiam.NewRole(stack, jsii.String("ecsTaskRole"), &awsiam.RoleProps{
		AssumedBy: awsiam.NewServicePrincipal(jsii.String("ecs-tasks.amazonaws.com"), &awsiam.ServicePrincipalOpts{}),
		ManagedPolicies: &[]awsiam.IManagedPolicy{
			awsiam.ManagedPolicy_FromAwsManagedPolicyName(jsii.String("service-role/AmazonECSTaskExecutionRolePolicy")),
		},
	})

	taskRole.AddToPolicy(awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
		Actions:   &[]*string{jsii.String("s3:*"), jsii.String("dynamodb:*")},
		Resources: &[]*string{jsii.String("arn:aws:s3:::*/*"), jsii.String("arn:aws:dynamodb:::*/*")},
	}))

	repository := awsecr.NewRepository(stack, jsii.String("MyRepository"), &awsecr.RepositoryProps{
		RepositoryName: jsii.String("my-ecr-repo"),
		RemovalPolicy:  awscdk.RemovalPolicy_DESTROY, // Change to RETAIN for production
	})

	// Create an ECS cluster
	cluster := awsecs.NewCluster(stack, jsii.String("EcsCluster"), &awsecs.ClusterProps{
		Vpc: vpc,
	})

	// Create a task definition
	taskDefinition := awsecs.NewTaskDefinition(stack, jsii.String("TaskDef"), &awsecs.TaskDefinitionProps{
		Compatibility: awsecs.Compatibility_FARGATE,
		Cpu:           jsii.String("256"), // 0.25 vCPU
		MemoryMiB:     jsii.String("512"), // 512 MiB
		TaskRole:      taskRole,           // Assign the created IAM role to the task
	})

	// Add a container to the task definition
	container := taskDefinition.AddContainer(jsii.String("WebContainer"), &awsecs.ContainerDefinitionOptions{
		Image:          awsecs.ContainerImage_FromEcrRepository(repository, jsii.String("latest")),
		MemoryLimitMiB: jsii.Number(512),
		Environment: &map[string]*string{
			"PORT": jsii.String("80"),
		},
	})

	// Map port 80 on the container to port 80 on the host
	container.AddPortMappings(&awsecs.PortMapping{
		ContainerPort: jsii.Number(80),
		Protocol:      awsecs.Protocol_TCP,
	})

	service := awsecs.NewFargateService(stack, jsii.String("Service"), &awsecs.FargateServiceProps{
		Cluster:        cluster,
		TaskDefinition: taskDefinition,
		DesiredCount:   jsii.Number(1),
	})

	hostedZone := awsroute53.NewHostedZone(stack, jsii.String("MyHostedZone"), &awsroute53.HostedZoneProps{
		ZoneName: jsii.String(constants.DOMAIN),
	})

	certificate := awscertificatemanager.NewCertificate(stack, jsii.String("MyCertificate"), &awscertificatemanager.CertificateProps{
		DomainName: jsii.String(constants.DOMAIN),
		Validation: awscertificatemanager.CertificateValidation_FromDns(hostedZone),
	})

	lb := awselasticloadbalancingv2.NewApplicationLoadBalancer(stack, jsii.String("LB"), &awselasticloadbalancingv2.ApplicationLoadBalancerProps{
		Vpc:            vpc,
		InternetFacing: jsii.Bool(true),
	})

	awsroute53.NewARecord(stack, jsii.String("DNSRecord"), &awsroute53.ARecordProps{
		Zone:       hostedZone,
		Target:     awsroute53.RecordTarget_FromAlias(awsroute53targets.NewLoadBalancerTarget(lb)),
		RecordName: jsii.String(constants.DOMAIN),
	})

	listener := lb.AddListener(jsii.String("HttpsListener"), &awselasticloadbalancingv2.BaseApplicationListenerProps{
		Port: jsii.Number(443),
		Certificates: &[]awselasticloadbalancingv2.IListenerCertificate{
			awselasticloadbalancingv2.ListenerCertificate_FromCertificateManager(certificate),
		},
	})

	listener.AddTargets(jsii.String("ECSTargets"), &awselasticloadbalancingv2.AddApplicationTargetsProps{
		Port: jsii.Number(80),
		Targets: &[]awselasticloadbalancingv2.IApplicationLoadBalancerTarget{
			service,
		},
	})

	return service
}
