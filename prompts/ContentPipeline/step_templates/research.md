# Research Step Template
# Generated for step: research

## Step Information
- **Name**: research
- **Type**: research
- **Description**: Step 1: research
- **Timeout**: 60 seconds
- **Retry Attempts**: 2

## Input Parameters

## Output Expectations

## Research Instructions

### Context
You are performing research as part of the research workflow. Your task is to gather comprehensive information on the given topic or query.

### Research Process
1. **Topic Analysis**: Understand the research topic and identify key areas to investigate
2. **Source Identification**: Find relevant and reliable sources of information
3. **Information Gathering**: Collect detailed information from multiple sources
4. **Data Validation**: Verify the accuracy and relevance of collected information
5. **Synthesis**: Organize and synthesize the research findings

### Research Guidelines
- Use multiple credible sources
- Cross-reference information for accuracy
- Document sources and citations
- Focus on recent and relevant information
- Consider different perspectives and viewpoints
- Maintain objectivity in analysis

### Output Format
Provide your research findings in the following structured format:

```json
{
  "research_summary": "Brief overview of research findings",
  "key_findings": [
    "Finding 1 with supporting evidence",
    "Finding 2 with supporting evidence",
    "Finding 3 with supporting evidence"
  ],
  "sources": [
    {
      "title": "Source title",
      "url": "Source URL",
      "author": "Author name",
      "date": "Publication date",
      "relevance": "Why this source is relevant"
    }
  ],
  "recommendations": [
    "Recommendation 1 based on findings",
    "Recommendation 2 based on findings"
  ],
  "next_steps": "Suggested next steps for the workflow"
}
```

### Quality Criteria
- **Completeness**: Research covers all relevant aspects of the topic
- **Accuracy**: Information is verified and reliable
- **Relevance**: Findings are directly applicable to the workflow
- **Clarity**: Results are presented in a clear, organized manner
- **Actionability**: Recommendations are specific and implementable

### Error Handling
If research cannot be completed due to:
- Insufficient information available
- Access restrictions to sources
- Time constraints
- Technical issues

Provide a detailed explanation of the issue and suggest alternative approaches or next steps.

## Step Completion Criteria
- Research findings are comprehensive and well-documented
- Sources are properly cited and verified
- Recommendations are actionable and relevant
- Output format is correct and complete
- All quality criteria are met
