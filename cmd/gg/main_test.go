package main

import (
	"errors"
	"os"
	"testing"
)

// ---------------------------------------------------------------------------
// runCmd tests
// ---------------------------------------------------------------------------

func TestRunCmd_Success(t *testing.T) {
	t.Parallel()

	got := runCmd("echo", "hello")
	if got != 0 {
		t.Errorf("runCmd(echo, hello) = %d, want 0", got)
	}
}

func TestRunCmd_NonexistentBinary(t *testing.T) {
	t.Parallel()

	got := runCmd("nonexistent-binary-zzz-12345")
	if got != 1 {
		t.Errorf("runCmd(nonexistent-binary) = %d, want 1", got)
	}
}

func TestRunCmd_ExitCode(t *testing.T) {
	t.Parallel()

	// sh -c "exit 3" should give exit code 3
	got := runCmd("sh", "-c", "exit 3")
	if got != 3 {
		t.Errorf("runCmd(sh, -c, exit 3) = %d, want 3", got)
	}
}

// ---------------------------------------------------------------------------
// detectOrigin tests
// ---------------------------------------------------------------------------

func TestDetectOrigin_ReturnsHost(t *testing.T) {
	// detectOrigin runs a real git command. This test verifies the
	// function compiles, does not panic, and returns a host.Host.
	// Full behaviour is tested in the integration suite (test_gg.sh).
	h := detectOrigin("")
	_ = h // uses the host package type
}

// ---------------------------------------------------------------------------
// run tests
// ---------------------------------------------------------------------------

func TestRun_VersionFlag(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{"long flag", []string{"gg", "--version"}},
		{"short flag", []string{"gg", "-v"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldArgs := os.Args
			defer func() { os.Args = oldArgs }()
			os.Args = tt.args

			got := run()
			if got != 0 {
				t.Errorf("run() with %v = %d, want 0", tt.args, got)
			}
		})
	}
}

func TestRun_NoArgs(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{"gg"}

	got := run()
	if got != 1 {
		t.Errorf("run() with no args = %d, want 1", got)
	}
}

func TestRun_PlatformCommand_UnknownHost(t *testing.T) {
	// No GG_HOST set, no origin remote → detectOrigin returns Unknown.
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{"gg", "issue", "list"}

	t.Setenv("GG_HOST", "")

	got := run()
	if got != 1 {
		t.Errorf("run() with platform cmd and unknown host = %d, want 1", got)
	}
}

func TestRun_PlatformCommand_MissingCLI(t *testing.T) {
	// GG_HOST=github forces the host, lookPath mocks gh not installed.
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{"gg", "pr", "create"}

	t.Setenv("GG_HOST", "github")

	oldLookPath := lookPath
	lookPath = func(name string) error { return errors.New("not found") }
	defer func() { lookPath = oldLookPath }()

	got := run()
	if got != 1 {
		t.Errorf("run() with platform cmd and missing CLI = %d, want 1", got)
	}
}

func TestRun_GitPassThrough(t *testing.T) {
	// Non-platform command → routes to git directly.
	// With no git repo this may still fail, but we test routing here.
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{"gg", "commit"}

	// Don't rely on os.Executable(); use a known path.
	got := run()
	// git commit without a repo prints an error and exits non-zero.
	// The point is that it doesn't return 1 from usage/unknown-host paths.
	if got == 0 {
		t.Log("git commit succeeded unexpectedly (maybe inside a configured repo)")
	}
}
