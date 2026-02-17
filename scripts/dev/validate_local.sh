#!/usr/bin/env bash
# Local validation script mirroring CI hard gates.
# See: https://github.com/wphillipmoore/standards-and-conventions/blob/develop/docs/repository/local-validation-scripts.md

set -euo pipefail

# --- Prerequisite checks ---

missing=()
for tool in go golangci-lint govulncheck; do
    if ! command -v "$tool" &>/dev/null; then
        missing+=("$tool")
    fi
done

if [[ ${#missing[@]} -gt 0 ]]; then
    echo "ERROR: Missing required tools: ${missing[*]}" >&2
    echo "Install with:" >&2
    for tool in "${missing[@]}"; do
        case "$tool" in
            go)              echo "  brew install go" >&2 ;;
            golangci-lint)   echo "  brew install golangci-lint" >&2 ;;
            govulncheck)     echo "  go install golang.org/x/vuln/cmd/govulncheck@latest" >&2 ;;
        esac
    done
    exit 1
fi

# --- Validation steps ---

run() {
    echo "Running: $*"
    "$@"
}

run go vet ./...
run golangci-lint run ./...
run go test -race -count=1 ./...
run govulncheck ./...
