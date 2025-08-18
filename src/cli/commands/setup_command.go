package commands

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"

	"github.com/contexis-cmp/contexis/src/cli/logger"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// GetSetupCommand returns the `setup` command.
//
// The `setup` command bootstraps a local-first Contexis project by:
//   1) Creating a Python virtual environment at `.venv/` (if missing)
//   2) Upgrading pip and installing required Python dependencies
//      (transformers, torch, sentence-transformers)
//   3) Optionally warming up a local model so the first inference is fast
//
// Flags:
//   - `--model-id` (default: `sshleifer/tiny-gpt2`) model to warm
//   - `--no-warmup` to skip the warmup step
//   - `--timeout-seconds` timeout for warmup execution
func GetSetupCommand() *cobra.Command {
	var (
		modelID     string
		skipWarmup  bool
		timeoutSecs int
	)

	cmd := &cobra.Command{
		Use:   "setup",
		Short: "Prepare local dev: create .venv, install Python deps, warm up local model",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			wd, err := os.Getwd()
			if err != nil {
				return err
			}

			venvPython := filepath.Join(wd, ".venv", "bin", "python")
			// 1) Create venv if missing
			if _, err := os.Stat(venvPython); err != nil {
				logger.LogInfo(ctx, "Creating Python virtual environment", zap.String("path", filepath.Join(wd, ".venv")))
				c := exec.Command("python3", "-m", "venv", ".venv")
				c.Stdout = os.Stdout
				c.Stderr = os.Stderr
				c.Dir = wd
				if err := c.Run(); err != nil {
					logger.LogErrorColored(ctx, "Failed to create virtualenv", err)
					return err
				}
			} else {
				logger.LogInfo(ctx, "Using existing virtual environment", zap.String("python", venvPython))
			}

			// 2) Upgrade pip
			logger.LogInfo(ctx, "Upgrading pip")
			cPipUp := exec.Command(venvPython, "-m", "pip", "install", "-q", "--upgrade", "pip")
			cPipUp.Stdout = os.Stdout
			cPipUp.Stderr = os.Stderr
			cPipUp.Dir = wd
			if err := cPipUp.Run(); err != nil {
				logger.LogErrorColored(ctx, "Failed to upgrade pip", err)
				return err
			}

			// 3) Install required deps
			logger.LogInfo(ctx, "Installing Python dependencies",
				zap.Strings("packages", []string{"transformers==4.55.0", "torch", "sentence-transformers==2.7.0"}))
			cPip := exec.Command(venvPython, "-m", "pip", "install", "-q",
				"transformers==4.55.0", "torch", "sentence-transformers==2.7.0")
			cPip.Stdout = os.Stdout
			cPip.Stderr = os.Stderr
			cPip.Dir = wd
			if err := cPip.Run(); err != nil {
				logger.LogErrorColored(ctx, "Failed to install Python dependencies", err)
				return err
			}

			// Export helper envs for subsequent commands
			_ = os.Setenv("CMP_PYTHON_BIN", venvPython)
			_ = os.Setenv("CMP_PROJECT_ROOT", wd)
			if os.Getenv("CMP_LOCAL_MODELS") == "" {
				_ = os.Setenv("CMP_LOCAL_MODELS", "true")
			}

			// 4) Warm up local model (optional)
			if !skipWarmup {
				logger.LogInfo(ctx, "Warming up local model (first run may take several minutes)",
					zap.String("model_id", modelID), zap.Int("timeout_seconds", timeoutSecs))
				exe, err := os.Executable()
				if err != nil {
					return err
				}
				env := os.Environ()
				env = append(env,
					"CMP_PYTHON_BIN="+venvPython,
					"CMP_PROJECT_ROOT="+wd,
					"CMP_LOCAL_MODELS=true",
					"CMP_LOCAL_TIMEOUT_SECONDS="+strconv.Itoa(timeoutSecs),
				)
				if modelID != "" {
					env = append(env, "CMP_LOCAL_MODEL_ID="+modelID)
				}
				warm := exec.Command(exe, "models", "warmup")
				warm.Stdout = os.Stdout
				warm.Stderr = os.Stderr
				warm.Dir = wd
				warm.Env = env
				if err := warm.Run(); err != nil {
					logger.LogErrorColored(ctx, "Model warmup failed", err)
					return err
				}
			}

			logger.LogSuccess(ctx, "Setup complete",
				zap.String("python", venvPython))
			fmt.Fprintln(cmd.OutOrStdout(), "Setup complete. You can now run: ctx serve --addr :8000")
			return nil
		},
	}

	cmd.Flags().StringVar(&modelID, "model-id", "sshleifer/tiny-gpt2", "Local model id to warm up (e.g., microsoft/Phi-3-mini-4k-instruct)")
	cmd.Flags().BoolVar(&skipWarmup, "no-warmup", false, "Skip model warmup step")
	cmd.Flags().IntVar(&timeoutSecs, "timeout-seconds", 1800, "Timeout for warmup (seconds)")
	return cmd
}
