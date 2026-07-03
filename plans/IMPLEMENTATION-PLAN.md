# Implementation Plan — gg Forgejo Actions Support

## Context

`gg` routes platform commands to `gh` (GitHub) or `tea` (Gitea/Forgejo) based on origin remote detection. `tea` is in maintenance mode and lacks Actions workflow management. This plan adds direct Gitea SDK calls for `workflow` and `run` commands on Forgejo repos.

## Phase 1 — Project Hygiene

### 1.1 Initialize git repo
- **Who**: git
- **Files**: `.git/`
- **Verify**: `git log` shows initial commit

### 1.2 Clean up runtime artifacts
- **Who**: build
- **Files**: remove `logs/`, `status/`; update `.gitignore`
- **Verify**: `git status` shows clean tree

### 1.3 Move misplaced document
- **Who**: build
- **Files**: move `plans/ARCHITECTURE-HANDOVER.md` → opencode project
- **Verify**: file no longer in omingit tree

### 1.4 Add README.md
- **Who**: docs
- **Files**: `README.md` (new)
- **Verify**: contains project description, build/install/test instructions

### 1.5 Add --version flag
- **Who**: build
- **Files**: `cmd/gg/main.go`
- **Verify**: `gg --version` prints version (from ldflags) or "dev"

### 1.6 Unit tests for main.go
- **Who**: test
- **Files**: `cmd/gg/main_test.go` (new)
- **Verify**: `go test ./cmd/gg/...` passes

## Phase 2 — Forgejo Actions (workflow + run via SDK)

### 2.1 Add SDK dependency
- **Who**: build
- **Files**: `go.mod`, `go.sum`
- **Verify**: `go mod tidy` succeeds, `go build` links cleanly

### 2.2 Implement internal/gitea package
- **Who**: build
- **Files**: `internal/gitea/client.go`, `internal/gitea/repo.go` (new)
- **Verify**: unit tests for client init, repo slug parsing

### 2.3 Implement workflow list/dispatch
- **Who**: build
- **Files**: `internal/gitea/workflow.go` (new)
- **Verify**: unit tests with mocked SDK client

### 2.4 Implement run list
- **Who**: build
- **Files**: `internal/gitea/run.go` (new)
- **Verify**: unit tests with mocked SDK client

### 2.5 Wire into dispatch
- **Who**: build
- **Files**: `cmd/gg/main.go`
- **Verify**: `GG_HOST=gitea` routes `workflow`/`run` to SDK, not tea

### 2.6 Integration tests
- **Who**: test
- **Files**: `test_gg.sh`
- **Verify**: mock SDK binary tests pass

## Auth Model

- `GITEA_TOKEN` environment variable (same pattern as `GG_HOST`)
- Token must have `read:repository` (list) or `write:repository` (dispatch)
- No stored credentials, no `gg login` command
- GitHub side unchanged (pure passthrough to `gh`)

## Deferred

- `gg variable` / `gg secret` — SDK supports it, not in this round
- `gg project` — blocked (not in SDK)
- `gg gist` — permanently blocked (no Forgejo feature)

## Verification

- `go test ./...` — existing + new tests pass
- `go vet ./...` — clean
- `./test_gg.sh` — all 20 existing tests + new integration tests pass
