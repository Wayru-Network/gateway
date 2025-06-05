package infra

import (
	"errors"

	"github.com/Wayru-Network/serve/middleware"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// ZapLoggerAdapter adapts zap.Logger to middleware.Logger interface
type ZapLoggerAdapter struct {
	logger *zap.Logger
}

func (l *ZapLoggerAdapter) Debug(msg string) {
	l.logger.Debug(msg)
}

func (l *ZapLoggerAdapter) Info(msg string) {
	l.logger.Info(msg)
}

func (l *ZapLoggerAdapter) Error(msg string) {
	l.logger.Error(msg)
}

// NewZapLoggerAdapter creates a new ZapLoggerAdapter from a zap.Logger
func NewZapLoggerAdapter(logger *zap.Logger) *ZapLoggerAdapter {
	return &ZapLoggerAdapter{logger: logger}
}

// InitLogger initializes and replaces the global zap logger based on the application environment.
// It returns the created *zap.Logger for further use.
func InitLogger(appEnv string) (*zap.Logger, error) {
	var cfg zap.Config
	switch appEnv {
	case "local":
		cfg = zap.NewDevelopmentConfig()
	case "dev":
		cfg = zap.NewProductionConfig()
		cfg.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	case "prod":
		cfg = zap.NewProductionConfig()
	default:
		return nil, errors.New("unknown APP_ENV: " + appEnv)
	}

	// Use ISO8601 for human-readable timestamps
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	logger, err := cfg.Build(zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	if err != nil {
		return nil, err
	}

	zap.ReplaceGlobals(logger)

	// If in local environment, also set the global Sugar logger
	if appEnv == "local" {
		zap.S() // This initializes the global sugared logger
	}

	return logger, nil
}

// ConfigureServeLogger configures the serve middleware to use our zap logger
func ConfigureServeLogger(logger *zap.Logger) {
	zapLogger := NewZapLoggerAdapter(logger)
	middleware.SetLogger(zapLogger)
}

// Sync flushes any buffered log entries; should be called before the process exits.
func Sync() {
	_ = zap.L().Sync()
}
