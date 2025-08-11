package commands

import (
	"context"
	"fmt"
	"os"
	"text/template"

	"github.com/contexis-cmp/contexis/src/cli/logger"
	"go.uber.org/zap"
)

// generateRAGTests creates test configuration for drift detection
func generateRAGTests(ctx context.Context, config RAGConfig) error {
	log := logger.WithContext(ctx)

	// Create drift detection test configuration
	testConfigPath := fmt.Sprintf("tests/%s/rag_drift_test.yaml", config.Name)

	testConfigTemplate := `# Drift Detection Tests for {{.Name}} RAG System

test_cases:
  - name: "basic_search_functionality"
    input: "What is the main topic?"
    expected_similarity: 0.85
    required_keywords: ["sample", "document"]
    forbidden_phrases: ["I don't know", "no information"]
    
  - name: "specific_information_retrieval"
    input: "How do I use this system?"
    expected_similarity: 0.80
    required_keywords: ["usage", "documents", "format"]
    
  - name: "no_results_handling"
    input: "completely unrelated query that should not match"
    expected_similarity: 0.30
    required_response: "couldn't find"
    
  - name: "response_format_consistency"
    input: "What is the key information?"
    expected_format: "markdown"
    required_sections: ["User Query", "Response", "Sources"]
    forbidden_sections: ["I apologize", "I'm not sure"]

drift_thresholds:
  similarity_threshold: 0.85
  response_time_threshold: 2000  # milliseconds
  token_count_threshold: 1000

business_rules:
  - name: "must_cite_sources"
    description: "All responses must include source citations"
    validation: "response_contains_sources"
    
  - name: "no_speculation"
    description: "Responses should not speculate beyond document content"
    validation: "no_speculative_phrases"
    
  - name: "consistent_format"
    description: "Responses should follow consistent markdown format"
    validation: "markdown_format_check"

performance_benchmarks:
  search_response_time: 1000  # milliseconds
  embedding_generation_time: 500  # milliseconds
  memory_usage: 512  # MB
`

	tmpl, err := template.New("test_config").Parse(testConfigTemplate)
	if err != nil {
		log.Error("failed to parse test config template", zap.Error(err))
		return fmt.Errorf("failed to parse test config template: %w", err)
	}

	file, err := os.Create(testConfigPath)
	if err != nil {
		log.Error("failed to create test config file", zap.String("path", testConfigPath), zap.Error(err))
		return fmt.Errorf("failed to create test config file: %w", err)
	}
	defer file.Close()

	if err := tmpl.Execute(file, config); err != nil {
		log.Error("failed to execute test config template", zap.Error(err))
		return fmt.Errorf("failed to execute test config template: %w", err)
	}

	// Create Python test script
	testScriptPath := fmt.Sprintf("tests/%s/test_rag.py", config.Name)

	testScriptTemplate := `#!/usr/bin/env python3
"""
Test script for {{.Name}} RAG system
"""

import sys
import os
import json
import yaml
from pathlib import Path

# Add tools to path
sys.path.append(os.path.join(os.path.dirname(__file__), '../../tools/{{.Name}}'))

from semantic_search import {{.Name}}SemanticSearch

def load_test_config():
    """Load test configuration"""
    config_path = Path(__file__).parent / "rag_drift_test.yaml"
    with open(config_path, 'r') as f:
        return yaml.safe_load(f)

def run_drift_tests():
    """Run drift detection tests"""
    config = load_test_config()
    search = {{.Name}}SemanticSearch()
    
    print(f"Running drift tests for {{.Name}} RAG system...")
    
    results = {
        'passed': 0,
        'failed': 0,
        'tests': []
    }
    
    for test_case in config['test_cases']:
        print(f"\\nTesting: {test_case['name']}")
        
        try:
            # Perform search
            search_results = search.search(test_case['input'])
            
            # Basic validation
            if not search_results:
                if 'no_results_handling' in test_case['name']:
                    results['passed'] += 1
                    results['tests'].append({
                        'name': test_case['name'],
                        'status': 'PASSED',
                        'reason': 'Correctly handled no results'
                    })
                else:
                    results['failed'] += 1
                    results['tests'].append({
                        'name': test_case['name'],
                        'status': 'FAILED',
                        'reason': 'No search results returned'
                    })
                continue
            
            # Check similarity threshold
            best_result = search_results[0]
            if best_result['similarity'] >= test_case.get('expected_similarity', 0.7):
                results['passed'] += 1
                results['tests'].append({
                    'name': test_case['name'],
                    'status': 'PASSED',
                    'similarity': best_result['similarity']
                })
            else:
                results['failed'] += 1
                results['tests'].append({
                    'name': test_case['name'],
                    'status': 'FAILED',
                    'reason': f"Similarity {best_result['similarity']} below threshold {test_case['expected_similarity']}"
                })
                
        except Exception as e:
            results['failed'] += 1
            results['tests'].append({
                'name': test_case['name'],
                'status': 'ERROR',
                'reason': str(e)
            })
    
    # Print results
    print(f"\\n=== Test Results ===")
    print(f"Passed: {results['passed']}")
    print(f"Failed: {results['failed']}")
    print(f"Total: {results['passed'] + results['failed']}")
    
    for test in results['tests']:
        status_icon = "" if test['status'] == 'PASSED' else ""
        print(f"{status_icon} {test['name']}: {test['status']}")
        if 'reason' in test:
            print(f"   Reason: {test['reason']}")
    
    return results['failed'] == 0

if __name__ == "__main__":
    success = run_drift_tests()
    sys.exit(0 if success else 1)
`

	tmpl, err = template.New("test_script").Parse(testScriptTemplate)
	if err != nil {
		log.Error("failed to parse test script template", zap.Error(err))
		return fmt.Errorf("failed to parse test script template: %w", err)
	}

	file, err = os.Create(testScriptPath)
	if err != nil {
		log.Error("failed to create test script file", zap.String("path", testScriptPath), zap.Error(err))
		return fmt.Errorf("failed to create test script file: %w", err)
	}
	defer file.Close()

	if err := tmpl.Execute(file, config); err != nil {
		log.Error("failed to execute test script template", zap.Error(err))
		return fmt.Errorf("failed to execute test script template: %w", err)
	}

	// Make the test script executable
	if err := os.Chmod(testScriptPath, 0755); err != nil {
		log.Error("failed to make test script executable", zap.Error(err))
	}

	log.Info("RAG tests generated", zap.String("test_config", testConfigPath))
	return nil
}
