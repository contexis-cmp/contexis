package runtimememory

import (
	"context"
	"testing"
)

func TestEpisodicStore_IngestAndSearch(t *testing.T) {
	root := t.TempDir()
	cfg := Config{Provider: "episodic", RootDir: root, ComponentName: "SupportBot", TenantID: "acme"}
	store, err := NewStore(cfg)
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}
	defer store.Close()

	docs := []string{"User asked about refund policy", "Follow-up: order 12345", "Resolved: provided RMA"}
	if _, err := store.IngestDocuments(context.Background(), docs); err != nil {
		t.Fatalf("IngestDocuments: %v", err)
	}
	res, err := store.Search(context.Background(), "refund policy", 3)
	if err != nil {
		t.Fatalf("Search: %v", err)
	}
	if len(res) == 0 {
		t.Fatalf("expected results, got none")
	}
}

func TestSQLiteVectorStore_IngestAndSearch_FileBackend(t *testing.T) {
	root := t.TempDir()
	cfg := Config{Provider: "sqlite", RootDir: root, ComponentName: "CustomerDocs", EmbeddingModel: "bge-small-en"}
	store, err := NewStore(cfg)
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}
	defer store.Close()

	docs := []string{"Returns are accepted within 30 days.", "Shipping takes 3-5 business days."}
	if _, err := store.IngestDocuments(context.Background(), docs); err != nil {
		t.Fatalf("IngestDocuments: %v", err)
	}
	res, err := store.Search(context.Background(), "return policy 30 days", 2)
	if err != nil {
		t.Fatalf("Search: %v", err)
	}
	if len(res) == 0 {
		t.Fatalf("expected results, got none")
	}
}
