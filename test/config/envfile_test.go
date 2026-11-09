package config_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sebastianneubert/tmdb/internal/config"
)

// withAppVersion sets the app version for one test and restores it afterwards.
func withAppVersion(t *testing.T, version string) {
	t.Helper()
	old := config.AppVersion()
	config.SetAppVersion(version)
	t.Cleanup(func() { config.SetAppVersion(old) })
}

func writeFile(t *testing.T, dir, content string) string {
	t.Helper()
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(dir, config.EnvFileName)
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestParseHeader(t *testing.T) {
	tests := []struct {
		name        string
		content     string
		wantVersion string
		wantFound   bool
	}{
		{"plain header", "# tmdb version 0.7.0\nTMDB_API_KEY=abc\n", "0.7.0", true},
		{"v prefix", "# tmdb version v1.2.3\n", "1.2.3", true},
		{"no space after hash", "#tmdb version 0.7.0\n", "0.7.0", true},
		{"windows line endings", "# tmdb version 0.7.0\r\nTMDB_API_KEY=abc\r\n", "0.7.0", true},
		{"header after other comments", "# my notes\n\n# tmdb version 0.8.1\nREGION=DE\n", "0.8.1", true},
		{"no header", "TMDB_API_KEY=abc\nREGION=DE\n", "", false},
		{"other tool's version header", "# some other tool version 1.0.0\n", "", false},
		{"not a full version", "# tmdb version 1.0\n", "", false},
		{"header too far down", strings.Repeat("# note\n", 12) + "# tmdb version 0.7.0\n", "", false},
		{"empty", "", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			version, found := config.ParseHeader(tt.content)
			if found != tt.wantFound || version != tt.wantVersion {
				t.Errorf("ParseHeader() = (%q, %v), want (%q, %v)", version, found, tt.wantVersion, tt.wantFound)
			}
		})
	}
}

func TestCheckCompatible(t *testing.T) {
	// MinConfigVersion is 0.7.0
	tests := []struct {
		name    string
		file    string
		current string
		wantErr string // empty = compatible
	}{
		{"same version", "0.7.0", "0.7.0", ""},
		{"older patch", "0.7.0", "0.7.3", ""},
		{"older minor, still supported", "0.7.0", "0.9.0", ""},
		{"release tag with suffix", "0.7.0", "0.7.1-rc1", ""},
		{"v prefix in running version", "0.7.0", "v0.7.0", ""},
		{"newer than running tmdb", "0.8.0", "0.7.0", "newer tmdb"},
		{"older than supported", "0.6.0", "0.7.0", "needs a config from"},
		{"different major", "1.0.0", "2.0.0", "not compatible"},
		{"dev build skips running version check", "0.9.0", "dev", ""},
		{"dev build still checks minimum", "0.5.0", "dev", "needs a config from"},
		{"invalid file version", "abc", "0.7.0", "invalid version"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := config.CheckCompatible(tt.file, tt.current)
			switch {
			case tt.wantErr == "" && err != nil:
				t.Errorf("unexpected error: %v", err)
			case tt.wantErr != "" && err == nil:
				t.Errorf("expected an error containing %q", tt.wantErr)
			case tt.wantErr != "" && !strings.Contains(err.Error(), tt.wantErr):
				t.Errorf("error %q should contain %q", err, tt.wantErr)
			}
		})
	}
}

func TestConfigDir(t *testing.T) {
	t.Run("uses XDG_CONFIG_HOME", func(t *testing.T) {
		xdg := t.TempDir()
		t.Setenv("XDG_CONFIG_HOME", xdg)

		dir, err := config.ConfigDir()
		if err != nil || dir != filepath.Join(xdg, "tmdb") {
			t.Errorf("ConfigDir() = (%q, %v), want %q", dir, err, filepath.Join(xdg, "tmdb"))
		}
	})

	t.Run("falls back to ~/.config/tmdb", func(t *testing.T) {
		home := t.TempDir()
		t.Setenv("XDG_CONFIG_HOME", "")
		t.Setenv("HOME", home)
		t.Setenv("USERPROFILE", home)

		path, err := config.ConfigFilePath()
		want := filepath.Join(home, ".config", "tmdb", ".env")
		if err != nil || path != want {
			t.Errorf("ConfigFilePath() = (%q, %v), want %q", path, err, want)
		}
	})
}

func TestLocate(t *testing.T) {
	withAppVersion(t, "0.7.2")

	valid := "# tmdb version 0.7.0\nTMDB_API_KEY=abc\n"
	foreign := "DATABASE_URL=postgres://localhost/app\n"
	tooNew := "# tmdb version 0.9.0\nTMDB_API_KEY=abc\n"
	tooOld := "# tmdb version 0.5.0\nTMDB_API_KEY=abc\n"

	tests := []struct {
		name         string
		local        string // content of ./.env, "" = no file
		home         string // content of the config folder's .env, "" = no file
		wantSource   string // "local", "home" or "" for no usable file
		wantSkipped  bool
		wantHome     bool
		wantProblems int
	}{
		{name: "local tmdb file wins", local: valid, home: valid, wantSource: "local", wantHome: false},
		{name: "local tmdb file without home file", local: valid, wantSource: "local"},
		{name: "foreign local file is skipped, home is used", local: foreign, home: valid, wantSource: "home", wantSkipped: true, wantHome: true},
		{name: "too new local file is skipped, home is used", local: tooNew, home: valid, wantSource: "home", wantSkipped: true, wantHome: true},
		{name: "home file only", home: valid, wantSource: "home", wantHome: true},
		{name: "nothing found", wantSource: ""},
		{name: "only foreign local file", local: foreign, wantSource: "", wantSkipped: true},
		{name: "foreign home file is not used", home: foreign, wantSource: "", wantHome: true, wantProblems: 1},
		{name: "too new home file is not used", home: tooNew, wantSource: "", wantHome: true, wantProblems: 1},
		{name: "too old home file is not used", home: tooOld, wantSource: "", wantHome: true, wantProblems: 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			workDir := t.TempDir()
			configDir := filepath.Join(t.TempDir(), "tmdb")

			var localPath, homePath string
			if tt.local != "" {
				localPath = writeFile(t, workDir, tt.local)
			}
			if tt.home != "" {
				homePath = writeFile(t, configDir, tt.home)
			}

			res := config.Locate(workDir, configDir)

			wantPath := ""
			switch tt.wantSource {
			case "local":
				wantPath = localPath
			case "home":
				wantPath = homePath
			}
			if res.Path != wantPath {
				t.Errorf("Path = %q, want %q", res.Path, wantPath)
			}
			if (res.SkippedLocal != "") != tt.wantSkipped {
				t.Errorf("SkippedLocal = %q, want skipped=%v", res.SkippedLocal, tt.wantSkipped)
			}
			if tt.wantSource == "" && res.HomeExists != tt.wantHome {
				t.Errorf("HomeExists = %v, want %v", res.HomeExists, tt.wantHome)
			}
			if len(res.Problems) != tt.wantProblems {
				t.Errorf("Problems = %v, want %d", res.Problems, tt.wantProblems)
			}
		})
	}
}

func TestLocateWithoutConfigDir(t *testing.T) {
	workDir := t.TempDir()
	writeFile(t, workDir, "# tmdb version 0.7.0\n")

	res := config.Locate(workDir, "")
	if res.Path == "" {
		t.Error("expected the local file to be used without a config dir")
	}

	res = config.Locate(t.TempDir(), "")
	if res.Path != "" || res.HomePath != "" || res.HomeExists {
		t.Errorf("expected nothing found, got %+v", res)
	}
}
