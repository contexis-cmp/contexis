# Conversational Agent Example

This example demonstrates how to build a conversational agent using the Contexis CMP framework.

## Overview

This conversational agent is designed to handle multi-turn conversations with context awareness and memory. It can:

- Maintain conversation history
- Remember user preferences
- Handle complex multi-step interactions
- Provide personalized responses

## Quick Start

```bash
# Initialize the project
ctx init conversation-agent
cd conversation-agent

# Generate the agent
ctx generate agent ConversationAgent --memory=episodic --tools=web_search,database

# Add conversation skills
ctx agent add-skill --name=greeting --context=./contexts/greeting.ctx
ctx agent add-skill --name=help --context=./contexts/help.ctx

# Test the agent
ctx test

# Start a conversation (example via server)
ctx serve --addr :8000
# curl -X POST http://localhost:8000/api/v1/chat -H 'Content-Type: application/json' \
#   -d '{"context":"ConversationAgent","component":"ConversationAgent","query":"hello","top_k":0,"data":{"user_input":"Hello, I need help with my account"}}'
```

## Generated Structure

```
conversation-agent/
├── contexts/
│   ├── conversation_agent.ctx    # Main agent definition
│   ├── greeting.ctx              # Greeting skill
│   └── help.ctx                  # Help skill
├── memory/
│   ├── conversations/            # Conversation history
│   └── user_preferences/         # User preferences
├── prompts/
│   ├── conversation_start.md     # Conversation initiation
│   ├── conversation_continue.md  # Conversation continuation
│   └── conversation_end.md       # Conversation ending
├── tools/
│   ├── conversation_tracker.py   # Conversation state management
│   └── user_preferences.py       # User preference storage
├── tests/
│   ├── conversation_flow.py      # Conversation flow tests
│   └── context_switching.py      # Context switching tests
└── context.lock.json             # Version locks
```

## Key Features

### Conversation Management
- **Multi-turn Support**: Handle extended conversations
- **Context Switching**: Switch between different conversation topics
- **Memory Persistence**: Remember conversation history across sessions
- **State Management**: Track conversation state and user intent

### Personalization
- **User Preferences**: Store and recall user preferences
- **Personalized Responses**: Tailor responses based on user history
- **Adaptive Behavior**: Learn from user interactions

### Integration
- **API Integration**: Connect to external services
- **Database Storage**: Persistent conversation storage
- **Real-time Updates**: Live conversation updates

## Configuration

### Agent Context

Edit `contexts/conversation_agent.ctx`:

```yaml
name: "Conversation Agent"
version: "1.0.0"
description: "Multi-turn conversational agent with memory"

role:
  persona: "Friendly, helpful conversational assistant"
  capabilities: ["multi_turn_conversation", "context_switching", "personalization"]
  limitations: ["no_personal_data", "no_harmful_content"]

tools:
  - name: "conversation_tracker"
    uri: "mcp://conversation.track"
    description: "Track conversation state and history"
  - name: "user_preferences"
    uri: "mcp://user.preferences"
    description: "Store and retrieve user preferences"
  - name: "context_switcher"
    uri: "mcp://conversation.switch"
    description: "Switch conversation context"

guardrails:
  tone: "friendly"
  format: "text"
  max_tokens: 1000
  temperature: 0.3
  
memory:
  episodic: true
  max_history: 50
  privacy: "user_isolated"

testing:
  drift_threshold: 0.85
  business_rules:
    - "always_maintain_context"
    - "respect_user_preferences"
    - "handle_context_switches_gracefully"
```

### Conversation Prompts

#### Conversation Start

Edit `prompts/conversation_start.md`:

```markdown
# Conversation Start Template

Hello! I'm your conversational assistant. How can I help you today?

**Previous Interactions**: {{ previous_interactions }}
**User Preferences**: {{ user_preferences }}

## Conversation Guidelines
- Be friendly and helpful
- Remember previous interactions
- Respect user preferences
- Ask clarifying questions when needed

## Response

{{ response_text }}

**Next Suggested Actions**: {{ next_actions }}
```

#### Conversation Continue

Edit `prompts/conversation_continue.md`:

```markdown
# Conversation Continue Template

**Current Context**: {{ current_context }}
**Conversation History**: {{ conversation_history }}
**User Intent**: {{ user_intent }}

## Response Guidelines
- Maintain conversation flow
- Reference previous context
- Provide helpful information
- Suggest next steps when appropriate

## Response

{{ response_text }}

**Context Updated**: {{ context_updated }}
**Suggested Follow-up**: {{ follow_up_questions }}
```

## Usage Examples

### Basic Conversation

```bash
# Start a conversation
ctx run conversation "Hi, I'm new here"

# Continue the conversation
ctx run conversation "Can you help me set up my account?"

# Ask follow-up questions
ctx run conversation "What are the security requirements?"
```

### Context Switching

```bash
# Switch to help context
ctx run conversation "I need help with my order" --context=help

# Switch back to general conversation
ctx run conversation "Thanks for the help!" --context=general
```

### Personalization

```bash
# Set user preferences
ctx user set-preferences --preferences='{"language": "en", "tone": "casual"}'

# Start personalized conversation
ctx run conversation "Hello" --user-id=user_123
```

## Testing

### Conversation Flow Testing

```bash
# Test conversation flow
ctx test --type=conversation --scenario=basic_flow

# Test context switching
ctx test --type=conversation --scenario=context_switch

# Test personalization
ctx test --type=conversation --scenario=personalization
```

### Automated Testing

The generated test suite includes:

```python
# tests/conversation_flow.py
def test_basic_conversation_flow():
    """Test basic conversation flow"""
    agent = ConversationAgent()
    
    # Test conversation start
    response1 = agent.start_conversation("Hi")
    assert "hello" in response1.lower()
    
    # Test conversation continue
    response2 = agent.continue_conversation("I need help")
    assert "help" in response2.lower()
    
    # Test conversation end
    response3 = agent.end_conversation("Goodbye")
    assert "goodbye" in response3.lower()

def test_context_switching():
    """Test context switching capabilities"""
    agent = ConversationAgent()
    
    # Start in general context
    response1 = agent.run("Hello", context="general")
    
    # Switch to help context
    response2 = agent.run("I need help", context="help")
    
    # Verify context switch
    assert response2.context == "help"
```

## Deployment

### Local Development

```bash
# Start development server
ctx dev --port=8080

# Test with curl
curl -X POST http://localhost:8080/conversation \
  -H "Content-Type: application/json" \
  -d '{"message": "Hello", "user_id": "user_123"}'
```

### Production Deployment

```bash
# Build for production
ctx build --environment=production

# Deploy to container
ctx deploy --target=docker --image=conversation-agent:latest

# Deploy to Kubernetes
ctx deploy --target=kubernetes --namespace=agents
```

## Monitoring

### Conversation Metrics

```bash
# View conversation metrics
ctx metrics --type=conversation

# Monitor conversation quality
ctx metrics --type=quality

# Track user engagement
ctx metrics --type=engagement
```

### Logs

```bash
# View conversation logs
ctx logs --type=conversation --user-id=user_123

# Monitor errors
ctx logs --level=error --type=conversation
```

## Customization

### Adding New Skills

1. **Create Skill Context**

```yaml
# contexts/booking.ctx
name: "Booking Skill"
version: "1.0.0"
description: "Handle booking and scheduling"

role:
  persona: "Efficient booking assistant"
  capabilities: ["schedule_appointments", "check_availability", "confirm_bookings"]

tools:
  - name: "calendar_integration"
    uri: "mcp://calendar.book"
    description: "Integrate with calendar systems"
```

2. **Add Skill to Agent**

```bash
ctx agent add-skill --name=booking --context=./contexts/booking.ctx
```

3. **Test New Skill**

```bash
ctx test --skill=booking
```

### Custom Tools

Create custom tools for specific integrations:

```python
# tools/custom_integration.py
class CustomIntegrationTool:
    def __init__(self, api_key):
        self.api_key = api_key
    
    async def call(self, parameters):
        # Custom integration logic
        return result
```

## Troubleshooting

### Common Issues

1. **Conversation State Loss**
   - Check memory configuration
   - Verify conversation persistence
   - Review state management

2. **Context Switching Issues**
   - Validate context definitions
   - Check context switching logic
   - Review conversation flow

3. **Personalization Not Working**
   - Verify user preference storage
   - Check preference retrieval
   - Review personalization logic

### Debug Mode

```bash
# Enable debug logging
ctx run conversation "test" --debug

# View detailed logs
ctx logs --level=debug --type=conversation
```

## Next Steps

- Add more conversation skills
- Implement advanced personalization
- Integrate with external services
- Add conversation analytics
- Implement multi-language support

## Resources

- [Conversation Design Guide](https://docs.contexis.dev/guides/conversation-design)
- [Personalization Best Practices](https://docs.contexis.dev/guides/personalization)
- [Context Switching Patterns](https://docs.contexis.dev/guides/context-switching)
- [Testing Strategies](https://docs.contexis.dev/guides/testing)
