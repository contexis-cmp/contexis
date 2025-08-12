# Changelog

All notable changes to the Contexis CMP Framework will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Phase 5: Advanced Features
  - Week 11: Enterprise Security & Compliance
    - Optional API key authentication, RBAC enforcement, and per-key/tenant/IP rate limiting (HTTP server)
    - Audit logging subsystem with JSONL sink and request context propagation
    - Prompt-injection guard (optional via `CMP_PI_ENFORCEMENT`): classify and block risky inputs (403)
    - Source-constrained answering (optional via `CMP_REQUIRE_CITATION`): `424` when no sources, `422` when missing citations
    - Episodic memory encryption-at-rest via AES-GCM (env key provider)
    - Helm hardening: securityContext (non-root, read-only FS), Ingress TLS/HSTS annotations
  - Week 12: Integration Ecosystem
    - Plugin system: local install from folder, remote install from ZIP URL or Git URL (#ref supported)
    - Plugin CLI: `ctx plugin list|info|install|remove`
    - Plugin scaffolder: `ctx generate plugin <name>`
    - Capability registry and manifest compatibility hooks
    - Templates: `templates/plugin/` example plugin
    - Integration test scaffold for plugin lifecycle
- Hugging Face Integration (Phase 4 follow-up):
  - Runtime provider layer with `HuggingFaceAPIProvider` (env-driven)
  - Optional inference in server after prompt rendering when `HF_TOKEN` and `HF_MODEL_ID` are set
  - CLI: `ctx hf test-model` for quick connectivity tests
  - Metrics: `cmp_hf_inference_latency_seconds`, `cmp_hf_inference_errors_total`
  - Tracing: OTel span `huggingface.generate` annotated with model id
  - Helm: values and templates to inject `HF_*` env vars via Secret/ExternalSecret
  - Docs: `docs/guides/hugging-face.md`, example `docs/examples/hugging-face.md`
- Phase 4 Deployment & Operations:
  - Containerization: root Dockerfile and docker-compose
  - CLI: `ctx build` (image build) and enhanced `ctx deploy` (docker/kubernetes)
  - Kubernetes: base manifests and Helm chart (`charts/contexis-app`) with optional ExternalSecret and Argo Rollouts
  - Server: health (`/healthz`), readiness (`/readyz`), version (`/version`), metrics (`/metrics`)
  - Worker: `ctx worker` process with `/healthz` and `/metrics`
  - Observability: Prometheus metrics for HTTP, prompt render, memory search; request logging and basic tracing
  - CI: `build-and-helm.yml` (build, helm lint/template) and `security.yml` (SBOM + Grype scan)
- Phase 3 Testing & Quality Infrastructure:
  - Drift Detection System:
    - YAML-driven specs under `tests/<Component>/rag_drift_test.yaml`
    - Baselines stored at `tests/<Component>/baselines/drift_baseline.json`
    - CLI: `ctx test --drift-detection [--component <Name>] [--semantic] [--update-baseline] [--junit] [--out <dir>]`
    - Reports: per-component JSON and aggregated index; optional JUnit XML
  - Go Testing Infrastructure:
    - CLI: `ctx test --all|--unit|--integration|--e2e|--category <name> [--coverage] [--junit] [--out <dir>]`
    - Coverage enforcement via `tests/test_config.yaml` thresholds
    - Reports: `tests/reports/go_<suite>.txt`, `go_tests.json`, optional `junit-go.xml`; coverage in `tests/coverage/*.out`
- Phase 2 Runtime Engine features:
  - Context Management System: tenant-aware resolution, inheritance/merge, JSON Schema validation; CLI `ctx context validate|reload`
  - Memory Management System: file-backed vector store ("sqlite") and episodic logs; CLI `ctx memory ingest|search|optimize`; component `memory_config.yaml` loader and encryption toggle
  - Prompt Management System: prompt engine with includes, helpers, format validation, simple token trimming; CLI `ctx prompt render|validate`, `ctx prompt-lint`
  - Guardrails: capability validation and format/max_tokens enforcement
  - Reproducibility: `ctx lock generate` producing `context.lock.json`
  - HTTP Server: `ctx serve` with `/api/v1/chat` endpoint for local experimentation
- **Complete Phase 1 Implementation**: All three core generators fully implemented
  - **RAG Generator**: Knowledge-based retrieval system generator
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

  - **Agent Generator**: Conversational agent generator with tools and memory
    - `ctx generate agent <name> --tools=<list> --memory=<type>` command
    - Support for multiple tool types: web_search, database, api, file_system, email
    - Support for memory types: episodic, none
    - Role-based context generation with personality and capabilities
    - Tool integration framework with MCP support
    - Memory management system for conversation persistence
    - Prompt templates for agent responses and interactions
    - Python tool implementations with error handling
    - Behavior consistency testing and validation
    - Auto-generated documentation and requirements files

  - **Workflow Generator**: Multi-step AI processing pipeline generator
    - `ctx generate workflow <name> --steps=<list>` command
    - Support for step types: research, write, review, extract, transform, load, analyze, generate, validate, deploy
    - Step dependency resolution and parallel execution support
    - State management system with checkpointing and recovery
    - Resource management with CPU, memory, storage, and network limits
    - Error handling with retry logic and failure recovery
    - Workflow orchestration with progress monitoring
    - Step-specific prompt templates and instructions
    - Integration testing framework for end-to-end validation
    - Auto-generated documentation and configuration files

- **Comprehensive Testing Framework**: Test-driven development implementation
  - **Unit Tests**: Complete test coverage for all generators and components
    - Agent generator validation and configuration tests
    - Workflow generator step and dependency tests
    - CLI command structure and flag validation tests
    - Helper functions and utility tests
  - **Integration Tests**: End-to-end testing for complete workflows
    - Agent generation with different tool combinations
    - Workflow generation with various step configurations
    - File generation and template processing validation
  - **Test Infrastructure**: Robust testing framework with utilities
    - Test fixtures and helper functions
    - Temporary directory management
    - File system assertions and validation
    - Mock and stub implementations for testing

- **CLI Framework Enhancements**:
  - New `generate` command with subcommand architecture
  - Support for multiple generator types: rag, agent, workflow
  - Comprehensive flag handling and validation
  - Structured logging with operation tracking
  - Input validation and security checks
  - Help system with examples and usage instructions
  - Command structure validation and error handling

- **Security & Compliance**:
  - Input validation for project names and configuration parameters
  - Secure file permissions (0750) for generated directories
  - Path traversal protection
  - Comprehensive error handling with user-friendly messages
  - Business rule validation and compliance checking

- **Documentation**:
  - Auto-generated README files for each generated component
  - Usage instructions and next steps guidance
  - Configuration examples and customization options
  - Testing and deployment instructions
  - Comprehensive API documentation and examples
  - Installation guides with local installation support

- **Installation & Deployment**:
  - Local installation support without sudo requirements
  - Automatic PATH configuration for different shells
  - Installation script with colored output and validation
  - Fallback installation options for different environments
  - Comprehensive installation testing and validation

### Changed
- README: documented security features and plugin commands
- Charts: deployment and ingress templates hardened; values support `security.authEnabled`
- README updated with new runtime commands and examples
- Added dependency `github.com/xeipuuv/gojsonschema` for schema validation
- **Repository Structure**: Updated to use correct GitHub repository URL `github.com/contexis-cmp/contexis`
- **Module Path**: Updated Go module path to match repository structure
- **Import Paths**: All Go files updated with correct import paths
- **Documentation URLs**: All documentation and configuration files updated with correct repository references
- **CLI Structure**: Refactored main CLI to support multiple generator types
- **Error Handling**: Enhanced error messages with structured logging
- **File Organization**: Modular architecture with separate files for each generator component
- **Code Quality**: Comprehensive code formatting and static analysis fixes
- **Build System**: Streamlined build process with automated quality checks

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
- `src/cli/commands/agent_generator.go` - Agent generator implementation
- `src/cli/commands/workflow_generator.go` - Workflow generator implementation
- `src/cli/commands/cli_commands.go` - Centralized CLI command definitions
- `tests/unit/agent_generator_test.go` - Agent generator unit tests
- `tests/unit/workflow_generator_test.go` - Workflow generator unit tests
- `tests/unit/cli_commands_test.go` - CLI command unit tests
- `tests/integration/agent_generator_integration_test.go` - Agent integration tests
- `tests/unit/helpers/helpers.go` - Test utilities and fixtures
- `tests/test_runner.go` - Test runner and orchestration
- `templates/agent/` - Agent generator templates
- `templates/workflow/` - Workflow generator templates
- `docs/testing.md` - Testing framework documentation
- `docs/testing_summary.md` - Testing implementation summary
- `WORKFLOW_GENERATOR_SUMMARY.md` - Workflow generator implementation summary
- `PHASE1_CLEANUP_SUMMARY.md` - Phase 1 cleanup and quality summary
- `scripts/install.sh` - Local installation script with automatic PATH setup



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

**RAG Generator Components:**
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

**Agent Generator Components:**
For each agent, the generator creates:

**Context Layer:**
- `contexts/<name>/<name>.ctx` - Defines agent role, personality, capabilities, guardrails
- `contexts/<name>/README.md` - Usage documentation

**Memory Layer:**
- `memory/<name>/memory_config.yaml` - Memory configuration for conversation persistence
- `memory/<name>/episodic/` - Episodic memory storage structure

**Prompt Layer:**
- `prompts/<name>/agent_response.md` - Agent response formatting template
- `prompts/<name>/interaction.md` - User interaction handling template

**Tools Layer:**
- `tools/<name>/<tool_name>.py` - Python tool implementations
- `tools/<name>/requirements.txt` - Python dependencies

**Testing Layer:**
- `tests/<name>/agent_behavior.yaml` - Behavior consistency test configuration
- `tests/<name>/test_agent.py` - Python test script

**Workflow Generator Components:**
For each workflow, the generator creates:

**Workflow Layer:**
- `workflows/<name>/<name>.yaml` - Workflow definition with steps and dependencies
- `workflows/<name>/README.md` - Usage documentation

**Context Layer:**
- `contexts/<name>/workflow_coordinator.ctx` - Orchestration logic and coordination

**Prompt Layer:**
- `prompts/<name>/step_templates/<step>.md` - Step-specific prompt templates

**Memory Layer:**
- `memory/<name>/workflow_state.yaml` - State persistence configuration

**Testing Layer:**
- `tests/<name>/workflow_integration.py` - End-to-end workflow testing

#### CLI Commands:
```bash
# Generate RAG system
ctx generate rag CustomerDocs --db=sqlite --embeddings=sentence-transformers

# Generate agent with tools and memory
ctx generate agent SupportBot --tools=web_search,database --memory=episodic

# Generate workflow with steps
ctx generate workflow ContentPipeline --steps=research,write,review
```

#### Configuration Options:
- **Database Types**: sqlite, postgres, chroma
- **Embedding Models**: sentence-transformers, openai, cohere, bge-small-en
- **Memory Types**: episodic, none (for agents)
- **Tools**: web_search, database, api, file_system, email (for agents)
- **Steps**: research, write, review, extract, transform, load, analyze, generate, validate, deploy (for workflows)

### Performance
- **Generation Time**: <5 seconds for complete system generation
- **Template Rendering**: Optimized Go template execution
- **Error Recovery**: Graceful handling of template parsing errors
- **Test Execution**: Fast test suite execution with parallel testing
- **Build Performance**: Optimized build process with dependency caching

### Security
- **Input Validation**: Regex-based project name validation
- **Path Security**: Prevention of directory traversal attacks
- **File Permissions**: Secure defaults (0750) for generated directories
- **Error Sanitization**: User-friendly error messages without sensitive data exposure

### Testing
- **Unit Tests**: Comprehensive test coverage for all generators and components
- **Integration Tests**: End-to-end generation workflow validation
- **Drift Detection**: Automated similarity testing configuration
- **Business Rules**: Compliance and consistency validation
- **Test Infrastructure**: Robust testing framework with utilities and fixtures
- **Code Coverage**: High test coverage with automated reporting

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

### Phase 1: Core Generators (Weeks 1-3)  COMPLETED
- [x] **Week 1**: RAG Generator Implementation 
- [x] **Week 2**: Agent Generator Implementation 
- [x] **Week 3**: Workflow Generator Implementation 

### Phase 2: Runtime Engine (Weeks 4-6)
- [x] **Week 4**: Context Management System
- [x] **Week 5**: Memory Management System
- [x] **Week 6**: Prompt Management System

### Phase 3: Testing & Quality Assurance (Weeks 7-8)
- [x] **Week 7**: Drift Detection System
- [x] **Week 8**: Testing Infrastructure

### Phase 4: Deployment & Operations (Weeks 9-10)
- [x] **Week 9**: Deployment System
- [x] **Week 10**: Monitoring & Observability

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
- **Unreleased**: Phase 1 Complete - All three core generators (RAG, Agent, Workflow) with comprehensive testing framework
