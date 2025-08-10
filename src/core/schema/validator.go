package schema

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/xeipuuv/gojsonschema"
	"gopkg.in/yaml.v3"
)

// ValidateContextYAML performs a lightweight validation of a .ctx YAML file.
// It ensures required top-level fields exist and have expected types.
// This is intentionally minimal and can be upgraded to full JSON Schema validation later.
func ValidateContextYAML(yamlBytes []byte) error {
	var m map[string]interface{}
	if err := yaml.Unmarshal(yamlBytes, &m); err != nil {
		return fmt.Errorf("invalid YAML: %w", err)
	}
	// Load schema relative to project root if possible
	schemaPath := filepath.Join("src", "core", "schema", "context_schema.json")
	if _, err := os.Stat(schemaPath); err != nil {
		// fallback: basic checks if schema missing
		if v, ok := m["name"].(string); !ok || v == "" {
			return fmt.Errorf("field 'name' is required and must be a non-empty string")
		}
		if v, ok := m["version"].(string); !ok || v == "" {
			return fmt.Errorf("field 'version' is required and must be a non-empty string")
		}
		role, ok := m["role"].(map[string]interface{})
		if !ok {
			return fmt.Errorf("field 'role' is required and must be an object")
		}
		if v, ok := role["persona"].(string); !ok || v == "" {
			return fmt.Errorf("field 'role.persona' is required and must be a non-empty string")
		}
		return nil
	}
	sl := gojsonschema.NewReferenceLoader("file://" + schemaPath)
	// convert YAML to JSON for validation
	by, err := json.Marshal(m)
	if err != nil {
		return err
	}
	dl := gojsonschema.NewBytesLoader(by)
	res, err := gojsonschema.Validate(sl, dl)
	if err != nil {
		return fmt.Errorf("schema validation error: %w", err)
	}
	if !res.Valid() {
		if len(res.Errors()) > 0 {
			return fmt.Errorf("schema invalid: %s", res.Errors()[0].String())
		}
		return fmt.Errorf("schema invalid")
	}
	return nil
}

// JSONMarshal is a small helper to marshal any value to JSON bytes.
func JSONMarshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}
