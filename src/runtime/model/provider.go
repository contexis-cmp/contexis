package model

import "context"

// Params defines model generation parameters that influence decoding.
// Fields may be ignored by providers that do not support them.
type Params struct {
	Temperature   float64 // Sampling temperature
	TopP          float64 // Nucleus sampling probability
	MaxNewTokens  int     // Maximum number of new tokens to generate
	RepetitionPen float64 // Repetition penalty
}

// Provider defines an inference provider capable of generating text based
// on a rendered prompt and parameter set.
type Provider interface {
	Generate(ctx context.Context, input string, params Params) (string, error)
}

// NewLocalProviderFromEnv returns a Provider backed by a local Python process
// when local-first variables are enabled. Implementation provided in local_provider.go.
func NewLocalProviderFromEnv() (Provider, error) {
	return newLocalPythonProviderFromEnv()
}
