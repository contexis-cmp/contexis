package context

import (
	"fmt"
	"strings"
)

func ExampleNew() {
	c := New("SupportBot", "1.0.0")
	fmt.Println(c.Name, c.Version)
	// Output:
	// SupportBot 1.0.0
}

func ExampleContext_Validate() {
	c := New("MyAgent", "0.1.0")
	c.Role.Persona = "Helpful assistant"
	fmt.Println(c.Validate() == nil)
	// Output:
	// true
}

func ExampleContext_GetSHA() {
	c := New("MyAgent", "0.1.0")
	sha, _ := c.GetSHA()
	fmt.Println(strings.HasPrefix(sha, "sha256:"))
	// Output:
	// true
}
