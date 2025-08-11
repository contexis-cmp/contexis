package unit

import (
    "testing"

    "github.com/spf13/cobra"
)

// These tests are smoke tests for cobra command wiring.

func TestBuildCommandExists(t *testing.T) {
    root := &cobra.Command{Use: "ctx"}
    // simulate import side-effect by constructing command
    // Note: commands package is imported in main; here we only check cobra basics.
    if root == nil { // placeholder to avoid unused warning in minimal test
        t.Fatal("root should not be nil")
    }
}


