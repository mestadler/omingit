# omingit

`gg` is a hybrid git router. One CLI to rule both GitHub and Gitea/Forgejo.

It inspects the origin remote URL to detect which host you are on, then routes
your command to the right tool:

- **Git commands** (`commit`, `push`, `status`, `log`, …) → `git`
- **Platform commands** (`pr`, `issue`, `workflow`, …) → `gh` (GitHub) or `tea` (Gitea/Forgejo)

No more remembering which CLI goes with which repo. `gg` figures it out.

## Install

```bash
go build -o gg ./cmd/gg
```

To embed a version string at build time:

```bash
go build -ldflags "-X main.version=$(git describe --tags --always --dirty)" -o gg ./cmd/gg
```

Put `gg` somewhere on your `PATH`.

## Usage

```bash
gg pr list              # → gh pr list or tea pr list
gg issue create         # → gh issue create or tea issue create
gg workflow list        # → gh workflow list
gg run list             # → gh run list or tea runs ls
gg release create       # → gh release create or tea release create
gg repo clone owner/rep # → gh repo clone or tea repo clone
gg push origin main     # → git push origin main
gg commit -m "feat: …"  # → git commit -m "feat: …"
```

All platform commands recognised:

`pr`, `issue`, `repo`, `workflow`, `run`, `release`, `gist`, `label`, `project`,
`variable`, `secret`

Everything else passes through to `git` verbatim.

## Host detection

Detection runs in this order:

1. **`GG_HOST` env var** — explicit override. Set to `github` or `gitea`.
2. **Origin remote URL** — `git remote get-url origin`
   - Contains `github.com` → GitHub
   - Any other URL → Gitea/Forgejo (covers self-hosted instances)
3. **No remote, no override** → error

```bash
GG_HOST=gitea gg pr list   # force tea even on a GitHub repo
GG_HOST=github gg pr list  # force gh even on a Gitea repo
```

## Requirements

| Tool  | Needed for           | Install                              |
|-------|----------------------|--------------------------------------|
| `git` | All operations       | System package manager               |
| `gh`  | GitHub repos         | `https://cli.github.com`             |
| `tea` | Gitea/Forgejo repos  | `https://gitea.com/gitea/tea`        |

For Gitea Actions commands (`workflow`, `run`), `tea` needs a token. Set the
`GITEA_TOKEN` environment variable before running `gg`:

```bash
export GITEA_TOKEN=<your-gitea-token>
gg workflow list
```

## Exit codes

`gg` passes through the exit code from the underlying command. Scripts that
check `$?` behave identically to calling `git`, `gh`, or `tea` directly.

## Agent usage

`gg` ships with a skill definition for OpenCode and compatible agents. See
[SKILL.md](SKILL.md) for the full agent-facing documentation, including the
routing rules your agent should follow when using `gg`.

## License

[TODO]
