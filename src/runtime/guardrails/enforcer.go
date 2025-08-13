package guardrails

import (
	"fmt"
	"strings"

	corectx "github.com/contexis-cmp/contexis/src/core/context"
	runtimeprompt "github.com/contexis-cmp/contexis/src/runtime/prompt"
)

// ValidateCapabilities returns an error if requestedAction is not permitted by
// the Context's role capabilities. An empty action is treated as allowed.
func ValidateCapabilities(ctx *corectx.Context, requestedAction string) error {
	if requestedAction == "" {
		return nil
	}
	for _, cap := range ctx.Role.Capabilities {
		if strings.EqualFold(cap, requestedAction) {
			return nil
		}
	}
	return fmt.Errorf("capability '%s' not allowed by context", requestedAction)
}

// EnforceGuardrails applies guardrail constraints to a candidate response.
//
// It enforces:
//   - Response format (json|markdown|text)
//   - Token length (MaxTokens)
//
// Tone and temperature are advisory and not enforced at this stage.
func EnforceGuardrails(ctx *corectx.Context, response string) (string, error) {
	gr := ctx.Guardrails
	if gr.MaxTokens > 0 {
		response = runtimeprompt.OptimizeTokens(response, gr.MaxTokens)
	}
	if err := runtimeprompt.ValidateFormat(gr.Format, response); err != nil {
		return "", err
	}
	// tone/temperature hooks: future implementation
	return response, nil
}
