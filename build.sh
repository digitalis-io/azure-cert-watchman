#!/bin/bash

set -e

APP_NAME="azure-cert-watchman"
VERSION="1.0.0"

echo "Building $APP_NAME v$VERSION..."

echo "Building for current platform..."
go build -ldflags="-s -w -X main.Version=$VERSION" -o $APP_NAME main.go

echo ""
echo "Building cross-platform binaries..."

echo "Building for Linux AMD64..."
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w -X main.Version=$VERSION" -o ${APP_NAME}-linux-amd64 main.go

echo "Building for Windows AMD64..."
GOOS=windows GOARCH=amd64 go build -ldflags="-s -w -X main.Version=$VERSION" -o ${APP_NAME}-windows-amd64.exe main.go

echo "Building for macOS AMD64..."
GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w -X main.Version=$VERSION" -o ${APP_NAME}-darwin-amd64 main.go

echo "Building for macOS ARM64 (M1/M2)..."
GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w -X main.Version=$VERSION" -o ${APP_NAME}-darwin-arm64 main.go

echo ""
echo "Build complete! Binaries created:"
ls -lh ${APP_NAME}*

echo ""
echo "Usage:"
echo "  ./$APP_NAME --help"
echo ""
echo "Or set environment variables:"
echo "  export AZURE_KEY_VAULT_URI=https://your-vault.vault.azure.net/"
echo "  export AZURE_CLIENT_ID=your-client-id"
echo "  export AZURE_TENANT_ID=your-tenant-id"
echo "  export AZURE_CLIENT_SECRET=your-client-secret"
echo "  ./$APP_NAME"