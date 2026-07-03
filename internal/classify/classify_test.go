package classify

import (
	"testing"
)

func TestIsPlatform(t *testing.T) {
	platformCommands := []string{
		"pr",
		"issue",
		"repo",
		"workflow",
		"run",
		"release",
		"gist",
		"label",
		"project",
		"variable",
		"secret",
	}

	gitCommands := []string{
		"commit",
		"push",
		"pull",
		"fetch",
		"clone",
		"init",
		"add",
		"checkout",
		"branch",
		"tag",
		"log",
		"status",
		"diff",
		"merge",
		"rebase",
		"reset",
		"revert",
		"stash",
		"config",
		"remote",
		"help",
		"version",
		"",
		"unknown-command",
	}

	for _, cmd := range platformCommands {
		t.Run("platform_"+cmd, func(t *testing.T) {
			if !IsPlatform(cmd) {
				t.Errorf("IsPlatform(%q) = false, want true", cmd)
			}
		})
	}

	for _, cmd := range gitCommands {
		t.Run("git_"+cmd, func(t *testing.T) {
			if IsPlatform(cmd) {
				t.Errorf("IsPlatform(%q) = true, want false", cmd)
			}
		})
	}
}
