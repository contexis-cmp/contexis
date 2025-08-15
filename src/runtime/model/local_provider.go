package model

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// localPythonProvider shells out to the Python LocalAIProvider to generate text.
type localPythonProvider struct {
	pythonBin  string
	scriptPath string
	timeout    time.Duration
}

type localReq struct {
	Prompt string `json:"prompt"`
	Params Params `json:"params"`
}

type localResp struct {
	Output string `json:"output"`
	Error  string `json:"error,omitempty"`
}

func newLocalPythonProviderFromEnv() (Provider, error) {
	py := os.Getenv("CMP_PYTHON_BIN")
	if py == "" {
		py = "python3"
	}
	// Resolve script path robustly for both repo root and generated project dirs
	if override := os.Getenv("CMP_PYTHON_SCRIPT"); override != "" {
		if _, err := os.Stat(override); err == nil {
			return &localPythonProvider{pythonBin: py, scriptPath: override, timeout: resolveLocalTimeout()}, nil
		}
	}
	candidates := []string{}
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
	for _, path := range candidates {
		if _, err := os.Stat(path); err == nil {
			return &localPythonProvider{pythonBin: py, scriptPath: path, timeout: resolveLocalTimeout()}, nil
		}
	}
	return nil, fmt.Errorf("local provider script not found in candidates: %v", candidates)
}

func resolveLocalTimeout() time.Duration {
	t := 600 * time.Second
	if v := os.Getenv("CMP_LOCAL_TIMEOUT_SECONDS"); v != "" {
		if dur, err := time.ParseDuration(v + "s"); err == nil {
			return dur
		}
	}
	return t
}

func (p *localPythonProvider) Generate(ctx context.Context, input string, params Params) (string, error) {
	req := localReq{Prompt: input, Params: params}
	payload, _ := json.Marshal(req)

	// Execute local Python script directly
	cmd := exec.CommandContext(ctx, p.pythonBin, p.scriptPath)
	cmd.Stdin = bytes.NewReader(payload)
	var out bytes.Buffer
	cmd.Stdout = &out
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	// Set working directory to project root if provided
	if cwd := os.Getenv("CMP_PROJECT_ROOT"); cwd != "" {
		if abs, err := filepath.Abs(cwd); err == nil {
			cmd.Dir = abs
		}
	}

	if err := cmd.Start(); err != nil {
		return "", fmt.Errorf("local provider start failed: %w", err)
	}

	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()

	select {
	case err := <-done:
		if err != nil {
			return "", fmt.Errorf("local provider error: %s", stderr.String())
		}
	case <-time.After(p.timeout):
		_ = cmd.Process.Kill()
		return "", fmt.Errorf("local provider timeout")
	}

	// Parse response
	scanner := bufio.NewScanner(&out)
	var resp localResp
	var buf bytes.Buffer
	for scanner.Scan() {
		buf.Write(scanner.Bytes())
	}
	if err := json.Unmarshal(buf.Bytes(), &resp); err != nil {
		return "", fmt.Errorf("failed to parse local provider output: %w", err)
	}
	if resp.Error != "" {
		return "", fmt.Errorf(resp.Error)
	}
	return resp.Output, nil
}
