// Package context defines the Contexis CMP Context domain model.
//
// A Context is a declarative specification of an AI agent: role/persona,
// capabilities, tool integrations, guardrails, memory behavior, and testing
// configuration. Contexts are authored as YAML (`.ctx`) files and validated
// against a JSON Schema at load time. The types in this package map 1:1 to the
// persisted representation and are used across the runtime and tooling.
package context
