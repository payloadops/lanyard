package main

import (
	"infra/pkg/stacks"
	"os"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/jsii-runtime-go"
	// "github.com/aws/aws-cdk-go/awscdk/v2/awssqs"
)

type InfraStackProps struct {
	awscdk.StackProps
}

func main() {
	stage := os.Getenv("STAGE")
	app := awscdk.NewApp(nil)

	vpc := stacks.CreateVPCStack(app, stage, *env().Region)
	stacks.CreateEcsFargateAPI(app, vpc, stage, "us-east-1")

	app.Synth(nil)
}

func env() *awscdk.Environment {
	account := os.Getenv("CDK_DEFAULT_ACCOUNT")
	region := os.Getenv("CDK_DEFAULT_REGION")
	return &awscdk.Environment{
		Account: jsii.String(account),
		Region:  jsii.String(region),
	}
}
