package logging

import (
	"testing"

	"github.com/payloadops/lanyard/app/config"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

// TestNewLoggerLocal tests the NewLogger function for the local environment.
func TestNewLoggerLocal(t *testing.T) {
	cfg := &config.Config{
		Environment: "local",
	}

	logger, err := NewLogger(cfg)
	assert.NoError(t, err)
	assert.NotNil(t, logger)
	assert.True(t, logger.Core().Enabled(zap.DebugLevel))

	logger.Debug("This is a debug log for local environment")
	logger.Info("This is an info log for local environment")
}

// TestNewLoggerDevelopment tests the NewLogger function for the development environment.
func TestNewLoggerDevelopment(t *testing.T) {
	cfg := &config.Config{
		Environment: "development",
	}

	logger, err := NewLogger(cfg)
	assert.NoError(t, err)
	assert.NotNil(t, logger)
	assert.True(t, logger.Core().Enabled(zap.DebugLevel))

	logger.Debug("This is a debug log for development environment")
	logger.Info("This is an info log for development environment")
}

// TestNewLoggerProduction tests the NewLogger function for the production environment.
func TestNewLoggerProduction(t *testing.T) {
	cfg := &config.Config{
		Environment: "production",
	}

	logger, err := NewLogger(cfg)
	assert.NoError(t, err)
	assert.NotNil(t, logger)
	assert.False(t, logger.Core().Enabled(zap.DebugLevel))
	assert.True(t, logger.Core().Enabled(zap.InfoLevel))

	// This is an expensive mistake.
	logger.Debug("This debug log should not appear in production environment")
	logger.Info("This is an info log for production environment")
}

// TestNewLoggerDefault tests the NewLogger function with an unspecified environment.
func TestNewLoggerDefault(t *testing.T) {
	cfg := &config.Config{
		Environment: "unspecified",
	}

	logger, err := NewLogger(cfg)
	assert.NoError(t, err)
	assert.NotNil(t, logger)
	assert.False(t, logger.Core().Enabled(zap.DebugLevel))
	assert.True(t, logger.Core().Enabled(zap.InfoLevel))

	logger.Debug("This debug log should not appear in default environment")
	logger.Info("This is an info log for default environment")
}
