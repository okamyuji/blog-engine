#!/usr/bin/env bash
set -e

echo "=================================="
echo "Running Go Quality Checks"
echo "=================================="

echo ""
echo "Step 1/4: Running go fmt..."
go fmt ./...

echo ""
echo "Step 2/4: Running go vet..."
go vet ./...

echo ""
echo "Step 3/4: Running golangci-lint..."
golangci-lint run ./...

echo ""
echo "Step 4/4: Running tests with shuffle..."
go test --shuffle=on --count=1 ./...

echo ""
echo "=================================="
echo "✓ All checks passed successfully!"
echo "=================================="

