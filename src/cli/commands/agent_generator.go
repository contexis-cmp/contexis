package commands

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	"github.com/contexis/cmp/src/cli/logger"
)

// generateAgent creates a conversational agent with tools and memory
func generateAgent(ctx context.Context, name, tools, memory string) error {
	log := logger.WithContext(ctx)

	log.Info("agent generator not yet implemented",
		zap.String("name", name),
		zap.String("tools", tools),
		zap.String("memory", memory))

	return fmt.Errorf("agent generator will be implemented in Week 2 of Phase 1")
}
