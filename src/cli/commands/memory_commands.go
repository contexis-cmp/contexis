package commands

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/contexis-cmp/contexis/src/cli/logger"
	runtimememory "github.com/contexis-cmp/contexis/src/runtime/memory"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// GetMemoryCommand returns the `memory` command with subcommands for ingest, seed, search, and optimize.
//
// Notable DX helpers:
//   - `ctx memory seed --component <Name>`: bulk-ingests all supported documents under memory/<Name>/documents
//   - `ctx memory ingest --all --component <Name>`: same as seed, inline flag
func GetMemoryCommand() *cobra.Command {
	memCmd := &cobra.Command{Use: "memory", Short: "Memory operations (ingest, search, optimize)"}
	memCmd.AddCommand(newMemoryIngestCmd())
	memCmd.AddCommand(newMemorySeedCmd())
	memCmd.AddCommand(newMemorySearchCmd())
	memCmd.AddCommand(newMemoryOptimizeCmd())
	return memCmd
}

// newMemoryIngestCmd returns the `ingest` subcommand which ingests documents from a file/stdin or all docs under a component.
func newMemoryIngestCmd() *cobra.Command {
	var (
		provider     string
		component    string
		tenant       string
		model        string
		inputPath    string
		allDocuments bool
	)
	cmd := &cobra.Command{
		Use:   "ingest",
		Short: "Ingest documents into a memory store",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			if component == "" {
				return fmt.Errorf("--component is required")
			}
			if provider == "" {
				provider = "sqlite"
			}

			logger.LogInfo(ctx, "Starting memory ingestion",
				zap.String("component", component),
				zap.String("provider", provider),
				zap.String("model", model))

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
				logger.LogErrorColored(ctx, "Failed to create memory store", err)
				return err
			}
			defer store.Close()

			var docs []string
			if allDocuments {
				// Read all files under memory/<component>/documents recursively
				docsDir := runtimememory.DerivePath(cfg.RootDir, component, tenant, "documents")
				logger.LogInfo(ctx, "Scanning documents directory", zap.String("path", docsDir))
				err := filepath.Walk(docsDir, func(path string, info os.FileInfo, err error) error {
					if err != nil {
						return err
					}
					if info.IsDir() {
						return nil
					}
					low := strings.ToLower(info.Name())
					if strings.HasSuffix(low, ".txt") || strings.HasSuffix(low, ".md") {
						b, rErr := os.ReadFile(path)
						if rErr != nil {
							return rErr
						}
						docs = append(docs, string(b))
					} else if strings.HasSuffix(low, ".pdf") {
						if bin, lookErr := exec.LookPath("pdftotext"); lookErr == nil {
							cmd := exec.Command(bin, "-layout", path, "-")
							out, cErr := cmd.Output()
							if cErr == nil {
								docs = append(docs, string(out))
							}
						} else {
							logger.WithContext(ctx).Info("skipping PDF (pdftotext not found)", zap.String("file", path))
						}
					}
					return nil
				})
				if err != nil {
					logger.LogErrorColored(ctx, "Failed to read documents directory", err)
					return fmt.Errorf("failed to read documents directory %s: %w", docsDir, err)
				}
				if len(docs) == 0 {
					logger.LogInfo(ctx, "No supported documents found (txt, md, pdf)")
				}
			} else {
				// load documents from file (one per line) or stdin
				var rErr error
				docs, rErr = readLines(inputPath)
				if rErr != nil {
					logger.LogErrorColored(ctx, "Failed to read documents", rErr)
					return rErr
				}
			}

			logger.LogInfo(ctx, "Ingesting documents", zap.Int("count", len(docs)))
			ver, err := store.IngestDocuments(context.Background(), docs)
			if err != nil {
				logger.LogErrorColored(ctx, "Failed to ingest documents", err)
				return err
			}

			logger.LogSuccess(ctx, "Memory ingestion completed",
				zap.String("version", ver),
				zap.Int("documents_ingested", len(docs)))

			fmt.Fprintf(cmd.OutOrStdout(), "ingested version: %s\n", ver)
			return nil
		},
	}
	cmd.Flags().StringVar(&provider, "provider", "sqlite", "Memory provider (sqlite, episodic)")
	cmd.Flags().StringVar(&component, "component", "", "Component name (e.g., CustomerDocs, SupportBot)")
	cmd.Flags().StringVar(&tenant, "tenant", "", "Tenant ID")
	cmd.Flags().StringVar(&model, "model", "bge-small-en", "Embedding model identifier")
	cmd.Flags().StringVar(&inputPath, "input", "", "Path to file with documents (one per line). If empty, read from stdin")
	cmd.Flags().BoolVar(&allDocuments, "all", false, "Ingest all documents under memory/<component>/documents (txt, md, pdf)")
	return cmd
}

// newMemorySeedCmd returns the `seed` subcommand which bulk-ingests all supported documents for a component.
// This is the DX-equivalent of Rails' db:seed for memory documents.
func newMemorySeedCmd() *cobra.Command {
	var (
		provider  string
		component string
		tenant    string
		model     string
	)
	cmd := &cobra.Command{
		Use:   "seed",
		Short: "Seed memory by ingesting all documents under memory/<component>/documents",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			if component == "" {
				return fmt.Errorf("--component is required")
			}
			if provider == "" {
				provider = "sqlite"
			}

			// Delegate to ingest --all
			logger.LogInfo(ctx, "Seeding memory for component",
				zap.String("component", component))
			ingest := newMemoryIngestCmd()
			ingest.Flags().Set("provider", provider)
			ingest.Flags().Set("component", component)
			if tenant != "" {
				ingest.Flags().Set("tenant", tenant)
			}
			if model != "" {
				ingest.Flags().Set("model", model)
			}
			ingest.Flags().Set("all", "true")
			return ingest.RunE(cmd, []string{})
		},
	}
	cmd.Flags().StringVar(&provider, "provider", "sqlite", "Memory provider (sqlite, episodic)")
	cmd.Flags().StringVar(&component, "component", "", "Component name (e.g., CustomerDocs, SupportBot)")
	cmd.Flags().StringVar(&tenant, "tenant", "", "Tenant ID")
	cmd.Flags().StringVar(&model, "model", "bge-small-en", "Embedding model identifier")
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
			ctx := cmd.Context()
			if component == "" {
				return fmt.Errorf("--component is required")
			}
			if query == "" {
				return fmt.Errorf("--query is required")
			}

			logger.LogInfo(ctx, "Starting memory search",
				zap.String("component", component),
				zap.String("provider", provider),
				zap.String("query", query),
				zap.Int("top_k", topK))

			cfg := runtimememory.Config{Provider: provider, RootDir: mustGetwd(), ComponentName: component, TenantID: tenant}
			store, err := runtimememory.NewStore(cfg)
			if err != nil {
				logger.LogErrorColored(ctx, "Failed to create memory store", err)
				return err
			}
			defer store.Close()

			results, err := store.Search(context.Background(), query, topK)
			if err != nil {
				logger.LogErrorColored(ctx, "Failed to search memory", err)
				return err
			}

			logger.LogSuccess(ctx, "Memory search completed",
				zap.Int("results_count", len(results)))

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
