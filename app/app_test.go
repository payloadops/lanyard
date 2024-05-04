package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"os"
	"testing"
)

func TestInitLogger(t *testing.T) {
	app := &App{}
	err := app.initLogger()

	assert.NoError(t, err, "initLogger should not return an error")
	assert.NotNil(t, app.logger, "Logger should be initialized")
}

func TestInitializeClients(t *testing.T) {
	// Mock AWS client initialization
	/*
		originalDynamoClient := awsclient.InitDynamoClient
		originalS3Client := awsclient.InitS3Client

		defer func() {
			awsclient.InitDynamoClient = originalDynamoClient
			awsclient.InitS3Client = originalS3Client
		}()

		awsclient.InitDynamoClient = func() (*awsclient.DynamoClient, error) {
			return nil, nil
		}
		awsclient.InitS3Client = func() (*awsclient.S3Client, error) {
			return nil, nil
		}
	*/

	app := &App{}
	err := app.initializeClients()

	assert.NoError(t, err, "initializeClients should not return an error")
}

func TestCreateRouter(t *testing.T) {
	app := &App{}
	router := app.createRouter()

	assert.NotNil(t, router, "Router should not be nil")
}

func TestStartHTTPServer_DefaultBindAddress(t *testing.T) {
	app := &App{
		logger: zap.NewNop(),
	}

	// Test with the default bind address
	os.Unsetenv("BIND_ADDRESS")
	defer func() {
		app.server.Close()
	}()

	assert.Nil(t, app.server, "HTTP server should initially be nil")

	app.startHTTPServer(chi.NewMux())
	assert.NotNil(t, app.server, "HTTP server should be initialized after startHTTPServer")
	assert.Equal(t, DefaultBindAddress, app.server.Addr, "Server address should default to DefaultBindAddress")
}
