#!/bin/bash

# go-find - Fast, minimal file explorer written in Go
# Usage: curl -sL https://go-find.sauravdev.in | bash

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Get target directory from argument or use current directory
TARGET_DIR="${1:-.}"

{
    echo -e "${GREEN}🚀 go-find${NC}"
    echo "─────────────────────────────────"

    # Check if Go is installed
    if ! command -v go &> /dev/null; then
        echo -e "${RED}❌ Go is not installed${NC}"
        echo "Please install Go from https://golang.org/dl/"
        exit 1
    fi

    # Create temporary directory for download and build
    TEMP_DIR=$(mktemp -d)
    trap "rm -rf $TEMP_DIR" EXIT

    echo -e "${YELLOW}📥 Downloading go-find...${NC}"
    cd "$TEMP_DIR"

    # Download the repository
    git clone --depth 1 "https://github.com/saurav-tiwari03/go-find.git" . 2>&1 | grep -v "Cloning into" || true

    # Build the application
    echo -e "${YELLOW}🔨 Building...${NC}"
    go build -o go-find . 2>&1 || {
        echo -e "${RED}❌ Failed to build go-find${NC}"
        exit 1
    }

    # Run the application
    echo -e "${YELLOW}▶️  Running...${NC}\n"
    ./go-find "$TARGET_DIR"

    echo -e "\n${GREEN}✅ Done!${NC}"
} 2>&1




