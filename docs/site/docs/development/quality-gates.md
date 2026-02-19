# Quality gates

--8<-- "development/quality-gates.md"

## Go-specific validation

The Go validation pipeline runs as individual commands:

```bash
go vet ./...
golangci-lint run
go test -race -count=1 -coverprofile=coverage.out ./...
```

This covers:

1. **go vet** — Standard Go static analysis
2. **golangci-lint** — Comprehensive linter suite
3. **go test -race** — Unit tests with race detector enabled
4. **Coverage enforcement** — 100% coverage after `go-test-coverage` exclusions
5. **govulncheck** — Dependency vulnerability scanning
6. **go-licenses** — License compliance checking

The CI matrix tests against Go 1.25 and 1.26. The package uses only the
Go standard library, so the dependency audit is minimal.

## External quality report

[![Go Report Card](https://goreportcard.com/badge/github.com/wphillipmoore/mq-rest-admin-go)](https://goreportcard.com/report/github.com/wphillipmoore/mq-rest-admin-go)

[Go Report Card](https://goreportcard.com/report/github.com/wphillipmoore/mq-rest-admin-go)
runs six automated checks against the codebase: go vet, gofmt, gocyclo,
ineffassign, license, and misspell. All of these overlap with checks
already enforced by the local validation pipeline and CI.
