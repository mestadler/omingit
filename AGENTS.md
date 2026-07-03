# Project Instructions — omingit

> This file extends `~/.config/opencode/AGENTS.md` (global conventions).
> Global rules (commits, security, languages) apply unless overridden below.

## Rules
{"precedence":["R3 scopes R2: simplify only code you author this turn; never refactor existing code to do it.","R5 overrides R3 on style when a session agreement exists.","Ask (R1) only when no strong criterion (R4) or prior agreement (R5) resolves it."],"rules":[{"id":1,"name":"Think First","do":["State assumptions; ask if unsure.","If ambiguous, give options; no silent choices.","Propose simpler paths; push back.","If unclear, name the confusion and stop."]},{"id":2,"name":"Simplicity First","do":["No unrequested features, abstractions, config, or flexibility.","No error handling for cases that can't occur.","Cut 200 lines to 50 where doable.","Test: too complex for a senior? Simplify."]},{"id":3,"name":"Surgical Edits","do":["Touch only what the request needs; every changed line traces to it.","Match existing style; don't refactor or reformat unbroken code.","Delete only orphans your edit created; report other dead code, don't touch it.","Surface adjacent contradictions before coding."]},{"id":4,"name":"Goal-Driven","do":["Recast tasks as verifiable goals: 'add validation'->test invalid inputs, then pass; 'fix bug'->write a test that reproduces it, then pass; 'refactor X'->tests pass before and after.","Strong criteria auto-loop; weak ones force asking.","Plan as Step->verify pairs."]},{"id":5,"name":"Memory","do":["Apply agreed conventions without re-asking.","Name which agreement you break, and why, before editing.","At session start, surface unresolved decisions."]}]}

## Identity
You are working on github, forgejo and git cli.

## Project Structure
- `cmd/gg/main.go` — entry point, arg parsing, dispatch
- `internal/host/` — host detection (env override + URL parsing)
- `internal/classify/` — command classification (git vs platform)
- `data/` — project-local data, excluded from git
- `plans/` — plan agent writes RFC drafts here (has edit permission for this directory)
- `SKILL.md` — AI agent skill for using the `gg` binary
- `test_gg.sh` — integration test suite

## Build
```bash
# Build the gg binary in the project root
go build -o gg ./cmd/gg
```

## Test
```bash
# Unit tests
go test ./...

# Integration tests (requires building first)
go build -o gg ./cmd/gg && ./test_gg.sh
```

## Conventions
- **Versioning**: [TODO: document versioning]
- **CI**: [TODO: CI platform]

## ai-memory

This project has an ai-memory namespace. When using ai-memory tools, always pass
`project: "omingit"` to scope queries to this project's wiki.
