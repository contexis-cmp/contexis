package commands

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

// GetModelsCommand provides model utilities like warmup for local models
func GetModelsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "models",
		Short: "Model utilities",
	}
	cmd.AddCommand(getWarmupCmd())
	return cmd
}

func getWarmupCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "warmup",
		Short: "Pre-download and initialize local models",
		RunE: func(cmd *cobra.Command, args []string) error {
			py := os.Getenv("CMP_PYTHON_BIN")
			if py == "" {
				py = "python3"
			}
			// Resolve script path similarly to runtime provider
			candidates := []string{}
			if override := os.Getenv("CMP_PYTHON_SCRIPT"); override != "" {
				candidates = append(candidates, override)
			}
			if root := os.Getenv("CMP_PROJECT_ROOT"); root != "" {
				candidates = append(candidates, filepath.Join(root, "src", "providers", "local_provider.py"))
			}
			if execPath, err := os.Executable(); err == nil {
				execDir := filepath.Dir(execPath)
				candidates = append(candidates, filepath.Join(execDir, "..", "src", "providers", "local_provider.py"))
			}
			if cwd, err := os.Getwd(); err == nil {
				candidates = append(candidates, filepath.Join(cwd, "..", "src", "providers", "local_provider.py"))
				candidates = append(candidates, filepath.Join(cwd, "src", "providers", "local_provider.py"))
			}
			var script string
			for _, p := range candidates {
				if _, err := os.Stat(p); err == nil {
					script = p
					break
				}
			}
			if script == "" {
				return fmt.Errorf("could not locate local_provider.py; tried: %v", candidates)
			}
			// Minimal no-op prompt to trigger model load
			payload := `{"prompt":"Hello","params":{"MaxNewTokens":1}}`
			c := exec.Command(py, script)
			c.Stdin = bytes.NewBufferString(payload)
			c.Stdout = os.Stdout
			c.Stderr = os.Stderr
			fmt.Println("Warming up local model (first run can take several minutes)...")
			return c.Run()
		},
	}
}
