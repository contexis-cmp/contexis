package model

import "context"

// Params defines model generation parameters.
type Params struct {
    Temperature   float64
    TopP          float64
    MaxNewTokens  int
    RepetitionPen float64
}

// Provider defines an inference provider.
type Provider interface {
    Generate(ctx context.Context, input string, params Params) (string, error)
}


