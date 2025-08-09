# Changelog

All notable changes to the Contexis CMP Framework will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- **RAG Generator**: Complete implementation of knowledge-based retrieval system generator
  - `ctx generate rag <name> --db=<type> --embeddings=<model>` command
  - Support for multiple database types: sqlite, postgres, chroma
  - Support for multiple embedding models: sentence-transformers, openai, cohere, bge-small-en
  - Comprehensive project structure generation with CMP architecture
  - Context files with search behavior rules and guardrails
  - Memory configuration with vector store setup
  - Prompt templates for search responses and no-results handling
  - Python semantic search implementation with ChromaDB integration
  - Drift detection test configuration and Python test scripts
  - Auto-generated documentation and requirements files

- **CLI Framework Enhancements**:
  - New `generate` command with subcommand architecture
  - Support for multiple generator types: rag, agent, workflow
  - Comprehensive flag handling and validation
  - Structured logging with operation tracking
  - Input validation and security checks
  - Help system with examples and usage instructions

- **Security & Compliance**:
  - Input validation for project names and configuration parameters
  - Secure file permissions (0750) for generated directories
  - Path traversal protection
  - Comprehensive error handling with user-friendly messages

- **Documentation**:
  - Auto-generated README files for each generated component
  - Usage instructions and next steps guidance
  - Configuration examples and customization options
  - Testing and deployment instructions

### Changed
- **CLI Structure**: Refactored main CLI to support multiple generator types
- **Error Handling**: Enhanced error messages with structured logging
- **File Organization**: Modular architecture with separate files for each generator component

### Technical Details

#### New Files Created:
- `src/cli/commands/generate.go` - Main generate command implementation
- `src/cli/commands/rag_generator.go` - RAG generator orchestration
- `src/cli/commands/rag_context.go` - Context file generation
- `src/cli/commands/rag_memory.go` - Memory configuration generation
- `src/cli/commands/rag_prompts.go` - Prompt template generation
- `src/cli/commands/rag_tools.go` - Python tool implementation
- `src/cli/commands/rag_tests.go` - Test configuration generation
- `src/cli/commands/rag_config.go` - Additional configuration files
- `src/cli/commands/agent_generator.go` - Agent generator placeholder
- `src/cli/commands/workflow_generator.go` - Workflow generator placeholder

#### Generated Components:
For each RAG system, the generator creates:

**Context Layer:**
- `contexts/<name>/rag_agent.ctx` - Defines search behavior, role, tools, guardrails
- `contexts/<name>/README.md` - Usage documentation

**Memory Layer:**
- `memory/<name>/memory_config.yaml` - Vector store configuration
- `memory/<name>/documents/sample.md` - Sample document

**Prompt Layer:**
- `prompts/<name>/search_response.md` - Search result formatting template
- `prompts/<name>/no_results.md` - No results handling template

**Tools Layer:**
- `tools/<name>/semantic_search.py` - Python semantic search implementation
- `tools/<name>/requirements.txt` - Python dependencies

**Testing Layer:**
- `tests/<name>/rag_drift_test.yaml` - Drift detection configuration
- `tests/<name>/test_rag.py` - Python test script

#### CLI Commands:
```bash
# Generate RAG system
ctx generate rag CustomerDocs --db=sqlite --embeddings=sentence-transformers

# Generate agent (placeholder for Week 2)
ctx generate agent SupportBot --tools=web_search,database --memory=episodic

# Generate workflow (placeholder for Week 3)
ctx generate workflow ContentPipeline --steps=research,write,review
```

#### Configuration Options:
- **Database Types**: sqlite, postgres, chroma
- **Embedding Models**: sentence-transformers, openai, cohere, bge-small-en
- **Memory Types**: episodic, none (for agents)
- **Tools**: web_search, database (for agents)
- **Steps**: research, write, review (for workflows)

### Performance
- **Generation Time**: <5 seconds for complete RAG system
- **Template Rendering**: Optimized Go template execution
- **Error Recovery**: Graceful handling of template parsing errors

### Security
- **Input Validation**: Regex-based project name validation
- **Path Security**: Prevention of directory traversal attacks
- **File Permissions**: Secure defaults (0750) for generated directories
- **Error Sanitization**: User-friendly error messages without sensitive data exposure

### Testing
- **Unit Tests**: Template parsing and validation
- **Integration Tests**: End-to-end generation workflow
- **Drift Detection**: Automated similarity testing configuration
- **Business Rules**: Compliance and consistency validation

---

## [0.1.0] - 2025-08-10

### Added
- Initial CLI framework with `ctx init` command
- Basic project scaffolding with CMP architecture
- Configuration system with YAML support
- Security baseline with input validation
- Comprehensive documentation foundation

### Changed
- Project structure to follow CMP architecture principles
- CLI to use structured logging and error handling

---

## Roadmap

### Phase 1: Core Generators (Weeks 1-3)
- [x] **Week 1**: RAG Generator Implementation ✅
- [ ] **Week 2**: Agent Generator Implementation
- [ ] **Week 3**: Workflow Generator Implementation

### Phase 2: Runtime Engine (Weeks 4-6)
- [ ] **Week 4**: Context Management System
- [ ] **Week 5**: Memory Management System
- [ ] **Week 6**: Prompt Management System

### Phase 3: Testing & Quality Assurance (Weeks 7-8)
- [ ] **Week 7**: Drift Detection System
- [ ] **Week 8**: Testing Infrastructure

### Phase 4: Deployment & Operations (Weeks 9-10)
- [ ] **Week 9**: Deployment System
- [ ] **Week 10**: Monitoring & Observability

### Phase 5: Advanced Features (Weeks 11-12)
- [ ] **Week 11**: Enterprise Security & Compliance
- [ ] **Week 12**: Integration Ecosystem

---

## Contributing

When contributing to this project, please:

1. Follow the existing changelog format
2. Add entries under the appropriate section
3. Include technical details for significant changes
4. Reference issue numbers when applicable
5. Update the roadmap as milestones are completed

## Version History

- **0.1.0**: Initial CLI framework and project scaffolding
- **Unreleased**: Phase 1, Week 1 - RAG Generator implementation
