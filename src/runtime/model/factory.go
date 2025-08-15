package model

import "os"

// FromEnv returns a Provider when environment variables are configured,
// or nil when no provider is configured. Supported variables:
//   - Local first via CMP_LOCAL_MODELS=true (uses local provider)
//   - HF_TOKEN, HF_MODEL_ID[, HF_ENDPOINT] for Hugging Face Inference API.
func FromEnv() (Provider, error) {
	if os.Getenv("CMP_LOCAL_MODELS") == "true" {
		if prov, err := NewLocalProviderFromEnv(); err == nil {
			return prov, nil
		}
	}
	if os.Getenv("HF_TOKEN") != "" && os.Getenv("HF_MODEL_ID") != "" {
		return NewHuggingFaceAPIProviderFromEnv()
	}
	return nil, nil
}
