#!/usr/bin/env python3
"""
Semantic Search Tool for CustomerDocs RAG System
"""

import os
import json
import logging
from typing import List, Dict, Any, Optional
from pathlib import Path

try:
    from sentence_transformers import SentenceTransformer
    import chromadb
    from chromadb.config import Settings
except ImportError:
    print("Installing required dependencies...")
    os.system("pip install sentence-transformers chromadb")
    from sentence_transformers import SentenceTransformer
    import chromadb
    from chromadb.config import Settings

logger = logging.getLogger(__name__)

class CustomerDocsSemanticSearch:
    """Semantic search implementation for CustomerDocs RAG system"""
    
    def __init__(self, db_type: str = "sqlite", embedding_model: str = "sentence-transformers"):
        self.db_type = db_type
        self.embedding_model = embedding_model
        self.model = None
        self.client = None
        self.collection = None
        
        self._initialize_components()
    
    def _initialize_components(self):
        """Initialize embedding model and vector database"""
        try:
            # Initialize embedding model
            if self.embedding_model == "sentence-transformers":
                self.model = SentenceTransformer('all-MiniLM-L6-v2')
            elif self.embedding_model == "bge-small-en":
                self.model = SentenceTransformer('BAAI/bge-small-en-v1.5')
            else:
                # Fallback to default
                self.model = SentenceTransformer('all-MiniLM-L6-v2')
            
            # Initialize vector database
            if self.db_type == "chroma":
                self.client = chromadb.PersistentClient(
                    path="./memory/CustomerDocs/vector_store",
                    settings=Settings(anonymized_telemetry=False)
                )
                self.collection = self.client.get_or_create_collection(
                    name="CustomerDocs_documents",
                    metadata={"description": "Document embeddings for CustomerDocs RAG system"}
                )
            
            logger.info(f"Initialized CustomerDocs semantic search with {self.embedding_model}")
            
        except Exception as e:
            logger.error(f"Failed to initialize semantic search: {e}")
            raise
    
    def search(self, query: str, top_k: int = 5, threshold: float = 0.7) -> List[Dict[str, Any]]:
        """
        Perform semantic search on documents
        
        Args:
            query: Search query
            top_k: Number of results to return
            threshold: Similarity threshold
            
        Returns:
            List of search results with metadata
        """
        try:
            # Generate query embedding
            query_embedding = self.model.encode([query])
            
            # Search in vector database
            if self.db_type == "chroma":
                results = self.collection.query(
                    query_embeddings=query_embedding.tolist(),
                    n_results=top_k,
                    include=["metadatas", "documents", "distances"]
                )
                
                # Format results
                formatted_results = []
                for i in range(len(results['ids'][0])):
                    distance = results['distances'][0][i]
                    similarity = 1 - distance  # Convert distance to similarity
                    
                    if similarity >= threshold:
                        formatted_results.append({
                            'id': results['ids'][0][i],
                            'content': results['documents'][0][i],
                            'metadata': results['metadatas'][0][i] or {},
                            'similarity': similarity,
                            'distance': distance
                        })
                
                return formatted_results
            
            else:
                # Fallback for other database types
                logger.warning(f"Database type {self.db_type} not fully implemented, using mock results")
                return self._mock_search(query, top_k)
                
        except Exception as e:
            logger.error(f"Search failed: {e}")
            return []
    
    def _mock_search(self, query: str, top_k: int) -> List[Dict[str, Any]]:
        """Mock search for testing purposes"""
        return [
            {
                'id': 'mock_1',
                'content': f'Mock result for query: {query}',
                'metadata': {'title': 'Mock Document', 'source': 'test'},
                'similarity': 0.85,
                'distance': 0.15
            }
        ]
    
    def add_documents(self, documents: List[Dict[str, Any]]) -> bool:
        """
        Add documents to the vector database
        
        Args:
            documents: List of documents with 'id', 'content', and 'metadata'
            
        Returns:
            Success status
        """
        try:
            if self.db_type == "chroma":
                # Prepare documents for insertion
                ids = [doc['id'] for doc in documents]
                contents = [doc['content'] for doc in documents]
                metadatas = [doc.get('metadata', {}) for doc in documents]
                
                # Generate embeddings
                embeddings = self.model.encode(contents)
                
                # Add to collection
                self.collection.add(
                    ids=ids,
                    embeddings=embeddings.tolist(),
                    documents=contents,
                    metadatas=metadatas
                )
                
                logger.info(f"Added {len(documents)} documents to CustomerDocs collection")
                return True
                
        except Exception as e:
            logger.error(f"Failed to add documents: {e}")
            return False
    
    def get_stats(self) -> Dict[str, Any]:
        """Get collection statistics"""
        try:
            if self.db_type == "chroma":
                count = self.collection.count()
                return {
                    'total_documents': count,
                    'database_type': self.db_type,
                    'embedding_model': self.embedding_model
                }
        except Exception as e:
            logger.error(f"Failed to get stats: {e}")
        
        return {'error': 'Unable to retrieve statistics'}

# Example usage
if __name__ == "__main__":
    # Initialize search
    search = CustomerDocsSemanticSearch()
    
    # Example search
    results = search.search("What is the main topic?")
    print(json.dumps(results, indent=2))
