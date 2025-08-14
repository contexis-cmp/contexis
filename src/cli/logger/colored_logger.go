package logger

import (
	"context"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Color constants for terminal output
const (
	// ANSI color codes
	ColorReset   = "\033[0m"
	ColorRed     = "\033[31m"
	ColorGreen   = "\033[32m"
	ColorYellow  = "\033[33m"
	ColorBlue    = "\033[34m"
	ColorMagenta = "\033[35m"
	ColorCyan    = "\033[36m"
	ColorWhite   = "\033[37m"

	// Bright colors
	ColorBrightRed    = "\033[91m"
	ColorBrightGreen  = "\033[92m"
	ColorBrightYellow = "\033[93m"
	ColorBrightBlue   = "\033[94m"
	ColorBrightCyan   = "\033[96m"

	// Background colors
	ColorBgRed     = "\033[41m"
	ColorBgGreen   = "\033[42m"
	ColorBgYellow  = "\033[43m"
	ColorBgBlue    = "\033[44m"
	ColorBgMagenta = "\033[45m"
	ColorBgCyan    = "\033[46m"
)

// ColoredEncoder provides Rails-like colored console output
type ColoredEncoder struct {
	zapcore.Encoder
	colored bool
}

// NewColoredEncoder creates a new colored encoder
func NewColoredEncoder(colored bool) zapcore.Encoder {
	config := zap.NewDevelopmentEncoderConfig()
	config.TimeKey = "timestamp"
	config.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("15:04:05"))
	}
	config.EncodeLevel = func(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
		if colored {
			enc.AppendString(getColoredLevel(l))
		} else {
			enc.AppendString(l.CapitalString())
		}
	}
	config.EncodeCaller = func(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
		if colored {
			enc.AppendString(ColorCyan + caller.TrimmedPath() + ColorReset)
		} else {
			enc.AppendString(caller.TrimmedPath())
		}
	}

	return zapcore.NewConsoleEncoder(config)
}

// getColoredLevel returns a colored level string
func getColoredLevel(level zapcore.Level) string {
	switch level {
	case zapcore.DebugLevel:
		return ColorBrightCyan + "DEBUG" + ColorReset
	case zapcore.InfoLevel:
		return ColorBrightGreen + "INFO " + ColorReset
	case zapcore.WarnLevel:
		return ColorBrightYellow + "WARN " + ColorReset
	case zapcore.ErrorLevel:
		return ColorBrightRed + "ERROR" + ColorReset
	case zapcore.FatalLevel:
		return ColorRed + ColorBgRed + ColorWhite + "FATAL" + ColorReset
	case zapcore.PanicLevel:
		return ColorRed + ColorBgRed + ColorWhite + "PANIC" + ColorReset
	default:
		return ColorWhite + level.CapitalString() + ColorReset
	}
}

// InitColoredLogger initializes a colored logger similar to Rails
func InitColoredLogger(level string) error {
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

	// Check if we should use colors (not in CI, not redirected to file)
	colored := shouldUseColors()

	encoder := NewColoredEncoder(colored)
	core := zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), zapLevel)

	globalLogger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	return nil
}

// shouldUseColors determines if we should use colored output
func shouldUseColors() bool {
	// Check if we're in a CI environment
	if os.Getenv("CI") != "" || os.Getenv("GITHUB_ACTIONS") != "" {
		return false
	}

	// Check if output is redirected to a file
	stat, err := os.Stdout.Stat()
	if err != nil {
		return false
	}

	// Use colors if it's a terminal
	return (stat.Mode() & os.ModeCharDevice) != 0
}

// LogOperationColored logs operations with colored status indicators
func LogOperationColored(ctx context.Context, operation string, fields ...zap.Field) func() {
	logger := WithContext(ctx).With(
		append(fields, zap.String("operation", operation))...,
	)

	// Start message with green status
	logger.Info("▶ " + operation + " started")

	start := time.Now()
	return func() {
		duration := time.Since(start)
		// Completion message with green status
		logger.Info("✓ "+operation+" completed", zap.Duration("duration", duration))
	}
}

// LogErrorColored logs errors with red status indicators
func LogErrorColored(ctx context.Context, msg string, err error, fields ...zap.Field) {
	logger := WithContext(ctx).With(fields...)
	logger.Error("✗ "+msg, zap.Error(err))
}

// LogSecurityColored logs security events with yellow status indicators
func LogSecurityColored(ctx context.Context, event string, fields ...zap.Field) {
	logger := WithContext(ctx).With(
		append(fields, zap.String("security_event", event))...,
	)
	logger.Warn("⚠ " + event + " detected")
}

// LogSuccess logs success messages with green status
func LogSuccess(ctx context.Context, msg string, fields ...zap.Field) {
	logger := WithContext(ctx).With(fields...)
	logger.Info("✓ "+msg, fields...)
}

// LogWarning logs warning messages with yellow status
func LogWarning(ctx context.Context, msg string, fields ...zap.Field) {
	logger := WithContext(ctx).With(fields...)
	logger.Warn("⚠ "+msg, fields...)
}

// LogInfo logs info messages with blue status
func LogInfo(ctx context.Context, msg string, fields ...zap.Field) {
	logger := WithContext(ctx).With(fields...)
	logger.Info("ℹ "+msg, fields...)
}

// LogDebug logs debug messages with cyan status
func LogDebug(ctx context.Context, msg string, fields ...zap.Field) {
	logger := WithContext(ctx).With(fields...)
	logger.Debug("🔍 "+msg, fields...)
}
