package guardrails

import (
	"fmt"

	corectx "github.com/contexis-cmp/contexis/src/core/context"
)

func ExampleValidateCapabilities() {
	ctx := &corectx.Context{Role: corectx.Role{Capabilities: []string{"search", "answer"}}}
	err := ValidateCapabilities(ctx, "search")
	fmt.Println(err == nil)
	// Output:
	// true
}

func ExampleEnforceGuardrails() {
	ctx := &corectx.Context{Guardrails: corectx.Guardrails{Format: "text", MaxTokens: 3}}
	out, err := EnforceGuardrails(ctx, "one two three four five")
	fmt.Println(err == nil && out == "one two three ...")
	// Output:
	// true
}
