# Write Step Template
# Generated for step: write

## Step Information
- **Name**: write
- **Type**: write
- **Description**: Step 2: write
- **Timeout**: 60 seconds
- **Retry Attempts**: 2

## Input Parameters

## Output Expectations

## Writing Instructions

### Context
You are performing content writing as part of the write workflow. Your task is to create high-quality, engaging content based on the research findings and requirements.

### Writing Process
1. **Content Planning**: Analyze requirements and research findings to plan content structure
2. **Outline Creation**: Develop a logical outline for the content
3. **Content Creation**: Write the main content following the outline
4. **Review and Refinement**: Review content for clarity, accuracy, and engagement
5. **Finalization**: Ensure content meets all requirements and quality standards

### Writing Guidelines
- Use clear, concise language appropriate for the target audience
- Structure content logically with proper headings and sections
- Incorporate research findings and supporting evidence
- Maintain consistent tone and style throughout
- Ensure accuracy and factual correctness
- Optimize for readability and engagement

### Content Structure
Organize your content using the following structure:

```markdown
# [Main Title]

## Executive Summary
Brief overview of the content and key points

## Introduction
Background information and context

## Main Content
### [Section 1]
Detailed content for section 1

### [Section 2]
Detailed content for section 2

### [Section 3]
Detailed content for section 3

## Key Takeaways
Summary of main points and insights

## Conclusion
Final thoughts and recommendations

## References
List of sources and citations
```

### Output Format
Provide your written content in the following structured format:

```json
{
  "content_title": "Main title of the content",
  "executive_summary": "Brief overview of the content",
  "main_content": "Full written content in markdown format",
  "key_takeaways": [
    "Key takeaway 1",
    "Key takeaway 2",
    "Key takeaway 3"
  ],
  "word_count": 1500,
  "reading_time": "5 minutes",
  "target_audience": "Description of target audience",
  "content_type": "Type of content (article, report, etc.)",
  "seo_keywords": [
    "keyword1",
    "keyword2",
    "keyword3"
  ],
  "next_steps": "Suggested next steps for the workflow"
}
```

### Quality Criteria
- **Clarity**: Content is clear and easy to understand
- **Engagement**: Content is interesting and holds reader attention
- **Accuracy**: All information is factually correct
- **Completeness**: Content covers all required topics
- **Structure**: Content is well-organized and logical
- **SEO**: Content is optimized for search engines (if applicable)

### Error Handling
If writing cannot be completed due to:
- Insufficient research data
- Unclear requirements
- Time constraints
- Technical issues

Provide a detailed explanation of the issue and suggest alternative approaches or next steps.

## Step Completion Criteria
- Content is complete and meets all requirements
- Writing quality meets established standards
- Content structure is logical and well-organized
- All quality criteria are met
- Output format is correct and complete
