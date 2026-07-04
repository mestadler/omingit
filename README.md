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

Or install directly with Go 1.26+:

```bash
go install github.com/mestadler/omingit/cmd/gg@latest
```

To embed a version string at build time:

```bash
go build -ldflags "-X main.version=$(git describe --tags --always --dirty)" -o gg ./cmd/gg
```

Put `gg` somewhere on your `PATH`.

## Usage

`gg --version` (or `gg -v`) prints the build version, set at build time via
`-ldflags` (defaults to `dev` otherwise).

```bash
gg pr list              # → gh pr list or tea pr list
gg issue create         # → gh issue create or tea issue create
gg workflow list        # → gh workflow list
gg run list             # → gh run list or Gitea SDK call
gg release create       # → gh release create or tea release create
gg repo clone owner/rep # → gh repo clone or tea repo clone
gg push origin main     # → git push origin main
gg commit -m "feat: …"  # → git commit -m "feat: …"
```

All platform commands recognised:

`pr`, `issue`, `repo`, `workflow`, `run`, `release`, `gist`, `label`, `project`,
`variable`, `secret`

On Forgejo repos:

- `pr`, `issue`, `repo`, `release`, `label` → `tea` (passthrough)
- `workflow`, `run` → gg's built-in Gitea SDK (direct API calls)
- `gist` → GitHub only (no Forgejo equivalent)
- `project` → unsupported on Forgejo

### `-C` flag

Run `gg` against a repo in a different directory. Must appear before the
subcommand:

```bash
gg -C /path/to/repo pr list
gg -C ~/other-project issue view 42
```

Everything else passes through to `git` verbatim. When `git` rejects an
unknown command, `gg` retries with `gh` or `tea` if either is available.
This makes `gg` forward-compatible with new `gh`/`tea` subcommands.

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

| Tool  | Needed for                                    | Install                              |
|-------|-----------------------------------------------|--------------------------------------|
| `git` | All operations                                | System package manager               |
| `gh`  | GitHub repos                                  | `https://cli.github.com`             |
| `tea` | Gitea/Forgejo (pr, issue, repo, release, label) | `https://gitea.com/gitea/tea`        |

Gitea `workflow` and `run` commands use gg's bundled Gitea SDK — no extra
binary needed.

For Gitea Actions commands (`workflow`, `run`), `gg` talks directly to the
Gitea API and needs `GITEA_TOKEN`. Set it before running `gg`:

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

## Shell completion

`gg` ships with bundled completions for bash, zsh, and fish. Generate and
source them:

```bash
# bash — add to ~/.bashrc
source <(gg completion bash)

# zsh — add to ~/.zshrc, or copy to a fpath directory
source <(gg completion zsh)

# fish — save to the completions directory
mkdir -p ~/.config/fish/completions
gg completion fish > ~/.config/fish/completions/gg.fish
```

## License

MIT — see [LICENSE](LICENSE) file.
