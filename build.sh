#!/bin/bash

# Cross-platform build script for all tools
# This script compiles all Go tools for Windows, macOS (x64 & ARM64), and Linux

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Project root directory
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
RELEASES_DIR="$PROJECT_ROOT/releases"

# Target platforms (format: "platform_dir:GOOS:GOARCH")
PLATFORMS="windows:windows:amd64
mac-arm64:darwin:arm64
mac-x64:darwin:amd64
linux:linux:amd64"

# Find all tool directories (directories containing go.mod)
find_tools() {
    for dir in "$PROJECT_ROOT"/*/; do
        if [[ -f "${dir}go.mod" ]]; then
            basename "$dir"
        fi
    done
}

# Print banner
print_banner() {
    echo -e "${BLUE}"
    echo "=============================================="
    echo "  Cross-Platform Build Script"
    echo "=============================================="
    echo -e "${NC}"
}

# Print step
print_step() {
    echo -e "${YELLOW}>>> $1${NC}"
}

# Print success
print_success() {
    echo -e "${GREEN}✓ $1${NC}"
}

# Print error
print_error() {
    echo -e "${RED}✗ $1${NC}"
}

# Clean releases directory
clean_releases() {
    print_step "Cleaning releases directory..."
    rm -rf "$RELEASES_DIR"
    mkdir -p "$RELEASES_DIR"
    print_success "Releases directory cleaned"
}

# Create platform directories
create_platform_dirs() {
    print_step "Creating platform directories..."
    echo "$PLATFORMS" | while IFS=: read -r platform_dir goos goarch; do
        mkdir -p "$RELEASES_DIR/$platform_dir"
    done
    print_success "Platform directories created"
}

# Build tool for a specific platform
build_tool() {
    local tool_name="$1"
    local platform_dir="$2"
    local goos="$3"
    local goarch="$4"
    
    local tool_dir="$PROJECT_ROOT/$tool_name"
    local output_name="$tool_name"
    
    # Add .exe extension for Windows
    if [[ "$goos" == "windows" ]]; then
        output_name="${tool_name}.exe"
    fi
    
    local output_path="$RELEASES_DIR/$platform_dir/$output_name"
    
    printf "  Building %-20s for %-12s ... " "$tool_name" "$platform_dir"
    
    cd "$tool_dir"
    
    if GOOS="$goos" GOARCH="$goarch" CGO_ENABLED=0 go build -ldflags="-s -w" -o "$output_path" . 2>/dev/null; then
        echo -e "${GREEN}done${NC}"
        return 0
    else
        echo -e "${RED}failed${NC}"
        return 1
    fi
}

# Build all tools for all platforms
build_all() {
    local tools=$(find_tools)
    local tool_count=$(echo "$tools" | wc -l | tr -d ' ')
    
    if [[ -z "$tools" ]]; then
        print_error "No tools found to build"
        exit 1
    fi
    
    echo -e "${BLUE}Found tools:${NC}"
    echo "$tools" | while read -r tool; do
        echo "  - $tool"
    done
    echo ""
    
    local platform_count=$(echo "$PLATFORMS" | wc -l | tr -d ' ')
    local total_builds=$((tool_count * platform_count))
    local successful_builds=0
    local failed_builds=0
    
    echo "$PLATFORMS" | while IFS=: read -r platform_dir goos goarch; do
        print_step "Building for $platform_dir ($goos/$goarch)..."
        
        echo "$tools" | while read -r tool; do
            build_tool "$tool" "$platform_dir" "$goos" "$goarch"
        done
        echo ""
    done
    
    # Print summary
    echo -e "${BLUE}=============================================="
    echo "  Build Summary"
    echo "==============================================${NC}"
    echo ""
    
    # List all built binaries
    print_step "Built binaries:"
    echo "$PLATFORMS" | while IFS=: read -r platform_dir goos goarch; do
        echo -e "  ${YELLOW}$platform_dir/${NC}"
        for file in "$RELEASES_DIR/$platform_dir"/*; do
            if [[ -f "$file" ]]; then
                local size=$(du -h "$file" | cut -f1)
                echo "    - $(basename "$file") ($size)"
            fi
        done
    done
}

# Main function
main() {
    print_banner
    
    # Check if Go is installed
    if ! command -v go &> /dev/null; then
        print_error "Go is not installed. Please install Go first."
        exit 1
    fi
    
    echo -e "Go version: $(go version)"
    echo ""
    
    clean_releases
    create_platform_dirs
    echo ""
    build_all
    
    echo ""
    print_success "Build completed! Binaries are in: $RELEASES_DIR"
}

# Run main function
main "$@"
