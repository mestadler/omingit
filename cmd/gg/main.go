package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/mestadler/omingit/internal/classify"
	"github.com/mestadler/omingit/internal/host"
)

var version = "dev"

func main() {
	os.Exit(run())
}

func run() int {
	if len(os.Args) > 1 {
		if os.Args[1] == "--version" || os.Args[1] == "-v" {
			fmt.Println(version)
			return 0
		}
	}

	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: gg <command> [args...]\n\n")
		fmt.Fprintf(os.Stderr, "gg routes commands to git, gh, or tea based on the remote's host.\n")
		fmt.Fprintf(os.Stderr, "See SKILL.md for full documentation.\n")
		return 1
	}

	cmd := os.Args[1]
	args := os.Args[1:] // pass through all args including the subcommand

	// Route platform-specific commands to gh or tea.
	if classify.IsPlatform(cmd) {
		h := detectOrigin()
		cli := h.CLI()
		if cli == "" {
			fmt.Fprintf(os.Stderr, "gg: unknown host. Set GG_HOST=github or GG_HOST=gitea.\n")
			return 1
		}
		if err := lookPath(cli); err != nil {
			fmt.Fprintf(os.Stderr, "gg: %s not found in PATH (required for %s operations)\n", cli, h)
			fmt.Fprintf(os.Stderr, "  Install: https://%s.com/cli\n", h)
			return 1
		}
		return runCmd(cli, args...)
	}

	// Default: pass through to git.
	return runCmd("git", args...)
}

// detectOrigin reads the origin remote URL to determine the host.
func detectOrigin() host.Host {
	out, err := exec.Command("git", "remote", "get-url", "origin").Output()
	if err != nil {
		return host.Detect("")
	}
	return host.Detect(string(out))
}

// runCmd executes the given command with args, connecting stdio directly.
// Returns the child process exit code.
func runCmd(name string, args ...string) int {
	cmd := exec.Command(name, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return exitErr.ExitCode()
		}
		return 1
	}
	return 0
}

// lookPath wraps exec.LookPath for clarity and testability.
var lookPath = func(name string) error {
	_, err := exec.LookPath(name)
	return err
}
