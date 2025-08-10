#!/bin/bash

# Contexis CMP Framework - Local Installation Script
# This script installs the CLI tool to the user's local directory and sets up PATH

set -e

echo "🚀 Installing Contexis CMP Framework..."

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

# Check if we're in the right directory
if [ ! -f "Makefile" ] || [ ! -f "go.mod" ]; then
    print_error "Please run this script from the contexis project root directory"
    exit 1
fi

# Build the CLI tool
print_status "Building the CLI tool..."
make build

# Create local bin directory
LOCAL_BIN="$HOME/.local/bin"
print_status "Creating local bin directory: $LOCAL_BIN"
mkdir -p "$LOCAL_BIN"

# Copy the CLI tool
print_status "Installing CLI tool to $LOCAL_BIN"
cp bin/ctx "$LOCAL_BIN/"

# Make it executable
chmod +x "$LOCAL_BIN/ctx"

print_success "CLI tool installed to $LOCAL_BIN/ctx"

# Detect shell and update PATH
SHELL_CONFIG=""
SHELL_NAME=""

if [ -n "$ZSH_VERSION" ]; then
    SHELL_CONFIG="$HOME/.zshrc"
    SHELL_NAME="zsh"
elif [ -n "$BASH_VERSION" ]; then
    SHELL_CONFIG="$HOME/.bashrc"
    SHELL_NAME="bash"
else
    # Try to detect from $SHELL
    case "$SHELL" in
        *zsh)
            SHELL_CONFIG="$HOME/.zshrc"
            SHELL_NAME="zsh"
            ;;
        *bash)
            SHELL_CONFIG="$HOME/.bashrc"
            SHELL_NAME="bash"
            ;;
        *)
            print_warning "Could not detect shell type. Please manually add to your shell configuration:"
            echo "export PATH=\"$LOCAL_BIN:\$PATH\""
            ;;
    esac
fi

# Update shell configuration if detected
if [ -n "$SHELL_CONFIG" ]; then
    print_status "Detected $SHELL_NAME shell, updating $SHELL_CONFIG"
    
    # Check if PATH is already set
    if grep -q "$LOCAL_BIN" "$SHELL_CONFIG" 2>/dev/null; then
        print_warning "PATH already contains $LOCAL_BIN in $SHELL_CONFIG"
    else
        # Add PATH export to shell config
        echo "" >> "$SHELL_CONFIG"
        echo "# Contexis CMP Framework" >> "$SHELL_CONFIG"
        echo "export PATH=\"$LOCAL_BIN:\$PATH\"" >> "$SHELL_CONFIG"
        print_success "Added PATH export to $SHELL_CONFIG"
    fi
fi

# Test the installation
print_status "Testing installation..."
if [ -f "$LOCAL_BIN/ctx" ]; then
    print_success "CLI tool is installed and executable"
    
    # Try to run version command
    if "$LOCAL_BIN/ctx" version >/dev/null 2>&1; then
        print_success "CLI tool is working correctly"
        echo ""
        echo "🎉 Installation completed successfully!"
        echo ""
        echo "To start using the CLI:"
        echo "1. Restart your terminal or run: source $SHELL_CONFIG"
        echo "2. Test the installation: ctx version"
        echo "3. Get help: ctx --help"
        echo ""
        echo "Quick start:"
        echo "  ctx init my-project"
        echo "  cd my-project"
        echo "  ctx generate rag CustomerDocs --db=sqlite --embeddings=sentence-transformers"
    else
        print_error "CLI tool installation test failed"
        exit 1
    fi
else
    print_error "CLI tool was not installed correctly"
    exit 1
fi

print_success "Installation completed! 🚀"
