## [develop-v1.1.1] - 2026-02-16

### 🚀 Features

- Go port of pymqrest (mq-rest-admin-go) (#1)
- Add CI/CD workflow, git hooks, and linter configuration (#2)
- Add Tier 1 security tooling (CodeQL, license compliance)
- Add Trivy and Semgrep CI jobs
- Add version constant set to 1.1.0 (#41)
- Replace go-ignore-cov with go-test-coverage (#46)
- Add publish workflow and versioned docs deployment (#52)

### 🐛 Bug Fixes

- Resolve golangci-lint issues in auth tests
- Disable MD041 for mkdocs snippet-include files
- Correct snippets base_path resolution for fragment includes (#33)
- Run mike from repo root so snippet base_path resolves in CI (#34)
- Propagate 4 missing mapping entries from canonical JSON (#39)
- Move coverage:ignore annotations to preceding line (#44)
- Set docs default version to latest on main deploy
- Install prepare_release.py and bump version to 1.1.1
- Allow commits on release/* branches in library-release model

### 🚜 Refactor

- Rename package from mqrest to mqrestadmin (#6)

### 📚 Documentation

- Normalize README formatting (#3)
- Expand README with installation, quick start, and API overview (#8)
- Update implementation progress for phases 6-7 (#10)
- Update coverage target from 100% to 99%
- Add MkDocs Material documentation site (#16)
- Fix hallucinated API references and correct ensure method count (#19)
- Address medium-severity documentation consistency findings (#21)
- Address cross-library documentation consistency nits (#25)
- Switch to shared fragment includes from common repo (#30)
- Add quality gates documentation page

### 🧪 Testing

- Add table-driven tests for all command wrapper methods
- Add ensure and sync wrapper coverage
- Add mapping edge case coverage
- Add HTTPTransport and buildClient coverage
- Cover auth, session init, errors, and remaining edge cases

### ⚙️ Miscellaneous Tasks

- Enforce 100% coverage gate with go-ignore-cov exclusion annotations (#14)
- Add workflow_dispatch trigger to docs workflow (#17)
- Restrict docs deploy to main branch and manual dispatch (#23)
- Add mike versioned documentation deployment (#27)
- Trigger dev docs deployment on push to develop (#28)
- Trigger CI re-run with updated PR body
- Trigger CI re-run with updated PR body
- Improve CLAUDE.md with accurate commands and workflow docs (#50)
- Re-trigger CI with issue linkage
