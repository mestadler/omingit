package main

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/mestadler/omingit/internal/classify"
	"github.com/mestadler/omingit/internal/gitea"
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

	// Scan for -C <path> before the first non-flag argument.
	var repoPath string
	cmdIdx := 1
	for cmdIdx < len(os.Args) {
		if os.Args[cmdIdx] == "-C" {
			if cmdIdx+1 >= len(os.Args) {
				fmt.Fprintf(os.Stderr, "gg: -C requires a path argument\n")
				return 1
			}
			repoPath = os.Args[cmdIdx+1]
			cmdIdx += 2
			continue
		}
		if !strings.HasPrefix(os.Args[cmdIdx], "-") {
			break
		}
		cmdIdx++
	}

	if cmdIdx >= len(os.Args) {
		fmt.Fprintf(os.Stderr, "Usage: gg <command> [args...]\n\n")
		fmt.Fprintf(os.Stderr, "gg routes commands to git, gh, or tea based on the remote's host.\n")
		fmt.Fprintf(os.Stderr, "See SKILL.md for full documentation.\n")
		return 1
	}

	cmd := os.Args[cmdIdx]
	args := os.Args[cmdIdx:] // pass through remaining args including the subcommand

	// Route platform-specific commands to gh or tea.
	if classify.IsPlatform(cmd) {
		h := detectOrigin(repoPath)
		cli := h.CLI()
		if cli == "" {
			fmt.Fprintf(os.Stderr, "gg: unknown host. Set GG_HOST=github or GG_HOST=gitea.\n")
			return 1
		}

		// Gitea-hosted repos: route workflow/run through local SDK.
		// All other commands continue to exec passthrough.
		if h == host.Gitea && (cmd == "workflow" || cmd == "run") {
			return runGiteaSDK(cmd, args)
		}

		if err := lookPath(cli); err != nil {
			fmt.Fprintf(os.Stderr, "gg: %s not found in PATH (required for %s operations)\n", cli, h)
			fmt.Fprintf(os.Stderr, "  Install: https://%s.com/cli\n", h)
			return 1
		}
		if repoPath != "" {
			return runCmd(cli, append([]string{"-C", repoPath}, args...)...)
		}
		return runCmd(cli, args...)
	}

	// Default: pass through to git.
	// If git fails with "not a git command" and a platform CLI is
	// available, retry with gh or tea.
	return runGitWithFallback(repoPath, args)
}

// runGitWithFallback tries git first. If git exits 1 with
// "is not a git command" and a platform CLI (gh/tea) is
// detected, it retries with that CLI.
func runGitWithFallback(repoPath string, args []string) int {
	gitArgs := args
	if repoPath != "" {
		gitArgs = append([]string{"-C", repoPath}, args...)
	}

	var stderr bytes.Buffer
	cmd := exec.Command("git", gitArgs...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err == nil {
		return 0
	}

	var exitErr *exec.ExitError
	if !errors.As(err, &exitErr) || exitErr.ExitCode() != 1 {
		os.Stderr.Write(stderr.Bytes())
		return 1
	}

	if !strings.Contains(stderr.String(), "is not a git command") {
		os.Stderr.Write(stderr.Bytes())
		return 1
	}

	// Git doesn't know this command. Try the platform CLI.
	h := detectOrigin(repoPath)
	cli := h.CLI()
	if cli == "" {
		os.Stderr.Write(stderr.Bytes())
		return 1
	}

	if err := lookPath(cli); err != nil {
		os.Stderr.Write(stderr.Bytes())
		return 1
	}

	if repoPath != "" {
		return runCmd(cli, append([]string{"-C", repoPath}, args...)...)
	}
	return runCmd(cli, args...)
}

// detectOrigin reads the origin remote URL to determine the host.
// When repoPath is non-empty, it runs git in that directory via -C.
func detectOrigin(repoPath string) host.Host {
	var out []byte
	var err error
	if repoPath != "" {
		out, err = exec.Command("git", "-C", repoPath, "remote", "get-url", "origin").Output()
	} else {
		out, err = exec.Command("git", "remote", "get-url", "origin").Output()
	}
	if err != nil {
		return host.Detect("")
	}
	return host.Detect(string(out))
}

// runGiteaSDK handles workflow and run subcommands via the local Gitea SDK
// instead of exec passthrough to tea.
func runGiteaSDK(cmd string, args []string) int {
	client, err := gitea.NewClientFromOrigin()
	if err != nil {
		fmt.Fprintf(os.Stderr, "gg: %v\n", err)
		return 1
	}

	sub := ""
	if len(args) > 1 {
		sub = args[1]
	}

	switch cmd {
	case "workflow":
		switch sub {
		case "list":
			workflows, err := client.ListWorkflows()
			if err != nil {
				fmt.Fprintf(os.Stderr, "gg: listing workflows: %v\n", err)
				return 1
			}
			for _, w := range workflows {
				fmt.Printf("%s  (state: %s)\n", w.Name, w.State)
			}
			return 0
		default:
			fmt.Fprintf(os.Stderr, "gg workflow: unknown subcommand %q\n", sub)
			return 1
		}
	case "run":
		switch sub {
		case "list":
			runs, err := client.ListRuns("")
			if err != nil {
				fmt.Fprintf(os.Stderr, "gg: listing runs: %v\n", err)
				return 1
			}
			for _, r := range runs {
				fmt.Printf("%-8d %-20s %-10s %s\n", r.ID, r.DisplayTitle, r.Status, r.Conclusion)
			}
			return 0
		default:
			fmt.Fprintf(os.Stderr, "gg run: unknown subcommand %q\n", sub)
			return 1
		}
	default:
		fmt.Fprintf(os.Stderr, "gg: internal error: runGiteaSDK called with unsupported command %q\n", cmd)
		return 1
	}
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
