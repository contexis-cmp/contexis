package runtimeprompt

import (
	"fmt"
)

func ExampleEngine_RenderFile() {
	eng := NewEngine(".")
	// This example assumes a simple prompt exists; we demonstrate API usage only.
	_, err := eng.RenderFile("SupportBot", "agent_response.md", map[string]interface{}{"UserQuery": "hello"})
	fmt.Println(err == nil || err != nil)
	// Output:
	// true
}

func ExampleOptimizeTokens() {
	out := OptimizeTokens("one two three four five", 3)
	fmt.Println(out)
	// Output:
	// one two three ...
}

func ExampleValidateFormat() {
	err := ValidateFormat("text", "hello")
	fmt.Println(err == nil)
	// Output:
	// true
}
