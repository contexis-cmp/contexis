// Package registry provides discovery and lifecycle management for Contexis plugins.
//
// Plugins are declared with a `plugin.json` manifest under `plugins/<name>/` and
// can advertise capabilities (e.g., memory backends, tools, providers). The
// Registry scans, installs, and removes plugin packages at runtime or via CLI.
package registry
