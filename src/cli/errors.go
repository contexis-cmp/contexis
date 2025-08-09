package main

import (
	"fmt"
)

// Structured errors for user-facing CLI
type ConfigError struct {
	Path   string
	Reason string
}

func (e ConfigError) Error() string {
	return fmt.Sprintf("config error in %s: %s", e.Path, e.Reason)
}

type ValidationError struct {
	Field   string
	Value   string
	Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("validation error for %s (%s): %s", e.Field, e.Value, e.Message)
}

type ProjectError struct {
	Name   string
	Action string
	Reason string
}

func (e ProjectError) Error() string {
	return fmt.Sprintf("project %s %s failed: %s", e.Name, e.Action, e.Reason)
}

type SecurityError struct {
	Operation string
	Reason    string
}

func (e SecurityError) Error() string {
	return fmt.Sprintf("security violation in %s: %s", e.Operation, e.Reason)
}
