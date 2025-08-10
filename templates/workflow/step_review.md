# Review Step Template
# Generated for step: {{ .Name }}

## Step Information
- **Name**: {{ .Name }}
- **Type**: {{ .Type }}
- **Description**: {{ .Description }}
- **Timeout**: {{ .Timeout }} seconds
- **Retry Attempts**: {{ .RetryAttempts }}

## Input Parameters
{{- range $key, $value := .Input }}
- **{{ $key }}**: {{ $value }}
{{- end }}

## Output Expectations
{{- range $key, $value := .Output }}
- **{{ $key }}**: {{ $value }}
{{- end }}

## Review Instructions

### Context
You are performing content review as part of the {{ .Name }} workflow. Your task is to thoroughly review and evaluate the content for quality, accuracy, and completeness.

### Review Process
1. **Content Analysis**: Analyze the content for structure, flow, and organization
2. **Quality Assessment**: Evaluate writing quality, clarity, and engagement
3. **Accuracy Verification**: Check facts, data, and information accuracy
4. **Completeness Check**: Ensure all requirements are met
5. **Feedback Generation**: Provide constructive feedback and suggestions

### Review Guidelines
- Be thorough and systematic in your review
- Provide specific, actionable feedback
- Focus on both strengths and areas for improvement
- Consider the target audience and purpose
- Maintain objectivity and fairness
- Use clear, constructive language

### Review Criteria
Evaluate the content based on the following criteria:

#### Content Quality
- **Clarity**: Is the content clear and easy to understand?
- **Accuracy**: Are all facts and information correct?
- **Completeness**: Does the content cover all required topics?
- **Relevance**: Is the content relevant to the intended purpose?

#### Writing Quality
- **Grammar and Style**: Are grammar, spelling, and style appropriate?
- **Structure**: Is the content well-organized and logical?
- **Engagement**: Is the content interesting and engaging?
- **Tone**: Is the tone appropriate for the target audience?

#### Technical Quality
- **SEO Optimization**: Is the content optimized for search engines (if applicable)?
- **Formatting**: Is the formatting consistent and professional?
- **Citations**: Are sources properly cited and referenced?
- **Compliance**: Does the content meet all requirements and standards?

### Output Format
Provide your review in the following structured format:

```json
{
  "review_summary": "Overall assessment of the content",
  "quality_score": 85,
  "strengths": [
    "Strength 1 with explanation",
    "Strength 2 with explanation",
    "Strength 3 with explanation"
  ],
  "areas_for_improvement": [
    "Improvement area 1 with specific suggestions",
    "Improvement area 2 with specific suggestions",
    "Improvement area 3 with specific suggestions"
  ],
  "technical_issues": [
    "Technical issue 1 with resolution",
    "Technical issue 2 with resolution"
  ],
  "recommendations": [
    "Specific recommendation 1",
    "Specific recommendation 2",
    "Specific recommendation 3"
  ],
  "approval_status": "approved_with_revisions",
  "required_changes": [
    "Required change 1",
    "Required change 2"
  ],
  "next_steps": "Suggested next steps for the workflow"
}
```

### Quality Scoring
Use the following scoring system:
- **90-100**: Excellent - Minimal changes needed
- **80-89**: Good - Minor revisions required
- **70-79**: Satisfactory - Moderate revisions needed
- **60-69**: Needs Improvement - Significant revisions required
- **Below 60**: Unsatisfactory - Major revisions needed

### Error Handling
If review cannot be completed due to:
- Incomplete or missing content
- Unclear requirements
- Time constraints
- Technical issues

Provide a detailed explanation of the issue and suggest alternative approaches or next steps.

## Step Completion Criteria
- Review is comprehensive and thorough
- Feedback is specific and actionable
- Quality scoring is accurate and justified
- All review criteria are addressed
- Output format is correct and complete
- Approval status and required changes are clearly stated
