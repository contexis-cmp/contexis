# Contributing to Contexis

Thank you for your interest in contributing to Contexis! This document provides guidelines and information for contributors.

## 🎯 Development Philosophy

Contexis follows a **Context-Memory-Prompt (CMP)** architecture that treats AI components as version-controlled, first-class citizens. We believe in:

- **Architectural Discipline**: Clear separation of concerns and well-defined interfaces
- **Reproducibility**: All AI components should be versioned and reproducible
- **Security First**: Security and privacy are non-negotiable
- **Developer Experience**: Tools should be intuitive and powerful

## 🏗️ Architecture Overview

### Core Components

- **Context**: Declarative instructions, agent roles, tool definitions
- **Memory**: Versioned knowledge stores, vector databases, logs
- **Prompt**: Pure templates hydrated at runtime

### Technology Stack

- **Go**: CLI and orchestration layer
- **Python**: AI/ML functionality and integrations
- **YAML**: Configuration files
- **JSON**: Data interchange and locking

## 🚀 Getting Started

### Prerequisites

- Go 1.21+
- Python 3.10+
- Git
- Make (optional, but recommended)

### Development Setup

```bash
# Clone the repository
git clone https://github.com/contexis/cmp.git
cd cmp

# Install dependencies
make setup

# Verify installation
go version
python --version
ctx version
```

### Development Workflow

1. **Fork and Clone**
   ```bash
   git fork https://github.com/contexis/cmp.git
   git clone https://github.com/YOUR_USERNAME/cmp.git
   cd cmp
   ```

2. **Create Feature Branch**
   ```bash
   git checkout -b feature/your-feature-name
   ```

3. **Make Changes**
   - Follow the coding standards below
   - Add tests for new functionality
   - Update documentation

4. **Test Your Changes**
   ```bash
   make test          # Run all tests
   make lint          # Check code quality
   make security      # Security checks
   ```

5. **Commit Your Changes**
   ```bash
   git add .
   git commit -m "feat: add new feature description"
   ```

6. **Push and Create PR**
   ```bash
   git push origin feature/your-feature-name
   # Create PR on GitHub
   ```

## 📝 Coding Standards

### Go Standards

- Follow [Effective Go](https://golang.org/doc/effective_go.html)
- Use structured error handling
- Implement proper logging with Zap
- Add comprehensive tests
- Use interfaces for dependency injection

```go
// Good example
func (s *Service) ProcessContext(ctx context.Context, req Request) error {
    logger := s.logger.With(
        zap.String("request_id", getRequestID(ctx)),
        zap.String("operation", "process_context"),
    )
    
    if err := s.validate(req); err != nil {
        logger.Error("validation failed", zap.Error(err))
        return fmt.Errorf("validate request: %w", err)
    }
    
    // Implementation
    return nil
}
```

### Python Standards

- Python 3.10+ required
- Use type hints everywhere
- Follow PEP 8 and Black formatting
- Use async/await for I/O operations
- Implement proper error handling

```python
# Good example
from typing import Protocol, TypedDict, List
import asyncio
from dataclasses import dataclass

@dataclass(frozen=True)
class EmbeddingRequest:
    text: str
    model: str
    normalize: bool = True

async def process_embeddings(
    requests: List[EmbeddingRequest]
) -> List[List[float]]:
    """Process embeddings asynchronously."""
    results = []
    async with asyncio.TaskGroup() as tg:
        for req in requests:
            task = tg.create_task(generate_embedding(req))
            results.append(task)
    
    return [await r for r in results]
```

### Security Guidelines

- **Never log sensitive data** (API keys, passwords, tokens)
- **Validate all inputs** with proper sanitization
- **Use secure defaults** for all configurations
- **Implement proper authentication** for multi-tenant features
- **Follow OWASP guidelines** for web components

```go
// Good security practice
func SafeJoin(base, path string) (string, error) {
    full := filepath.Join(base, filepath.Clean(path))
    
    if !strings.HasPrefix(full, base) {
        return "", errors.New("path escapes base directory")
    }
    
    return full, nil
}
```

## 🧪 Testing Guidelines

### Test Structure

```
tests/
├── unit/              # Unit tests
├── integration/       # Integration tests
├── e2e/              # End-to-end tests
└── fixtures/         # Test data and fixtures
```

### Test Requirements

- **Unit Tests**: 80%+ coverage for core functionality
- **Integration Tests**: All major workflows
- **E2E Tests**: Critical user journeys
- **Performance Tests**: For AI operations

### Running Tests

```bash
# Run all tests
make test

# Run specific test suites
go test ./... -v
pytest tests/ -v

# Run with coverage
make test-coverage

# Run performance tests
make test-performance
```

## 📚 Documentation Standards

### Documentation Structure

```
docs/
├── api/              # API documentation
├── guides/           # User guides
├── examples/         # Code examples
└── technical/        # Technical specifications
```

### Writing Documentation

- Use clear, concise language
- Include code examples
- Add diagrams for complex concepts
- Keep documentation up-to-date
- Use Markdown with proper formatting

### Documentation Requirements

- **README.md**: Project overview and quick start
- **API docs**: Complete API reference
- **Guides**: Step-by-step tutorials
- **Technical docs**: Architecture and design decisions

## 🔄 Pull Request Process

### PR Requirements

1. **Clear Description**: What, why, and how
2. **Tests**: New functionality must have tests
3. **Documentation**: Update docs if needed
4. **Security Review**: Security implications considered
5. **Performance**: No performance regressions

### PR Template

```markdown
## Description
Brief description of the changes

## Type of Change
- [ ] Bug fix
- [ ] New feature
- [ ] Breaking change
- [ ] Documentation update

## Testing
- [ ] Unit tests pass
- [ ] Integration tests pass
- [ ] E2E tests pass
- [ ] Performance tests pass

## Security
- [ ] No sensitive data logged
- [ ] Input validation added
- [ ] Security implications considered

## Documentation
- [ ] README updated
- [ ] API docs updated
- [ ] Guides updated
```

## 🏷️ Versioning

We follow [Semantic Versioning](https://semver.org/) (SemVer):

- **MAJOR**: Breaking changes
- **MINOR**: New features, backward compatible
- **PATCH**: Bug fixes, backward compatible

## 🚀 Release Process

1. **Feature Freeze**: Stop adding new features
2. **Testing**: Comprehensive testing on staging
3. **Documentation**: Update all documentation
4. **Release**: Tag and release
5. **Announcement**: Communicate changes

## 🤝 Community Guidelines

### Code of Conduct

- Be respectful and inclusive
- Welcome new contributors
- Provide constructive feedback
- Follow project conventions

### Communication

- **Issues**: Use GitHub issues for bugs and feature requests
- **Discussions**: Use GitHub discussions for questions
- **Security**: Email security@contexis.dev for security issues

## 🎯 Contribution Areas

### Priority Areas

1. **Core Framework**: CMP architecture improvements
2. **Security**: Security enhancements and audits
3. **Performance**: Performance optimization
4. **Testing**: Test coverage and quality
5. **Documentation**: Documentation improvements

### Good First Issues

- Documentation updates
- Test additions
- Bug fixes
- Small feature improvements

## 📞 Getting Help

- **Documentation**: Check the docs first
- **Issues**: Search existing issues
- **Discussions**: Ask questions in discussions
- **Email**: contact@contexis.dev

## 🙏 Acknowledgments

Thank you for contributing to Contexis! Your contributions help make AI applications more reliable, secure, and accessible.

---

**Contexis** - Bringing architectural discipline to AI applications 🚀
