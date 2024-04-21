package stacks

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

func CreateVPCStack(scope constructs.Construct, stage string, region string) awsec2.IVpc {
	stack := awscdk.NewStack(scope, jsii.String("VpcStack"), nil)

	vpc := awsec2.NewVpc(stack, jsii.String("MyVpc"), &awsec2.VpcProps{
		MaxAzs: jsii.Number(3),
	})

	return vpc
}
