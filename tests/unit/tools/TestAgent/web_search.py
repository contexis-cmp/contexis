#!/usr/bin/env python3
"""
Web Search Tool for CMP Agents
Provides web search capabilities using DuckDuckGo API
"""

import requests
import json
from typing import List, Dict, Optional
from dataclasses import dataclass
import logging

logger = logging.getLogger(__name__)

@dataclass
class SearchResult:
    """Represents a single search result"""
    title: str
    url: str
    snippet: str
    source: str
    relevance_score: float

class WebSearchTool:
    """Web search tool for agents"""
    
    def __init__(self, api_key: Optional[str] = None):
        self.api_key = api_key
        self.base_url = "https://api.duckduckgo.com/"
        
    def search(self, query: str, max_results: int = 5) -> List[SearchResult]:
        """
        Perform web search using DuckDuckGo API
        
        Args:
            query: Search query string
            max_results: Maximum number of results to return
            
        Returns:
            List of SearchResult objects
        """
        try:
            logger.info(f"Performing web search for: {query}")
            
            # Use DuckDuckGo Instant Answer API
            params = {
                'q': query,
                'format': 'json',
                'no_html': '1',
                'skip_disambig': '1'
            }
            
            response = requests.get(self.base_url, params=params, timeout=10)
            response.raise_for_status()
            
            data = response.json()
            results = []
            
            # Extract results from response
            if 'AbstractURL' in data and data['AbstractURL']:
                results.append(SearchResult(
                    title=data.get('Abstract', 'No title'),
                    url=data['AbstractURL'],
                    snippet=data.get('Abstract', ''),
                    source='DuckDuckGo',
                    relevance_score=0.9
                ))
            
            # Add related topics if available
            if 'RelatedTopics' in data:
                for topic in data['RelatedTopics'][:max_results-1]:
                    if 'Text' in topic and 'FirstURL' in topic:
                        results.append(SearchResult(
                            title=topic.get('Text', 'No title')[:100],
                            url=topic['FirstURL'],
                            snippet=topic.get('Text', ''),
                            source='DuckDuckGo',
                            relevance_score=0.7
                        ))
            
            logger.info(f"Found {len(results)} search results")
            return results[:max_results]
            
        except requests.RequestException as e:
            logger.error(f"Web search failed: {e}")
            return []
        except Exception as e:
            logger.error(f"Unexpected error in web search: {e}")
            return []
    
    def search_news(self, query: str, max_results: int = 5) -> List[SearchResult]:
        """
        Search for recent news articles
        
        Args:
            query: News search query
            max_results: Maximum number of results
            
        Returns:
            List of news SearchResult objects
        """
        try:
            logger.info(f"Performing news search for: {query}")
            
            # Add news-specific parameters
            params = {
                'q': f"{query} news",
                'format': 'json',
                'no_html': '1'
            }
            
            response = requests.get(self.base_url, params=params, timeout=10)
            response.raise_for_status()
            
            data = response.json()
            results = []
            
            # Extract news results
            if 'AbstractURL' in data and data['AbstractURL']:
                results.append(SearchResult(
                    title=data.get('Abstract', 'No title'),
                    url=data['AbstractURL'],
                    snippet=data.get('Abstract', ''),
                    source='DuckDuckGo News',
                    relevance_score=0.85
                ))
            
            logger.info(f"Found {len(results)} news results")
            return results[:max_results]
            
        except Exception as e:
            logger.error(f"News search failed: {e}")
            return []

def main():
    """Test the web search tool"""
    tool = WebSearchTool()
    
    # Test basic search
    results = tool.search("Python programming")
    print("Basic search results:")
    for result in results:
        print(f"- {result.title}: {result.url}")
    
    # Test news search
    news_results = tool.search_news("artificial intelligence")
    print("\nNews search results:")
    for result in news_results:
        print(f"- {result.title}: {result.url}")

if __name__ == "__main__":
    main()
