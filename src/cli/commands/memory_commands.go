package commands

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	runtimememory "github.com/contexis-cmp/contexis/src/runtime/memory"
	"github.com/spf13/cobra"
)

// GetMemoryCommand returns the parent memory command with subcommands.
func GetMemoryCommand() *cobra.Command {
	memCmd := &cobra.Command{Use: "memory", Short: "Memory operations (ingest, search, optimize)"}
	memCmd.AddCommand(newMemoryIngestCmd())
	memCmd.AddCommand(newMemorySearchCmd())
	memCmd.AddCommand(newMemoryOptimizeCmd())
	return memCmd
}

func newMemoryIngestCmd() *cobra.Command {
	var (
		provider  string
		component string
		tenant    string
		model     string
		inputPath string
	)
	cmd := &cobra.Command{
		Use:   "ingest",
		Short: "Ingest documents into a memory store",
		RunE: func(cmd *cobra.Command, args []string) error {
			if component == "" {
				return fmt.Errorf("--component is required")
			}
			if provider == "" {
				provider = "sqlite"
			}
			cfg := runtimememory.Config{
				Provider:       provider,
				RootDir:        mustGetwd(),
				ComponentName:  component,
				EmbeddingModel: model,
				Settings:       map[string]string{},
				TenantID:       tenant,
			}
			store, err := runtimememory.NewStore(cfg)
			if err != nil {
				return err
			}
			defer store.Close()
			// load documents from file (one per line) or stdin
			docs, err := readLines(inputPath)
			if err != nil {
				return err
			}
			ver, err := store.IngestDocuments(context.Background(), docs)
			if err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "ingested version: %s\n", ver)
			return nil
		},
	}
	cmd.Flags().StringVar(&provider, "provider", "sqlite", "Memory provider (sqlite, episodic)")
	cmd.Flags().StringVar(&component, "component", "", "Component name (e.g., CustomerDocs, SupportBot)")
	cmd.Flags().StringVar(&tenant, "tenant", "", "Tenant ID")
	cmd.Flags().StringVar(&model, "model", "bge-small-en", "Embedding model identifier")
	cmd.Flags().StringVar(&inputPath, "input", "", "Path to file with documents (one per line). If empty, read from stdin")
	return cmd
}

func newMemorySearchCmd() *cobra.Command {
	var (
		provider  string
		component string
		tenant    string
		query     string
		topK      int
	)
	cmd := &cobra.Command{
		Use:   "search",
		Short: "Search a memory store",
		RunE: func(cmd *cobra.Command, args []string) error {
			if component == "" {
				return fmt.Errorf("--component is required")
			}
			if query == "" {
				return fmt.Errorf("--query is required")
			}
			cfg := runtimememory.Config{Provider: provider, RootDir: mustGetwd(), ComponentName: component, TenantID: tenant}
			store, err := runtimememory.NewStore(cfg)
			if err != nil {
				return err
			}
			defer store.Close()
			results, err := store.Search(context.Background(), query, topK)
			if err != nil {
				return err
			}
			for _, r := range results {
				fmt.Fprintf(cmd.OutOrStdout(), "%.3f\t%s\n", r.Score, strings.ReplaceAll(strings.TrimSpace(r.Content), "\n", " "))
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&provider, "provider", "sqlite", "Memory provider (sqlite, episodic)")
	cmd.Flags().StringVar(&component, "component", "", "Component name (e.g., CustomerDocs, SupportBot)")
	cmd.Flags().StringVar(&tenant, "tenant", "", "Tenant ID")
	cmd.Flags().StringVar(&query, "query", "", "Search query")
	cmd.Flags().IntVar(&topK, "top-k", 5, "Number of results to return")
	return cmd
}

func newMemoryOptimizeCmd() *cobra.Command {
	var (
		provider  string
		component string
		tenant    string
		version   string
	)
	cmd := &cobra.Command{
		Use:   "optimize",
		Short: "Optimize a memory store",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := runtimememory.Config{Provider: provider, RootDir: mustGetwd(), ComponentName: component, TenantID: tenant}
			store, err := runtimememory.NewStore(cfg)
			if err != nil {
				return err
			}
			defer store.Close()
			return store.Optimize(context.Background(), version)
		},
	}
	cmd.Flags().StringVar(&provider, "provider", "sqlite", "Memory provider (sqlite, episodic)")
	cmd.Flags().StringVar(&component, "component", "", "Component name (e.g., CustomerDocs, SupportBot)")
	cmd.Flags().StringVar(&tenant, "tenant", "", "Tenant ID")
	cmd.Flags().StringVar(&version, "version", "", "Memory version identifier to optimize (optional)")
	return cmd
}

func readLines(path string) ([]string, error) {
	var f *os.File
	var err error
	if path == "" {
		f = os.Stdin
	} else {
		f, err = os.Open(path)
		if err != nil {
			return nil, err
		}
		defer f.Close()
	}
	var lines []string
	s := bufio.NewScanner(f)
	for s.Scan() {
		lines = append(lines, s.Text())
	}
	if err := s.Err(); err != nil {
		return nil, err
	}
	return lines, nil
}

func mustGetwd() string {
	wd, _ := os.Getwd()
	return wd
}
