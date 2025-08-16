# LangChain Integration Guide

This guide shows how to integrate Contexis with LangChain for building sophisticated AI applications and agents.

## Overview

LangChain is a framework for developing applications powered by language models. Contexis integrates with LangChain through:
- **Custom Tools**: Contexis components as LangChain tools
- **Memory Integration**: Contexis memory as LangChain memory backends
- **Agent Integration**: Contexis contexts as LangChain agents
- **Chain Composition**: Contexis workflows as LangChain chains

## Prerequisites

- Python 3.8+ with LangChain installed
- Contexis server running with API access
- API key for Contexis (if authentication enabled)

## Installation

```bash
# Install LangChain and dependencies
pip install langchain langchain-community langchain-core

# Install Contexis Python SDK (if available)
pip install contexis-sdk
```

## Basic Integration

### 1. Contexis Tool for LangChain

Create a custom LangChain tool that wraps Contexis chat functionality:

```python
from langchain.tools import BaseTool
from typing import Optional, Type
from pydantic import BaseModel, Field
import requests
import json

class ContexisChatInput(BaseModel):
    query: str = Field(description="The question or query to send to Contexis")
    context: str = Field(default="CustomerDocs", description="The Contexis context to use")
    tenant_id: Optional[str] = Field(default=None, description="Tenant ID for multi-tenancy")

class ContexisChatTool(BaseTool):
    name = "contexis_chat"
    description = "Query Contexis AI system for information and responses"
    args_schema: Type[BaseModel] = ContexisChatInput
    
    def __init__(self, base_url: str = "http://localhost:8080", api_key: Optional[str] = None):
        super().__init__()
        self.base_url = base_url
        self.api_key = api_key
        self.headers = {
            "Content-Type": "application/json"
        }
        if api_key:
            self.headers["Authorization"] = f"Bearer {api_key}"
    
    def _run(self, query: str, context: str = "CustomerDocs", tenant_id: Optional[str] = None) -> str:
        """Execute the tool."""
        payload = {
            "context": context,
            "query": query,
            "top_k": 5
        }
        if tenant_id:
            payload["tenant_id"] = tenant_id
        
        try:
            response = requests.post(
                f"{self.base_url}/api/v1/chat",
                headers=self.headers,
                json=payload,
                timeout=30
            )
            response.raise_for_status()
            return response.json()["rendered"]
        except requests.exceptions.RequestException as e:
            return f"Error calling Contexis: {str(e)}"
    
    async def _arun(self, query: str, context: str = "CustomerDocs", tenant_id: Optional[str] = None) -> str:
        """Execute the tool asynchronously."""
        import aiohttp
        
        payload = {
            "context": context,
            "query": query,
            "top_k": 5
        }
        if tenant_id:
            payload["tenant_id"] = tenant_id
        
        try:
            async with aiohttp.ClientSession() as session:
                async with session.post(
                    f"{self.base_url}/api/v1/chat",
                    headers=self.headers,
                    json=payload,
                    timeout=aiohttp.ClientTimeout(total=30)
                ) as response:
                    response.raise_for_status()
                    result = await response.json()
                    return result["rendered"]
        except Exception as e:
            return f"Error calling Contexis: {str(e)}"

# Usage example
contexis_tool = ContexisChatTool(
    base_url="http://localhost:8080",
    api_key="your_api_key_here"
)
```

### 2. Contexis Memory Integration

Create a LangChain memory backend using Contexis memory:

```python
from langchain.memory import BaseMemory
from typing import Dict, List, Any, Optional
import requests

class ContexisMemory(BaseMemory):
    """LangChain memory backend using Contexis memory system."""
    
    def __init__(self, base_url: str, component: str, api_key: Optional[str] = None):
        super().__init__()
        self.base_url = base_url
        self.component = component
        self.api_key = api_key
        self.headers = {"Content-Type": "application/json"}
        if api_key:
            self.headers["Authorization"] = f"Bearer {api_key}"
    
    @property
    def memory_variables(self) -> List[str]:
        """Return the memory variables."""
        return ["chat_history", "current_context"]
    
    def load_memory_variables(self, inputs: Dict[str, Any]) -> Dict[str, Any]:
        """Load memory variables."""
        # Get conversation history from Contexis
        try:
            response = requests.get(
                f"{self.base_url}/api/v1/memory/episodic/{self.component}",
                headers=self.headers
            )
            if response.status_code == 200:
                history = response.json().get("history", [])
                return {
                    "chat_history": history,
                    "current_context": self._get_current_context(inputs)
                }
        except Exception as e:
            print(f"Error loading memory: {e}")
        
        return {
            "chat_history": [],
            "current_context": self._get_current_context(inputs)
        }
    
    def save_context(self, inputs: Dict[str, Any], outputs: Dict[str, str]) -> None:
        """Save context to memory."""
        # Save conversation to Contexis episodic memory
        conversation_data = {
            "input": inputs.get("input", ""),
            "output": outputs.get("output", ""),
            "timestamp": str(datetime.now())
        }
        
        try:
            requests.post(
                f"{self.base_url}/api/v1/memory/episodic/{self.component}",
                headers=self.headers,
                json=conversation_data
            )
        except Exception as e:
            print(f"Error saving to memory: {e}")
    
    def clear(self) -> None:
        """Clear memory."""
        try:
            requests.delete(
                f"{self.base_url}/api/v1/memory/episodic/{self.component}",
                headers=self.headers
            )
        except Exception as e:
            print(f"Error clearing memory: {e}")
    
    def _get_current_context(self, inputs: Dict[str, Any]) -> str:
        """Extract current context from inputs."""
        return inputs.get("input", "")

# Usage example
contexis_memory = ContexisMemory(
    base_url="http://localhost:8080",
    component="CustomerDocs",
    api_key="your_api_key_here"
)
```

## Advanced Integration

### 1. Contexis Agent Integration

Create a LangChain agent that uses Contexis as the primary reasoning engine:

```python
from langchain.agents import AgentExecutor, create_openai_functions_agent
from langchain.prompts import ChatPromptTemplate, MessagesPlaceholder
from langchain.schema import SystemMessage, HumanMessage
from langchain.tools import Tool
import requests

class ContexisAgent:
    """LangChain agent that uses Contexis for reasoning and tool execution."""
    
    def __init__(self, base_url: str, context: str, api_key: Optional[str] = None):
        self.base_url = base_url
        self.context = context
        self.api_key = api_key
        self.headers = {"Content-Type": "application/json"}
        if api_key:
            self.headers["Authorization"] = f"Bearer {api_key}"
    
    def run(self, query: str, tools: List[Tool] = None) -> str:
        """Run the agent with the given query and tools."""
        
        # Prepare the system prompt with tool descriptions
        system_prompt = f"""You are a helpful AI assistant using the {self.context} context.
        
        Available tools:
        {self._format_tools(tools) if tools else "No additional tools available."}
        
        Use the tools when needed to answer the user's question. Always provide helpful and accurate responses."""
        
        # Create the prompt template
        prompt = ChatPromptTemplate.from_messages([
            ("system", system_prompt),
            MessagesPlaceholder(variable_name="chat_history"),
            ("human", "{input}"),
            MessagesPlaceholder(variable_name="agent_scratchpad"),
        ])
        
        # Execute with Contexis
        try:
            response = requests.post(
                f"{self.base_url}/api/v1/chat",
                headers=self.headers,
                json={
                    "context": self.context,
                    "query": query,
                    "data": {
                        "system_prompt": system_prompt,
                        "tools": [tool.dict() for tool in tools] if tools else []
                    }
                }
            )
            response.raise_for_status()
            return response.json()["rendered"]
        except Exception as e:
            return f"Error running agent: {str(e)}"
    
    def _format_tools(self, tools: List[Tool]) -> str:
        """Format tools for the system prompt."""
        if not tools:
            return ""
        
        tool_descriptions = []
        for tool in tools:
            tool_descriptions.append(f"- {tool.name}: {tool.description}")
        
        return "\n".join(tool_descriptions)

# Usage example
agent = ContexisAgent(
    base_url="http://localhost:8080",
    context="CustomerDocs",
    api_key="your_api_key_here"
)

# Add LangChain tools
tools = [
    Tool(
        name="web_search",
        func=lambda q: f"Search results for: {q}",
        description="Search the web for current information"
    ),
    Tool(
        name="calculator",
        func=lambda expr: eval(expr),
        description="Perform mathematical calculations"
    )
]

response = agent.run("What is the return policy and how much is 15% off $100?", tools)
print(response)
```

### 2. Contexis Chain Integration

Create a LangChain chain that orchestrates multiple Contexis components:

```python
from langchain.chains import LLMChain
from langchain.prompts import PromptTemplate
from typing import Dict, Any
import requests

class ContexisChain:
    """LangChain chain that orchestrates multiple Contexis components."""
    
    def __init__(self, base_url: str, api_key: Optional[str] = None):
        self.base_url = base_url
        self.api_key = api_key
        self.headers = {"Content-Type": "application/json"}
        if api_key:
            self.headers["Authorization"] = f"Bearer {api_key}"
    
    def run(self, steps: List[Dict[str, Any]], initial_input: str) -> Dict[str, Any]:
        """Run a chain of Contexis operations."""
        results = {"initial_input": initial_input}
        current_input = initial_input
        
        for i, step in enumerate(steps):
            step_name = step.get("name", f"step_{i}")
            context = step.get("context", "CustomerDocs")
            query = step.get("query", current_input)
            
            try:
                response = requests.post(
                    f"{self.base_url}/api/v1/chat",
                    headers=self.headers,
                    json={
                        "context": context,
                        "query": query,
                        "data": step.get("data", {})
                    }
                )
                response.raise_for_status()
                step_result = response.json()["rendered"]
                
                results[step_name] = step_result
                current_input = step_result  # Pass result to next step
                
            except Exception as e:
                results[step_name] = f"Error in step {step_name}: {str(e)}"
                break
        
        return results

# Usage example
chain = ContexisChain(
    base_url="http://localhost:8080",
    api_key="your_api_key_here"
)

# Define the chain steps
steps = [
    {
        "name": "research",
        "context": "ResearchBot",
        "query": "Research the latest AI trends in 2024"
    },
    {
        "name": "analyze",
        "context": "AnalysisBot",
        "query": "Analyze the research findings and identify key insights"
    },
    {
        "name": "summarize",
        "context": "SummaryBot",
        "query": "Create a concise summary of the analysis"
    }
]

results = chain.run(steps, "Generate a report on AI trends")
print(results)
```

### 3. Contexis Tool for Memory Operations

Create a LangChain tool for Contexis memory operations:

```python
from langchain.tools import BaseTool
from typing import Optional, Type
from pydantic import BaseModel, Field
import requests

class ContexisMemoryInput(BaseModel):
    operation: str = Field(description="Memory operation: search, ingest, or clear")
    query: Optional[str] = Field(default=None, description="Search query or content to ingest")
    component: str = Field(default="CustomerDocs", description="Contexis component name")
    top_k: Optional[int] = Field(default=5, description="Number of results to return")

class ContexisMemoryTool(BaseTool):
    name = "contexis_memory"
    description = "Search, ingest, or clear Contexis memory"
    args_schema: Type[BaseModel] = ContexisMemoryInput
    
    def __init__(self, base_url: str = "http://localhost:8080", api_key: Optional[str] = None):
        super().__init__()
        self.base_url = base_url
        self.api_key = api_key
        self.headers = {"Content-Type": "application/json"}
        if api_key:
            self.headers["Authorization"] = f"Bearer {api_key}"
    
    def _run(self, operation: str, query: Optional[str] = None, component: str = "CustomerDocs", top_k: int = 5) -> str:
        """Execute the memory operation."""
        
        if operation == "search":
            if not query:
                return "Error: Query required for search operation"
            
            try:
                response = requests.post(
                    f"{self.base_url}/api/v1/memory/search",
                    headers=self.headers,
                    json={
                        "provider": "sqlite",
                        "component": component,
                        "query": query,
                        "top_k": top_k
                    }
                )
                response.raise_for_status()
                results = response.json()
                return f"Search results: {results}"
                
            except Exception as e:
                return f"Error searching memory: {str(e)}"
        
        elif operation == "ingest":
            if not query:
                return "Error: Content required for ingest operation"
            
            try:
                response = requests.post(
                    f"{self.base_url}/api/v1/memory/ingest",
                    headers=self.headers,
                    json={
                        "provider": "sqlite",
                        "component": component,
                        "input": query
                    }
                )
                response.raise_for_status()
                return "Content ingested successfully"
                
            except Exception as e:
                return f"Error ingesting content: {str(e)}"
        
        elif operation == "clear":
            try:
                response = requests.delete(
                    f"{self.base_url}/api/v1/memory/clear",
                    headers=self.headers,
                    json={
                        "provider": "sqlite",
                        "component": component
                    }
                )
                response.raise_for_status()
                return "Memory cleared successfully"
                
            except Exception as e:
                return f"Error clearing memory: {str(e)}"
        
        else:
            return f"Unknown operation: {operation}"

# Usage example
memory_tool = ContexisMemoryTool(
    base_url="http://localhost:8080",
    api_key="your_api_key_here"
)
```

## Complete Example: Customer Support Agent

Here's a complete example of a customer support agent using LangChain and Contexis:

```python
from langchain.agents import initialize_agent, AgentType
from langchain.llms import OpenAI
from langchain.memory import ConversationBufferMemory
from langchain.tools import Tool
import os

# Initialize Contexis tools
contexis_chat = ContexisChatTool(
    base_url="http://localhost:8080",
    api_key=os.getenv("CONTEXIS_API_KEY")
)

contexis_memory = ContexisMemoryTool(
    base_url="http://localhost:8080",
    api_key=os.getenv("CONTEXIS_API_KEY")
)

# Additional tools
def search_orders(customer_id: str) -> str:
    """Search for customer orders."""
    return f"Found orders for customer {customer_id}: Order #12345, Order #12346"

def update_order_status(order_id: str, status: str) -> str:
    """Update order status."""
    return f"Updated order {order_id} status to {status}"

# Create tools list
tools = [
    Tool(
        name="contexis_knowledge",
        func=contexis_chat._run,
        description="Search company knowledge base and get AI-powered responses"
    ),
    Tool(
        name="contexis_memory",
        func=contexis_memory._run,
        description="Search or update customer conversation history"
    ),
    Tool(
        name="search_orders",
        func=search_orders,
        description="Search for customer orders by customer ID"
    ),
    Tool(
        name="update_order_status",
        func=update_order_status,
        description="Update the status of an order"
    )
]

# Initialize the agent
llm = OpenAI(temperature=0)
memory = ConversationBufferMemory(memory_key="chat_history", return_messages=True)

agent = initialize_agent(
    tools,
    llm,
    agent=AgentType.CONVERSATIONAL_REACT_DESCRIPTION,
    memory=memory,
    verbose=True
)

# Example conversation
response = agent.run("Hi, I need help with my order #12345. What's the status?")
print(response)

response = agent.run("I want to return it. What's your return policy?")
print(response)
```

## Error Handling and Retry Logic

### 1. Retry Decorator

```python
import time
from functools import wraps
from typing import Callable, Any

def retry_on_failure(max_retries: int = 3, delay: float = 1.0):
    """Decorator to retry function calls on failure."""
    def decorator(func: Callable) -> Callable:
        @wraps(func)
        def wrapper(*args, **kwargs) -> Any:
            last_exception = None
            
            for attempt in range(max_retries):
                try:
                    return func(*args, **kwargs)
                except Exception as e:
                    last_exception = e
                    if attempt < max_retries - 1:
                        time.sleep(delay * (2 ** attempt))  # Exponential backoff
                    continue
            
            raise last_exception
        
        return wrapper
    return decorator

# Apply to Contexis tools
@retry_on_failure(max_retries=3, delay=1.0)
def robust_contexis_call(base_url: str, payload: dict, headers: dict) -> dict:
    """Make a robust call to Contexis with retry logic."""
    response = requests.post(f"{base_url}/api/v1/chat", headers=headers, json=payload)
    response.raise_for_status()
    return response.json()
```

### 2. Circuit Breaker Pattern

```python
from enum import Enum
import time
from typing import Optional

class CircuitState(Enum):
    CLOSED = "closed"
    OPEN = "open"
    HALF_OPEN = "half_open"

class CircuitBreaker:
    """Simple circuit breaker implementation."""
    
    def __init__(self, failure_threshold: int = 5, timeout: float = 60.0):
        self.failure_threshold = failure_threshold
        self.timeout = timeout
        self.state = CircuitState.CLOSED
        self.failure_count = 0
        self.last_failure_time: Optional[float] = None
    
    def call(self, func: Callable, *args, **kwargs) -> Any:
        """Execute function with circuit breaker protection."""
        
        if self.state == CircuitState.OPEN:
            if time.time() - self.last_failure_time > self.timeout:
                self.state = CircuitState.HALF_OPEN
            else:
                raise Exception("Circuit breaker is OPEN")
        
        try:
            result = func(*args, **kwargs)
            self._on_success()
            return result
        except Exception as e:
            self._on_failure()
            raise e
    
    def _on_success(self):
        """Handle successful call."""
        self.failure_count = 0
        self.state = CircuitState.CLOSED
    
    def _on_failure(self):
        """Handle failed call."""
        self.failure_count += 1
        self.last_failure_time = time.time()
        
        if self.failure_count >= self.failure_threshold:
            self.state = CircuitState.OPEN

# Usage with Contexis
circuit_breaker = CircuitBreaker(failure_threshold=3, timeout=30.0)

def safe_contexis_call(query: str) -> str:
    """Make a safe call to Contexis with circuit breaker protection."""
    return circuit_breaker.call(
        lambda: contexis_chat._run(query),
        query
    )
```

## Best Practices

### 1. Tool Design

- **Clear Descriptions**: Provide detailed descriptions for tools
- **Error Handling**: Implement robust error handling in tools
- **Type Safety**: Use Pydantic models for input validation
- **Async Support**: Implement async versions for better performance

### 2. Memory Management

- **Context Window**: Be mindful of memory context window limits
- **Cleanup**: Regularly clean up old conversations
- **Persistence**: Use persistent memory for important conversations
- **Privacy**: Implement proper data privacy controls

### 3. Performance Optimization

- **Caching**: Cache frequently accessed information
- **Batch Operations**: Use batch operations when possible
- **Connection Pooling**: Reuse HTTP connections
- **Async Operations**: Use async/await for I/O operations

### 4. Security

- **API Key Management**: Store API keys securely
- **Input Validation**: Validate all inputs before processing
- **Rate Limiting**: Implement rate limiting for API calls
- **Audit Logging**: Log all operations for security

## Troubleshooting

### Common Issues

1. **Connection Errors**
   - Check Contexis server status
   - Verify network connectivity
   - Check firewall settings

2. **Authentication Errors**
   - Verify API key is correct
   - Check API key permissions
   - Ensure proper header format

3. **Memory Issues**
   - Check available memory
   - Implement proper cleanup
   - Use streaming for large responses

### Debug Tools

```python
import logging

# Enable debug logging
logging.basicConfig(level=logging.DEBUG)

# Debug wrapper for Contexis calls
def debug_contexis_call(func):
    """Debug wrapper for Contexis API calls."""
    def wrapper(*args, **kwargs):
        print(f"Calling {func.__name__} with args: {args}, kwargs: {kwargs}")
        try:
            result = func(*args, **kwargs)
            print(f"Success: {result}")
            return result
        except Exception as e:
            print(f"Error: {e}")
            raise
    return wrapper

# Apply debug wrapper
contexis_chat._run = debug_contexis_call(contexis_chat._run)
```

## Support

For LangChain integration support:
- **Documentation**: [docs.contexis.dev/integrations/langchain](https://docs.contexis.dev/integrations/langchain)
- **Community**: [Discord Integration Channel](https://discord.gg/contexis)
- **Issues**: [GitHub Issues](https://github.com/contexis-cmp/contexis/issues)
