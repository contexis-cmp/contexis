package schema

import (
    "encoding/json"
    "fmt"

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

    // name (string)
    if v, ok := m["name"].(string); !ok || v == "" {
        return fmt.Errorf("field 'name' is required and must be a non-empty string")
    }
    // version (string)
    if v, ok := m["version"].(string); !ok || v == "" {
        return fmt.Errorf("field 'version' is required and must be a non-empty string")
    }
    // role.persona (string)
    role, ok := m["role"].(map[string]interface{})
    if !ok {
        return fmt.Errorf("field 'role' is required and must be an object")
    }
    if v, ok := role["persona"].(string); !ok || v == "" {
        return fmt.Errorf("field 'role.persona' is required and must be a non-empty string")
    }
    return nil
}

// JSONMarshal is a small helper to marshal any value to JSON bytes.
func JSONMarshal(v interface{}) ([]byte, error) {
    return json.Marshal(v)
}


