package config

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

const (
	// EnvFileName is the name of the config file, both in the current folder and in the config folder.
	EnvFileName = ".env"

	// MinConfigVersion is the oldest tmdb version whose config files this build still understands.
	// Raise it when a release changes the config format in an incompatible way.
	MinConfigVersion = "0.7.0"

	// devVersion is used by builds without an injected version (go run, plain go build).
	devVersion = "dev"

	// headerScanLines is the number of leading lines that are searched for the version header.
	headerScanLines = 10
)

var (
	appVersion = devVersion

	// headerRe matches the header that marks a file as tmdb config: "# tmdb version 0.7.0".
	headerRe  = regexp.MustCompile(`^#\s*tmdb version\s+v?(\d+\.\d+\.\d+)\s*$`)
	versionRe = regexp.MustCompile(`^v?(\d+)\.(\d+)\.(\d+)`)
)

// SetAppVersion sets the version of the running tmdb binary (injected at build time).
// Empty values are ignored, which keeps the "dev" default.
func SetAppVersion(v string) {
	if v = strings.TrimSpace(v); v != "" {
		appVersion = v
	}
}

// AppVersion returns the version of the running tmdb binary, or "dev".
func AppVersion() string {
	return appVersion
}

type semver [3]int

func parseVersion(s string) (semver, bool) {
	m := versionRe.FindStringSubmatch(strings.TrimSpace(s))
	if m == nil {
		return semver{}, false
	}

	var v semver
	for i := range v {
		n, err := strconv.Atoi(m[i+1])
		if err != nil {
			return semver{}, false
		}
		v[i] = n
	}
	return v, true
}

func (v semver) less(o semver) bool {
	for i := range v {
		if v[i] != o[i] {
			return v[i] < o[i]
		}
	}
	return false
}

func (v semver) String() string {
	return fmt.Sprintf("%d.%d.%d", v[0], v[1], v[2])
}

// headerVersion is the version written into the header of new config files.
func headerVersion() string {
	if v, ok := parseVersion(appVersion); ok {
		return v.String()
	}
	return MinConfigVersion
}

// ParseHeader looks for the "# tmdb version x.y.z" header in the first lines of a file
// and returns the version. It reports false if the file has no such header, which means
// it is not a tmdb config file.
func ParseHeader(content string) (string, bool) {
	for i, line := range strings.Split(content, "\n") {
		if i >= headerScanLines {
			break
		}
		if m := headerRe.FindStringSubmatch(strings.TrimSpace(line)); m != nil {
			return m[1], true
		}
	}
	return "", false
}

// CheckCompatible reports whether a config file written by tmdb fileVersion can be read by
// the tmdb version current. Config files are compatible if they are not older than
// MinConfigVersion, have the same major version and are not newer than the running tmdb.
// Builds without a version ("dev") skip the checks against the running version.
func CheckCompatible(fileVersion, current string) error {
	fv, ok := parseVersion(fileVersion)
	if !ok {
		return fmt.Errorf("invalid version %q in the tmdb header", fileVersion)
	}

	oldest, _ := parseVersion(MinConfigVersion)
	if fv.less(oldest) {
		return fmt.Errorf("it was written by tmdb %s, but this tmdb needs a config from %s or newer", fv, oldest)
	}

	cv, ok := parseVersion(current)
	if !ok {
		return nil
	}

	if fv[0] != cv[0] {
		return fmt.Errorf("it was written by tmdb %s, which is not compatible with tmdb %s", fv, cv)
	}
	if cv.less(fv) {
		return fmt.Errorf("it was written by a newer tmdb (%s), this is tmdb %s", fv, cv)
	}
	return nil
}

// ConfigDir returns the folder of the user's tmdb config: $XDG_CONFIG_HOME/tmdb if
// XDG_CONFIG_HOME is set, otherwise ~/.config/tmdb.
func ConfigDir() (string, error) {
	if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
		return filepath.Join(xdg, "tmdb"), nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config", "tmdb"), nil
}

// ConfigFilePath returns the path of the user's tmdb config file.
func ConfigFilePath() (string, error) {
	dir, err := ConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, EnvFileName), nil
}

// ActiveFile returns the config file that is used right now, or "" if there is none.
// It only looks at the files and has no side effects, so it is safe for help output.
func ActiveFile() string {
	workDir, err := os.Getwd()
	if err != nil {
		workDir = "."
	}

	configDir, err := ConfigDir()
	if err != nil {
		configDir = ""
	}

	return Locate(workDir, configDir).Path
}

// Located is the result of looking for a config file.
type Located struct {
	// Path is the config file to load, empty if none was found.
	Path string
	// LocalPath and HomePath are the two places that were searched.
	LocalPath string
	HomePath  string
	// SkippedLocal describes why a .env in the current folder was not used
	// (not a tmdb file or incompatible). Empty if there is none or it was used.
	SkippedLocal string
	// HomeExists is true if there is a file at HomePath, even if it is not usable.
	HomeExists bool
	// Problems lists reasons why the file at HomePath is not usable.
	Problems []string
}

// Locate looks for the tmdb config file. A .env in workDir wins; it is only used if it carries
// the "# tmdb version x.y.z" header and is compatible, so foreign .env files are never read.
// Otherwise the .env in configDir is used. An empty configDir disables that location.
func Locate(workDir, configDir string) Located {
	res := Located{LocalPath: filepath.Join(workDir, EnvFileName)}
	if configDir != "" {
		res.HomePath = filepath.Join(configDir, EnvFileName)
	}

	if content, err := os.ReadFile(res.LocalPath); err == nil {
		if version, found := ParseHeader(string(content)); !found {
			res.SkippedLocal = fmt.Sprintf("%s has no \"# tmdb version x.y.z\" header and is not used as tmdb config", res.LocalPath)
		} else if err := CheckCompatible(version, AppVersion()); err != nil {
			res.SkippedLocal = fmt.Sprintf("%s is not used: %v", res.LocalPath, err)
		} else {
			res.Path = res.LocalPath
			return res
		}
	}

	if res.HomePath == "" {
		return res
	}

	content, err := os.ReadFile(res.HomePath)
	if errors.Is(err, fs.ErrNotExist) {
		return res
	}

	res.HomeExists = true
	if err != nil {
		res.Problems = append(res.Problems, fmt.Sprintf("cannot read %s: %v", res.HomePath, err))
		return res
	}

	version, found := ParseHeader(string(content))
	if !found {
		res.Problems = append(res.Problems, fmt.Sprintf("%s has no \"# tmdb version x.y.z\" header and is not used as tmdb config", res.HomePath))
		return res
	}
	if err := CheckCompatible(version, AppVersion()); err != nil {
		res.Problems = append(res.Problems, fmt.Sprintf("%s is not used: %v", res.HomePath, err))
		return res
	}

	res.Path = res.HomePath
	return res
}
