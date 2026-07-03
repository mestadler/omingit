// Package host detects the Git hosting platform from the origin remote URL
// or the GG_HOST environment variable.
package host

import (
	"os"
	"strings"
)

// Host represents a supported Git hosting platform.
type Host int

const (
	Unknown Host = iota
	GitHub
	Gitea
)

func (h Host) String() string {
	switch h {
	case GitHub:
		return "github"
	case Gitea:
		return "gitea"
	default:
		return "unknown"
	}
}

// CLI returns the platform-specific CLI binary name for the host.
// Returns empty string for unknown hosts.
func (h Host) CLI() string {
	switch h {
	case GitHub:
		return "gh"
	case Gitea:
		return "tea"
	default:
		return ""
	}
}

// Detect determines the hosting platform.
// Priority:
//  1. GG_HOST environment variable override
//  2. Origin URL containing "github.com" → GitHub
//  3. Any other non-empty origin URL → Gitea (covers self-hosted Forgejo/Gitea)
//  4. Empty origin URL → Unknown
func Detect(originURL string) Host {
	switch strings.ToLower(os.Getenv("GG_HOST")) {
	case "github":
		return GitHub
	case "gitea":
		return Gitea
	}

	u := strings.ToLower(strings.TrimSpace(originURL))
	if u == "" {
		return Unknown
	}

	if strings.Contains(u, "github.com") {
		return GitHub
	}

	return Gitea
}
