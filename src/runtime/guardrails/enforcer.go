package guardrails

import (
    "fmt"
    "strings"

    corectx "github.com/contexis-cmp/contexis/src/core/context"
    runtimeprompt "github.com/contexis-cmp/contexis/src/runtime/prompt"
)

// Capability validation: ensure requested action is allowed by context.
func ValidateCapabilities(ctx *corectx.Context, requestedAction string) error {
    if requestedAction == "" { return nil }
    for _, cap := range ctx.Role.Capabilities {
        if strings.EqualFold(cap, requestedAction) {
            return nil
        }
    }
    return fmt.Errorf("capability '%s' not allowed by context", requestedAction)
}

// EnforceGuardrails performs simple enforcement based on context guardrails.
// - Ensures response format (json/markdown/text)
// - Trims tokens if MaxTokens specified
// - Temperature and tone are advisory here (no model call), so not enforced beyond bounds check.
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


