# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/)
and this project adheres to [Semantic Versioning](https://semver.org/).

## [1.1.6] - 2026-02-19

### Bug fixes

- sync prepare_release.py ruff lint fixes from canonical (#98)
- sync prepare_release.py empty changelog abort from canonical (#100)
- sync shared tooling to v1.0.2
- sync hook and lint scripts from standards-and-conventions (#113)

### CI

- auto-add issues to GitHub Project (#109)

### Documentation

- rename mq-dev-environment references to mq-rest-admin-dev-environment (#114)

### Features

- add canonical local validation script (#104)
- add shared tooling sync from standard-tooling v1.0.0
- apply gofmt -s simplifications for goreportcard.com compliance (#115)
- reduce mqscCommand cyclomatic complexity for goreportcard.com (#116)
- add tools.go for development tool dependency management (#118)

## [1.1.5] - 2026-02-17

### Bug fixes

- truncate docs version to major.minor (#89)

### Features

- use GitHub App token for bump PR to trigger CI (#91)

## [1.1.4] - 2026-02-16

### Bug fixes

- sync prepare_release.py with canonical version
- sync prepare_release.py merge message fix from canonical (#74)
- add cliff.toml for markdownlint-compliant changelog generation (#78)
- sync prepare_release.py changelog conflict fix from canonical (#81)

## [1.1.3] - 2026-02-16

### Bug fixes

- remove PR_BUMP_TOKEN and add issue linkage to bump PR (#66)

## [1.1.2] - 2026-02-16

### Bug fixes

- remove duplicate CI runs and relax coverage threshold (#60)

## [1.1.1] - 2026-02-16

### Bug fixes

- resolve golangci-lint issues in auth tests
- disable MD041 for mkdocs snippet-include files
- correct snippets base_path resolution for fragment includes (#33)
- run mike from repo root so snippet base_path resolves in CI (#34)
- propagate 4 missing mapping entries from canonical JSON (#39)
- move coverage:ignore annotations to preceding line (#44)
- set docs default version to latest on main deploy
- install prepare_release.py and bump version to 1.1.1
- allow commits on release/* branches in library-release model
- rename integration-test job to integration-tests (#56)

### Documentation

- normalize README formatting (#3)
- expand README with installation, quick start, and API overview (#8)
- update implementation progress for phases 6-7 (#10)
- update coverage target from 100% to 99%
- add MkDocs Material documentation site (#16)
- fix hallucinated API references and correct ensure method count (#19)
- address medium-severity documentation consistency findings (#21)
- address cross-library documentation consistency nits (#25)
- switch to shared fragment includes from common repo (#30)
- add quality gates documentation page

### Features

- Go port of pymqrest (mq-rest-admin-go) (#1)
- add CI/CD workflow, git hooks, and linter configuration (#2)
- add Tier 1 security tooling (CodeQL, license compliance)
- add Trivy and Semgrep CI jobs
- add version constant set to 1.1.0 (#41)
- replace go-ignore-cov with go-test-coverage (#46)
- add publish workflow and versioned docs deployment (#52)

### Refactoring

- rename package from mqrest to mqrestadmin (#6)

### Testing

- add table-driven tests for all command wrapper methods
- add ensure and sync wrapper coverage
- add mapping edge case coverage
- add HTTPTransport and buildClient coverage
- cover auth, session init, errors, and remaining edge cases
