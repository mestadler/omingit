package gitea

import (
	"testing"
)

func TestNewClient_EmptyToken(t *testing.T) {
	_, err := NewClient("https://git.example.com", "", "owner", "repo")
	if err == nil {
		t.Error("expected error for empty token, got nil")
	}
}

func TestRepoOwner(t *testing.T) {
	tests := []struct {
		name  string
		owner string
		repo  string
		want  string
	}{
		{"simple", "alice", "project", "alice/project"},
		{"with dots", "user.name", "my.repo", "user.name/my.repo"},
		{"with dashes", "org-name", "repo-name", "org-name/repo-name"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{owner: tt.owner, repo: tt.repo}
			if got := c.RepoOwner(); got != tt.want {
				t.Errorf("RepoOwner() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestParseOriginURL_HTTPS(t *testing.T) {
	tests := []struct {
		name      string
		url       string
		wantBase  string
		wantOwner string
		wantRepo  string
	}{
		{
			name:      "standard HTTPS with .git",
			url:       "https://github.com/owner/repo.git",
			wantBase:  "https://github.com",
			wantOwner: "owner",
			wantRepo:  "repo",
		},
		{
			name:      "HTTPS without .git",
			url:       "https://git.example.com/org/project",
			wantBase:  "https://git.example.com",
			wantOwner: "org",
			wantRepo:  "project",
		},
		{
			name:      "HTTPS with trailing newline",
			url:       "https://git.example.com/org/repo.git\n",
			wantBase:  "https://git.example.com",
			wantOwner: "org",
			wantRepo:  "repo",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			base, owner, repo, err := parseOriginURL(tt.url)
			if err != nil {
				t.Fatalf("parseOriginURL(%q) unexpected error: %v", tt.url, err)
			}
			if base != tt.wantBase {
				t.Errorf("base = %q, want %q", base, tt.wantBase)
			}
			if owner != tt.wantOwner {
				t.Errorf("owner = %q, want %q", owner, tt.wantOwner)
			}
			if repo != tt.wantRepo {
				t.Errorf("repo = %q, want %q", repo, tt.wantRepo)
			}
		})
	}
}

func TestParseOriginURL_GitSSH(t *testing.T) {
	tests := []struct {
		name      string
		url       string
		wantBase  string
		wantOwner string
		wantRepo  string
	}{
		{
			name:      "standard SSH",
			url:       "git@github.com:owner/repo.git",
			wantBase:  "https://github.com",
			wantOwner: "owner",
			wantRepo:  "repo",
		},
		{
			name:      "SSH without .git",
			url:       "git@git.example.com:org/project",
			wantBase:  "https://git.example.com",
			wantOwner: "org",
			wantRepo:  "project",
		},
		{
			name:      "SSH with trailing newline",
			url:       "git@forgejo.example.com:user/repo.git\n",
			wantBase:  "https://forgejo.example.com",
			wantOwner: "user",
			wantRepo:  "repo",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			base, owner, repo, err := parseOriginURL(tt.url)
			if err != nil {
				t.Fatalf("parseOriginURL(%q) unexpected error: %v", tt.url, err)
			}
			if base != tt.wantBase {
				t.Errorf("base = %q, want %q", base, tt.wantBase)
			}
			if owner != tt.wantOwner {
				t.Errorf("owner = %q, want %q", owner, tt.wantOwner)
			}
			if repo != tt.wantRepo {
				t.Errorf("repo = %q, want %q", repo, tt.wantRepo)
			}
		})
	}
}

func TestParseOriginURL_Invalid(t *testing.T) {
	tests := []struct {
		name string
		url  string
	}{
		{"empty", ""},
		{"SSH missing colon", "git@host"},
		{"SSH colon but no path", "git@host:"},
		{"SSH single path component", "git@host:only-owner"},
		{"HTTPS no path", "https://host"},
		{"HTTPS host component only", "https://host/"},
		{"HTTPS single path component", "https://host/only-owner"},
		{"no scheme", "not-a-url"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, _, err := parseOriginURL(tt.url)
			if err == nil {
				t.Errorf("parseOriginURL(%q) expected error, got nil", tt.url)
			}
		})
	}
}
