# Agent Response Template

## Conversation Context
- User ID: [USER_ID]
- Session ID: [SESSION_ID]
- Conversation History: [CONVERSATION_HISTORY]
- Current Context: [CURRENT_CONTEXT]

## Tool Usage
- **Tool**: [TOOL_NAME]
- **Input**: [TOOL_INPUT]
- **Output**: [TOOL_OUTPUT]
- **Status**: [TOOL_STATUS]

## Response Guidelines
- **Tone**: Professional and helpful
- **Format**: JSON
- **Max Tokens**: 500

## Response Template

```json
{
  "response": "[RESPONSE_TEXT]",
  "confidence": [CONFIDENCE_SCORE],
  "tools_used": ["[TOOLS_USED]"],
  "next_actions": ["[NEXT_ACTIONS]"],
  "memory_updated": [MEMORY_UPDATED],
  "context_sha": "[CONTEXT_SHA]"
}
```

## Example Response

For a support inquiry:
```json
{
  "response": "I understand you're having trouble with your account. Let me look up your information and help you resolve this issue.",
  "confidence": 0.92,
  "tools_used": ["user_lookup", "account_status"],
  "next_actions": ["verify_identity", "check_account_status"],
  "memory_updated": true,
  "context_sha": "sha256:def456..."
}
```

## Template Variables

This template uses the following placeholder variables that will be populated at runtime:

- `[USER_ID]` - Unique user identifier
- `[SESSION_ID]` - Current session identifier
- `[CONVERSATION_HISTORY]` - Previous conversation messages
- `[CURRENT_CONTEXT]` - Current conversation context
- `[TOOL_NAME]` - Name of the tool being used
- `[TOOL_INPUT]` - Input provided to the tool
- `[TOOL_OUTPUT]` - Output from the tool
- `[TOOL_STATUS]` - Status of tool execution
- `[RESPONSE_TEXT]` - The actual response text
- `[CONFIDENCE_SCORE]` - Confidence score (0.0-1.0)
- `[TOOLS_USED]` - List of tools used
- `[NEXT_ACTIONS]` - Suggested next actions
- `[MEMORY_UPDATED]` - Whether memory was updated
- `[CONTEXT_SHA]` - Context hash for tracking
