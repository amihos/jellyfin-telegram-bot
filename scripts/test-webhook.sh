#!/bin/bash
# Test script for webhook handler

set -e

echo "Running webhook handler tests..."
cd /home/huso/jellyfin-telegram-bot

# Run webhook tests only
go test -v -count=1 ./internal/handlers/

echo ""
echo "All webhook tests passed!"
