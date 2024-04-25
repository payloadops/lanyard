package stacks

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigatewayv2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsecr"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsecs"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

func CreateEcsFargateAPI(scope constructs.Construct, vpc awsec2.IVpc, stage string, region string) awsecs.IFargateService {
	stack := awscdk.NewStack(scope, jsii.String("ECSApiGatewayStack"), nil)
	cluster := awsecs.NewCluster(stack, jsii.String("EcsCluster"), &awsecs.ClusterProps{})

	repository := awsecr.NewRepository(stack, jsii.String("MyRepository"), &awsecr.RepositoryProps{
		RepositoryName: jsii.String("my-ecr-repo"),
		RemovalPolicy:  awscdk.RemovalPolicy_DESTROY, // Change to RETAIN for production
	})

	// awsec2.NewInterfaceVpcEndpoint(stack, jsii.String("EcrApiVpcEndpoint"), &awsec2.InterfaceVpcEndpointProps{
	// 	Service:           awsec2.InterfaceVpcEndpointAwsService_ECR(),
	// 	PrivateDnsEnabled: jsii.Bool(true),
	// 	Subnets:           &awsec2.SubnetSelection{SubnetType: awsec2.SubnetType_PRIVATE_ISOLATED},
	// })

	// // ECR VPC Endpoint for Docker image pulling (ECR_DOCKER)
	// awsec2.NewInterfaceVpcEndpoint(stack, jsii.String("EcrDockerVpcEndpoint"), &awsec2.InterfaceVpcEndpointProps{
	// 	Service:           awsec2.InterfaceVpcEndpointAwsService_ECR_DOCKER(),
	// 	PrivateDnsEnabled: jsii.Bool(true),
	// 	Subnets:           &awsec2.SubnetSelection{SubnetType: awsec2.SubnetType_PRIVATE_ISOLATED},
	// })

	// Define the task definition with a single container using an image from ECR
	taskDefinition := awsecs.NewFargateTaskDefinition(stack, jsii.String("TaskDef"), nil)
	container := taskDefinition.AddContainer(jsii.String("web"), &awsecs.ContainerDefinitionOptions{
		Image:          awsecs.ContainerImage_FromEcrRepository(repository, jsii.String("latest")), // Assuming 'latest' tag is used
		MemoryLimitMiB: jsii.Number(512),
	})
	container.AddPortMappings(&awsecs.PortMapping{
		ContainerPort: jsii.Number(80),
	})

	service := awsecs.NewFargateService(stack, jsii.String("Service"), &awsecs.FargateServiceProps{
		Cluster:        cluster,
		TaskDefinition: taskDefinition,
		DesiredCount:   jsii.Number(1),
	})

	// lb := awselasticloadbalancingv2.NewApplicationLoadBalancer(stack, jsii.String("LB"), &awselasticloadbalancingv2.ApplicationLoadBalancerProps{
	// 	InternetFacing: jsii.Bool(true),
	// })

	// listener := lb.AddListener(jsii.String("Listener"), &awselasticloadbalancingv2.BaseApplicationListenerProps{Port: jsii.Number(80)})
	// listener.AddTargets(jsii.String("Target"), &awselasticloadbalancingv2.AddApplicationTargetsProps{
	// 	Targets: &[]awselasticloadbalancingv2.IApplicationLoadBalancerTarget{service},
	// 	Port:    jsii.Number(80),
	// })

	awsapigatewayv2.NewCfnApi(stack, jsii.String("HttpApi"), &awsapigatewayv2.CfnApiProps{
		Name:         jsii.String("MyAPI"),
		ProtocolType: jsii.String("HTTP"),
	})

	// integration := awsapigatewayv2.NewCfnIntegration(stack, jsii.String("ApiIntegration"), &awsapigatewayv2.CfnIntegrationProps{
	// 	ApiId:                api.AttrApiId(),
	// 	IntegrationType:      jsii.String("HTTP_PROXY"),
	// 	IntegrationUri:       listener.ListenerArn(),
	// 	IntegrationMethod:    jsii.String("ANY"),
	// 	PayloadFormatVersion: jsii.String("1.0"),
	// })

	// awsapigatewayv2.NewCfnRoute(stack, jsii.String("ApiRoute"), &awsapigatewayv2.CfnRouteProps{
	// 	ApiId:    api.AttrApiId(),
	// 	RouteKey: jsii.String("ANY /{proxy+}"),
	// 	Target:   jsii.String("integrations/" + *integration.ApiId()),
	// })

	// awsapigatewayv2.NewCfnStage(stack, jsii.String("ApiStage"), &awsapigatewayv2.CfnStageProps{
	// 	ApiId:      api.AttrApiId(),
	// 	StageName:  jsii.String(stage),
	// 	AutoDeploy: jsii.Bool(true),
	// 	DefaultRouteSettings: &awsapigatewayv2.CfnStage_RouteSettingsProperty{
	// 		ThrottlingRateLimit:  jsii.Number(10), // max requests per second
	// 		ThrottlingBurstLimit: jsii.Number(20), // max concurrent requests
	// 	},
	// })

	return service
}
