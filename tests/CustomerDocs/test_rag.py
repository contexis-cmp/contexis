#!/usr/bin/env python3
"""
Test script for CustomerDocs RAG system
"""

import sys
import os
import json
import yaml
from pathlib import Path

# Add tools to path
sys.path.append(os.path.join(os.path.dirname(__file__), '../../tools/CustomerDocs'))

from semantic_search import CustomerDocsSemanticSearch

def load_test_config():
    """Load test configuration"""
    config_path = Path(__file__).parent / "rag_drift_test.yaml"
    with open(config_path, 'r') as f:
        return yaml.safe_load(f)

def run_drift_tests():
    """Run drift detection tests"""
    config = load_test_config()
    search = CustomerDocsSemanticSearch()
    
    print(f"Running drift tests for CustomerDocs RAG system...")
    
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
