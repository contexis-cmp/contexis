package model

import "os"

// FromEnv returns a Provider when environment is configured, else nil.
func FromEnv() (Provider, error) {
    if os.Getenv("HF_TOKEN") != "" && os.Getenv("HF_MODEL_ID") != "" {
        return NewHuggingFaceAPIProviderFromEnv()
    }
    return nil, nil
}


