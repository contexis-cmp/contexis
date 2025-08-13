// Package model defines the model provider abstraction used by the runtime.
//
// Providers implement text generation over rendered prompts. A factory reads
// environment variables to select a concrete implementation (e.g. Hugging Face
// Inference API), or returns nil when no provider is configured.
package model
