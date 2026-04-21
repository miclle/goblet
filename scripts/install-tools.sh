#!/usr/bin/env bash
# Install/update development tools: reflex, staticcheck, golangci-lint.
# Idempotent: existing binaries are left alone.
set -euo pipefail

GOLANGCI_LINT_VERSION="v2.11.0"

which reflex > /dev/null 2>&1 || go install github.com/cespare/reflex@latest
which staticcheck > /dev/null 2>&1 || go install honnef.co/go/tools/cmd/staticcheck@latest

GOBIN="$(go env GOPATH)/bin"
mkdir -p "$GOBIN"
if [ ! -x "$GOBIN/golangci-lint" ]; then
  GOBIN="$GOBIN" GO111MODULE=on go install "github.com/golangci/golangci-lint/v2/cmd/golangci-lint@${GOLANGCI_LINT_VERSION}"
fi
