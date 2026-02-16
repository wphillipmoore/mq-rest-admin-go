# AI-assisted engineering

--8<-- "ai-engineering.md"

## Go-specific quality standards

**Test coverage**: All production code is covered by unit tests. Coverage
is enforced as a CI hard gate.

**Race detection**: All tests run with `go test -race` to catch
concurrent access issues.

**Static analysis**: golangci-lint runs with a comprehensive linter
configuration. Both `go vet` and additional linters run as CI hard gates.

**Zero dependencies**: The package uses only the Go standard library,
eliminating supply-chain risk and keeping the dependency audit trivial.
