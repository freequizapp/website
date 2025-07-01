#!/bin/bash
set -e

echo "Building Go binary for AWS Lambda..."
GOOS=linux GOARCH=amd64 go build -v -o main ./main.go

echo "Build complete"
