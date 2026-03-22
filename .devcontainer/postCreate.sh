#!/usr/bin/env bash
set -euo pipefail

echo "==> Installing Go tools..."

# Hot reload
go install github.com/air-verse/air@latest

# Linter
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# DB migrations
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Mock generator
go install github.com/vektra/mockery/v2@latest

echo "==> Installing project dependencies..."
make setup

echo "==> Done! Run 'make dev' to start the development environment."
