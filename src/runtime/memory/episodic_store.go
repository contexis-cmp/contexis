package runtimememory

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type episodicStore struct {
	logPath string
	encrypt bool
}

func newEpisodicStore(cfg Config) (MemoryStore, error) {
	logPath := DerivePath(cfg.RootDir, cfg.ComponentName, cfg.TenantID, "episodic/episodes.log")
	if err := os.MkdirAll(filepath.Dir(logPath), 0o755); err != nil {
		return nil, err
	}
	// Ensure file exists
	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		f, err := os.Create(logPath)
		if err != nil {
			return nil, err
		}
		_ = f.Close()
	}
	enc := false
	if v, ok := cfg.Settings["episodic_encryption"]; ok && v == "true" {
		enc = true
	}
	return &episodicStore{logPath: logPath, encrypt: enc}, nil
}

func (e *episodicStore) Close() error { return nil }

func (e *episodicStore) IngestDocuments(ctx context.Context, documents []string) (string, error) {
	if len(documents) == 0 {
		return "", fmt.Errorf("no episodic entries to ingest")
	}
	f, err := os.OpenFile(e.logPath, os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return "", err
	}
	defer f.Close()
	for _, d := range documents {
		line := strings.TrimSpace(d)
		if e.encrypt {
			// placeholder: reversible base64-like marker, not real crypto; replace with AES-GCM later
			line = "enc:" + line
		}
		if _, err := fmt.Fprintf(f, "%s\n", line); err != nil {
			return "", err
		}
	}
	// version based on content and count; simple timestamp-less hash via contentSHA
	return contentSHA(documents, "episodic"), nil
}

func (e *episodicStore) Search(ctx context.Context, query string, topK int) ([]SearchResult, error) {
	if topK <= 0 {
		topK = 5
	}
	q := strings.ToLower(strings.TrimSpace(query))
	file, err := os.Open(e.logPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	results := make([]SearchResult, 0, topK)
	idx := 0
	for scanner.Scan() {
		line := scanner.Text()
		score := simpleMatchScore(strings.ToLower(line), q)
		if score > 0 {
			results = append(results, SearchResult{ID: fmt.Sprintf("%d", idx), Content: line, Score: score})
		}
		idx++
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	// sort by score desc
	for i := 0; i < len(results); i++ {
		max := i
		for j := i + 1; j < len(results); j++ {
			if results[j].Score > results[max].Score {
				max = j
			}
		}
		results[i], results[max] = results[max], results[i]
	}
	if len(results) > topK {
		results = results[:topK]
	}
	return results, nil
}

func (e *episodicStore) Optimize(ctx context.Context, _ string) error {
	// No-op for simple file-backed store
	return nil
}

func simpleMatchScore(text, query string) float64 {
	if query == "" {
		return 0
	}
	// crude score: proportion of query terms present
	qterms := strings.Fields(query)
	if len(qterms) == 0 {
		return 0
	}
	present := 0
	for _, t := range qterms {
		if strings.Contains(text, t) {
			present++
		}
	}
	return float64(present) / float64(len(qterms))
}
