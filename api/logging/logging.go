package logging

import (
	"github.com/payloadops/plato/api/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// NewLogger initializes a Zap logging suitable for the given environment.
func NewLogger(cfg *config.Config) (*zap.Logger, error) {
	level := zap.InfoLevel
	switch cfg.Environment {
	case "local", "development":
		level = zap.DebugLevel
	}

	config := zap.Config{
		Encoding:         "json",
		Level:            zap.NewAtomicLevelAt(level),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig:    getEncoderConfig(zapcore.CapitalLevelEncoder),
	}

	logger, err := config.Build()
	if err != nil {
		return nil, err
	}

	return logger, nil
}

// getEncoderConfig returns a common encoder configuration.
func getEncoderConfig(encodeLevel zapcore.LevelEncoder) zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logging",
		CallerKey:      "caller",
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    encodeLevel,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}
