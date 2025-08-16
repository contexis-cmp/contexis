package commands

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/contexis-cmp/contexis/src/cli/logger"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// RunRequest represents the request structure for the run command
type RunRequest struct {
	TenantID   string                 `json:"tenant_id"`
	Context    string                 `json:"context"`
	Component  string                 `json:"component"`
	Query      string                 `json:"query"`
	TopK       int                    `json:"top_k"`
	Data       map[string]interface{} `json:"data"`
	PromptFile string                 `json:"prompt_file"`
}

// RunResponse represents the response structure for the run command
type RunResponse struct {
	Rendered string `json:"rendered"`
}

// GetRunCommand returns the run command for direct query execution
func GetRunCommand() *cobra.Command {
	var (
		addr       string
		tenantID   string
		component  string
		topK       int
		promptFile string
		debug      bool
		timeout    int
	)

	cmd := &cobra.Command{
		Use:   "run [context] [query]",
		Short: "Run a query against a context directly",
		Long: `Execute a query against a context without manually starting the server.
This command will temporarily start the server, send the query, and return the response.

Examples:
  ctx run SupportBot "What is your return policy?"
  ctx run CustomerDocs "How do I reset my password?" --component CustomerDocs
  ctx run WorkflowProcessor "Process data" --data '{"action":"process"}'`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			contextName := args[0]
			query := args[1]

			// Validate inputs
			if contextName == "" {
				return fmt.Errorf("context name is required")
			}
			if query == "" {
				return fmt.Errorf("query is required")
			}

			// Get project root
			projectRoot, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("failed to get project root: %w", err)
			}

			// Set default component if not provided
			if component == "" {
				component = contextName
			}

			// Set default prompt file if not provided
			if promptFile == "" {
				promptFile = "agent_response.md"
			}

			// Prepare request data
			data := make(map[string]interface{})
			if cmd.Flags().Changed("data") {
				dataStr, _ := cmd.Flags().GetString("data")
				if dataStr != "" {
					if err := json.Unmarshal([]byte(dataStr), &data); err != nil {
						return fmt.Errorf("invalid JSON in --data flag: %w", err)
					}
				}
			}

			// Add user input to data
			data["user_input"] = query

			// Create request
			req := RunRequest{
				TenantID:   tenantID,
				Context:    contextName,
				Component:  component,
				Query:      query,
				TopK:       topK,
				Data:       data,
				PromptFile: promptFile,
			}

			// Execute the query
			return executeQuery(cmd.Context(), projectRoot, addr, req, debug, timeout)
		},
	}

	// Add flags
	cmd.Flags().StringVar(&addr, "addr", ":8000", "Server address to use")
	cmd.Flags().StringVar(&tenantID, "tenant", "", "Tenant ID for multi-tenant setups")
	cmd.Flags().StringVar(&component, "component", "", "Component name (defaults to context name)")
	cmd.Flags().IntVar(&topK, "top-k", 5, "Number of memory results to retrieve")
	cmd.Flags().StringVar(&promptFile, "prompt-file", "", "Prompt template file to use")
	cmd.Flags().String("data", "{}", "Additional JSON data to include in the request")
	cmd.Flags().BoolVar(&debug, "debug", false, "Enable debug output")
	cmd.Flags().IntVar(&timeout, "timeout", 30, "Timeout in seconds for the request")

	return cmd
}

// executeQuery handles the actual query execution
func executeQuery(ctx context.Context, projectRoot, addr string, req RunRequest, debug bool, timeout int) error {
	// Check if server is already running
	if isServerRunning(addr) {
		if debug {
			logger.LogInfo(ctx, "Server already running, using existing instance", zap.String("addr", addr))
		}
		return sendQuery(ctx, addr, req, debug, timeout)
	}

	// Start server in background
	if debug {
		logger.LogInfo(context.Background(), "Starting server", zap.String("addr", addr))
	}

	// Change to project root directory
	originalDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	if err := os.Chdir(projectRoot); err != nil {
		return fmt.Errorf("failed to change to project root: %w", err)
	}
	defer os.Chdir(originalDir)

	// Start server process
	executable, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}

	serverCmd := exec.CommandContext(ctx, executable, "serve", "--addr", addr)
	serverCmd.Stdout = os.Stdout
	serverCmd.Stderr = os.Stderr

	if err := serverCmd.Start(); err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}

	// Wait for server to be ready
	if err := waitForServer(addr, 10*time.Second); err != nil {
		serverCmd.Process.Kill()
		return fmt.Errorf("server failed to start: %w", err)
	}

	if debug {
		logger.LogSuccess(ctx, "Server started successfully")
	}

	// Send query
	err = sendQuery(ctx, addr, req, debug, timeout)

	// Clean up server
	if debug {
		logger.LogInfo(ctx, "Stopping server")
	}
	serverCmd.Process.Kill()
	serverCmd.Wait()

	return err
}

// isServerRunning checks if a server is already running on the given address
func isServerRunning(addr string) bool {
	// Normalize address
	if !strings.Contains(addr, ":") {
		addr = "localhost:" + addr
	}
	if strings.HasPrefix(addr, ":") {
		addr = "localhost" + addr
	}

	// Try to connect
	client := &http.Client{Timeout: 2 * time.Second}
	resp, err := client.Get("http://" + addr + "/healthz")
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}

// waitForServer waits for the server to be ready
func waitForServer(addr string, timeout time.Duration) error {
	// Normalize address
	if !strings.Contains(addr, ":") {
		addr = "localhost:" + addr
	}
	if strings.HasPrefix(addr, ":") {
		addr = "localhost" + addr
	}

	client := &http.Client{Timeout: 1 * time.Second}
	deadline := time.Now().Add(timeout)

	for time.Now().Before(deadline) {
		resp, err := client.Get("http://" + addr + "/healthz")
		if err == nil {
			resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				return nil
			}
		}
		time.Sleep(100 * time.Millisecond)
	}

	return fmt.Errorf("server not ready after %v", timeout)
}

// sendQuery sends the query to the server
func sendQuery(ctx context.Context, addr string, req RunRequest, debug bool, timeout int) error {
	// Normalize address
	if !strings.Contains(addr, ":") {
		addr = "localhost:" + addr
	}
	if strings.HasPrefix(addr, ":") {
		addr = "localhost" + addr
	}

	// Prepare request
	jsonData, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	if debug {
		logger.LogInfo(ctx, "Sending request", zap.String("url", "http://"+addr+"/api/v1/chat"))
		logger.LogDebugWithContext(ctx, "Request payload", zap.String("data", string(jsonData)))
	}

	// Create HTTP request
	httpReq, err := http.NewRequestWithContext(ctx, "POST", "http://"+addr+"/api/v1/chat", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	// Send request
	client := &http.Client{Timeout: time.Duration(timeout) * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	if debug {
		logger.LogInfo(ctx, "Response received", zap.Int("status", resp.StatusCode))
		logger.LogDebugWithContext(ctx, "Response body", zap.String("body", string(body)))
	}

	// Handle errors
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned error %d: %s", resp.StatusCode, string(body))
	}

	// Parse response
	var runResp RunResponse
	if err := json.Unmarshal(body, &runResp); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	// Output the response
	fmt.Println(runResp.Rendered)

	return nil
}
