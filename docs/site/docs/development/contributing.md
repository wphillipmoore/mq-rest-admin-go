# Contributing

This project welcomes contributions from humans working with or without
AI assistance. AI tooling is available but not required.

## Branching and workflow

All contributors follow the same branching model:

- Branch from `develop` using `feature/*`, `bugfix/*`, `hotfix/*`, or
  `chore/*` prefixes.
- Commit messages follow
  [conventional commits](https://www.conventionalcommits.org/) and are
  validated by CI.
- Feature PRs: squash merge to `develop`.
- Release PRs: regular merge to `main` (preserves shared ancestry).

## Commit conventions

Commits must follow the [Conventional Commits](https://www.conventionalcommits.org/)
format:

```text
<type>: <description>

[optional body]

[optional footer(s)]
```

Types: `feat`, `fix`, `docs`, `style`, `refactor`, `test`, `chore`

## Code quality requirements

All code must pass quality checks before merging:

- **Static analysis**: `go vet ./...`
- **Tests**: `go test -race ./...`
- **Build**: `go build ./...`

Run all checks locally before pushing:

```bash
go build ./... && go vet ./... && go test -race ./...
```

## Pull request process

1. Create a `feature/*` branch from `develop`
2. Make changes and ensure all checks pass
3. Open a PR targeting `develop`
4. CI runs the full validation pipeline
5. After review and approval, squash-merge into `develop`

## For human contributors

- Run quality checks before pushing to catch issues early.
- Reference `docs/repository-standards.md` for the full standards
  specification.
- The `CLAUDE.md` and `AGENTS.md` files document architecture,
  patterns, and key design decisions. They are useful as reference
  material even when not using an AI agent.

## For AI agent contributors

### Agent entry points

- **Claude Code**: reads `CLAUDE.md`, which loads repository standards
  via include directives.
- **Codex and other agents**: reads `AGENTS.md`, which loads the same
  standards plus shared skills from the `standards-and-conventions`
  repository.

### Quality expectations

AI-generated code must pass all the same validation gates listed
above. There are no exceptions.

### Co-author trailers

AI agents add co-author trailers to commits automatically when
following the repository standards.
