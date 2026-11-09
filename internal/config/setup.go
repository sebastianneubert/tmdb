package config

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/sebastianneubert/tmdb"
)

const (
	apiKeyName        = "TMDB_API_KEY"
	apiKeyPlaceholder = "your_api_key_here"
)

var (
	// settingRe matches an active setting line of the dist file: KEY=value
	settingRe = regexp.MustCompile(`^([A-Za-z_][A-Za-z0-9_]*)=(.*)$`)
	// commentedSettingRe matches an optional, commented-out setting: # KEY=value
	commentedSettingRe = regexp.MustCompile(`^#\s*[A-Z][A-Z0-9_]*=`)
)

// ErrSetupAborted is returned when the input ends before the setup is finished.
var ErrSetupAborted = errors.New("setup aborted")

// IsInteractive reports whether tmdb runs in a terminal, i.e. whether it may ask questions.
func IsInteractive() bool {
	for _, f := range []*os.File{os.Stdin, os.Stdout, os.Stderr} {
		info, err := f.Stat()
		if err != nil || info.Mode()&os.ModeCharDevice == 0 {
			return false
		}
	}
	return true
}

// readLine prints a prompt and reads one trimmed line. It returns ErrSetupAborted at end of input.
func readLine(in *bufio.Reader, out io.Writer, prompt string) (string, error) {
	fmt.Fprint(out, prompt)

	line, err := in.ReadString('\n')
	if err != nil && line == "" {
		fmt.Fprintln(out)
		return "", ErrSetupAborted
	}
	return strings.TrimSpace(line), nil
}

// Confirm asks a yes/no question. An empty answer selects the default.
func Confirm(in *bufio.Reader, out io.Writer, question string, defaultYes bool) (bool, error) {
	hint := "[y/N]"
	if defaultYes {
		hint = "[Y/n]"
	}

	for {
		answer, err := readLine(in, out, fmt.Sprintf("%s %s ", question, hint))
		if err != nil {
			return false, err
		}

		switch strings.ToLower(answer) {
		case "":
			return defaultYes, nil
		case "y", "yes", "j", "ja":
			return true, nil
		case "n", "no", "nein":
			return false, nil
		}
		fmt.Fprintln(out, "Please answer with y or n.")
	}
}

// validateLike checks that value has the same kind (integer, number or text) as the default.
func validateLike(def, value string) error {
	if strings.ContainsAny(value, "\"\n\r") {
		return errors.New("the value must not contain double quotes")
	}

	if _, err := strconv.Atoi(def); err == nil {
		if _, err := strconv.Atoi(value); err != nil {
			return errors.New("please enter a whole number")
		}
	} else if _, err := strconv.ParseFloat(def, 64); err == nil {
		if _, err := strconv.ParseFloat(value, 64); err != nil {
			return errors.New("please enter a number, e.g. 7.5")
		}
	}
	return nil
}

func formatValue(value string, quoted bool) string {
	if quoted || strings.ContainsAny(value, " \t#'") {
		return `"` + value + `"`
	}
	return value
}

// Setup creates the config file at path interactively. It walks through the active settings of
// the template (the dist file .env.example), asks for each value with the template value as
// default, and keeps comments and optional (commented-out) settings as they are. The file gets
// the "# tmdb version x.y.z" header and is only readable by the user, as it contains the API key.
func Setup(in *bufio.Reader, out io.Writer, template, path string) error {
	fmt.Fprintf(out, "Creating the tmdb config %s\n", path)
	fmt.Fprintln(out, "Press Enter to accept the value in brackets.")
	fmt.Fprintln(out)

	var content strings.Builder
	fmt.Fprintf(&content, "# tmdb version %s\n", headerVersion())
	content.WriteString("# Created by tmdb. Edit the values or run `tmdb init` to create this file again.\n")

	var help []string
	for _, line := range strings.Split(strings.TrimRight(template, "\r\n"), "\n") {
		line = strings.TrimRight(line, "\r")

		if headerRe.MatchString(strings.TrimSpace(line)) {
			continue
		}

		m := settingRe.FindStringSubmatch(line)
		if m == nil {
			content.WriteString(line + "\n")
			switch {
			case strings.HasPrefix(line, "#") && !commentedSettingRe.MatchString(line):
				help = append(help, strings.TrimSpace(strings.TrimPrefix(line, "#")))
			default:
				help = nil
			}
			continue
		}

		key, def := m[1], strings.TrimSpace(m[2])
		quoted := len(def) >= 2 && strings.HasPrefix(def, `"`) && strings.HasSuffix(def, `"`)
		if quoted {
			def = def[1 : len(def)-1]
		}

		required := key == apiKeyName
		if required && def == apiKeyPlaceholder {
			def = ""
		}

		value, err := askSetting(in, out, key, def, help, required)
		if err != nil {
			return err
		}
		help = nil

		content.WriteString(key + "=" + formatValue(value, quoted) + "\n")
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return fmt.Errorf("cannot create %s: %w", filepath.Dir(path), err)
	}
	if err := os.WriteFile(path, []byte(content.String()), 0o600); err != nil {
		return fmt.Errorf("cannot write %s: %w", path, err)
	}
	if err := os.Chmod(path, 0o600); err != nil {
		return fmt.Errorf("cannot protect %s: %w", path, err)
	}

	fmt.Fprintf(out, "\nSaved %s\n", path)
	return nil
}

// askSetting asks for the value of one setting until it is valid.
func askSetting(in *bufio.Reader, out io.Writer, key, def string, help []string, required bool) (string, error) {
	for _, h := range help {
		if h != "" {
			fmt.Fprintf(out, "%s\n", h)
		}
	}
	if key == apiKeyName {
		fmt.Fprintln(out, "Get an API key at https://www.themoviedb.org/settings/api")
	}

	prompt := key + ": "
	if def != "" {
		prompt = fmt.Sprintf("%s [%s]: ", key, def)
	}

	for {
		answer, err := readLine(in, out, prompt)
		if err != nil {
			return "", err
		}

		if answer == "" {
			answer = def
		}
		if answer == "" && required {
			fmt.Fprintf(out, "%s is required.\n", key)
			continue
		}
		if answer != "" {
			if err := validateLike(def, answer); err != nil {
				fmt.Fprintln(out, err)
				continue
			}
		}

		fmt.Fprintln(out)
		return answer, nil
	}
}

// SetupFromDist runs Setup with the dist file .env.example that is embedded in the binary.
func SetupFromDist(in *bufio.Reader, out io.Writer, path string) error {
	return Setup(in, out, tmdb.EnvExample, path)
}
