#!/bin/bash
# Run tests with race detection enabled

echo "Running tests with race detection..."
go test -race ./...

echo "Test with race detection completed"
