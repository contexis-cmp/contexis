package commands

import (
	"context"
	"fmt"
	"os"

	"github.com/contexis-cmp/contexis/src/cli/logger"
	runtimemodel "github.com/contexis-cmp/contexis/src/runtime/model"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func GetHFCommand() *cobra.Command {
	hfCmd := &cobra.Command{
		Use:   "hf",
		Short: "Hugging Face utilities",
	}

	testCmd := &cobra.Command{
		Use:   "test-model <prompt>",
		Short: "Test the configured HF model with a prompt",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			if os.Getenv("HF_TOKEN") == "" || os.Getenv("HF_MODEL_ID") == "" {
				return fmt.Errorf("HF_TOKEN and HF_MODEL_ID must be set")
			}

			logger.LogInfo(ctx, "Testing Hugging Face model",
				zap.String("model", os.Getenv("HF_MODEL_ID")),
				zap.String("prompt", args[0]))

			prov, err := runtimemodel.NewHuggingFaceAPIProviderFromEnv()
			if err != nil {
				logger.LogErrorColored(ctx, "Failed to create HF provider", err)
				return err
			}

			out, err := prov.Generate(context.Background(), args[0], runtimemodel.Params{MaxNewTokens: 128})
			if err != nil {
				logger.LogErrorColored(ctx, "Failed to generate response", err)
				return err
			}

			logger.LogSuccess(ctx, "HF model test completed",
				zap.Int("response_length", len(out)),
				zap.String("model", os.Getenv("HF_MODEL_ID")))

			fmt.Println(out)
			return nil
		},
	}
	hfCmd.AddCommand(testCmd)
	return hfCmd
}

// zapField is a small helper to avoid importing zap in this file's signature
func zapField(key string, val interface{}) interface{} {
	return struct {
		K string
		V interface{}
	}{K: key, V: val}
}
