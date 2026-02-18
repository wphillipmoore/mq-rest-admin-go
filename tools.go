//go:build tools

package tools

import (
	_ "github.com/fzipp/gocyclo/cmd/gocyclo"
	_ "github.com/vladopajic/go-test-coverage/v2"
	_ "golang.org/x/vuln/cmd/govulncheck"
)
