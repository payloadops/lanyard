package main

import (
	awsclient "plato/app/pkg/client/aws"
	dbClient "plato/app/pkg/client/db"
	"plato/app/pkg/router"
)

func main() {
	initializeClients()
	router.Init()
	defer cleanupResources()
}

func initializeClients() {
	dbClient.Init()
	awsclient.InitS3Client()
}

func cleanupResources() {
	dbClient.CleanUp()
}
