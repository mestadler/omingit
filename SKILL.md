# gg — Hybrid Git Router

`gg` is a smart CLI router that standardizes developer workflows across
multiple Git hosting providers. It inspects the origin remote URL to detect
the host (GitHub vs. Gitea/Forgejo) and routes commands to the appropriate
CLI binary.

## How It Works

```
gg <command> [args...]
```

**Git commands** are passed directly to `git`. Everything works: `commit`,
`push`, `pull`, `log`, `status`, `diff`, `branch`, `tag`, `config`, etc.

**Platform commands** are routed to `gh` (GitHub) or `tea` (Gitea/Forgejo):

| Command   | Description              |
|-----------|--------------------------|
| `gg pr`   | Pull request operations  |
| `gg issue`| Issue tracking           |
| `gg repo` | Repository management    |
| `gg run`  | CI/CD runs               |
| `gg release` | Release management   |
| `gg workflow` | Workflow management (GitHub) |
| `gg gist` | Gist management          |
| `gg label`| Label management         |
| `gg project` | Project board management |
| `gg variable` | CI/CD variables     |
| `gg secret` | CI/CD secrets        |

## Host Detection

Detection priority:

1. **`GG_HOST` env var** — explicit override (`github` or `gitea`)
2. **Origin remote URL** — `git remote get-url origin`
   - Contains `github.com` → GitHub (`gh`)
   - Any other URL → Gitea/Forgejo (`tea`)
3. **No remote, no override** → error with setup instructions

### Examples

```bash
# In a GitHub-hosted repo
gg pr create          # → gh pr create

# In a Gitea-hosted repo
gg pr create          # → tea pr create

# Override detection
GG_HOST=gitea gg pr create  # → tea pr create (even on GitHub)

# Git commands always pass through
gg push origin main   # → git push origin main
gg commit -m "feat"   # → git commit -m "feat"
```

## Requirements

| Binary | Required for | Installation                  |
|--------|-------------|-------------------------------|
| `git`  | Everything  | System package manager        |
| `gh`   | GitHub ops  | `https://cli.github.com`      |
| `tea`  | Gitea/Forgejo ops | `https://gitea.com/gitea/tea` |

## Environment Variables

| Variable     | Description                                   |
|-------------|-----------------------------------------------|
| `GG_HOST`   | Override host detection: `github` or `gitea`  |

## Exit Codes

`gg` propagates the exit code from the underlying command verbatim.
Scripts that check `$?` work identically to calling git/gh/tea directly.
