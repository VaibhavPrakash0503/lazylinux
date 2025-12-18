#!/usr/bin/env bash

set -e

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m'

BINARY_NAME="lazylinux"
INSTALL_DIR="$HOME/.local/bin"

# Help message
show_help() {
  echo "LazyLinux Setup Script"
  echo ""
  echo "Usage: ./setup.sh [option]"
  echo ""
  echo "Options:"
  echo "  --install, -i    Build and install lazylinux to ~/.local/bin"
  echo "  --uninstall, -u  Remove lazylinux from ~/.local/bin"
  echo "  --help, -h       Show this help message"
  echo ""
}

# Check if Go is installed
check_go() {
  if ! command -v go &>/dev/null; then
    echo -e "${RED}Error: Go is not installed. Please install Go 1.21+ first.${NC}"
    exit 1
  fi
}

# Check if PATH contains ~/.local/bin
check_path() {
  if [[ ":$PATH:" != *":$INSTALL_DIR:"* ]]; then
    echo ""
    echo -e "${YELLOW}⚠️  Warning: $INSTALL_DIR is not in your PATH${NC}"
    echo ""
    echo "Add it by running:"
    echo ""
    
    SHELL_NAME=$(basename "$SHELL")
    case "$SHELL_NAME" in
      bash)
        echo "  echo 'export PATH=\"\$HOME/.local/bin:\$PATH\"' >> ~/.bashrc"
        echo "  source ~/.bashrc"
        ;;
      zsh)
        echo "  echo 'export PATH=\"\$HOME/.local/bin:\$PATH\"' >> ~/.zshrc"
        echo "  source ~/.zshrc"
        ;;
      fish)
        echo "  fish_add_path ~/.local/bin"
        ;;
      *)
        echo "  echo 'export PATH=\"\$HOME/.local/bin:\$PATH\"' >> ~/.bashrc"
        echo "  source ~/.bashrc"
        ;;
    esac
  fi
}

# Install function
install_lazylinux() {
  echo -e "${BLUE}Installing LazyLinux...${NC}"
  echo ""

  # Check Go
  check_go

  # Build binary
  echo "Building binary..."
  go build -o "$BINARY_NAME" main.go

  if [ ! -f "$BINARY_NAME" ]; then
    echo -e "${RED}Error: Build failed${NC}"
    exit 1
  fi

  # Create install directory
  mkdir -p "$INSTALL_DIR"

  # Move binary
  echo "Installing to $INSTALL_DIR..."
  mv "$BINARY_NAME" "$INSTALL_DIR/"

  echo -e "${GREEN}✅ Installation complete!${NC}"
  echo ""

  # Check PATH
  check_path

  echo ""
  echo "Next steps:"
  echo "  1. Run 'lazylinux setup' to configure your package manager"
  echo "  2. Run 'lazylinux --help' to get started"
  echo ""
  echo -e "${YELLOW}💡 Tip: Add a shorter alias${NC}"
  echo ""
  echo "alias lzl='lazylinux'"
  echo ""
  echo "Then use: lzl install firefox"
}

# Uninstall function
uninstall_lazylinux() {
  echo -e "${BLUE}Uninstalling LazyLinux...${NC}"
  echo ""

  BINARY_PATH="$INSTALL_DIR/$BINARY_NAME"

  if [ -f "$BINARY_PATH" ]; then
    rm "$BINARY_PATH"
    echo -e "${GREEN}✅ Removed $BINARY_PATH${NC}"
  else
    echo -e "${YELLOW}LazyLinux is not installed at $BINARY_PATH${NC}"
  fi

  # Ask about config
  CONFIG_DIR="$HOME/.config/lazylinux"
  if [ -d "$CONFIG_DIR" ]; then
    echo ""
    read -p "Remove config directory ($CONFIG_DIR)? [y/N]: " remove_config
    if [[ "$remove_config" =~ ^([yY][eE][sS]|[yY])$ ]]; then
      rm -rf "$CONFIG_DIR"
      echo -e "${GREEN}✅ Removed config directory${NC}"
    else
      echo "Config directory kept"
    fi
  fi

  echo ""
  echo -e "${GREEN}Uninstall complete!${NC}"
}

# Main script logic
case "${1:-}" in
--install | -i)
  install_lazylinux
  ;;
--uninstall | -u)
  uninstall_lazylinux
  ;;
--help | -h)
  show_help
  ;;
"")
  echo -e "${RED}Error: No option specified${NC}"
  echo ""
  show_help
  exit 1
  ;;
*)
  echo -e "${RED}Error: Unknown option '$1'${NC}"
  echo ""
  show_help
  exit 1
  ;;
esac
