package cli

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/git-pkgs/forge/internal/config"
)

func TestAPIBaseURLKnownDomains(t *testing.T) {
	tests := []struct {
		name      string
		domain    string
		forgeType string
		want      string
	}{
		{
			name:   "github",
			domain: "github.com",
			want:   "https://api.github.com",
		},
		{
			name:   "gitlab",
			domain: "gitlab.com",
			want:   "https://gitlab.com/api/v4",
		},
		{
			name:   "codeberg",
			domain: "codeberg.org",
			want:   "https://codeberg.org/api/v1",
		},
		{
			name:   "bitbucket",
			domain: "bitbucket.org",
			want:   "https://api.bitbucket.org/2.0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := apiBaseURL(tt.domain, tt.forgeType)
			if got != tt.want {
				t.Fatalf("apiBaseURL(%q, %q) = %q, want %q", tt.domain, tt.forgeType, got, tt.want)
			}
		})
	}
}

func TestAPIBaseURLUnknownDomainWithForgeType(t *testing.T) {
	tests := []struct {
		name      string
		forgeType string
		want      string
	}{
		{
			name:      "github",
			forgeType: "github",
			want:      "https://api.mylab.example.org",
		},
		{
			name:      "gitlab",
			forgeType: "gitlab",
			want:      "https://mylab.example.org/api/v4",
		},
		{
			name:      "gitea",
			forgeType: "gitea",
			want:      "https://mylab.example.org/api/v1",
		},
		{
			name:      "forgejo",
			forgeType: "forgejo",
			want:      "https://mylab.example.org/api/v1",
		},
		{
			name:      "bitbucket",
			forgeType: "bitbucket",
			want:      "https://api.bitbucket.org/2.0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := apiBaseURL("mylab.example.org", tt.forgeType)
			if got != tt.want {
				t.Fatalf("apiBaseURL(%q, %q) = %q, want %q", "mylab.example.org", tt.forgeType, got, tt.want)
			}
		})
	}
}

func TestAPIBaseURLUnknownDomainWithConfigType(t *testing.T) {
	config.ResetCache()
	defer config.ResetCache()

	xdgConfigHome := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", xdgConfigHome)

	configDir := filepath.Join(xdgConfigHome, "forge")
	if err := os.MkdirAll(configDir, 0700); err != nil {
		t.Fatal(err)
	}
	configPath := filepath.Join(configDir, "config")
	if err := os.WriteFile(configPath, []byte(`[mylab.example.org]
type = gitlab
`), 0600); err != nil {
		t.Fatal(err)
	}

	got := apiBaseURL("mylab.example.org", "")
	want := "https://mylab.example.org/api/v4"
	if got != want {
		t.Fatalf("apiBaseURL(%q, %q) = %q, want %q", "mylab.example.org", "", got, want)
	}
}

func TestAPIBaseURLUnknownDomainWithoutForgeTypeOrConfig(t *testing.T) {
	config.ResetCache()
	defer config.ResetCache()

	t.Setenv("XDG_CONFIG_HOME", t.TempDir())

	got := apiBaseURL("mylab.example.org", "")
	want := "https://mylab.example.org/api/v1"
	if got != want {
		t.Fatalf("apiBaseURL(%q, %q) = %q, want %q", "mylab.example.org", "", got, want)
	}
}
