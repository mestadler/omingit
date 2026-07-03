package host

import (
	"testing"
)

func TestDetect_GGHostOverride(t *testing.T) {
	t.Setenv("GG_HOST", "github")
	if got := Detect("git@gitlab.com:user/repo.git"); got != GitHub {
		t.Errorf("with GG_HOST=github, Detect(...) = %v, want %v", got, GitHub)
	}

	t.Setenv("GG_HOST", "gitea")
	if got := Detect("git@github.com:user/repo.git"); got != Gitea {
		t.Errorf("with GG_HOST=gitea, Detect(...) = %v, want %v", got, Gitea)
	}
}

func TestDetect_GGHostCaseInsensitive(t *testing.T) {
	t.Setenv("GG_HOST", "GitHub")
	if got := Detect(""); got != GitHub {
		t.Errorf("with GG_HOST=GitHub, Detect(...) = %v, want %v", got, GitHub)
	}

	t.Setenv("GG_HOST", "GITEA")
	if got := Detect(""); got != Gitea {
		t.Errorf("with GG_HOST=GITEA, Detect(...) = %v, want %v", got, Gitea)
	}
}

func TestDetect_EmptyEnv(t *testing.T) {
	t.Setenv("GG_HOST", "")
	tests := []struct {
		name string
		url  string
		want Host
	}{
		{"GitHub SSH", "git@github.com:org/repo.git", GitHub},
		{"GitHub HTTPS", "https://github.com/org/repo.git", GitHub},
		{"GitHub with trailing newline", "git@github.com:org/repo.git\n", GitHub},
		{"Gitea SSH", "git@forgejo.example.com:org/repo.git", Gitea},
		{"Gitea HTTPS", "https://git.example.com/org/repo.git", Gitea},
		{"Gitea localhost", "git@git.internal:org/repo.git", Gitea},
		{"empty URL", "", Unknown},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Detect(tt.url); got != tt.want {
				t.Errorf("Detect(%q) = %v, want %v", tt.url, got, tt.want)
			}
		})
	}
}

func TestHost_CLI(t *testing.T) {
	tests := []struct {
		host Host
		want string
	}{
		{GitHub, "gh"},
		{Gitea, "tea"},
		{Unknown, ""},
	}
	for _, tt := range tests {
		if got := tt.host.CLI(); got != tt.want {
			t.Errorf("%v.CLI() = %q, want %q", tt.host, got, tt.want)
		}
	}
}

func TestHost_String(t *testing.T) {
	tests := []struct {
		host Host
		want string
	}{
		{GitHub, "github"},
		{Gitea, "gitea"},
		{Unknown, "unknown"},
	}
	for _, tt := range tests {
		if got := tt.host.String(); got != tt.want {
			t.Errorf("%v.String() = %q, want %q", tt.host, got, tt.want)
		}
	}
}
