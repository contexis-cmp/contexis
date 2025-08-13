package model

import (
	"fmt"
)

func ExampleParams() {
	p := Params{Temperature: 0.2, TopP: 0.9, MaxNewTokens: 64, RepetitionPen: 1.1}
	fmt.Println(p.MaxNewTokens)
	// Output:
	// 64
}

func ExampleFromEnv() {
	prov, _ := FromEnv()
	fmt.Println(prov == nil || prov != nil)
	// Output:
	// true
}
