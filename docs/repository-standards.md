# Repository Standards

## Table of Contents

- [Pre-flight checklist](#pre-flight-checklist)
- [AI co-authors](#ai-co-authors)
- [Repository profile](#repository-profile)
- [Local validation](#local-validation)
- [Tooling requirement](#tooling-requirement)
- [Merge strategy override](#merge-strategy-override)
- [Approved domain abbreviations](#approved-domain-abbreviations)
- [Accepted naming deviations](#accepted-naming-deviations)

## Pre-flight checklist

- Before modifying any files, check the current branch with `git status -sb`.
- If on `develop`, create a short-lived `feature/*` branch or ask for explicit approval to proceed on `develop`.
- If approval is granted to work on `develop`, call it out in the response and proceed only for that user-approved scope.
- Enable repository git hooks before committing: `git config core.hooksPath scripts/git-hooks`.

## AI co-authors

- Co-Authored-By: wphillipmoore-codex <255923655+wphillipmoore-codex@users.noreply.github.com>
- Co-Authored-By: wphillipmoore-claude <255925739+wphillipmoore-claude@users.noreply.github.com>

## Repository profile

- repository_type: library
- versioning_scheme: library
- branching_model: library-release
- release_model: artifact-publishing
- supported_release_lines: current and previous
- canonical_local_validation_command: scripts/dev/validate_local.sh

## Local validation

Canonical local validation command:

```bash
scripts/dev/validate_local.sh
```

Individual checks (run by the validation script):

```bash
go vet ./...                    # Static analysis
golangci-lint run ./...         # Lint checks
go test -race -count=1 ./...   # Unit tests with race detection
govulncheck ./...               # Vulnerability scanning
```

Integration tests (require MQ environment, not included in validation script):

```bash
go test -race -count=1 -tags=integration ./...
```

## Tooling requirement

Required for daily workflow:

- Go 1.22+ (`brew install go` or <https://go.dev/dl/>)
- `golangci-lint` (`brew install golangci-lint`)
- `govulncheck` (`go install golang.org/x/vuln/cmd/govulncheck@latest`)
- `markdownlint` (required for docs validation and PR pre-submission)

## Merge strategy override

- Feature, bugfix, and chore PRs targeting `develop` use squash merges (`--squash`).
- Release PRs targeting `main` use regular merges (`--merge`) to preserve shared
  ancestry between `main` and `develop`.
- Auto-merge commands:
  - Feature PRs: `gh pr merge --auto --squash --delete-branch`
  - Release PRs: `gh pr merge --auto --merge --delete-branch`

## Approved domain abbreviations

Domain-specific abbreviations that are well-established in the IBM MQ ecosystem
and may be used in identifiers without expansion:

- `qmgr` — queue manager (established MQSC domain term)

## Accepted naming deviations

None yet.
