# CrewAI Integration Guide

This guide shows how to integrate Contexis with CrewAI for building sophisticated multi-agent systems and collaborative AI workflows.

## Overview

CrewAI is a framework for orchestrating role-playing autonomous AI agents. Contexis integrates with CrewAI through:
- **Custom Tools**: Contexis components as CrewAI tools
- **Agent Integration**: Contexis contexts as CrewAI agents
- **Memory Sharing**: Contexis memory as shared knowledge base
- **Task Delegation**: Contexis workflows as CrewAI tasks

## Prerequisites

- Python 3.8+ with CrewAI installed
- Contexis server running with API access
- API key for Contexis (if authentication enabled)

## Installation

```bash
# Install CrewAI and dependencies
pip install crewai langchain langchain-openai

# Install Contexis Python SDK (if available)
pip install contexis-sdk
```

## Basic Integration

### 1. Contexis Tool for CrewAI

Create a custom CrewAI tool that wraps Contexis functionality:

```python
from crewai.tools import BaseTool
from typing import Optional
import requests
import json

class ContexisChatTool(BaseTool):
    name: str = "contexis_chat"
    description: str = "Query Contexis AI system for information and responses"
    
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
```

### 2. Contexis Agent for CrewAI

Create a CrewAI agent that uses Contexis as its knowledge base:

```python
from crewai import Agent, Task, Crew
from langchain_openai import ChatOpenAI

class ContexisAgent:
    """CrewAI agent that integrates with Contexis."""
    
    def __init__(self, base_url: str = "http://localhost:8080", api_key: Optional[str] = None):
        self.base_url = base_url
        self.api_key = api_key
        self.llm = ChatOpenAI(model="gpt-4", temperature=0.1)
    
    def create_researcher_agent(self) -> Agent:
        """Create a research agent that uses Contexis for knowledge retrieval."""
        return Agent(
            role="Research Analyst",
            goal="Conduct thorough research using Contexis knowledge base",
            backstory="""You are an expert research analyst with access to a comprehensive 
            knowledge base through Contexis. You excel at finding relevant information 
            and synthesizing insights.""",
            tools=[ContexisChatTool(self.base_url, self.api_key)],
            llm=self.llm,
            verbose=True
        )
    
    def create_writer_agent(self) -> Agent:
        """Create a writer agent that uses Contexis for content creation."""
        return Agent(
            role="Content Writer",
            goal="Create high-quality content based on research findings",
            backstory="""You are a skilled content writer who creates engaging and 
            informative content. You use Contexis to ensure accuracy and relevance.""",
            tools=[ContexisChatTool(self.base_url, self.api_key)],
            llm=self.llm,
            verbose=True
        )
    
    def create_reviewer_agent(self) -> Agent:
        """Create a reviewer agent that uses Contexis for quality assurance."""
        return Agent(
            role="Quality Reviewer",
            goal="Review and validate content quality and accuracy",
            backstory="""You are a meticulous quality reviewer who ensures content 
            meets high standards. You use Contexis to verify facts and check consistency.""",
            tools=[ContexisChatTool(self.base_url, self.api_key)],
            llm=self.llm,
            verbose=True
        )
```

### 3. Complete Workflow Example

Create a complete content creation workflow using CrewAI and Contexis:

```python
from crewai import Task, Crew
from typing import List, Dict, Any

class ContexisCrewAIWorkflow:
    """Complete workflow combining CrewAI and Contexis."""
    
    def __init__(self, base_url: str = "http://localhost:8080", api_key: Optional[str] = None):
        self.contexis_agent = ContexisAgent(base_url, api_key)
    
    def create_content_workflow(self, topic: str, requirements: Dict[str, Any]) -> str:
        """Create content using a multi-agent workflow."""
        
        # Create agents
        researcher = self.contexis_agent.create_researcher_agent()
        writer = self.contexis_agent.create_writer_agent()
        reviewer = self.contexis_agent.create_reviewer_agent()
        
        # Define tasks
        research_task = Task(
            description=f"""Research the topic: {topic}
            
            Requirements:
            - Use Contexis knowledge base for accurate information
            - Find relevant examples and case studies
            - Identify key points and insights
            
            Deliver a comprehensive research summary.""",
            agent=researcher,
            expected_output="Detailed research findings with key insights and examples"
        )
        
        writing_task = Task(
            description=f"""Create content based on the research findings.
            
            Topic: {topic}
            Requirements: {requirements}
            
            Use the research summary to create engaging, informative content.
            Ensure accuracy by cross-referencing with Contexis knowledge base.""",
            agent=writer,
            expected_output="High-quality content that meets all requirements",
            context=[research_task]
        )
        
        review_task = Task(
            description="""Review the created content for quality and accuracy.
            
            Check for:
            - Factual accuracy (verify with Contexis)
            - Content quality and engagement
            - Meeting all requirements
            - Grammar and style
            
            Provide feedback and suggest improvements.""",
            agent=reviewer,
            expected_output="Comprehensive review with feedback and improvements",
            context=[writing_task]
        )
        
        # Create and run the crew
        crew = Crew(
            agents=[researcher, writer, reviewer],
            tasks=[research_task, writing_task, review_task],
            verbose=True
        )
        
        result = crew.kickoff()
        return result
```

## Advanced Integration

### 1. Multi-Tenant Support

Extend the integration to support multi-tenancy:

```python
class MultiTenantContexisTool(BaseTool):
    name: str = "multi_tenant_contexis_chat"
    description: str = "Query Contexis AI system with tenant isolation"
    
    def __init__(self, base_url: str = "http://localhost:8080", api_key: Optional[str] = None):
        super().__init__()
        self.base_url = base_url
        self.api_key = api_key
        self.headers = {
            "Content-Type": "application/json"
        }
        if api_key:
            self.headers["Authorization"] = f"Bearer {api_key}"
    
    def _run(self, query: str, context: str = "CustomerDocs", tenant_id: str = "default") -> str:
        """Execute the tool with tenant isolation."""
        payload = {
            "context": context,
            "query": query,
            "tenant_id": tenant_id,
            "top_k": 5
        }
        
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
```

### 2. Memory Integration

Integrate Contexis memory with CrewAI for persistent knowledge:

```python
class ContexisMemoryTool(BaseTool):
    name: str = "contexis_memory_search"
    description: str = "Search Contexis memory for relevant information"
    
    def __init__(self, base_url: str = "http://localhost:8080", api_key: Optional[str] = None):
        super().__init__()
        self.base_url = base_url
        self.api_key = api_key
        self.headers = {
            "Content-Type": "application/json"
        }
        if api_key:
            self.headers["Authorization"] = f"Bearer {api_key}"
    
    def _run(self, query: str, component: str = "CustomerDocs", top_k: int = 5) -> str:
        """Search Contexis memory."""
        payload = {
            "component": component,
            "query": query,
            "top_k": top_k
        }
        
        try:
            response = requests.post(
                f"{self.base_url}/api/v1/memory/search",
                headers=self.headers,
                json=payload,
                timeout=30
            )
            response.raise_for_status()
            results = response.json()["results"]
            
            # Format results for CrewAI
            formatted_results = []
            for result in results:
                formatted_results.append(f"Content: {result['content']}\nRelevance: {result['score']}")
            
            return "\n\n".join(formatted_results)
        except requests.exceptions.RequestException as e:
            return f"Error searching Contexis memory: {str(e)}"
```

### 3. Workflow Orchestration

Create complex workflows that combine CrewAI and Contexis:

```python
class ContexisWorkflowOrchestrator:
    """Orchestrates complex workflows using CrewAI and Contexis."""
    
    def __init__(self, base_url: str = "http://localhost:8080", api_key: Optional[str] = None):
        self.base_url = base_url
        self.api_key = api_key
        self.contexis_agent = ContexisAgent(base_url, api_key)
    
    def create_customer_support_workflow(self, customer_query: str, tenant_id: str = "default") -> str:
        """Create a customer support workflow."""
        
        # Create specialized agents
        support_agent = Agent(
            role="Customer Support Specialist",
            goal="Provide excellent customer support using Contexis knowledge base",
            backstory="""You are an experienced customer support specialist with 
            access to comprehensive product and policy information through Contexis.""",
            tools=[
                ContexisChatTool(self.base_url, self.api_key),
                ContexisMemoryTool(self.base_url, self.api_key)
            ],
            llm=ChatOpenAI(model="gpt-4", temperature=0.1),
            verbose=True
        )
        
        escalation_agent = Agent(
            role="Support Escalation Specialist",
            goal="Handle complex customer issues that require escalation",
            backstory="""You are a senior support specialist who handles complex 
            cases and escalations. You have deep knowledge of company policies.""",
            tools=[
                ContexisChatTool(self.base_url, self.api_key),
                ContexisMemoryTool(self.base_url, self.api_key)
            ],
            llm=ChatOpenAI(model="gpt-4", temperature=0.1),
            verbose=True
        )
        
        # Define tasks
        initial_support_task = Task(
            description=f"""Handle the customer query: {customer_query}
            
            Use Contexis to:
            - Search for relevant policies and procedures
            - Find similar cases and solutions
            - Provide accurate and helpful information
            
            If the issue is complex, escalate to the escalation specialist.""",
            agent=support_agent,
            expected_output="Initial response to customer query"
        )
        
        escalation_task = Task(
            description="""Handle escalated customer issues.
            
            Review the initial response and:
            - Provide more detailed solutions
            - Suggest next steps
            - Ensure customer satisfaction
            
            Use Contexis for comprehensive policy and procedure information.""",
            agent=escalation_agent,
            expected_output="Detailed escalation response with next steps",
            context=[initial_support_task]
        )
        
        # Create and run the crew
        crew = Crew(
            agents=[support_agent, escalation_agent],
            tasks=[initial_support_task, escalation_task],
            verbose=True
        )
        
        result = crew.kickoff()
        return result
```

## Testing

### Unit Tests

```python
import pytest
from unittest.mock import Mock, patch
from crewai import Agent, Task, Crew

class TestContexisCrewAI:
    """Test cases for Contexis-CrewAI integration."""
    
    @patch('requests.post')
    def test_contexis_tool_integration(self, mock_post):
        """Test Contexis tool integration with CrewAI."""
        # Mock response
        mock_response = Mock()
        mock_response.json.return_value = {"rendered": "Test response from Contexis"}
        mock_response.raise_for_status.return_value = None
        mock_post.return_value = mock_response
        
        # Create tool
        tool = ContexisChatTool("http://localhost:8080")
        
        # Test tool execution
        result = tool._run("test query", "CustomerDocs")
        
        assert result == "Test response from Contexis"
        mock_post.assert_called_once()
    
    def test_agent_creation(self):
        """Test CrewAI agent creation with Contexis tools."""
        contexis_agent = ContexisAgent()
        
        researcher = contexis_agent.create_researcher_agent()
        writer = contexis_agent.create_writer_agent()
        reviewer = contexis_agent.create_reviewer_agent()
        
        assert isinstance(researcher, Agent)
        assert isinstance(writer, Agent)
        assert isinstance(reviewer, Agent)
        
        assert len(researcher.tools) > 0
        assert len(writer.tools) > 0
        assert len(reviewer.tools) > 0
```

### Integration Tests

```python
class TestContexisCrewAIWorkflow:
    """Integration tests for Contexis-CrewAI workflows."""
    
    @patch('requests.post')
    def test_content_workflow(self, mock_post):
        """Test complete content creation workflow."""
        # Mock responses
        mock_responses = [
            Mock(json=lambda: {"rendered": "Research findings"}),
            Mock(json=lambda: {"rendered": "Written content"}),
            Mock(json=lambda: {"rendered": "Review feedback"})
        ]
        
        for response in mock_responses:
            response.raise_for_status.return_value = None
        
        mock_post.side_effect = mock_responses
        
        # Create workflow
        workflow = ContexisCrewAIWorkflow()
        
        # Test workflow execution
        result = workflow.create_content_workflow(
            "AI Trends 2024",
            {"format": "blog post", "length": "1000 words"}
        )
        
        assert "Research findings" in result
        assert "Written content" in result
        assert "Review feedback" in result
```

## Deployment

### Docker Configuration

```dockerfile
# Dockerfile.crewai
FROM python:3.9-slim

WORKDIR /app

# Install dependencies
COPY requirements.txt .
RUN pip install -r requirements.txt

# Copy CrewAI integration
COPY tools/crewai/ ./tools/crewai/

# Set environment variables
ENV PYTHONPATH=/app
ENV CREWAI_VERBOSE=true

# Expose port for Contexis
EXPOSE 8000

CMD ["python", "-m", "contexis.server"]
```

### Kubernetes Deployment

```yaml
# k8s/crewai-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: contexis-crewai
spec:
  replicas: 2
  selector:
    matchLabels:
      app: contexis-crewai
  template:
    metadata:
      labels:
        app: contexis-crewai
    spec:
      containers:
      - name: contexis-crewai
        image: contexis/crewai:latest
        ports:
        - containerPort: 8000
        env:
        - name: OPENAI_API_KEY
          valueFrom:
            secretKeyRef:
              name: contexis-secrets
              key: openai-api-key
        - name: CONTEXIS_API_KEY
          valueFrom:
            secretKeyRef:
              name: contexis-secrets
              key: contexis-api-key
        resources:
          requests:
            memory: "1Gi"
            cpu: "500m"
          limits:
            memory: "2Gi"
            cpu: "1000m"
```

## Best Practices

### 1. Agent Design

- **Clear Roles**: Define specific, non-overlapping roles for each agent
- **Tool Selection**: Choose appropriate Contexis tools for each agent's needs
- **Memory Management**: Use Contexis memory for persistent knowledge sharing

### 2. Workflow Design

- **Task Decomposition**: Break complex tasks into smaller, manageable subtasks
- **Dependencies**: Clearly define task dependencies and execution order
- **Error Handling**: Implement robust error handling for both CrewAI and Contexis

### 3. Performance Optimization

- **Parallel Execution**: Use CrewAI's parallel task execution where possible
- **Caching**: Implement caching for frequently accessed Contexis data
- **Resource Management**: Monitor and optimize resource usage

### 4. Security

- **API Key Management**: Securely manage API keys for both CrewAI and Contexis
- **Input Validation**: Validate all inputs to prevent injection attacks
- **Access Control**: Implement proper access control for multi-tenant scenarios

## Troubleshooting

### Common Issues

1. **Agent Communication**
   - Check Contexis server connectivity
   - Verify API key configuration
   - Review agent tool configurations

2. **Memory Issues**
   - Monitor memory usage during workflow execution
   - Implement proper cleanup for large workflows
   - Use streaming for large responses

3. **Performance Issues**
   - Optimize task dependencies
   - Use parallel execution where possible
   - Implement caching strategies

### Debug Commands

```bash
# Enable debug logging
export CREWAI_VERBOSE=true
export CONTEXIS_DEBUG=true

# Test Contexis connectivity
curl -X POST http://localhost:8080/api/v1/chat \
  -H "Content-Type: application/json" \
  -d '{"context":"CustomerDocs","query":"test"}'

# Run CrewAI with debug output
python -m crewai.main --debug
```

## Resources

- **Documentation**: [CrewAI Documentation](https://docs.crewai.com/)
- **Examples**: [CrewAI Examples](https://github.com/joaomdmoura/crewAI/tree/main/examples)
- **Community**: [CrewAI Discord](https://discord.gg/crewai)
- **Issues**: [GitHub Issues](https://github.com/joaomdmoura/crewAI/issues)