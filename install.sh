#!/usr/bin/env bash

set -e

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo "Installing LazyLinux..."

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo -e "${RED}Error: Go is not installed. Please install Go 1.21+ first.${NC}"
    exit 1
fi

# Build the binary
echo "Building binary..."
go build -o lazylinux ./cmd/lazylinux

# Install to /usr/local/bin
echo "Installing to /usr/local/bin (requires sudo)..."
sudo mv lazylinux /usr/local/bin/

echo -e "${GREEN}✅ Installation complete!${NC}"
echo ""
echo "Next steps:"
echo "  1. Run 'lazylinux setup' to configure your package manager"
echo "  2. Run 'lazylinux --help' to get started"
echo ""
echo -e "${YELLOW}💡 Tip: Add a shorter alias for convenience${NC}"
echo ""
