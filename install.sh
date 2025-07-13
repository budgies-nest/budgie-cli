#!/bin/bash

# Budgie CLI Installation Script

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print colored output
print_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if binary exists
if [ ! -f "budgie" ]; then
    print_error "Binary 'budgie' not found in current directory"
    print_info "Please build the binary first with: go build -o budgie main.go"
    exit 1
fi

# Check if budgie.config.json exists
if [ ! -f "budgie.config.json" ]; then
    print_error "Config file 'budgie.config.json' not found in current directory"
    exit 1
fi

# Check if budgie.system.md exists
if [ ! -f "budgie.system.md" ]; then
    print_warning "System file 'budgie.system.md' not found in current directory"
    print_info "You may need to create a system instructions file"
fi

# Installation options
echo "Budgie CLI Installation"
echo "======================"
echo "1. Install to /usr/local/bin (system-wide, requires sudo)"
echo "2. Install to ~/bin (user-only)"
echo "3. Create symlink in /usr/local/bin (requires sudo)"
echo "4. Cancel installation"
echo ""

read -p "Choose installation option (1-4): " choice

case $choice in
    1)
        print_info "Installing to /usr/local/bin (system-wide)..."
        
        # Check if /usr/local/bin exists
        if [ ! -d "/usr/local/bin" ]; then
            print_error "/usr/local/bin directory does not exist"
            exit 1
        fi
        
        # Copy binary and config
        sudo cp budgie /usr/local/bin/
        sudo cp budgie.config.json /usr/local/bin/
        
        # Copy budgie.system.md if it exists
        if [ -f "budgie.system.md" ]; then
            sudo cp budgie.system.md /usr/local/bin/
        fi
        
        # Set permissions
        sudo chmod +x /usr/local/bin/budgie
        
        print_info "Installation completed successfully!"
        print_info "You can now run 'budgie' from anywhere"
        ;;
        
    2)
        print_info "Installing to ~/bin (user-only)..."
        
        # Create ~/bin if it doesn't exist
        mkdir -p ~/bin
        
        # Copy binary and config
        cp budgie ~/bin/
        cp budgie.config.json ~/bin/
        
        # Copy budgie.system.md if it exists
        if [ -f "budgie.system.md" ]; then
            cp budgie.system.md ~/bin/
        fi
        
        # Set permissions
        chmod +x ~/bin/budgie
        
        # Check if ~/bin is in PATH
        if [[ ":$PATH:" != *":$HOME/bin:"* ]]; then
            print_warning "~/bin is not in your PATH"
            print_info "Add the following line to your shell profile (~/.bashrc, ~/.zshrc, etc.):"
            echo "export PATH=\"\$HOME/bin:\$PATH\""
            print_info "Then restart your terminal or run: source ~/.bashrc (or ~/.zshrc)"
        fi
        
        print_info "Installation completed successfully!"
        ;;
        
    3)
        print_info "Creating symlink in /usr/local/bin..."
        
        # Check if /usr/local/bin exists
        if [ ! -d "/usr/local/bin" ]; then
            print_error "/usr/local/bin directory does not exist"
            exit 1
        fi
        
        # Get absolute path
        CURRENT_DIR=$(pwd)
        
        # Create symlink
        sudo ln -sf "$CURRENT_DIR/budgie" /usr/local/bin/budgie
        
        print_info "Symlink created successfully!"
        print_warning "Note: budgie.config.json and budgie.system.md will be read from the current directory"
        print_info "Make sure to run 'budgie' from this directory or use absolute paths for --config and --system flags"
        ;;
        
    4)
        print_info "Installation cancelled"
        exit 0
        ;;
        
    *)
        print_error "Invalid option. Please choose 1-4"
        exit 1
        ;;
esac

# Test installation
echo ""
print_info "Testing installation..."

if command -v budgie >/dev/null 2>&1; then
    print_info "✓ 'budgie' command is available"
    
    # Test help command
    if budgie --help >/dev/null 2>&1; then
        print_info "✓ Help command works"
    else
        print_warning "Help command test failed"
    fi
else
    print_warning "'budgie' command not found in PATH"
    print_info "You may need to restart your terminal or update your PATH"
fi

echo ""
print_info "Installation process completed!"
print_info "Example usage: budgie ask --question \"What is Go programming language?\""