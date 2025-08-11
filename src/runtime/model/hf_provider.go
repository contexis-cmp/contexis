package model

import (
    "bytes"
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "os"
    "time"
)

// HuggingFaceAPIProvider calls the HF Inference API for text generation.
type HuggingFaceAPIProvider struct {
    client   *http.Client
    token    string
    endpoint string
    modelID  string
}

func NewHuggingFaceAPIProviderFromEnv() (*HuggingFaceAPIProvider, error) {
    token := os.Getenv("HF_TOKEN")
    modelID := os.Getenv("HF_MODEL_ID")
    endpoint := os.Getenv("HF_ENDPOINT")
    if endpoint == "" {
        endpoint = "https://api-inference.huggingface.co/models"
    }
    if token == "" || modelID == "" {
        return nil, fmt.Errorf("HF_TOKEN and HF_MODEL_ID are required")
    }
    return &HuggingFaceAPIProvider{
        client:   &http.Client{Timeout: 60 * time.Second},
        token:    token,
        endpoint: endpoint,
        modelID:  modelID,
    }, nil
}

type hfRequest struct {
    Inputs string                 `json:"inputs"`
    Params map[string]interface{} `json:"parameters,omitempty"`
}

type hfResponse []struct {
    GeneratedText string `json:"generated_text"`
}

func (p *HuggingFaceAPIProvider) Generate(ctx context.Context, input string, params Params) (string, error) {
    body := hfRequest{Inputs: input}
    prm := map[string]interface{}{}
    if params.MaxNewTokens > 0 {
        prm["max_new_tokens"] = params.MaxNewTokens
    }
    if params.Temperature > 0 {
        prm["temperature"] = params.Temperature
    }
    if params.TopP > 0 {
        prm["top_p"] = params.TopP
    }
    if len(prm) > 0 {
        body.Params = prm
    }
    by, _ := json.Marshal(body)
    req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/%s", p.endpoint, p.modelID), bytes.NewReader(by))
    if err != nil {
        return "", err
    }
    req.Header.Set("Authorization", "Bearer "+p.token)
    req.Header.Set("Content-Type", "application/json")
    resp, err := p.client.Do(req)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()
    if resp.StatusCode >= 300 {
        return "", fmt.Errorf("hf api error: %s", resp.Status)
    }
    var out hfResponse
    if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
        return "", err
    }
    if len(out) == 0 {
        return "", fmt.Errorf("empty response from hf")
    }
    return out[0].GeneratedText, nil
}


