package main

import (
	"context"
	"fmt"
	"go.uber.org/zap/zapcore"
	"net/http"
	"os"
	"os/signal"
	"plato/app/pkg/router"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	awsclient "plato/app/pkg/client/aws"
)

// DefaultBindAddress represents the default host and port to host the server on.
const DefaultBindAddress = ":8080"

// GracefulShutdownTimeout represents the maximum amount of time the server will wait for requests to complete.
const GracefulShutdownTimeout = 5 * time.Second

// App holds the application's runtime context, including logging and services.
type App struct {
	logger *zap.Logger
	server *http.Server
}

func main() {
	app := &App{}
	if err := app.initLogger(); err != nil {
		fmt.Printf("Error initializing zap logger: %v\n", err)
		os.Exit(1)
	}

	defer app.cleanupResources()
	if err := app.initializeClients(); err != nil {
		fmt.Printf("Error initializing AWS clients: %v\n", err)
		os.Exit(1)
	}

	r := app.createRouter()
	app.startHTTPServer(r)
	app.gracefulShutdown()
}

func (app *App) initLogger() error {
	config := zap.NewProductionConfig()
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.OutputPaths = []string{"stdout"}
	config.ErrorOutputPaths = []string{"stderr"}

	logger, err := config.Build()
	if err != nil {
		return err
	}

	app.logger = logger
	return nil
}

func (app *App) initializeClients() error {
	if _, err := awsclient.InitDynamoClient(); err != nil {
		return fmt.Errorf("failed to initialize DynamoDB client: %w", err)
	}

	if _, err := awsclient.InitS3Client(); err != nil {
		return fmt.Errorf("failed to initialize S3 client: %w", err)
	}

	return nil
}

func (app *App) cleanupResources() {
	app.logger.Info("Cleaning up resources...")
	_ = app.logger.Sync()
}

func (app *App) createRouter() *chi.Mux {
	return router.NewRouter()
}

func (app *App) startHTTPServer(r *chi.Mux) {
	bindAddress := os.Getenv("BIND_ADDRESS")
	if bindAddress == "" {
		bindAddress = DefaultBindAddress
	}

	server := &http.Server{
		Addr:    bindAddress,
		Handler: r,
	}

	go func() {
		app.logger.Info("Starting server on port 8080")
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			app.logger.Fatal("Error starting HTTP server", zap.Error(err))
		}
	}()

	app.server = server
}

func (app *App) gracefulShutdown() {
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

	<-stopChan // wait for SIGINT
	app.logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), GracefulShutdownTimeout)
	defer cancel()

	if err := app.server.Shutdown(ctx); err != nil {
		app.logger.Fatal("Server shutdown failed", zap.Error(err))
	}

	app.logger.Info("Server gracefully stopped")
}
