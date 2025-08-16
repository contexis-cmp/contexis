# CMP Framework Testing Guide

This directory contains comprehensive tests for the CMP (Context-Memory-Prompt) framework, following Test-Driven Development (TDD) principles.

## Test Architecture

### Test Structure
```
tests/
├── unit/                    # Unit tests for individual components
│   ├── helpers/            # Test utilities and fixtures
│   ├── agent_generator_test.go
│   └── cli_commands_test.go
├── integration/            # Integration tests for component interactions
│   └── agent_generator_integration_test.go
├── e2e/                   # End-to-end tests for complete workflows
├── fixtures/              # Test data and fixtures
├── coverage/              # Coverage reports
├── reports/               # Test reports
├── temp/                  # Temporary test files
├── test_config.yaml       # Test configuration
├── test_runner.go         # Test runner utilities
└── README.md              # This file
```

### Test Categories

1. **Unit Tests**: Test individual functions and components in isolation
2. **Integration Tests**: Test interactions between components
3. **End-to-End Tests**: Test complete user workflows
4. **Performance Tests**: Test system performance under load
5. **Security Tests**: Test security and privacy features

##  Running Tests

### Prerequisites
```bash
# Install Go testing dependencies
go get github.com/stretchr/testify/assert
go get github.com/stretchr/testify/require
go get github.com/stretchr/testify/suite
```

### Using ctx test (recommended)
```bash
# All Go suites with coverage and JUnit
ctx test --all --coverage --junit --out tests/reports

# Specific suite
ctx test --unit --coverage

# By category from tests/test_config.yaml
ctx test --category core_components --coverage

# Drift detection for all components
ctx test --drift-detection --out tests/reports

# Drift detection for a single component with semantic similarity and baseline update
ctx test --drift-detection --component CustomerDocs --semantic --update-baseline --junit --out tests/reports
```

Artifacts:
- `tests/reports/go_<suite>.txt`, `tests/reports/go_tests.json`, optional `junit-go.xml`
- `tests/reports/drift_<Component>.json`, `tests/reports/drift_index.json`, optional `junit-drift.xml`
- Coverage profiles in `tests/coverage/*.out`

### Direct Go Test Commands
```bash
# Run unit tests
go test ./tests/unit/... -v

# Run integration tests
go test ./tests/integration/... -v

# Run with coverage
go test ./tests/... -v -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html

# Run specific test file
go test ./tests/unit/agent_generator_test.go -v

# Run specific test function
go test ./tests/unit/ -v -run TestAgentConfigValidation
```

##  Test Coverage Requirements

### Minimum Coverage Targets
- **Unit Tests**: 80% coverage
- **Integration Tests**: 70% coverage
- **End-to-End Tests**: 60% coverage

### Coverage Reports
Coverage reports are generated in `tests/coverage/`:
- `coverage.html` - HTML coverage report
- `unit.out` - Unit test coverage data
- `integration.out` - Integration test coverage data
- `e2e.out` - E2E test coverage data

##  Test-Driven Development (TDD) Workflow

### 1. Write Failing Test (Red)
```go
func TestNewFeature(t *testing.T) {
    // Write test for feature that doesn't exist yet
    result := NewFeature()
    assert.Equal(t, expected, result)
}
```

### 2. Write Minimal Implementation (Green)
```go
func NewFeature() string {
    return "expected" // Minimal implementation to make test pass
}
```

### 3. Refactor (Refactor)
```go
func NewFeature() string {
    // Refactor to improve code quality while keeping tests green
    return calculateExpectedValue()
}
```

### 4. Repeat
Continue the cycle for each new feature or bug fix.

##  Test Utilities and Helpers

### TestFixtures
Provides common test data:
```go
fixtures := &helpers.TestFixtures{}
validNames := fixtures.ValidAgentNames()
validTools := fixtures.ValidTools()
validConfigs := fixtures.ValidAgentConfigs()
```

### TestUtils
Provides utility functions:
```go
utils := &helpers.TestUtils{}
utils.CreateTestProjectStructure(t, basePath)
utils.AssertFileContains(t, filePath, expectedContent)
utils.ValidateYAMLContent(t, content)
```

### Test Helpers
Common test functions:
```go
ctx := helpers.TestContext(t)
tempDir := helpers.CreateTempDir(t)
helpers.AssertFileExists(t, filePath)
helpers.AssertFileContent(t, filePath, expectedContent)
```

##  Test Configuration

### Test Configuration File
Tests are configured via `tests/test_config.yaml`:

```yaml
test_suites:
  unit:
    enabled: true
    timeout: 30s
    parallel: true
    coverage_threshold: 80
```

### Environment Variables
```bash
CMP_ENV=test              # Test environment
CMP_LOG_LEVEL=debug       # Log level for tests
CMP_TEMP_DIR=tests/temp   # Temporary directory
```

##  Test Categories

### Agent Generator Tests
- **Unit Tests**: Test individual agent generation functions
- **Integration Tests**: Test complete agent generation workflow
- **Validation Tests**: Test input validation and error handling

### CLI Command Tests
- **Command Structure**: Test command hierarchy and flags
- **Argument Validation**: Test input validation
- **Execution Flow**: Test command execution and output

### Core Component Tests
- **Context Management**: Test context loading and validation
- **Memory Management**: Test memory operations
- **Prompt Management**: Test template processing

### Security Tests
- **Input Validation**: Test malicious input handling
- **Access Control**: Test permission enforcement
- **Data Privacy**: Test data isolation and encryption

### Performance Tests
- **Load Testing**: Test system under load
- **Memory Usage**: Test memory consumption
- **Response Time**: Test response latency

##  Debugging Tests

### Verbose Output
```bash
go test -v ./tests/unit/
```

### Test Debugging
```bash
# Run single test with debug output
go test -v -run TestSpecificFunction ./tests/unit/

# Run with race detection
go test -race ./tests/unit/

# Run with memory profiling
go test -memprofile=mem.out ./tests/unit/
```

### Test Logs
Tests use structured logging with Zap:
```go
logger := helpers.TestLogger(t)
logger.Info("Test step completed", zap.String("step", "validation"))
```

##  Continuous Integration

### CI/CD Pipeline
Tests are automatically run in CI/CD:
- Unit tests on every commit
- Integration tests on pull requests
- E2E tests on merge to main
- Performance tests on release candidates

### Test Reports
Test results are published to:
- `tests/reports/` - Detailed test reports
- Coverage reports in HTML format
- Performance benchmarks

##  Best Practices

### Test Naming
- Use descriptive test names: `TestAgentConfigValidation`
- Group related tests in test suites
- Use table-driven tests for multiple scenarios

### Test Organization
- One test file per component
- Group related tests in test suites
- Use helper functions for common operations

### Test Data
- Use fixtures for test data
- Clean up test data after tests
- Use temporary directories for file operations

### Assertions
- Use specific assertions: `assert.Equal`, `assert.Contains`
- Provide meaningful error messages
- Test both positive and negative cases

### Error Handling
- Test error conditions
- Verify error messages
- Test edge cases and boundary conditions

##  Common Issues

### Import Path Issues
If you encounter import path issues:
```bash
# Ensure you're in the project root
cd /path/to/contexis

# Update Go modules
go mod tidy

# Check module path
go list -m
```

### Test Environment Issues
If tests fail due to environment issues:
```bash
# Clean test artifacts
make clean

# Rebuild test environment
make setup

# Run tests with fresh environment
make test
```

### Coverage Issues
If coverage is below targets:
```bash
# Generate detailed coverage report
make test-coverage

# Review uncovered code
open tests/coverage/coverage.html

# Add tests for uncovered code
# Focus on critical paths first
```

##  Additional Resources

- [Go Testing Package](https://golang.org/pkg/testing/)
- [Testify Documentation](https://github.com/stretchr/testify)
- [Test-Driven Development](https://en.wikipedia.org/wiki/Test-driven_development)
- [Go Testing Best Practices](https://golang.org/doc/tutorial/testing)

##  Contributing

When adding new tests:
1. Follow TDD principles
2. Add tests for new functionality
3. Update this documentation
4. Ensure tests pass in CI/CD
5. Maintain test coverage targets
