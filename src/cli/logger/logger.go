package logger

import (
	"context"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var globalLogger *zap.Logger

// InitLogger initializes the global logger
func InitLogger(level string, format string) error {
	var zapLevel zapcore.Level
	switch level {
	case "debug":
		zapLevel = zapcore.DebugLevel
	case "info":
		zapLevel = zapcore.InfoLevel
	case "warn":
		zapLevel = zapcore.WarnLevel
	case "error":
		zapLevel = zapcore.ErrorLevel
	default:
		zapLevel = zapcore.InfoLevel
	}

	var encoderConfig zapcore.EncoderConfig
	if format == "json" {
		encoderConfig = zap.NewProductionEncoderConfig()
	} else {
		encoderConfig = zap.NewDevelopmentEncoderConfig()
	}

	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	var encoder zapcore.Encoder
	if format == "json" {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	core := zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), zapLevel)
	globalLogger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	return nil
}

// GetLogger returns the global logger
func GetLogger() *zap.Logger {
	if globalLogger == nil {
		// Initialize with defaults if not already done
		InitLogger("info", "json")
	}
	return globalLogger
}

// WithContext returns a logger with context fields
func WithContext(ctx context.Context) *zap.Logger {
	logger := GetLogger()
	
	// Add request ID if available
	if requestID := getRequestID(ctx); requestID != "" {
		logger = logger.With(zap.String("request_id", requestID))
	}
	
	// Add tenant ID if available
	if tenantID := getTenantID(ctx); tenantID != "" {
		logger = logger.With(zap.String("tenant_id", tenantID))
	}
	
	return logger
}

// WithFields returns a logger with additional fields
func WithFields(fields ...zap.Field) *zap.Logger {
	return GetLogger().With(fields...)
}

// getRequestID extracts request ID from context
func getRequestID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if requestID, ok := ctx.Value("request_id").(string); ok {
		return requestID
	}
	return ""
}

// getTenantID extracts tenant ID from context
func getTenantID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if tenantID, ok := ctx.Value("tenant_id").(string); ok {
		return tenantID
	}
	return ""
}

// LogOperation logs the start and completion of an operation
func LogOperation(ctx context.Context, operation string, fields ...zap.Field) func() {
	logger := WithContext(ctx).With(
		append(fields, zap.String("operation", operation))...,
	)
	
	logger.Info("operation started")
	
	start := time.Now()
	return func() {
		duration := time.Since(start)
		logger.Info("operation completed", zap.Duration("duration", duration))
	}
}

// LogError logs an error with context
func LogError(ctx context.Context, msg string, err error, fields ...zap.Field) {
	logger := WithContext(ctx).With(fields...)
	logger.Error(msg, zap.Error(err))
}

// LogSecurity logs security-related events
func LogSecurity(ctx context.Context, event string, fields ...zap.Field) {
	logger := WithContext(ctx).With(
		append(fields, zap.String("security_event", event))...,
	)
	logger.Warn("security event detected")
}
