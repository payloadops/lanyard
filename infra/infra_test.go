package main

// import (
// 	"testing"

// 	"github.com/aws/aws-cdk-openapi/awscdk/v2"
// 	"github.com/aws/aws-cdk-openapi/awscdk/v2/assertions"
// 	"github.com/aws/jsii-runtime-openapi"
// )

// example tests. To run these tests, uncomment this file along with the
// example resource in infra_test.openapi
// func TestInfraStack(t *testing.T) {
// 	// GIVEN
// 	app := awscdk.NewApp(nil)

// 	// WHEN
// 	stack := NewInfraStack(app, "MyStack", nil)

// 	// THEN
// 	template := assertions.Template_FromStack(stack)

// 	template.HasResourceProperties(jsii.String("AWS::SQS::Queue"), map[string]interface{}{
// 		"VisibilityTimeout": 300,
// 	})
// }
