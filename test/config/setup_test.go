package config_test

import (
	"bufio"
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/sebastianneubert/tmdb"
	"github.com/sebastianneubert/tmdb/internal/config"
)

const testTemplate = `TMDB_API_KEY=your_api_key_here
# use the tmdb providers command to get a list of available providers
PROVIDERS="Netflix,Disney Plus"
REGION=DE
MIN_RATING=7.5
API_TIMEOUT_SECONDS=20
# hide old movies, 0 = off
# MAX_AGE=12
# DEBUG=true
`

func runSetup(t *testing.T, template, input string) (path string, output string, err error) {
	t.Helper()

	path = filepath.Join(t.TempDir(), "config", "tmdb", ".env")
	var out bytes.Buffer
	err = config.Setup(bufio.NewReader(strings.NewReader(input)), &out, template, path)
	return path, out.String(), err
}

func TestSetupWritesConfigWithHeader(t *testing.T) {
	withAppVersion(t, "0.7.4")

	// API key, then Enter for PROVIDERS, a new REGION, Enter for the rest
	path, output, err := runSetup(t, testTemplate, "my-secret-key\n\nAT\n\n\n")
	if err != nil {
		t.Fatalf("Setup failed: %v\noutput:\n%s", err, output)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	content := string(data)

	if !strings.HasPrefix(content, "# tmdb version 0.7.4\n") {
		t.Errorf("file must start with the version header, got:\n%s", content)
	}
	for _, want := range []string{
		"TMDB_API_KEY=my-secret-key\n",
		"PROVIDERS=\"Netflix,Disney Plus\"\n", // default kept, quoting kept
		"REGION=AT\n",
		"MIN_RATING=7.5\n",
		"API_TIMEOUT_SECONDS=20\n",
		"# use the tmdb providers command to get a list of available providers\n", // comments kept
		"# MAX_AGE=12\n", // optional settings stay commented out
		"# DEBUG=true\n",
	} {
		if !strings.Contains(content, want) {
			t.Errorf("file should contain %q, got:\n%s", want, content)
		}
	}

	// the generated file has to be accepted by the checks of this version
	version, found := config.ParseHeader(content)
	if !found {
		t.Fatal("generated file has no header")
	}
	if err := config.CheckCompatible(version, config.AppVersion()); err != nil {
		t.Errorf("generated file is not compatible: %v", err)
	}
	if res := config.Locate(t.TempDir(), filepath.Dir(path)); res.Path != path {
		t.Errorf("Locate() should find the generated file, got %+v", res)
	}

	// the user is told where the help text for a setting comes from
	if !strings.Contains(output, "use the tmdb providers command") {
		t.Errorf("prompt should show the comment above PROVIDERS, got:\n%s", output)
	}
	if !strings.Contains(output, "PROVIDERS [Netflix,Disney Plus]: ") {
		t.Errorf("prompt should show the default, got:\n%s", output)
	}
}

func TestSetupFilePermissions(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("unix permissions")
	}

	path, _, err := runSetup(t, testTemplate, "key\n\n\n\n\n")
	if err != nil {
		t.Fatal(err)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if perm := info.Mode().Perm(); perm != 0o600 {
		t.Errorf("config file permissions = %o, want 600 (it contains the API key)", perm)
	}
	dirInfo, err := os.Stat(filepath.Dir(path))
	if err != nil {
		t.Fatal(err)
	}
	if perm := dirInfo.Mode().Perm(); perm != 0o700 {
		t.Errorf("config folder permissions = %o, want 700", perm)
	}
}

func TestSetupAsksAgainForInvalidInput(t *testing.T) {
	// empty API key (required) -> asked again; MIN_RATING "abc" is no number -> asked again
	input := "\nmy-key\n\n\nabc\n8.1\n\n"
	path, output, err := runSetup(t, testTemplate, input)
	if err != nil {
		t.Fatalf("Setup failed: %v\noutput:\n%s", err, output)
	}

	data, _ := os.ReadFile(path)
	content := string(data)
	if !strings.Contains(content, "TMDB_API_KEY=my-key\n") {
		t.Errorf("API key missing in:\n%s", content)
	}
	if !strings.Contains(content, "MIN_RATING=8.1\n") {
		t.Errorf("MIN_RATING should be 8.1 in:\n%s", content)
	}
	if !strings.Contains(output, "TMDB_API_KEY is required.") {
		t.Errorf("empty API key should be rejected, output:\n%s", output)
	}
	if !strings.Contains(output, "please enter a number") {
		t.Errorf("invalid number should be rejected, output:\n%s", output)
	}
}

func TestSetupQuotesValuesWithSpaces(t *testing.T) {
	path, _, err := runSetup(t, testTemplate, "key\n\nDE DE\n\n\n")
	if err != nil {
		t.Fatal(err)
	}

	data, _ := os.ReadFile(path)
	if !strings.Contains(string(data), "REGION=\"DE DE\"\n") {
		t.Errorf("values with spaces must be quoted, got:\n%s", data)
	}
}

func TestSetupRejectsQuotesInValues(t *testing.T) {
	_, output, err := runSetup(t, testTemplate, "key\n\"bad\"\nGood,Value\n\n\n\n")
	if err != nil {
		t.Fatalf("Setup failed: %v", err)
	}
	if !strings.Contains(output, "must not contain double quotes") {
		t.Errorf("quotes should be rejected, output:\n%s", output)
	}
}

func TestSetupAbortsAtEndOfInput(t *testing.T) {
	path, _, err := runSetup(t, testTemplate, "my-key\n")
	if !errors.Is(err, config.ErrSetupAborted) {
		t.Fatalf("err = %v, want ErrSetupAborted", err)
	}
	if _, statErr := os.Stat(path); !os.IsNotExist(statErr) {
		t.Error("an aborted setup must not leave a config file behind")
	}
}

func TestSetupReplacesHeaderOfTemplate(t *testing.T) {
	template := "# tmdb version 0.0.1\n" + testTemplate
	path, _, err := runSetup(t, template, "key\n\n\n\n\n")
	if err != nil {
		t.Fatal(err)
	}

	data, _ := os.ReadFile(path)
	if n := strings.Count(string(data), "# tmdb version"); n != 1 {
		t.Errorf("expected exactly one version header, found %d in:\n%s", n, data)
	}
}

func TestSetupOverwritesExistingFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "tmdb", ".env")
	writeFile(t, filepath.Dir(path), "# tmdb version 0.7.0\nTMDB_API_KEY=old\n")

	var out bytes.Buffer
	in := bufio.NewReader(strings.NewReader("new-key\n\n\n\n\n"))
	if err := config.Setup(in, &out, testTemplate, path); err != nil {
		t.Fatal(err)
	}

	data, _ := os.ReadFile(path)
	if strings.Contains(string(data), "TMDB_API_KEY=old") || !strings.Contains(string(data), "TMDB_API_KEY=new-key") {
		t.Errorf("file should be replaced, got:\n%s", data)
	}
}

// The embedded dist file has to work with the setup, otherwise `tmdb init` is broken.
func TestEmbeddedDistFile(t *testing.T) {
	if !strings.Contains(tmdb.EnvExample, "TMDB_API_KEY=") {
		t.Fatal("the embedded .env.example must contain TMDB_API_KEY")
	}

	path := filepath.Join(t.TempDir(), "tmdb", ".env")
	var out bytes.Buffer
	// the API key, then Enter for every other setting
	in := bufio.NewReader(strings.NewReader("my-key\n" + strings.Repeat("\n", 30)))
	if err := config.SetupFromDist(in, &out, path); err != nil {
		t.Fatalf("SetupFromDist failed: %v\noutput:\n%s", err, out.String())
	}

	data, _ := os.ReadFile(path)
	content := string(data)
	if !strings.HasPrefix(content, "# tmdb version ") {
		t.Errorf("missing header in:\n%s", content)
	}
	if !strings.Contains(content, "TMDB_API_KEY=my-key\n") || !strings.Contains(content, "REGION=") {
		t.Errorf("unexpected content:\n%s", content)
	}
	if strings.Contains(content, "your_api_key_here") {
		t.Errorf("the placeholder API key must not end up in the config:\n%s", content)
	}
}

func TestConfirm(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		defaultYes bool
		want       bool
	}{
		{"empty uses default yes", "\n", true, true},
		{"empty uses default no", "\n", false, false},
		{"yes", "y\n", false, true},
		{"YES", "YES\n", false, true},
		{"ja", "ja\n", false, true},
		{"no", "n\n", true, false},
		{"nein", "nein\n", true, false},
		{"invalid answer is asked again", "maybe\ny\n", false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var out bytes.Buffer
			got, err := config.Confirm(bufio.NewReader(strings.NewReader(tt.input)), &out, "Create?", tt.defaultYes)
			if err != nil || got != tt.want {
				t.Errorf("Confirm() = (%v, %v), want (%v, nil)", got, err, tt.want)
			}
		})
	}

	t.Run("end of input aborts", func(t *testing.T) {
		var out bytes.Buffer
		_, err := config.Confirm(bufio.NewReader(strings.NewReader("")), &out, "Create?", true)
		if !errors.Is(err, config.ErrSetupAborted) {
			t.Errorf("err = %v, want ErrSetupAborted", err)
		}
	})
}
