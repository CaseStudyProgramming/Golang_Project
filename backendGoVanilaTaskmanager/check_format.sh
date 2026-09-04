#!/bin/bash
# Check if Go code is properly formatted

echo "Checking Go code formatting..."
UNFORMATTED=$(gofmt -l .)

if [ -n "$UNFORMATTED" ]; then
    echo "The following files are not properly formatted:"
    echo "$UNFORMATTED"
    echo "Run 'gofmt -w .' to format all files"
    exit 1
else
    echo "All Go files are properly formatted ✅"
    exit 0
fi
