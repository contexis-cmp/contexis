# LangGraph Integration Guide

This guide shows how to integrate Contexis with LangGraph for building complex, stateful AI workflows and multi-agent systems.

## Overview

LangGraph is a library for building stateful, multi-actor applications with LLMs. Contexis integrates with LangGraph through:
- **Custom Nodes**: Contexis components as LangGraph nodes
- **State Management**: Contexis memory as persistent state
- **Multi-Agent Workflows**: Contexis contexts as specialized agents
- **Conditional Routing**: Dynamic workflow paths based on Contexis responses

## Prerequisites

- Python 3.8+ with LangGraph installed
- Contexis server running with API access
- API key for Contexis (if authentication enabled)

## Installation

```bash
# Install LangGraph and dependencies
pip install langgraph langchain langchain-core

# Install Contexis Python SDK (if available)
pip install contexis-sdk
```

## Basic Integration

### 1. Contexis Node for LangGraph

Create a custom LangGraph node that wraps Contexis functionality:

```python
from langgraph.graph import StateGraph, END
from typing import Dict, Any, TypedDict, Annotated
import requests
from datetime import datetime

# Define the state schema
class WorkflowState(TypedDict):
    query: str
    context: str
    contexis_response: str
    memory_results: list
    current_step: str
    error: str
    metadata: Dict[str, Any]

class ContexisNode:
    """LangGraph node that integrates with Contexis."""
    
    def __init__(self, base_url: str = "http://localhost:8080", api_key: str = None):
        self.base_url = base_url
        self.api_key = api_key
        self.headers = {"Content-Type": "application/json"}
        if api_key:
            self.headers["Authorization"] = f"Bearer {api_key}"
    
    def chat_node(self, state: WorkflowState) -> WorkflowState:
        """Node for Contexis chat operations."""
        try:
            response = requests.post(
                f"{self.base_url}/api/v1/chat",
                headers=self.headers,
                json={
                    "context": state["context"],
                    "query": state["query"],
                    "top_k": 5
                },
                timeout=30
            )
            response.raise_for_status()
            
            return {
                **state,
                "contexis_response": response.json()["rendered"],
                "current_step": "chat_completed",
                "metadata": {
                    **state.get("metadata", {}),
                    "chat_timestamp": datetime.now().isoformat(),
                    "model_used": response.json().get("model", "unknown")
                }
            }
        except Exception as e:
            return {
                **state,
                "error": f"Chat error: {str(e)}",
                "current_step": "error"
            }
    
    def memory_search_node(self, state: WorkflowState) -> WorkflowState:
        """Node for Contexis memory search operations."""
        try:
            response = requests.post(
                f"{self.base_url}/api/v1/memory/search",
                headers=self.headers,
                json={
                    "provider": "sqlite",
                    "component": state["context"],
                    "query": state["query"],
                    "top_k": 5
                }
            )
            response.raise_for_status()
            
            return {
                **state,
                "memory_results": response.json().get("results", []),
                "current_step": "memory_search_completed",
                "metadata": {
                    **state.get("metadata", {}),
                    "search_timestamp": datetime.now().isoformat(),
                    "results_count": len(response.json().get("results", []))
                }
            }
        except Exception as e:
            return {
                **state,
                "error": f"Memory search error: {str(e)}",
                "current_step": "error"
            }
    
    def memory_ingest_node(self, state: WorkflowState) -> WorkflowState:
        """Node for Contexis memory ingest operations."""
        try:
            content = state.get("content_to_ingest", state["query"])
            response = requests.post(
                f"{self.base_url}/api/v1/memory/ingest",
                headers=self.headers,
                json={
                    "provider": "sqlite",
                    "component": state["context"],
                    "input": content
                }
            )
            response.raise_for_status()
            
            return {
                **state,
                "current_step": "memory_ingest_completed",
                "metadata": {
                    **state.get("metadata", {}),
                    "ingest_timestamp": datetime.now().isoformat(),
                    "ingest_success": True
                }
            }
        except Exception as e:
            return {
                **state,
                "error": f"Memory ingest error: {str(e)}",
                "current_step": "error"
            }

# Usage example
contexis_node = ContexisNode(
    base_url="http://localhost:8080",
    api_key="your_api_key_here"
)
```

### 2. Basic Workflow Graph

Create a simple workflow that uses Contexis for chat and memory operations:

```python
from langgraph.graph import StateGraph, END
from typing import Dict, Any

def create_basic_workflow() -> StateGraph:
    """Create a basic workflow with Contexis integration."""
    
    # Initialize Contexis node
    contexis = ContexisNode(
        base_url="http://localhost:8080",
        api_key="your_api_key_here"
    )
    
    # Create the workflow graph
    workflow = StateGraph(WorkflowState)
    
    # Add nodes
    workflow.add_node("chat", contexis.chat_node)
    workflow.add_node("memory_search", contexis.memory_search_node)
    workflow.add_node("memory_ingest", contexis.memory_ingest_node)
    
    # Define the workflow
    workflow.set_entry_point("chat")
    workflow.add_edge("chat", "memory_search")
    workflow.add_edge("memory_search", END)
    
    return workflow.compile()

# Usage
workflow = create_basic_workflow()

# Run the workflow
initial_state = {
    "query": "What is your return policy?",
    "context": "CustomerDocs",
    "contexis_response": "",
    "memory_results": [],
    "current_step": "start",
    "error": "",
    "metadata": {}
}

result = workflow.invoke(initial_state)
print(result)
```

## Advanced Integration

### 1. Multi-Step Research Workflow

Create a complex workflow that orchestrates multiple Contexis components:

```python
from langgraph.graph import StateGraph, END
from typing import Dict, Any, List
import json

class ResearchWorkflowState(TypedDict):
    topic: str
    research_query: str
    research_results: str
    analysis_query: str
    analysis_results: str
    summary_query: str
    final_summary: str
    current_step: str
    error: str
    metadata: Dict[str, Any]

def create_research_workflow() -> StateGraph:
    """Create a research workflow using multiple Contexis contexts."""
    
    contexis = ContexisNode(
        base_url="http://localhost:8080",
        api_key="your_api_key_here"
    )
    
    # Define custom nodes for the research workflow
    def research_node(state: ResearchWorkflowState) -> ResearchWorkflowState:
        """Research phase using Contexis ResearchBot."""
        try:
            response = requests.post(
                f"{contexis.base_url}/api/v1/chat",
                headers=contexis.headers,
                json={
                    "context": "ResearchBot",
                    "query": f"Research the topic: {state['topic']}",
                    "data": {
                        "research_depth": "comprehensive",
                        "include_sources": True
                    }
                }
            )
            response.raise_for_status()
            
            return {
                **state,
                "research_results": response.json()["rendered"],
                "current_step": "research_completed"
            }
        except Exception as e:
            return {
                **state,
                "error": f"Research error: {str(e)}",
                "current_step": "error"
            }
    
    def analysis_node(state: ResearchWorkflowState) -> ResearchWorkflowState:
        """Analysis phase using Contexis AnalysisBot."""
        try:
            analysis_query = f"Analyze the following research findings and identify key insights: {state['research_results']}"
            
            response = requests.post(
                f"{contexis.base_url}/api/v1/chat",
                headers=contexis.headers,
                json={
                    "context": "AnalysisBot",
                    "query": analysis_query,
                    "data": {
                        "analysis_type": "comprehensive",
                        "include_recommendations": True
                    }
                }
            )
            response.raise_for_status()
            
            return {
                **state,
                "analysis_results": response.json()["rendered"],
                "current_step": "analysis_completed"
            }
        except Exception as e:
            return {
                **state,
                "error": f"Analysis error: {str(e)}",
                "current_step": "error"
            }
    
    def summary_node(state: ResearchWorkflowState) -> ResearchWorkflowState:
        """Summary phase using Contexis SummaryBot."""
        try:
            summary_query = f"Create a concise executive summary based on the analysis: {state['analysis_results']}"
            
            response = requests.post(
                f"{contexis.base_url}/api/v1/chat",
                headers=contexis.headers,
                json={
                    "context": "SummaryBot",
                    "query": summary_query,
                    "data": {
                        "summary_format": "executive",
                        "max_length": 500
                    }
                }
            )
            response.raise_for_status()
            
            return {
                **state,
                "final_summary": response.json()["rendered"],
                "current_step": "summary_completed"
            }
        except Exception as e:
            return {
                **state,
                "error": f"Summary error: {str(e)}",
                "current_step": "error"
            }
    
    def error_handler(state: ResearchWorkflowState) -> ResearchWorkflowState:
        """Handle errors in the workflow."""
        return {
            **state,
            "final_summary": f"Error occurred during processing: {state['error']}",
            "current_step": "error_handled"
        }
    
    # Create the workflow graph
    workflow = StateGraph(ResearchWorkflowState)
    
    # Add nodes
    workflow.add_node("research", research_node)
    workflow.add_node("analysis", analysis_node)
    workflow.add_node("summary", summary_node)
    workflow.add_node("error_handler", error_handler)
    
    # Define the workflow
    workflow.set_entry_point("research")
    workflow.add_edge("research", "analysis")
    workflow.add_edge("analysis", "summary")
    workflow.add_edge("summary", END)
    
    # Add error handling
    workflow.add_conditional_edges(
        "research",
        lambda state: "error_handler" if state.get("error") else "analysis"
    )
    workflow.add_conditional_edges(
        "analysis",
        lambda state: "error_handler" if state.get("error") else "summary"
    )
    workflow.add_conditional_edges(
        "summary",
        lambda state: "error_handler" if state.get("error") else END
    )
    workflow.add_edge("error_handler", END)
    
    return workflow.compile()

# Usage
research_workflow = create_research_workflow()

initial_state = {
    "topic": "Artificial Intelligence trends in 2024",
    "research_query": "",
    "research_results": "",
    "analysis_query": "",
    "analysis_results": "",
    "summary_query": "",
    "final_summary": "",
    "current_step": "start",
    "error": "",
    "metadata": {}
}

result = research_workflow.invoke(initial_state)
print(result["final_summary"])
```

### 2. Customer Support Workflow

Create a customer support workflow with dynamic routing:

```python
from langgraph.graph import StateGraph, END
from typing import Dict, Any

class SupportWorkflowState(TypedDict):
    customer_id: str
    customer_query: str
    query_type: str
    knowledge_response: str
    order_info: str
    final_response: str
    escalation_needed: bool
    current_step: str
    error: str
    metadata: Dict[str, Any]

def create_support_workflow() -> StateGraph:
    """Create a customer support workflow with dynamic routing."""
    
    contexis = ContexisNode(
        base_url="http://localhost:8080",
        api_key="your_api_key_here"
    )
    
    def classify_query(state: SupportWorkflowState) -> SupportWorkflowState:
        """Classify the customer query type."""
        try:
            response = requests.post(
                f"{contexis.base_url}/api/v1/chat",
                headers=contexis.headers,
                json={
                    "context": "QueryClassifier",
                    "query": f"Classify this customer query: {state['customer_query']}",
                    "data": {
                        "classification_options": ["order_status", "return_request", "product_info", "technical_support", "billing"]
                    }
                }
            )
            response.raise_for_status()
            
            # Extract classification from response
            classification = response.json()["rendered"].lower()
            if "order" in classification:
                query_type = "order_status"
            elif "return" in classification:
                query_type = "return_request"
            elif "product" in classification:
                query_type = "product_info"
            elif "technical" in classification:
                query_type = "technical_support"
            elif "billing" in classification:
                query_type = "billing"
            else:
                query_type = "general"
            
            return {
                **state,
                "query_type": query_type,
                "current_step": "query_classified"
            }
        except Exception as e:
            return {
                **state,
                "error": f"Classification error: {str(e)}",
                "current_step": "error"
            }
    
    def knowledge_lookup(state: SupportWorkflowState) -> SupportWorkflowState:
        """Look up information in knowledge base."""
        try:
            response = requests.post(
                f"{contexis.base_url}/api/v1/chat",
                headers=contexis.headers,
                json={
                    "context": "CustomerDocs",
                    "query": state["customer_query"],
                    "data": {
                        "query_type": state["query_type"]
                    }
                }
            )
            response.raise_for_status()
            
            return {
                **state,
                "knowledge_response": response.json()["rendered"],
                "current_step": "knowledge_lookup_completed"
            }
        except Exception as e:
            return {
                **state,
                "error": f"Knowledge lookup error: {str(e)}",
                "current_step": "error"
            }
    
    def order_lookup(state: SupportWorkflowState) -> SupportWorkflowState:
        """Look up order information."""
        # Simulate order lookup
        order_info = f"Order status for customer {state['customer_id']}: Processing"
        
        return {
            **state,
            "order_info": order_info,
            "current_step": "order_lookup_completed"
        }
    
    def generate_response(state: SupportWorkflowState) -> SupportWorkflowState:
        """Generate final response to customer."""
        try:
            context = f"""
            Customer Query: {state['customer_query']}
            Query Type: {state['query_type']}
            Knowledge Base Response: {state['knowledge_response']}
            Order Information: {state.get('order_info', 'N/A')}
            """
            
            response = requests.post(
                f"{contexis.base_url}/api/v1/chat",
                headers=contexis.headers,
                json={
                    "context": "ResponseGenerator",
                    "query": f"Generate a helpful response to the customer based on this context: {context}",
                    "data": {
                        "tone": "professional",
                        "include_order_info": state["query_type"] == "order_status"
                    }
                }
            )
            response.raise_for_status()
            
            return {
                **state,
                "final_response": response.json()["rendered"],
                "current_step": "response_generated"
            }
        except Exception as e:
            return {
                **state,
                "error": f"Response generation error: {str(e)}",
                "current_step": "error"
            }
    
    def check_escalation(state: SupportWorkflowState) -> SupportWorkflowState:
        """Check if escalation is needed."""
        try:
            escalation_query = f"Based on this customer interaction, does this need escalation? Query: {state['customer_query']}, Response: {state['final_response']}"
            
            response = requests.post(
                f"{contexis.base_url}/api/v1/chat",
                headers=contexis.headers,
                json={
                    "context": "EscalationChecker",
                    "query": escalation_query,
                    "data": {
                        "escalation_criteria": ["complex_technical", "billing_dispute", "urgent_issue"]
                    }
                }
            )
            response.raise_for_status()
            
            escalation_needed = "yes" in response.json()["rendered"].lower()
            
            return {
                **state,
                "escalation_needed": escalation_needed,
                "current_step": "escalation_checked"
            }
        except Exception as e:
            return {
                **state,
                "error": f"Escalation check error: {str(e)}",
                "current_step": "error"
            }
    
    # Create the workflow graph
    workflow = StateGraph(SupportWorkflowState)
    
    # Add nodes
    workflow.add_node("classify", classify_query)
    workflow.add_node("knowledge_lookup", knowledge_lookup)
    workflow.add_node("order_lookup", order_lookup)
    workflow.add_node("generate_response", generate_response)
    workflow.add_node("check_escalation", check_escalation)
    
    # Define the workflow with conditional routing
    workflow.set_entry_point("classify")
    
    # Route based on query type
    def route_by_query_type(state: SupportWorkflowState) -> str:
        if state["query_type"] == "order_status":
            return "order_lookup"
        else:
            return "knowledge_lookup"
    
    workflow.add_conditional_edges("classify", route_by_query_type)
    
    # Continue workflow
    workflow.add_edge("knowledge_lookup", "generate_response")
    workflow.add_edge("order_lookup", "generate_response")
    workflow.add_edge("generate_response", "check_escalation")
    workflow.add_edge("check_escalation", END)
    
    return workflow.compile()

# Usage
support_workflow = create_support_workflow()

initial_state = {
    "customer_id": "CUST123",
    "customer_query": "I need help with my order #12345. What's the status?",
    "query_type": "",
    "knowledge_response": "",
    "order_info": "",
    "final_response": "",
    "escalation_needed": False,
    "current_step": "start",
    "error": "",
    "metadata": {}
}

result = support_workflow.invoke(initial_state)
print(f"Final Response: {result['final_response']}")
print(f"Escalation Needed: {result['escalation_needed']}")
```

### 3. Stateful Memory Integration

Create a workflow that maintains state across multiple interactions:

```python
from langgraph.graph import StateGraph, END
from typing import Dict, Any, List

class ConversationState(TypedDict):
    session_id: str
    user_messages: List[str]
    ai_responses: List[str]
    conversation_context: str
    current_query: str
    memory_results: List[Dict[str, Any]]
    current_step: str
    error: str
    metadata: Dict[str, Any]

def create_conversation_workflow() -> StateGraph:
    """Create a conversation workflow with persistent memory."""
    
    contexis = ContexisNode(
        base_url="http://localhost:8080",
        api_key="your_api_key_here"
    )
    
    def update_conversation_context(state: ConversationState) -> ConversationState:
        """Update conversation context with new message."""
        context = f"""
        Previous conversation:
        {chr(10).join([f"User: {msg}" for msg in state['user_messages']])}
        {chr(10).join([f"AI: {resp}" for resp in state['ai_responses']])}
        
        Current query: {state['current_query']}
        """
        
        return {
            **state,
            "conversation_context": context,
            "current_step": "context_updated"
        }
    
    def search_memory(state: ConversationState) -> ConversationState:
        """Search conversation memory for relevant context."""
        try:
            response = requests.post(
                f"{contexis.base_url}/api/v1/memory/search",
                headers=contexis.headers,
                json={
                    "provider": "sqlite",
                    "component": "ConversationMemory",
                    "query": state["current_query"],
                    "top_k": 3
                }
            )
            response.raise_for_status()
            
            return {
                **state,
                "memory_results": response.json().get("results", []),
                "current_step": "memory_searched"
            }
        except Exception as e:
            return {
                **state,
                "error": f"Memory search error: {str(e)}",
                "current_step": "error"
            }
    
    def generate_response(state: ConversationState) -> ConversationState:
        """Generate AI response with conversation context."""
        try:
            memory_context = ""
            if state["memory_results"]:
                memory_context = f"Relevant previous context: {chr(10).join([result.get('content', '') for result in state['memory_results']])}"
            
            full_context = f"""
            {state['conversation_context']}
            
            {memory_context}
            
            Generate a helpful response to: {state['current_query']}
            """
            
            response = requests.post(
                f"{contexis.base_url}/api/v1/chat",
                headers=contexis.headers,
                json={
                    "context": "ConversationBot",
                    "query": full_context,
                    "data": {
                        "conversation_history": state["conversation_context"],
                        "memory_context": memory_context
                    }
                }
            )
            response.raise_for_status()
            
            ai_response = response.json()["rendered"]
            
            return {
                **state,
                "ai_responses": state["ai_responses"] + [ai_response],
                "current_step": "response_generated"
            }
        except Exception as e:
            return {
                **state,
                "error": f"Response generation error: {str(e)}",
                "current_step": "error"
            }
    
    def save_to_memory(state: ConversationState) -> ConversationState:
        """Save conversation to memory."""
        try:
            conversation_data = {
                "session_id": state["session_id"],
                "user_message": state["current_query"],
                "ai_response": state["ai_responses"][-1],
                "timestamp": datetime.now().isoformat()
            }
            
            response = requests.post(
                f"{contexis.base_url}/api/v1/memory/ingest",
                headers=contexis.headers,
                json={
                    "provider": "sqlite",
                    "component": "ConversationMemory",
                    "input": json.dumps(conversation_data)
                }
            )
            response.raise_for_status()
            
            return {
                **state,
                "current_step": "memory_saved"
            }
        except Exception as e:
            return {
                **state,
                "error": f"Memory save error: {str(e)}",
                "current_step": "error"
            }
    
    # Create the workflow graph
    workflow = StateGraph(ConversationState)
    
    # Add nodes
    workflow.add_node("update_context", update_conversation_context)
    workflow.add_node("search_memory", search_memory)
    workflow.add_node("generate_response", generate_response)
    workflow.add_node("save_memory", save_to_memory)
    
    # Define the workflow
    workflow.set_entry_point("update_context")
    workflow.add_edge("update_context", "search_memory")
    workflow.add_edge("search_memory", "generate_response")
    workflow.add_edge("generate_response", "save_memory")
    workflow.add_edge("save_memory", END)
    
    return workflow.compile()

# Usage
conversation_workflow = create_conversation_workflow()

# Initial state
initial_state = {
    "session_id": "session_123",
    "user_messages": ["Hello, I need help with my account"],
    "ai_responses": ["Hello! I'd be happy to help you with your account. What specific issue are you experiencing?"],
    "conversation_context": "",
    "current_query": "I can't log in to my account",
    "memory_results": [],
    "current_step": "start",
    "error": "",
    "metadata": {}
}

result = conversation_workflow.invoke(initial_state)
print(f"AI Response: {result['ai_responses'][-1]}")
```

## Error Handling and Monitoring

### 1. Error Handling Nodes

```python
def create_error_handling_workflow() -> StateGraph:
    """Create a workflow with comprehensive error handling."""
    
    def error_handler(state: Dict[str, Any]) -> Dict[str, Any]:
        """Handle errors in the workflow."""
        error_message = state.get("error", "Unknown error")
        
        # Log error
        print(f"Error in workflow: {error_message}")
        
        # Try to recover or provide fallback
        fallback_response = f"I apologize, but I encountered an error: {error_message}. Please try again or contact support."
        
        return {
            **state,
            "final_response": fallback_response,
            "current_step": "error_handled",
            "error": error_message
        }
    
    def retry_node(state: Dict[str, Any]) -> Dict[str, Any]:
        """Retry a failed operation."""
        max_retries = state.get("metadata", {}).get("retry_count", 0)
        
        if max_retries < 3:
            return {
                **state,
                "metadata": {
                    **state.get("metadata", {}),
                    "retry_count": max_retries + 1
                },
                "current_step": "retry"
            }
        else:
            return {
                **state,
                "error": "Max retries exceeded",
                "current_step": "max_retries_exceeded"
            }
    
    # Create workflow with error handling
    workflow = StateGraph(Dict[str, Any])
    
    workflow.add_node("error_handler", error_handler)
    workflow.add_node("retry", retry_node)
    
    # Add conditional edges for error handling
    workflow.add_conditional_edges(
        "main_node",
        lambda state: "retry" if state.get("error") and state.get("metadata", {}).get("retry_count", 0) < 3 else "error_handler" if state.get("error") else END
    )
    
    return workflow.compile()
```

### 2. Monitoring and Logging

```python
import logging
from datetime import datetime

class WorkflowMonitor:
    """Monitor and log workflow execution."""
    
    def __init__(self):
        self.logger = logging.getLogger("contexis_workflow")
        self.logger.setLevel(logging.INFO)
        
        # Add file handler
        handler = logging.FileHandler("workflow.log")
        formatter = logging.Formatter('%(asctime)s - %(name)s - %(levelname)s - %(message)s')
        handler.setFormatter(formatter)
        self.logger.addHandler(handler)
    
    def log_step(self, step_name: str, state: Dict[str, Any]):
        """Log workflow step execution."""
        self.logger.info(f"Step: {step_name}, State: {state.get('current_step', 'unknown')}")
        
        if state.get("error"):
            self.logger.error(f"Error in step {step_name}: {state['error']}")
    
    def log_completion(self, workflow_name: str, final_state: Dict[str, Any]):
        """Log workflow completion."""
        duration = final_state.get("metadata", {}).get("duration", 0)
        self.logger.info(f"Workflow {workflow_name} completed in {duration}s")
    
    def log_error(self, error: str, context: Dict[str, Any]):
        """Log workflow errors."""
        self.logger.error(f"Workflow error: {error}, Context: {context}")

# Usage in workflow nodes
monitor = WorkflowMonitor()

def monitored_node(state: Dict[str, Any]) -> Dict[str, Any]:
    """Node with monitoring."""
    start_time = datetime.now()
    
    try:
        # Node logic here
        result = some_operation(state)
        
        # Log success
        monitor.log_step("node_name", result)
        
        return result
    except Exception as e:
        # Log error
        monitor.log_error(str(e), state)
        return {
            **state,
            "error": str(e),
            "current_step": "error"
        }
```

## Best Practices

### 1. State Management

- **Immutable Updates**: Always create new state objects
- **Type Safety**: Use TypedDict for state schemas
- **Validation**: Validate state at each step
- **Persistence**: Save important state to external storage

### 2. Error Handling

- **Graceful Degradation**: Provide fallback responses
- **Retry Logic**: Implement exponential backoff
- **Circuit Breakers**: Prevent cascading failures
- **Monitoring**: Log all errors and performance metrics

### 3. Performance Optimization

- **Async Operations**: Use async nodes for I/O operations
- **Caching**: Cache frequently accessed data
- **Batch Processing**: Process multiple items together
- **Resource Management**: Clean up resources properly

### 4. Security

- **Input Validation**: Validate all inputs
- **Authentication**: Use secure API keys
- **Rate Limiting**: Implement rate limiting
- **Audit Logging**: Log all operations

## Troubleshooting

### Common Issues

1. **State Persistence**
   - Check state serialization
   - Verify state updates are immutable
   - Monitor memory usage

2. **Node Failures**
   - Implement proper error handling
   - Add retry logic
   - Monitor node performance

3. **Memory Issues**
   - Clean up old state
   - Implement pagination
   - Monitor memory usage

### Debug Tools

```python
def debug_workflow(workflow: StateGraph, initial_state: Dict[str, Any]):
    """Debug workflow execution."""
    
    def debug_node(state: Dict[str, Any]) -> Dict[str, Any]:
        """Debug wrapper for nodes."""
        print(f"Entering node with state: {state}")
        
        try:
            # Execute node logic
            result = node_logic(state)
            print(f"Node completed successfully: {result}")
            return result
        except Exception as e:
            print(f"Node failed with error: {e}")
            raise
    
    # Add debug wrapper to all nodes
    for node_name in workflow.nodes:
        original_node = workflow.nodes[node_name]
        workflow.nodes[node_name] = debug_node
    
    # Run workflow with debugging
    return workflow.invoke(initial_state)
```

## Support

For LangGraph integration support:
- **Documentation**: [docs.contexis.dev/integrations/langgraph](https://docs.contexis.dev/integrations/langgraph)
- **Community**: [Discord Integration Channel](https://discord.gg/contexis)
- **Issues**: [GitHub Issues](https://github.com/contexis-cmp/contexis/issues)
