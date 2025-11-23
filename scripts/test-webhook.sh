#!/bin/bash
# Test script for webhook handler

set -e

# Get the directory where this script is located
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# Navigate to project root (one level up from scripts/)
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

echo "Running webhook handler tests..."
cd "$PROJECT_ROOT"

# Run webhook tests only
go test -v -count=1 ./internal/handlers/

echo ""
echo "All webhook tests passed!"
