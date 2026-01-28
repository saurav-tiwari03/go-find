#!/bin/bash

# go-find - Fast, minimal file explorer written in Go
# Usage: curl -sL https://go-find.sauravdev.in | bash

# Get target directory from argument or use current directory
TARGET_DIR="${1:-.}"

echo "🚀 go-find"
echo "─────────────────────────────────"

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed"
    echo "Please install Go from https://golang.org/dl/"
    exit 1
fi

# Create temporary directory for download and build
TEMP_DIR=$(mktemp -d)
cleanup() {
    rm -rf "$TEMP_DIR"
}
trap cleanup EXIT

echo "📥 Downloading go-find..."
cd "$TEMP_DIR"

# Download the repository
if ! git clone --depth 1 "https://github.com/saurav-tiwari03/go-find.git" . > /dev/null 2>&1; then
    echo "❌ Failed to download go-find"
    exit 1
fi

# Build the application
echo "🔨 Building..."
if ! go build -o go-find . 2>&1; then
    echo "❌ Failed to build go-find"
    exit 1
fi

# Run the application
echo "▶️  Running..."
echo ""
./go-find "$TARGET_DIR"

echo ""
echo "✅ Done!"





