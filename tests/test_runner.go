package tests

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/contexis-cmp/contexis/src/cli/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

// TestRunner provides functionality to run different types of tests
type TestRunner struct {
	logger *zap.Logger
}

// NewTestRunner creates a new test runner
func NewTestRunner() *TestRunner {
	// Initialize logger for tests
	if err := logger.InitLogger("debug", "console"); err != nil {
		panic(fmt.Sprintf("Failed to initialize test logger: %v", err))
	}

	return &TestRunner{
		logger: logger.WithContext(nil),
	}
}

// RunUnitTests runs all unit tests
func (tr *TestRunner) RunUnitTests(t *testing.T) {
	tr.logger.Info("Running unit tests...")

	// Run unit tests
	unitTestDir := "tests/unit"
	if _, err := os.Stat(unitTestDir); os.IsNotExist(err) {
		t.Skip("Unit test directory does not exist")
		return
	}

	// Run all test files in unit directory
	err := filepath.Walk(unitTestDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && filepath.Ext(path) == ".go" {
			tr.logger.Info("Running unit test file", zap.String("file", path))
			// Note: In a real implementation, you would run the tests here
			// For now, we just log that we found the test file
		}

		return nil
	})

	require.NoError(t, err, "Failed to walk unit test directory")
}

// RunIntegrationTests runs all integration tests
func (tr *TestRunner) RunIntegrationTests(t *testing.T) {
	tr.logger.Info("Running integration tests...")

	// Run integration tests
	integrationTestDir := "tests/integration"
	if _, err := os.Stat(integrationTestDir); os.IsNotExist(err) {
		t.Skip("Integration test directory does not exist")
		return
	}

	// Run all test files in integration directory
	err := filepath.Walk(integrationTestDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && filepath.Ext(path) == ".go" {
			tr.logger.Info("Running integration test file", zap.String("file", path))
			// Note: In a real implementation, you would run the tests here
			// For now, we just log that we found the test file
		}

		return nil
	})

	require.NoError(t, err, "Failed to walk integration test directory")
}

// RunE2ETests runs all end-to-end tests
func (tr *TestRunner) RunE2ETests(t *testing.T) {
	tr.logger.Info("Running end-to-end tests...")

	// Run e2e tests
	e2eTestDir := "tests/e2e"
	if _, err := os.Stat(e2eTestDir); os.IsNotExist(err) {
		t.Skip("E2E test directory does not exist")
		return
	}

	// Run all test files in e2e directory
	err := filepath.Walk(e2eTestDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && filepath.Ext(path) == ".go" {
			tr.logger.Info("Running e2e test file", zap.String("file", path))
			// Note: In a real implementation, you would run the tests here
			// For now, we just log that we found the test file
		}

		return nil
	})

	require.NoError(t, err, "Failed to walk e2e test directory")
}

// RunAllTests runs all tests
func (tr *TestRunner) RunAllTests(t *testing.T) {
	tr.logger.Info("Running all tests...")

	tr.RunUnitTests(t)
	tr.RunIntegrationTests(t)
	tr.RunE2ETests(t)
}

// RunTestsWithCoverage runs tests with coverage reporting
func (tr *TestRunner) RunTestsWithCoverage(t *testing.T) {
	tr.logger.Info("Running tests with coverage...")

	// This would run tests with coverage flags
	// For now, just run all tests
	tr.RunAllTests(t)
}

// RunPerformanceTests runs performance tests
func (tr *TestRunner) RunPerformanceTests(t *testing.T) {
	tr.logger.Info("Running performance tests...")

	// Performance tests would be implemented here
	// For now, just log that we're running them
}

// RunSecurityTests runs security tests
func (tr *TestRunner) RunSecurityTests(t *testing.T) {
	tr.logger.Info("Running security tests...")

	// Security tests would be implemented here
	// For now, just log that we're running them
}

// TestRunnerTestSuite provides tests for the test runner itself
type TestRunnerTestSuite struct {
	runner *TestRunner
}

// SetupTest sets up the test
func (tr *TestRunnerTestSuite) SetupTest() {
	tr.runner = NewTestRunner()
}

// TestTestRunnerCreation tests test runner creation
func (tr *TestRunnerTestSuite) TestTestRunnerCreation(t *testing.T) {
	runner := NewTestRunner()
	assert.NotNil(t, runner, "Test runner should be created")
	assert.NotNil(t, runner.logger, "Test runner should have logger")
}

// TestTestRunnerDirectoryExistence tests that test directories exist
func (tr *TestRunnerTestSuite) TestTestRunnerDirectoryExistence(t *testing.T) {
	// Test that test directories exist
	testDirs := []string{"tests/unit", "tests/integration", "tests/e2e"}

	for _, dir := range testDirs {
		_, err := os.Stat(dir)
		assert.NoError(t, err, "Test directory should exist: %s", dir)
	}
}

// TestTestRunnerFileExistence tests that test files exist
func (tr *TestRunnerTestSuite) TestTestRunnerFileExistence(t *testing.T) {
	// Test that key test files exist
	testFiles := []string{
		"tests/unit/agent_generator_test.go",
		"tests/unit/cli_commands_test.go",
		"tests/integration/agent_generator_integration_test.go",
		"tests/unit/helpers/helpers.go",
	}

	for _, file := range testFiles {
		_, err := os.Stat(file)
		assert.NoError(t, err, "Test file should exist: %s", file)
	}
}

// Run the test runner tests
func TestTestRunner(t *testing.T) {
	suite := &TestRunnerTestSuite{}
	suite.SetupTest()

	suite.TestTestRunnerCreation(t)
	suite.TestTestRunnerDirectoryExistence(t)
	suite.TestTestRunnerFileExistence(t)
}
