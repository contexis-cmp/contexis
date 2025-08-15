#!/bin/bash

# Contexis Dev-Only Setup Script
# This script sets up a local development environment with Phi-3.5-Mini and Chroma(SQLite)
# No external API calls required - everything runs locally

set -e

echo "🚀 Setting up Contexis Dev-Only Environment"
echo "=========================================="

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if Python 3.10+ is available
print_status "Checking Python version..."
python_version=$(python3 --version 2>&1 | grep -oE '[0-9]+\.[0-9]+' | head -1)
if [[ $(echo "$python_version >= 3.10" | bc -l) -eq 1 ]]; then
    print_success "Python $python_version found"
else
    print_error "Python 3.10+ is required. Found: $python_version"
    exit 1
fi

# Check if pip is available
if ! command -v pip3 &> /dev/null; then
    print_error "pip3 is not installed"
    exit 1
fi

# Create virtual environment if it doesn't exist
if [ ! -d "venv" ]; then
    print_status "Creating virtual environment..."
    python3 -m venv venv
    print_success "Virtual environment created"
else
    print_status "Virtual environment already exists"
fi

# Activate virtual environment
print_status "Activating virtual environment..."
source venv/bin/activate

# Upgrade pip
print_status "Upgrading pip..."
pip install --upgrade pip

# Install dev-only requirements
print_status "Installing dev-only dependencies..."
if [ -f "requirements-dev.txt" ]; then
    pip install -r requirements-dev.txt
    print_success "Dev-only dependencies installed"
else
    print_error "requirements-dev.txt not found"
    exit 1
fi

# Create data directories
print_status "Creating data directories..."
mkdir -p data/models
mkdir -p data/embeddings
mkdir -p data/chroma
mkdir -p data/development

# Set environment variables for local development
print_status "Setting up environment variables..."
cat > .env.dev << EOF
# Contexis Dev-Only Environment
# No external API keys needed - everything runs locally

# Environment
CMP_ENV=development-dev

# Local Model Settings
CMP_LOCAL_MODELS=true
CMP_OFFLINE_MODE=true

# Model Cache
CMP_MODEL_CACHE_DIR=./data/models

# Database
CMP_DB_PROVIDER=sqlite
CMP_DB_PATH=./data/development/development.db

# Vector Database
CMP_VECTOR_DB_PROVIDER=chroma
CMP_VECTOR_DB_PATH=./data/embeddings
CMP_CHROMA_PERSIST_DIR=./data/chroma

# Logging
CMP_LOG_LEVEL=debug
CMP_LOG_FORMAT=json

# Development Features
CMP_HOT_RELOAD=true
CMP_DEBUG_MODE=true
CMP_MOCK_PROVIDERS=false
CMP_ENABLE_TELEMETRY=false
EOF

print_success "Environment file created: .env.dev"

# Create a simple test script
print_status "Creating test script..."
cat > test-dev.py << 'EOF'
#!/usr/bin/env python3
"""
Test script for Contexis Dev-Only Environment
"""

import os
import sys
from pathlib import Path

# Add the project root to Python path
project_root = Path(__file__).parent
sys.path.insert(0, str(project_root))

def test_local_provider():
    """Test the local AI provider."""
    try:
        from src.providers.local_provider import create_local_ai_provider
        
        config = {
            "model": "microsoft/Phi-3-mini-4k-instruct",
            "temperature": 0.1,
            "max_tokens": 100,
            "device": "cpu",  # Use CPU for testing
            "model_cache": {"directory": "./data/models"}
        }
        
        print("🧪 Testing local AI provider...")
        provider = create_local_ai_provider(config)
        
        # Test generation
        prompt = "Hello! How are you today?"
        response = provider.generate(prompt)
        
        print(f"✅ Local AI provider test successful!")
        print(f"📝 Prompt: {prompt}")
        print(f"🤖 Response: {response}")
        
        return True
        
    except Exception as e:
        print(f"❌ Local AI provider test failed: {e}")
        return False

def test_local_embeddings():
    """Test the local embeddings provider."""
    try:
        from src.providers.local_provider import create_local_embeddings_provider
        
        config = {
            "model": "all-MiniLM-L6-v2",
            "device": "cpu",
            "model_cache": {"directory": "./data/models"}
        }
        
        print("🧪 Testing local embeddings provider...")
        provider = create_local_embeddings_provider(config)
        
        # Test embedding
        text = "This is a test sentence for embeddings."
        embedding = provider.embed_single(text)
        
        print(f"✅ Local embeddings provider test successful!")
        print(f"📝 Text: {text}")
        print(f"🔢 Embedding dimensions: {len(embedding)}")
        
        return True
        
    except Exception as e:
        print(f"❌ Local embeddings provider test failed: {e}")
        return False

def main():
    """Run all tests."""
    print("🚀 Contexis Dev-Only Environment Test")
    print("=====================================")
    
    # Set environment
    os.environ["CMP_ENV"] = "development-dev"
    
    # Run tests
    ai_success = test_local_provider()
    embeddings_success = test_local_embeddings()
    
    print("\n📊 Test Results:")
    print(f"   AI Provider: {'✅ PASS' if ai_success else '❌ FAIL'}")
    print(f"   Embeddings: {'✅ PASS' if embeddings_success else '❌ FAIL'}")
    
    if ai_success and embeddings_success:
        print("\n🎉 All tests passed! Your dev-only environment is ready.")
        print("\n📝 Next steps:")
        print("   1. Run: ctx init MyProject")
        print("   2. cd MyProject")
        print("   3. cp ../.env.dev .env")
        print("   4. ctx generate rag MyRAG --db=sqlite --embeddings=local")
        print("   5. ctx run MyRAG 'Hello, world!'")
    else:
        print("\n❌ Some tests failed. Please check the error messages above.")
        sys.exit(1)

if __name__ == "__main__":
    main()
EOF

chmod +x test-dev.py
print_success "Test script created: test-dev.py"

# Create a quick start guide
print_status "Creating quick start guide..."
cat > QUICKSTART-DEV.md << 'EOF'
# Contexis Dev-Only Quick Start

This guide helps you get started with Contexis using local models (Phi-3.5-Mini) and Chroma(SQLite) without any external API calls.

## 🚀 Quick Start

1. **Activate the virtual environment:**
   ```bash
   source venv/bin/activate
   ```

2. **Test the environment:**
   ```bash
   python test-dev.py
   ```

3. **Create a new project:**
   ```bash
   ctx init MyProject
   cd MyProject
   ```

4. **Set up the dev environment:**
   ```bash
   cp ../.env.dev .env
   ```

5. **Create your first RAG system:**
   ```bash
   ctx generate rag MyRAG --db=sqlite --embeddings=local
   ```

6. **Add some knowledge:**
   ```bash
   echo "Contexis is a Context-Memory-Prompt framework for AI applications." > memory/documents/intro.txt
   ctx memory ingest --provider=sqlite --component=MyRAG --input=memory/documents/intro.txt
   ```

7. **Test your system:**
   ```bash
   ctx run MyRAG "What is Contexis?"
   ```

## 🔧 Configuration

The dev-only environment uses:
- **AI Model**: Microsoft Phi-3-mini-4k-instruct (local)
- **Embeddings**: all-MiniLM-L6-v2 (local)
- **Vector DB**: Chroma with SQLite backend
- **Database**: SQLite

## 📁 Directory Structure

```
data/
├── models/          # Downloaded models
├── embeddings/      # Vector embeddings
├── chroma/         # Chroma database
└── development/    # SQLite database
```

## 🛠️ Troubleshooting

### Memory Issues
If you encounter memory issues, try:
1. Use CPU instead of GPU: Set `device: "cpu"` in config
2. Enable quantization: Set `load_in_8bit: true` or `load_in_4bit: true`

### Model Download Issues
Models are automatically downloaded on first use. If download fails:
1. Check internet connection
2. Clear cache: `rm -rf data/models`
3. Try again

### Performance Tips
- Use GPU if available (CUDA)
- Enable quantization for memory optimization
- Use smaller models for faster inference

## 📚 Next Steps

- Read the main README.md for advanced features
- Check out examples/ for sample projects
- Explore the documentation in docs/
EOF

print_success "Quick start guide created: QUICKSTART-DEV.md"

print_status "Setup complete!"
print_success "Your Contexis dev-only environment is ready!"
echo ""
print_status "Next steps:"
echo "  1. Activate virtual environment: source venv/bin/activate"
echo "  2. Test the environment: python test-dev.py"
echo "  3. Create a project: ctx init MyProject"
echo "  4. Follow the guide: cat QUICKSTART-DEV.md"
echo ""
print_success "🎉 Happy coding with Contexis!"
