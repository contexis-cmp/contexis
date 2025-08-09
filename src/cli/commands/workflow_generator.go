package commands

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	"github.com/contexis/cmp/src/cli/logger"
)

// generateWorkflow creates a multi-step AI processing pipeline
func generateWorkflow(ctx context.Context, name, steps string) error {
	log := logger.WithContext(ctx)

	log.Info("workflow generator not yet implemented",
		zap.String("name", name),
		zap.String("steps", steps))

	return fmt.Errorf("workflow generator will be implemented in Week 3 of Phase 1")
}
