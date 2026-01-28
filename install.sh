#!/bin/bash

# go-find installer script
# Usage: curl -sL https://go-find.sauravdev.in | bash

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configuration
REPO="saurav-tiwari03/go-find"
INSTALL_DIR="${HOME}/.local/bin"
BINARY_NAME="go-find"

echo -e "${GREEN}🚀 go-find Installer${NC}"
echo "─────────────────────────────────"

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo -e "${RED}❌ Go is not installed${NC}"
    echo "Please install Go from https://golang.org/dl/"
    exit 1
fi

echo -e "${GREEN}✓ Go detected${NC}"

# Create install directory if it doesn't exist
mkdir -p "$INSTALL_DIR"

# Clone or update the repository
TEMP_DIR=$(mktemp -d)
trap "rm -rf $TEMP_DIR" EXIT

echo -e "${YELLOW}📥 Downloading go-find...${NC}"
cd "$TEMP_DIR"
git clone --depth 1 "https://github.com/${REPO}.git" . 2>/dev/null || {
    echo -e "${RED}❌ Failed to download go-find${NC}"
    exit 1
}

# Build the binary
echo -e "${YELLOW}🔨 Building...${NC}"
go build -o "$INSTALL_DIR/$BINARY_NAME" .

# Make it executable
chmod +x "$INSTALL_DIR/$BINARY_NAME"

# Check if ~/.local/bin is in PATH
if [[ ":$PATH:" != *":$INSTALL_DIR:"* ]]; then
    echo -e "${YELLOW}⚠️  Note: $INSTALL_DIR is not in your PATH${NC}"
    echo "Add this line to your shell configuration file (~/.bashrc, ~/.zshrc, etc):"
    echo -e "${YELLOW}export PATH=\"\$HOME/.local/bin:\$PATH\"${NC}"
fi

echo "─────────────────────────────────"
echo -e "${GREEN}✅ Installation complete!${NC}"
echo "Run: $BINARY_NAME"
