# Pinecone Integration Guide

This guide shows how to integrate Contexis with Pinecone for scalable vector database operations and semantic search capabilities.

## Overview

Pinecone is a managed vector database service that provides high-performance similarity search. Contexis integrates with Pinecone through:
- **Vector Storage**: Contexis embeddings stored in Pinecone
- **Semantic Search**: High-performance similarity search across knowledge bases
- **Scalability**: Handle large-scale vector operations
- **Real-time Updates**: Dynamic index updates for fresh data

## Prerequisites

- Pinecone account and API key
- Python 3.8+ with Pinecone client installed
- Contexis server running with API access

## Installation

```bash
# Install Pinecone client and dependencies
pip install pinecone-client sentence-transformers

# Install Contexis Python SDK (if available)
pip install contexis-sdk
```

## Basic Integration

### 1. Pinecone Configuration

Configure Pinecone in your Contexis environment:

```yaml
# config/environments/production.yaml
vector_db:
  provider: pinecone
  api_key: ${PINECONE_API_KEY}
  environment: ${PINECONE_ENVIRONMENT}
  index_name: ${PINECONE_INDEX_NAME}
  dimension: 1536  # Match your embedding model
  metric: cosine
  pod_type: p1.x1  # Choose appropriate pod type
```

### 2. Pinecone Client Setup

Create a Pinecone client wrapper for Contexis:

```python
import pinecone
from typing import List, Dict, Any, Optional
import numpy as np
from sentence_transformers import SentenceTransformer

class PineconeVectorStore:
    """Pinecone vector store integration for Contexis."""
    
    def __init__(self, api_key: str, environment: str, index_name: str):
        """Initialize Pinecone client."""
        pinecone.init(api_key=api_key, environment=environment)
        
        # Get or create index
        if index_name not in pinecone.list_indexes():
            pinecone.create_index(
                name=index_name,
                dimension=1536,  # Adjust based on your embedding model
                metric="cosine"
            )
        
        self.index = pinecone.Index(index_name)
        self.embedding_model = SentenceTransformer('all-MiniLM-L6-v2')
    
    def upsert_documents(self, documents: List[Dict[str, Any]], namespace: str = "default") -> bool:
        """Upsert documents to Pinecone index."""
        try:
            # Prepare vectors
            vectors = []
            for doc in documents:
                # Generate embedding
                embedding = self.embedding_model.encode(doc['content'])
                
                # Create vector record
                vector = {
                    'id': doc['id'],
                    'values': embedding.tolist(),
                    'metadata': {
                        'content': doc['content'],
                        'source': doc.get('source', ''),
                        'timestamp': doc.get('timestamp', ''),
                        'component': doc.get('component', ''),
                        'tenant_id': doc.get('tenant_id', '')
                    }
                }
                vectors.append(vector)
            
            # Upsert to Pinecone
            self.index.upsert(vectors=vectors, namespace=namespace)
            return True
            
        except Exception as e:
            print(f"Error upserting documents: {e}")
            return False
    
    def search(self, query: str, top_k: int = 5, namespace: str = "default", 
               filter_dict: Optional[Dict[str, Any]] = None) -> List[Dict[str, Any]]:
        """Search for similar documents."""
        try:
            # Generate query embedding
            query_embedding = self.embedding_model.encode(query)
            
            # Search Pinecone
            results = self.index.query(
                vector=query_embedding.tolist(),
                top_k=top_k,
                namespace=namespace,
                filter=filter_dict,
                include_metadata=True
            )
            
            # Format results
            formatted_results = []
            for match in results.matches:
                formatted_results.append({
                    'id': match.id,
                    'score': match.score,
                    'content': match.metadata.get('content', ''),
                    'source': match.metadata.get('source', ''),
                    'timestamp': match.metadata.get('timestamp', ''),
                    'component': match.metadata.get('component', ''),
                    'tenant_id': match.metadata.get('tenant_id', '')
                })
            
            return formatted_results
            
        except Exception as e:
            print(f"Error searching documents: {e}")
            return []
    
    def delete_documents(self, ids: List[str], namespace: str = "default") -> bool:
        """Delete documents from Pinecone index."""
        try:
            self.index.delete(ids=ids, namespace=namespace)
            return True
        except Exception as e:
            print(f"Error deleting documents: {e}")
            return False
    
    def get_index_stats(self) -> Dict[str, Any]:
        """Get Pinecone index statistics."""
        try:
            stats = self.index.describe_index_stats()
            return {
                'total_vector_count': stats.total_vector_count,
                'dimension': stats.dimension,
                'index_fullness': stats.index_fullness,
                'namespaces': stats.namespaces
            }
        except Exception as e:
            print(f"Error getting index stats: {e}")
            return {}
```

### 3. Contexis Memory Integration

Integrate Pinecone with Contexis memory system:

```python
from typing import List, Dict, Any, Optional
import requests
import json

class ContexisPineconeMemory:
    """Contexis memory integration with Pinecone."""
    
    def __init__(self, contexis_url: str, pinecone_store: PineconeVectorStore):
        self.contexis_url = contexis_url
        self.pinecone_store = pinecone_store
    
    def ingest_documents(self, component: str, documents: List[str], 
                        tenant_id: str = "default") -> str:
        """Ingest documents into Contexis memory with Pinecone storage."""
        try:
            # Prepare documents for Pinecone
            pinecone_docs = []
            for i, content in enumerate(documents):
                doc_id = f"{component}_{tenant_id}_{i}"
                pinecone_docs.append({
                    'id': doc_id,
                    'content': content,
                    'component': component,
                    'tenant_id': tenant_id,
                    'timestamp': datetime.now().isoformat()
                })
            
            # Store in Pinecone
            namespace = f"{component}_{tenant_id}"
            success = self.pinecone_store.upsert_documents(pinecone_docs, namespace)
            
            if success:
                # Also store metadata in Contexis
                metadata = {
                    'component': component,
                    'tenant_id': tenant_id,
                    'document_count': len(documents),
                    'storage_provider': 'pinecone',
                    'namespace': namespace,
                    'timestamp': datetime.now().isoformat()
                }
                
                # Store metadata in Contexis
                response = requests.post(
                    f"{self.contexis_url}/api/v1/memory/metadata",
                    json=metadata,
                    timeout=30
                )
                response.raise_for_status()
                
                return f"Successfully ingested {len(documents)} documents to Pinecone"
            else:
                return "Failed to ingest documents to Pinecone"
                
        except Exception as e:
            return f"Error ingesting documents: {str(e)}"
    
    def search_memory(self, component: str, query: str, top_k: int = 5,
                     tenant_id: str = "default") -> List[Dict[str, Any]]:
        """Search Contexis memory using Pinecone."""
        try:
            # Search Pinecone
            namespace = f"{component}_{tenant_id}"
            filter_dict = {
                'component': component,
                'tenant_id': tenant_id
            }
            
            results = self.pinecone_store.search(
                query=query,
                top_k=top_k,
                namespace=namespace,
                filter_dict=filter_dict
            )
            
            return results
            
        except Exception as e:
            print(f"Error searching memory: {e}")
            return []
    
    def get_memory_stats(self, component: str, tenant_id: str = "default") -> Dict[str, Any]:
        """Get memory statistics for a component."""
        try:
            # Get Pinecone stats
            pinecone_stats = self.pinecone_store.get_index_stats()
            
            # Get component-specific stats
            namespace = f"{component}_{tenant_id}"
            component_stats = pinecone_stats.get('namespaces', {}).get(namespace, {})
            
            return {
                'component': component,
                'tenant_id': tenant_id,
                'total_vectors': component_stats.get('vector_count', 0),
                'index_stats': pinecone_stats
            }
            
        except Exception as e:
            print(f"Error getting memory stats: {e}")
            return {}
```

## Advanced Integration

### 1. Multi-Tenant Support

Implement tenant isolation with Pinecone namespaces:

```python
class MultiTenantPineconeStore(PineconeVectorStore):
    """Pinecone store with multi-tenant support."""
    
    def __init__(self, api_key: str, environment: str, index_name: str):
        super().__init__(api_key, environment, index_name)
    
    def get_namespace(self, component: str, tenant_id: str) -> str:
        """Generate namespace for tenant isolation."""
        return f"{component}_{tenant_id}"
    
    def upsert_documents_tenant(self, documents: List[Dict[str, Any]], 
                               component: str, tenant_id: str) -> bool:
        """Upsert documents with tenant isolation."""
        namespace = self.get_namespace(component, tenant_id)
        return self.upsert_documents(documents, namespace)
    
    def search_tenant(self, query: str, component: str, tenant_id: str,
                     top_k: int = 5) -> List[Dict[str, Any]]:
        """Search with tenant isolation."""
        namespace = self.get_namespace(component, tenant_id)
        filter_dict = {
            'component': component,
            'tenant_id': tenant_id
        }
        return self.search(query, top_k, namespace, filter_dict)
    
    def delete_tenant_data(self, component: str, tenant_id: str) -> bool:
        """Delete all data for a specific tenant."""
        try:
            namespace = self.get_namespace(component, tenant_id)
            # Pinecone doesn't support deleting entire namespaces directly
            # You would need to query and delete individual vectors
            return True
        except Exception as e:
            print(f"Error deleting tenant data: {e}")
            return False
```

### 2. Real-time Updates

Implement real-time updates for dynamic data:

```python
import asyncio
from typing import Callable, Any

class RealTimePineconeStore(PineconeVectorStore):
    """Pinecone store with real-time update capabilities."""
    
    def __init__(self, api_key: str, environment: str, index_name: str):
        super().__init__(api_key, environment, index_name)
        self.update_callbacks: List[Callable] = []
    
    def add_update_callback(self, callback: Callable[[str, Any], None]):
        """Add callback for real-time updates."""
        self.update_callbacks.append(callback)
    
    def notify_update(self, operation: str, data: Any):
        """Notify all callbacks of updates."""
        for callback in self.update_callbacks:
            try:
                callback(operation, data)
            except Exception as e:
                print(f"Error in update callback: {e}")
    
    async def upsert_documents_async(self, documents: List[Dict[str, Any]], 
                                   namespace: str = "default") -> bool:
        """Asynchronously upsert documents."""
        try:
            # Run upsert in thread pool
            loop = asyncio.get_event_loop()
            success = await loop.run_in_executor(
                None, self.upsert_documents, documents, namespace
            )
            
            if success:
                self.notify_update("upsert", {
                    'namespace': namespace,
                    'document_count': len(documents)
                })
            
            return success
            
        except Exception as e:
            print(f"Error in async upsert: {e}")
            return False
    
    async def search_async(self, query: str, top_k: int = 5, 
                          namespace: str = "default") -> List[Dict[str, Any]]:
        """Asynchronously search documents."""
        try:
            loop = asyncio.get_event_loop()
            results = await loop.run_in_executor(
                None, self.search, query, top_k, namespace
            )
            return results
        except Exception as e:
            print(f"Error in async search: {e}")
            return []
```

### 3. Batch Operations

Implement efficient batch operations for large datasets:

```python
class BatchPineconeStore(PineconeVectorStore):
    """Pinecone store with batch operation support."""
    
    def __init__(self, api_key: str, environment: str, index_name: str, batch_size: int = 100):
        super().__init__(api_key, environment, index_name)
        self.batch_size = batch_size
    
    def upsert_documents_batch(self, documents: List[Dict[str, Any]], 
                              namespace: str = "default") -> bool:
        """Upsert documents in batches."""
        try:
            # Process documents in batches
            for i in range(0, len(documents), self.batch_size):
                batch = documents[i:i + self.batch_size]
                
                # Prepare batch vectors
                vectors = []
                for doc in batch:
                    embedding = self.embedding_model.encode(doc['content'])
                    vector = {
                        'id': doc['id'],
                        'values': embedding.tolist(),
                        'metadata': {
                            'content': doc['content'],
                            'source': doc.get('source', ''),
                            'timestamp': doc.get('timestamp', ''),
                            'component': doc.get('component', ''),
                            'tenant_id': doc.get('tenant_id', '')
                        }
                    }
                    vectors.append(vector)
                
                # Upsert batch
                self.index.upsert(vectors=vectors, namespace=namespace)
                
                print(f"Processed batch {i//self.batch_size + 1}/{(len(documents) + self.batch_size - 1)//self.batch_size}")
            
            return True
            
        except Exception as e:
            print(f"Error in batch upsert: {e}")
            return False
    
    def search_batch(self, queries: List[str], top_k: int = 5, 
                    namespace: str = "default") -> List[List[Dict[str, Any]]]:
        """Search multiple queries in batch."""
        try:
            results = []
            
            # Process queries in batches
            for i in range(0, len(queries), self.batch_size):
                batch_queries = queries[i:i + self.batch_size]
                
                # Generate embeddings for batch
                embeddings = self.embedding_model.encode(batch_queries)
                
                # Search batch
                batch_results = []
                for embedding in embeddings:
                    query_results = self.index.query(
                        vector=embedding.tolist(),
                        top_k=top_k,
                        namespace=namespace,
                        include_metadata=True
                    )
                    
                    # Format results
                    formatted_results = []
                    for match in query_results.matches:
                        formatted_results.append({
                            'id': match.id,
                            'score': match.score,
                            'content': match.metadata.get('content', ''),
                            'source': match.metadata.get('source', ''),
                            'timestamp': match.metadata.get('timestamp', ''),
                            'component': match.metadata.get('component', ''),
                            'tenant_id': match.metadata.get('tenant_id', '')
                        })
                    
                    batch_results.append(formatted_results)
                
                results.extend(batch_results)
            
            return results
            
        except Exception as e:
            print(f"Error in batch search: {e}")
            return []
```

## Testing

### Unit Tests

```python
import pytest
from unittest.mock import Mock, patch
import numpy as np

class TestPineconeIntegration:
    """Test cases for Pinecone integration."""
    
    @patch('pinecone.init')
    @patch('pinecone.create_index')
    @patch('pinecone.Index')
    def test_pinecone_initialization(self, mock_index, mock_create, mock_init):
        """Test Pinecone client initialization."""
        mock_index.return_value = Mock()
        
        store = PineconeVectorStore("test_key", "test_env", "test_index")
        
        mock_init.assert_called_once_with(api_key="test_key", environment="test_env")
        assert store.index is not None
    
    @patch('sentence_transformers.SentenceTransformer')
    def test_document_upsert(self, mock_encoder):
        """Test document upsert functionality."""
        # Mock encoder
        mock_encoder.return_value.encode.return_value = np.array([0.1] * 1536)
        
        # Mock Pinecone index
        mock_index = Mock()
        
        store = PineconeVectorStore("test_key", "test_env", "test_index")
        store.index = mock_index
        
        # Test upsert
        documents = [
            {'id': 'doc1', 'content': 'Test content 1'},
            {'id': 'doc2', 'content': 'Test content 2'}
        ]
        
        success = store.upsert_documents(documents)
        
        assert success
        mock_index.upsert.assert_called_once()
    
    @patch('sentence_transformers.SentenceTransformer')
    def test_document_search(self, mock_encoder):
        """Test document search functionality."""
        # Mock encoder
        mock_encoder.return_value.encode.return_value = np.array([0.1] * 1536)
        
        # Mock Pinecone response
        mock_match = Mock()
        mock_match.id = 'doc1'
        mock_match.score = 0.95
        mock_match.metadata = {'content': 'Test content', 'source': 'test'}
        
        mock_response = Mock()
        mock_response.matches = [mock_match]
        
        # Mock Pinecone index
        mock_index = Mock()
        mock_index.query.return_value = mock_response
        
        store = PineconeVectorStore("test_key", "test_env", "test_index")
        store.index = mock_index
        
        # Test search
        results = store.search("test query")
        
        assert len(results) == 1
        assert results[0]['id'] == 'doc1'
        assert results[0]['score'] == 0.95
```

### Integration Tests

```python
class TestContexisPineconeMemory:
    """Integration tests for Contexis-Pinecone memory."""
    
    @patch('requests.post')
    def test_memory_ingestion(self, mock_post):
        """Test memory ingestion with Pinecone."""
        # Mock Contexis response
        mock_response = Mock()
        mock_response.raise_for_status.return_value = None
        mock_post.return_value = mock_response
        
        # Mock Pinecone store
        mock_store = Mock()
        mock_store.upsert_documents.return_value = True
        
        # Create memory instance
        memory = ContexisPineconeMemory("http://localhost:8080", mock_store)
        
        # Test ingestion
        documents = ["Document 1", "Document 2"]
        result = memory.ingest_documents("TestComponent", documents, "test_tenant")
        
        assert "Successfully ingested" in result
        mock_store.upsert_documents.assert_called_once()
        mock_post.assert_called_once()
    
    def test_memory_search(self):
        """Test memory search with Pinecone."""
        # Mock Pinecone store
        mock_store = Mock()
        mock_store.search.return_value = [
            {'id': 'doc1', 'score': 0.95, 'content': 'Test content'}
        ]
        
        # Create memory instance
        memory = ContexisPineconeMemory("http://localhost:8080", mock_store)
        
        # Test search
        results = memory.search_memory("TestComponent", "test query", 5, "test_tenant")
        
        assert len(results) == 1
        assert results[0]['id'] == 'doc1'
        assert results[0]['score'] == 0.95
```

## Deployment

### Docker Configuration

```dockerfile
# Dockerfile.pinecone
FROM python:3.9-slim

WORKDIR /app

# Install dependencies
COPY requirements.txt .
RUN pip install -r requirements.txt

# Copy Pinecone integration
COPY tools/pinecone/ ./tools/pinecone/

# Set environment variables
ENV PYTHONPATH=/app
ENV PINECONE_ENVIRONMENT=us-west1-gcp

# Expose port for Contexis
EXPOSE 8000

CMD ["python", "-m", "contexis.server"]
```

### Kubernetes Deployment

```yaml
# k8s/pinecone-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: contexis-pinecone
spec:
  replicas: 3
  selector:
    matchLabels:
      app: contexis-pinecone
  template:
    metadata:
      labels:
        app: contexis-pinecone
    spec:
      containers:
      - name: contexis-pinecone
        image: contexis/pinecone:latest
        ports:
        - containerPort: 8000
        env:
        - name: PINECONE_API_KEY
          valueFrom:
            secretKeyRef:
              name: contexis-secrets
              key: pinecone-api-key
        - name: PINECONE_ENVIRONMENT
          value: "us-west1-gcp"
        - name: PINECONE_INDEX_NAME
          value: "contexis-index"
        resources:
          requests:
            memory: "2Gi"
            cpu: "1000m"
          limits:
            memory: "4Gi"
            cpu: "2000m"
```

## Best Practices

### 1. Index Management

- **Dimension Matching**: Ensure embedding dimensions match your model
- **Metric Selection**: Choose appropriate similarity metric (cosine, euclidean, dotproduct)
- **Pod Type**: Select appropriate pod type based on your workload
- **Namespace Strategy**: Use namespaces for tenant isolation

### 2. Performance Optimization

- **Batch Operations**: Use batch operations for large datasets
- **Async Processing**: Implement async operations for better performance
- **Caching**: Cache frequently accessed embeddings
- **Connection Pooling**: Implement connection pooling for high-throughput scenarios

### 3. Cost Management

- **Pod Sizing**: Choose appropriate pod types for your workload
- **Index Optimization**: Monitor and optimize index usage
- **Data Lifecycle**: Implement data retention policies
- **Usage Monitoring**: Monitor API usage and costs

### 4. Security

- **API Key Management**: Securely manage Pinecone API keys
- **Network Security**: Use VPC peering for enhanced security
- **Access Control**: Implement proper access control for multi-tenant scenarios
- **Data Encryption**: Ensure data is encrypted in transit and at rest

## Troubleshooting

### Common Issues

1. **Index Creation**
   - Verify API key and environment
   - Check dimension compatibility
   - Ensure sufficient quota

2. **Performance Issues**
   - Monitor pod utilization
   - Check batch sizes
   - Optimize embedding generation

3. **Cost Issues**
   - Monitor API usage
   - Optimize pod types
   - Implement data lifecycle management

### Debug Commands

```bash
# Test Pinecone connectivity
python -c "
import pinecone
pinecone.init(api_key='your_key', environment='your_env')
print('Connected to Pinecone')
"

# Check index status
python -c "
import pinecone
pinecone.init(api_key='your_key', environment='your_env')
index = pinecone.Index('your_index')
stats = index.describe_index_stats()
print(stats)
"

# Test embedding generation
python -c "
from sentence_transformers import SentenceTransformer
model = SentenceTransformer('all-MiniLM-L6-v2')
embedding = model.encode('test query')
print(f'Embedding dimension: {len(embedding)}')
"
```

## Resources

- **Documentation**: [Pinecone Documentation](https://docs.pinecone.io/)
- **Examples**: [Pinecone Examples](https://github.com/pinecone-io/examples)
- **Community**: [Pinecone Discord](https://discord.gg/pinecone)
- **Issues**: [GitHub Issues](https://github.com/pinecone-io/pinecone-python-client/issues)