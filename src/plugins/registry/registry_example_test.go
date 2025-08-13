package registry

import (
	"fmt"
)

func ExampleNewRegistry() {
	r := NewRegistry(".")
	list, err := r.List()
	fmt.Println(err == nil && (len(list) >= 0))
	// Output:
	// true
}

func ExampleRegistry_Info() {
	r := NewRegistry(".")
	_, err := r.Info("nonexistent")
	fmt.Println(err != nil)
	// Output:
	// true
}
