package model

import "os"

// FromEnv returns a Provider when environment variables are configured,
// or nil when no provider is configured. Supported variables:
//   - HF_TOKEN, HF_MODEL_ID[, HF_ENDPOINT] for Hugging Face Inference API.
func FromEnv() (Provider, error) {
	if os.Getenv("HF_TOKEN") != "" && os.Getenv("HF_MODEL_ID") != "" {
		return NewHuggingFaceAPIProviderFromEnv()
	}
	return nil, nil
}
