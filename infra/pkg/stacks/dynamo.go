package stacks

import (
	"infra/pkg/util"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsdynamodb"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

func CreateDynamoDBStack(scope constructs.Construct, vpc awsec2.IVpc, stage string, region string) awscdk.Resource {
	stack := awscdk.NewStack(scope, jsii.String(util.StageRegionDisambiguator("DynamoDBStack", stage, region)), &awscdk.StackProps{})

	table := awsdynamodb.NewTable(stack, jsii.String("OrganizationsTeamsUsersTable"), &awsdynamodb.TableProps{
		TableName:   jsii.String("OrganizationsTeamsUsers"),
		BillingMode: awsdynamodb.BillingMode_PAY_PER_REQUEST, // Use on-demand scaling
		PartitionKey: &awsdynamodb.Attribute{
			Name: jsii.String("PK"),
			Type: awsdynamodb.AttributeType_STRING,
		},
		SortKey: &awsdynamodb.Attribute{
			Name: jsii.String("SK"),
			Type: awsdynamodb.AttributeType_STRING,
		},
		RemovalPolicy: awscdk.RemovalPolicy_DESTROY, // Useful for development, change to RETAIN for production
	})

	// Create a Global Secondary Index for querying user memberships across teams
	table.AddGlobalSecondaryIndex(&awsdynamodb.GlobalSecondaryIndexProps{
		IndexName: jsii.String("UserMembershipIndex"),
		PartitionKey: &awsdynamodb.Attribute{
			Name: jsii.String("PK"),
			Type: awsdynamodb.AttributeType_STRING,
		},
		SortKey: &awsdynamodb.Attribute{
			Name: jsii.String("SK"),
			Type: awsdynamodb.AttributeType_STRING,
		},
		ProjectionType: awsdynamodb.ProjectionType_ALL,
	})

	return table
}
