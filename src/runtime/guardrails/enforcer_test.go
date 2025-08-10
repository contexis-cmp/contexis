package guardrails

import (
	"testing"

	corectx "github.com/contexis-cmp/contexis/src/core/context"
)

func TestValidateCapabilities(t *testing.T) {
	ctx := &corectx.Context{Role: corectx.Role{Capabilities: []string{"answer_questions", "search"}}}
	if err := ValidateCapabilities(ctx, "search"); err != nil {
		t.Fatalf("unexpected: %v", err)
	}
	if err := ValidateCapabilities(ctx, "write"); err == nil {
		t.Fatalf("expected error for disallowed capability")
	}
}

func TestEnforceGuardrails(t *testing.T) {
	ctx := &corectx.Context{Guardrails: corectx.Guardrails{Format: "markdown", MaxTokens: 5}}
	out, err := EnforceGuardrails(ctx, "# Title\nBody with more than five tokens")
	if err != nil {
		t.Fatalf("unexpected: %v", err)
	}
	if out == "" {
		t.Fatalf("expected content")
	}
}
