#!/usr/bin/env bash
# Run golangci-lint with development build tag.
# Prefers the binary in GOPATH/bin; falls back to whatever is on PATH.
set -euo pipefail

LINT_BIN="$(go env GOPATH)/bin/golangci-lint"
if [ ! -x "$LINT_BIN" ]; then
  LINT_BIN=golangci-lint
fi
"$LINT_BIN" run --build-tags=development ./...
