#!/usr/bin/env bash
#
# test_gg.sh — Integration tests for the gg CLI.
#
# Tests host detection and command routing using mock gh/tea binaries
# in a temporary PATH, so tests are deterministic and self-contained.
#
set -euo pipefail

BIN="$(cd "$(dirname "$0")" && pwd)/gg"
PASS=0
FAIL=0
COUNT=0

ok()      { PASS=$((PASS+1)); echo "  ok  $COUNT - $1"; }
not_ok()  { FAIL=$((FAIL+1)); echo "  not ok  $COUNT - $1"; }

# --- setup ---
TMPDIR="$(mktemp -d)"
trap 'rm -rf "$TMPDIR"' EXIT

MOCKBIN="$TMPDIR/mockbin"
mkdir -p "$MOCKBIN"

cat > "$MOCKBIN/gh" <<'MOCKGH'
#!/usr/bin/env bash
echo "gh called with: $@"
exit 0
MOCKGH
chmod +x "$MOCKBIN/gh"

cat > "$MOCKBIN/tea" <<'MOCKTEA'
#!/usr/bin/env bash
echo "tea called with: $@"
exit 0
MOCKTEA
chmod +x "$MOCKBIN/tea"

# --- repo helpers ---
make_github_repo() {
	git -C "$1" init --quiet
	git -C "$1" remote add origin "git@github.com:org/repo.git"
}

make_gitea_repo() {
	git -C "$1" init --quiet
	git -C "$1" remote add origin "https://git.example.com/org/repo.git"
}

make_noremote_repo() {
	git -C "$1" init --quiet
}

# --- test helpers ---

# assert_output: run gg in dir with optional env vars, check output matches.
assert_output() {
	local name="$1" pattern="$2" dir="$3"
	shift 3
	COUNT=$((COUNT+1))
	local output
	output=$(cd "$dir" && PATH="$MOCKBIN:$PATH" "$BIN" "$@" 2>&1 || true)
	if echo "$output" | grep -q "$pattern"; then
		ok "$name"
	else
		not_ok "$name"
		echo "#       expected: $pattern"
		echo "#       got: $(echo "$output" | head -c 300)"
	fi
}

# assert_output_env: like assert_output but with extra env vars.
assert_output_env() {
	local name="$1" pattern="$2" dir="$3" envvars="$4"
	shift 4
	COUNT=$((COUNT+1))
	local output
	output=$(cd "$dir" && env $envvars PATH="$MOCKBIN:$PATH" "$BIN" "$@" 2>&1 || true)
	if echo "$output" | grep -q "$pattern"; then
		ok "$name"
	else
		not_ok "$name"
		echo "#       expected: $pattern"
		echo "#       got: $(echo "$output" | head -c 300)"
	fi
}

# assert_exit: run gg in dir, check exit code is non-zero.
assert_exit_fail() {
	local name="$1" dir="$2"
	shift 2
	COUNT=$((COUNT+1))
	if ! (cd "$dir" && "$BIN" "$@") >/dev/null 2>&1; then
		ok "$name"
	else
		not_ok "$name"
		echo "#       expected non-zero exit but got 0"
	fi
}

# --- 1. binary sanity ---
echo "# binary sanity"

COUNT=$((COUNT+1))
if [[ -x "$BIN" ]]; then
	ok "binary exists at $BIN"
else
	not_ok "binary not found — build with 'go build -o $BIN ./cmd/gg'"
	exit 1
fi

assert_output "no args prints usage" "Usage:" "$TMPDIR"
assert_output "gg help passes through to git" "usage: git" "$TMPDIR" help

# --- 2. native git pass-through ---
echo "# native git pass-through"

NOREMOTE="$(mktemp -d)"
make_noremote_repo "$NOREMOTE"

assert_output "gg status works" "On branch" "$NOREMOTE" status
assert_output "gg status exits 0" "" "$NOREMOTE" status
assert_exit_fail "gg log with invalid flag exits non-zero" "$NOREMOTE" log --invalid-flag
assert_output "gg commit --allow-empty" "root-commit" "$NOREMOTE" commit --allow-empty -m "test"

# --- 3. platform routing: GitHub remote → gh ---
echo "# platform routing (GitHub remote → gh)"

GHREPO="$(mktemp -d)"
make_github_repo "$GHREPO"

assert_output "pr create"    "gh called with: pr create"    "$GHREPO" pr create
assert_output "issue list"   "gh called with: issue list"  "$GHREPO" issue list
assert_output "repo view"    "gh called with: repo view"   "$GHREPO" repo view
assert_output "workflow run" "gh called with: workflow run" "$GHREPO" workflow run
assert_output "release create" "gh called with: release create" "$GHREPO" release create
assert_output "gist list"    "gh called with: gist list"   "$GHREPO" gist list

# --- 4. platform routing: Gitea remote → tea ---
echo "# platform routing (Gitea remote → tea)"

GITREPO="$(mktemp -d)"
make_gitea_repo "$GITREPO"

assert_output "pr create"    "tea called with: pr create"    "$GITREPO" pr create
assert_output "issue list"   "tea called with: issue list"  "$GITREPO" issue list
assert_output "release create" "tea called with: release create" "$GITREPO" release create

# --- 5. GG_HOST override ---
echo "# GG_HOST override"

assert_output_env "GG_HOST=github overrides Gitea → gh" \
	"gh called with: pr list" "$GITREPO" "GG_HOST=github" pr list

assert_output_env "GG_HOST=gitea overrides GitHub → tea" \
	"tea called with: pr list" "$GHREPO" "GG_HOST=gitea" pr list

# --- 6. error handling ---
echo "# error handling"

NOREMOTE2="$(mktemp -d)"
make_noremote_repo "$NOREMOTE2"

assert_output "platform command with no remote" \
	"unknown host" "$NOREMOTE2" pr list

# Test missing CLI by clearing PATH of everything except gg's location.
# GG_HOST=gitea forces tea lookup, which won't be found.
COUNT=$((COUNT+1))
output=$(cd "$GITREPO" && env GG_HOST=gitea PATH="$(dirname "$BIN")" "$BIN" pr list 2>&1 || true)
if echo "$output" | grep -qi "not found"; then
	ok "missing platform CLI shows error"
else
	not_ok "missing platform CLI shows error"
	echo "#       expected: 'not found' in output"
	echo "#       got: $(echo "$output" | head -c 300)"
fi

# --- Gitea SDK workflow/run commands ---
echo "# Gitea SDK workflow/run commands"

GITEASDKREPO="$(mktemp -d)"
make_gitea_repo "$GITEASDKREPO"

# gg workflow list: routed to SDK; fails because there is no real server.
COUNT=$((COUNT+1))
tmpout="$TMPDIR/gg_wf_list"
if (cd "$GITEASDKREPO" && env GG_HOST=gitea GITEA_TOKEN=fake-token PATH="$MOCKBIN:$PATH" "$BIN" workflow list) >"$tmpout" 2>&1; then
	not_ok "gg_workflow_list_gitea (SDK dispatch)"
	echo "#       expected non-zero exit but got 0"
else
	output=$(cat "$tmpout")
	if echo "$output" | grep -qi "not a git command"; then
		not_ok "gg_workflow_list_gitea (SDK dispatch)"
		echo "#       misrouted to git (saw 'not a git command')"
	elif echo "$output" | grep -qi "gg:"; then
		ok "gg_workflow_list_gitea (SDK dispatch)"
	else
		not_ok "gg_workflow_list_gitea (SDK dispatch)"
		echo "#       unexpected output (no gg: prefix)"
		echo "#       got: $(echo "$output" | head -c 300)"
	fi
fi

# gg run list: routed to SDK; fails because there is no real server.
COUNT=$((COUNT+1))
tmpout="$TMPDIR/gg_run_list"
if (cd "$GITEASDKREPO" && env GG_HOST=gitea GITEA_TOKEN=fake-token PATH="$MOCKBIN:$PATH" "$BIN" run list) >"$tmpout" 2>&1; then
	not_ok "gg_run_list_gitea (SDK dispatch)"
	echo "#       expected non-zero exit but got 0"
else
	output=$(cat "$tmpout")
	if echo "$output" | grep -qi "tea called"; then
		not_ok "gg_run_list_gitea (SDK dispatch)"
		echo "#       misrouted to tea (saw 'tea called')"
	elif echo "$output" | grep -qi "usage: git"; then
		not_ok "gg_run_list_gitea (SDK dispatch)"
		echo "#       misrouted to git (saw 'usage: git')"
	elif echo "$output" | grep -qi "gg:"; then
		ok "gg_run_list_gitea (SDK dispatch)"
	else
		not_ok "gg_run_list_gitea (SDK dispatch)"
		echo "#       unexpected output"
		echo "#       got: $(echo "$output" | head -c 300)"
	fi
fi

# gg pr list: should still route to tea, not the SDK.
assert_output_env "gg_pr_list_gitea (tea passthrough)" \
	"tea called with: pr list" "$GITEASDKREPO" "GG_HOST=gitea" pr list

# --- results ---
echo ""
echo "# $COUNT tests, $PASS pass, $FAIL fail"
if [[ $FAIL -gt 0 ]]; then
	exit 1
fi
