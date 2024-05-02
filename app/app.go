package main

import (
	awsclient "plato/app/pkg/client/aws"
	"plato/app/pkg/router"
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
