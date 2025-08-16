// Package logger provides structured logging functionality for the Contexis CLI.
//
// The logger package implements Rails-inspired colored console output with
// structured logging capabilities. It provides both colored and JSON logging
// formats with configurable levels and output destinations.
//
// Key Features:
//   - Rails-like colored console output
//   - Structured JSON logging for production
//   - Configurable log levels and formats
//   - Colored log levels and timestamps
//   - Helper functions for common log operations
//
// Example Usage:
//
//	// Initialize colored logger
//	err := logger.InitColoredLogger("info")
//
//	// Get logger instance
//	log := logger.GetLogger()
//
//	// Log with different levels
//	log.Info("Application started")
//	log.Error("An error occurred", zap.Error(err))
//
//	// Use helper functions
//	logger.LogSuccess("Operation completed successfully")
//	logger.LogWarning("Deprecated feature used")
package logger

import (
	"context"
	"os"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Color constants for terminal output
const (
	// ANSI color codes
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorPurple = "\033[35m"
	ColorCyan   = "\033[36m"
	ColorWhite  = "\033[37m"
	ColorGray   = "\033[90m"

	// Bright colors
	ColorBrightRed    = "\033[91m"
	ColorBrightGreen  = "\033[92m"
	ColorBrightYellow = "\033[93m"
	ColorBrightBlue   = "\033[94m"
	ColorBrightPurple = "\033[95m"
	ColorBrightCyan   = "\033[96m"
	ColorBrightWhite  = "\033[97m"

	// Background colors
	ColorBgRed    = "\033[41m"
	ColorBgGreen  = "\033[42m"
	ColorBgYellow = "\033[43m"
	ColorBgBlue   = "\033[44m"
	ColorBgPurple = "\033[45m"
	ColorBgCyan   = "\033[46m"
	ColorBgWhite  = "\033[47m"
)

// ColoredEncoder provides Rails-like colored console output.
// It implements zapcore.Encoder to provide colored log output
// similar to Rails' console logging with colored levels and timestamps.
type ColoredEncoder struct {
	zapcore.Encoder
}

// NewColoredEncoder creates a new colored encoder.
// It wraps a console encoder with color formatting for
// log levels and other log components.
//
// Returns:
//   - zapcore.Encoder: A colored encoder instance
func NewColoredEncoder() zapcore.Encoder {
	encoderConfig := zap.NewDevelopmentEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	
	encoder := zapcore.NewConsoleEncoder(encoderConfig)
	return &ColoredEncoder{Encoder: encoder}
}

// getColoredLevel returns a colored level string.
// It maps log levels to appropriate colors for better
// visual distinction in console output.
//
// Parameters:
//   - level: The log level to color
//
// Returns:
//   - string: Colored level string with ANSI color codes
func getColoredLevel(level zapcore.Level) string {
	switch level {
	case zapcore.DebugLevel:
		return ColorGray + "DEBUG" + ColorReset
	case zapcore.InfoLevel:
		return ColorBlue + "INFO" + ColorReset
	case zapcore.WarnLevel:
		return ColorYellow + "WARN" + ColorReset
	case zapcore.ErrorLevel:
		return ColorRed + "ERROR" + ColorReset
	case zapcore.DPanicLevel:
		return ColorBrightRed + "DPANIC" + ColorReset
	case zapcore.PanicLevel:
		return ColorBrightRed + "PANIC" + ColorReset
	case zapcore.FatalLevel:
		return ColorBrightRed + "FATAL" + ColorReset
	default:
		return ColorWhite + "UNKNOWN" + ColorReset
	}
}

// InitColoredLogger initializes a colored logger similar to Rails.
// It sets up a logger with colored console output and configurable
// log level. The logger will use colors when outputting to a terminal
// and fall back to plain text when redirected to a file.
//
// Parameters:
//   - level: Log level string (debug, info, warn, error)
//
// Returns:
//   - error: Any error that occurred during initialization
func InitColoredLogger(level string) error {
	// Parse log level
	var zapLevel zapcore.Level
	switch strings.ToLower(level) {
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

	// Check if we should use colors (not in CI, not redirected to file)
	useColors := shouldUseColors()

	var encoder zapcore.Encoder
	if useColors {
		encoder = NewColoredEncoder()
	} else {
		encoderConfig := zap.NewDevelopmentEncoderConfig()
		encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	core := zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), zapLevel)
	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	
	// Replace the global logger
	zap.ReplaceGlobals(logger)
	
	return nil
}

// shouldUseColors determines if we should use colored output.
// It checks various conditions to determine if colored output
// is appropriate, such as terminal type and environment variables.
//
// Returns:
//   - bool: True if colors should be used, false otherwise
func shouldUseColors() bool {
	// Check if we're in a CI environment
	if os.Getenv("CI") != "" {
		return false
	}

	// Check if output is redirected to a file
	stat, err := os.Stdout.Stat()
	if err != nil {
		return false
	}

	// Check if it's a terminal
	return (stat.Mode() & os.ModeCharDevice) != 0
}

// LogOperationColored logs an operation with colored output.
// It provides a convenient way to log operations with consistent
// formatting and color coding.
//
// Parameters:
//   - ctx: Context for request tracking
//   - operation: The operation being performed
//   - fields: Additional structured fields
func LogOperationColored(ctx context.Context, operation string, fields ...zap.Field) func() {
	log := WithContext(ctx)
	log.Info("🔄 "+operation, fields...)
	
	start := time.Now()
	return func() {
		duration := time.Since(start)
		log.Info("✅ "+operation+" completed", zap.Duration("duration", duration))
	}
}

// LogErrorColored logs an error with colored output.
// It provides a convenient way to log errors with consistent
// formatting and color coding.
//
// Parameters:
//   - ctx: Context for request tracking
//   - message: Error message
//   - err: The error object
//   - fields: Additional structured fields
func LogErrorColored(ctx context.Context, message string, err error, fields ...zap.Field) {
	log := WithContext(ctx)
	allFields := append(fields, zap.Error(err))
	log.Error("❌ "+message, allFields...)
}

// LogSecurityColored logs security-related events with colored output.
// It provides a convenient way to log security events with consistent
// formatting and color coding.
//
// Parameters:
//   - event: The security event
//   - details: Additional details about the event
func LogSecurityColored(event, details string) {
	logger := GetLogger()
	logger.Warn("🔒 "+event, zap.String("details", details))
}

// LogSuccess logs a success message with green color.
// It provides a convenient way to log successful operations
// with consistent formatting.
//
// Parameters:
//   - ctx: Context for request tracking
//   - message: Success message
//   - fields: Additional structured fields
func LogSuccess(ctx context.Context, message string, fields ...zap.Field) {
	log := WithContext(ctx)
	log.Info("✅ "+message, fields...)
}

// LogWarning logs a warning message with yellow color.
// It provides a convenient way to log warnings with consistent
// formatting.
//
// Parameters:
//   - message: Warning message
func LogWarning(message string) {
	logger := GetLogger()
	logger.Warn("⚠️ " + message)
}

// LogInfo logs an info message with blue color.
// It provides a convenient way to log informational messages
// with consistent formatting.
//
// Parameters:
//   - ctx: Context for request tracking
//   - message: Info message  
//   - fields: Additional structured fields
func LogInfo(ctx context.Context, message string, fields ...zap.Field) {
	log := WithContext(ctx)
	log.Info("ℹ️ "+message, fields...)
}

// LogDebug logs a debug message with gray color.
// It provides a convenient way to log debug messages with
// consistent formatting.
//
// Parameters:
//   - message: Debug message
func LogDebug(message string) {
	logger := GetLogger()
	logger.Debug("🔍 " + message)
}

// LogDebugWithContext logs a debug message with context support.
func LogDebugWithContext(ctx context.Context, message string, fields ...zap.Field) {
	log := WithContext(ctx)
	log.Debug("🔍 "+message, fields...)
}
