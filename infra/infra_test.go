package main

// import (
// 	"testing"

// 	"github.com/client/client-cdk-go/awscdk/v2"
// 	"github.com/client/client-cdk-go/awscdk/v2/assertions"
// 	"github.com/client/jsii-runtime-go"
// )

// example tests. To run these tests, uncomment this file along with the
// example resource in infra_test.go
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
