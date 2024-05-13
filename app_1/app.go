package main

import (
	"plato/app/pkg/router"
	awsclient "plato/app_1/pkg/client/aws"
)

func main() {
	initializeClients()
	router.Init()
	defer cleanupResources()
}

func initializeClients() {
	awsclient.InitDynamoClient()
	awsclient.InitS3Client()
}

func cleanupResources() {
}
