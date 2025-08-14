package logger

import (
	"context"
	"testing"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestColoredLogging(t *testing.T) {
	// Initialize colored logger
	err := InitColoredLogger("debug")
	if err != nil {
		t.Fatalf("Failed to initialize colored logger: %v", err)
	}

	ctx := context.Background()

	// Test different log levels and status indicators
	t.Run("success_logging", func(t *testing.T) {
		LogSuccess(ctx, "Operation completed successfully", zap.String("operation", "test"))
	})

	t.Run("info_logging", func(t *testing.T) {
		LogInfo(ctx, "Processing request", zap.String("request_id", "req_123"))
	})

	t.Run("warning_logging", func(t *testing.T) {
		LogWarning(ctx, "High memory usage detected", zap.Float64("usage_percent", 85.5))
	})

	t.Run("debug_logging", func(t *testing.T) {
		LogDebug(ctx, "Detailed operation trace", zap.Duration("duration", 150*time.Millisecond))
	})

	t.Run("error_logging", func(t *testing.T) {
		LogErrorColored(ctx, "Failed to connect to database",
			&testError{msg: "connection timeout"},
			zap.String("database", "postgres"))
	})

	t.Run("security_logging", func(t *testing.T) {
		LogSecurityColored(ctx, "Unauthorized access attempt",
			zap.String("ip", "192.168.1.100"),
			zap.String("user_agent", "curl/7.68.0"))
	})

	t.Run("operation_logging", func(t *testing.T) {
		done := LogOperationColored(ctx, "Data processing",
			zap.String("dataset", "customer_data"),
			zap.Int("records", 1000))

		// Simulate some work
		time.Sleep(10 * time.Millisecond)
		done()
	})
}

// testError implements error interface for testing
type testError struct {
	msg string
}

func (e *testError) Error() string {
	return e.msg
}

func TestColorDetection(t *testing.T) {
	// Test color detection logic
	t.Run("should_use_colors", func(t *testing.T) {
		// This should return true in a terminal environment
		colored := shouldUseColors()
		t.Logf("Color detection result: %v", colored)
	})
}

func TestColoredLevels(t *testing.T) {
	// Test level coloring
	levels := []zapcore.Level{
		zapcore.DebugLevel,
		zapcore.InfoLevel,
		zapcore.WarnLevel,
		zapcore.ErrorLevel,
		zapcore.FatalLevel,
	}

	for _, level := range levels {
		colored := getColoredLevel(level)
		t.Logf("Level %s: %s", level.String(), colored)
	}
}
