package runtimememory

import (
	"bufio"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"
)

type sqliteVectorStore struct {
	filePath     string
	embeddingDim int
	model        string
}

type vecRecord struct {
	ID      string `json:"id"`
	Content string `json:"content"`
	Vector  string `json:"vector_b64"`
}

// item represents a scored vector match used internally for ranking
type item struct {
	id, content string
	score       float64
}

func newSQLiteVectorStore(cfg Config) (MemoryStore, error) {
	dim := 384
	if v, ok := cfg.Settings["embedding_dim"]; ok && v != "" {
		fmt.Sscanf(v, "%d", &dim)
	}
	filePath := DerivePath(cfg.RootDir, cfg.ComponentName, cfg.TenantID, "vector_store.jsonl")
	if err := os.MkdirAll(filepath.Dir(filePath), 0o755); err != nil {
		return nil, err
	}
	// ensure file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		if err := os.WriteFile(filePath, nil, 0o644); err != nil {
			return nil, err
		}
	}
	return &sqliteVectorStore{filePath: filePath, embeddingDim: dim, model: cfg.EmbeddingModel}, nil
}

func (s *sqliteVectorStore) Close() error { return nil }

func (s *sqliteVectorStore) IngestDocuments(ctx context.Context, documents []string) (string, error) {
	if len(documents) == 0 {
		return "", fmt.Errorf("no documents to ingest")
	}
	version := contentSHA(documents, s.model)
	f, err := os.OpenFile(s.filePath, os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return "", err
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	for i, doc := range documents {
		id := fmt.Sprintf("%s_%d", version, i)
		vec := naiveEmbed(doc, s.embeddingDim)
		rec := vecRecord{ID: id, Content: doc, Vector: base64.StdEncoding.EncodeToString(float64sToBytes(vec))}
		by, _ := json.Marshal(rec)
		if _, err := w.Write(by); err != nil {
			return "", err
		}
		if _, err := w.WriteString("\n"); err != nil {
			return "", err
		}
	}
	if err := w.Flush(); err != nil {
		return "", err
	}
	return version, nil
}

func (s *sqliteVectorStore) Search(ctx context.Context, query string, topK int) ([]SearchResult, error) {
	if topK <= 0 {
		topK = 5
	}
	qvec := naiveEmbed(query, s.embeddingDim)
	f, err := os.Open(s.filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	items := make([]item, 0, 64)
	scan := bufio.NewScanner(f)
	for scan.Scan() {
		var rec vecRecord
		if err := json.Unmarshal([]byte(scan.Text()), &rec); err != nil {
			continue
		}
		vb, err := base64.StdEncoding.DecodeString(rec.Vector)
		if err != nil {
			continue
		}
		v := bytesToFloat64s(vb)
		if len(v) != len(qvec) {
			continue
		}
		score := cosine(qvec, v)
		items = append(items, item{id: rec.ID, content: rec.Content, score: score})
	}
	if err := scan.Err(); err != nil {
		return nil, err
	}
	selectTopK(items, topK)
	results := make([]SearchResult, 0, min(topK, len(items)))
	for i := 0; i < min(topK, len(items)); i++ {
		it := items[i]
		results = append(results, SearchResult{ID: it.id, Content: it.content, Score: it.score})
	}
	return results, nil
}

func (s *sqliteVectorStore) Optimize(ctx context.Context, _ string) error { return nil }

// --- helpers ---

func naiveEmbed(text string, dim int) []float64 {
	vec := make([]float64, dim)
	var h uint64 = 1469598103934665603
	const prime uint64 = 1099511628211
	for i := 0; i < len(text); i++ {
		h ^= uint64(text[i])
		h *= prime
		idx := int(h % uint64(dim))
		vec[idx] += 1.0
	}
	var norm float64
	for _, v := range vec {
		norm += v * v
	}
	norm = math.Sqrt(norm)
	if norm == 0 {
		norm = 1
	}
	for i := range vec {
		vec[i] /= norm
	}
	return vec
}

func cosine(a, b []float64) float64 {
	var dot, na, nb float64
	for i := range a {
		dot += a[i] * b[i]
		na += a[i] * a[i]
		nb += b[i] * b[i]
	}
	if na == 0 || nb == 0 {
		return 0
	}
	return dot / (math.Sqrt(na) * math.Sqrt(nb))
}

func selectTopK(items []item, k int) {
	n := len(items)
	if k > n {
		k = n
	}
	for i := 0; i < k; i++ {
		maxIdx := i
		for j := i + 1; j < n; j++ {
			if items[j].score > items[maxIdx].score {
				maxIdx = j
			}
		}
		items[i], items[maxIdx] = items[maxIdx], items[i]
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func float64sToBytes(v []float64) []byte {
	b := make([]byte, 8*len(v))
	for i, f := range v {
		u := math.Float64bits(f)
		for j := 0; j < 8; j++ {
			b[i*8+j] = byte(u >> (8 * j))
		}
	}
	return b
}

func bytesToFloat64s(b []byte) []float64 {
	if len(b)%8 != 0 {
		return nil
	}
	n := len(b) / 8
	v := make([]float64, n)
	for i := 0; i < n; i++ {
		var u uint64
		for j := 0; j < 8; j++ {
			u |= uint64(b[i*8+j]) << (8 * j)
		}
		v[i] = math.Float64frombits(u)
	}
	return v
}
