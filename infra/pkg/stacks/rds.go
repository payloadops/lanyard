package stacks

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsrds"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

func CreateRdsStack(scope constructs.Construct, vpc awsec2.IVpc, stage string, region string) {

	stack := awscdk.NewStack(scope, jsii.String("PostgresRDSStack"), &awscdk.StackProps{})

	dbInstance := awsrds.NewDatabaseInstance(stack, jsii.String("PostgresDB"), &awsrds.DatabaseInstanceProps{
		Engine: awsrds.DatabaseInstanceEngine_Postgres(&awsrds.PostgresInstanceEngineProps{
			Version: awsrds.PostgresEngineVersion_VER_12_5(),
		}),
		InstanceType: awsec2.InstanceType_Of(awsec2.InstanceClass_BURSTABLE2, awsec2.InstanceSize_MICRO),
		Vpc:          vpc,
		VpcSubnets: &awsec2.SubnetSelection{
			SubnetType: awsec2.SubnetType_PUBLIC,
		},
		MultiAz:                jsii.Bool(false),
		AllocatedStorage:       jsii.Number(20),
		MaxAllocatedStorage:    jsii.Number(100),
		BackupRetention:        awscdk.Duration_Days(jsii.Number(7)),
		DeleteAutomatedBackups: jsii.Bool(false),
		DeletionProtection:     jsii.Bool(true), // Be careful with this in production!
	})

	awscdk.NewCfnOutput(stack, jsii.String("DBEndpoint"), &awscdk.CfnOutputProps{
		Value: dbInstance.DbInstanceEndpointAddress(),
	})
}
