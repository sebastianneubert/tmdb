package config

import (
	"bufio"
	"fmt"
	"os"

	"github.com/spf13/viper"
)

const (
	DefaultRegion       = "DE"
	DefaultProviders    = "Netflix,DisneyPlus,Wow,RtlPlus,AmazonPrime"
	DefaultTimeout      = 20
	DefaultMinRating    = 7.5
	DefaultMinVotes     = 1000
	DefaultMaxAge       = 0 // 0 = age filter off
	DefaultDebug        = false
	MaxPagesToSearch    = 5
	MaxResultsToDisplay = 40
)

type Config struct {
	APIKey    string  `mapstructure:"TMDB_API_KEY"`
	Region    string  `mapstructure:"REGION"`
	Providers string  `mapstructure:"PROVIDERS"`
	MinRating float64 `mapstructure:"MIN_RATING"`
	MinVotes  int     `mapstructure:"MIN_VOTES"`
	Timeout   int     `mapstructure:"API_TIMEOUT_SECONDS"`
	DEBUG     bool    `mapstructure:"DEBUG"`

	// MaxAge hides movies with an age recommendation above this age (0 = off).
	MaxAge int `mapstructure:"MAX_AGE"`
	// IncludeUnrated also shows movies without an age recommendation when MaxAge is set.
	IncludeUnrated bool `mapstructure:"INCLUDE_UNRATED"`
}

var AppConfig Config

// InitOptions controls how Init treats the config file.
type InitOptions struct {
	// Quiet suppresses all notices. Used for commands like completion and help, where
	// extra output would get in the way of the shell.
	Quiet bool
	// AllowSetup lets Init offer to create the config file interactively if there is none.
	// It only asks if tmdb runs in a terminal.
	AllowSetup bool
}

// Init loads the config without asking any questions.
func Init() {
	InitWith(InitOptions{})
}

// InitWith loads the config. The config file is searched in this order:
//  1. .env in the current folder
//  2. .env in the config folder (~/.config/tmdb, see ConfigDir)
//
// A file is only used if it starts with the "# tmdb version x.y.z" header and is compatible
// with this tmdb version. Environment variables always override values of the file.
func InitWith(opts InitOptions) {
	viper.AutomaticEnv()

	viper.SetDefault("REGION", DefaultRegion)
	viper.SetDefault("PROVIDERS", DefaultProviders)
	viper.SetDefault("MIN_RATING", DefaultMinRating)
	viper.SetDefault("MIN_VOTES", DefaultMinVotes)
	viper.SetDefault("API_TIMEOUT_SECONDS", DefaultTimeout)
	viper.SetDefault("MAX_AGE", DefaultMaxAge)
	viper.SetDefault("INCLUDE_UNRATED", false)
	viper.SetDefault("DEBUG", DefaultDebug)
	viper.SetDefault("TMDB_API_KEY", "")

	path := resolveConfigFile(opts)
	if path != "" {
		viper.SetConfigFile(path)
		viper.SetConfigType("env")

		if err := viper.ReadInConfig(); err != nil {
			path = ""
			if !opts.Quiet {
				fmt.Fprintf(os.Stderr, "ERROR: cannot read config: %v\n", err)
			}
		}
	}

	if err := viper.Unmarshal(&AppConfig); err != nil {
		fmt.Fprintf(os.Stderr, "FATAL: cannot convert config: %v\n", err)
		os.Exit(1) // critical error, the program has to stop
	}

	if AppConfig.DEBUG && path != "" && !opts.Quiet {
		fmt.Fprintf(os.Stderr, "Using config file: %s\n", path)
	}
}

// resolveConfigFile returns the config file to load, or "" if there is none. Problems are
// reported on stderr, and the user may be asked to create a config if there is none.
func resolveConfigFile(opts InitOptions) string {
	workDir, err := os.Getwd()
	if err != nil {
		workDir = "."
	}

	configDir, err := ConfigDir()
	if err != nil {
		configDir = ""
	}

	loc := Locate(workDir, configDir)
	if loc.Path != "" || opts.Quiet {
		return loc.Path
	}

	// A file exists in the config folder, but it cannot be used: never overwrite it silently.
	if loc.HomeExists {
		for _, problem := range loc.Problems {
			fmt.Fprintf(os.Stderr, "WARNING: %s\n", problem)
		}
		fmt.Fprintln(os.Stderr, "Run `tmdb init` to create a new config file.")
		return ""
	}

	// Setups that only use environment variables do not need a file.
	if os.Getenv(apiKeyName) != "" {
		return ""
	}

	if loc.SkippedLocal != "" {
		fmt.Fprintf(os.Stderr, "NOTE: %s\n", loc.SkippedLocal)
	}

	notFound := fmt.Sprintf("No tmdb config found (looked in %s).", loc.LocalPath)
	if loc.HomePath != "" {
		notFound = fmt.Sprintf("No tmdb config found (looked in %s and %s).", loc.LocalPath, loc.HomePath)
	}

	if opts.AllowSetup && loc.HomePath != "" && IsInteractive() {
		fmt.Fprintln(os.Stderr, notFound)
		if path := offerSetup(loc.HomePath); path != "" {
			return path
		}
		fmt.Fprintln(os.Stderr, "Continuing with environment variables and defaults. Run `tmdb init` to create the config later.")
		return ""
	}

	fmt.Fprintf(os.Stderr, "%s Using environment variables and defaults.", notFound)
	if loc.HomePath != "" {
		fmt.Fprint(os.Stderr, " Run `tmdb init` to create a config.")
	}
	fmt.Fprintln(os.Stderr)
	return ""
}

// offerSetup asks the user to create the config file at path.
// It returns path if the file was created, or "" if nothing was created.
func offerSetup(path string) string {
	in := bufio.NewReader(os.Stdin)

	create, err := Confirm(in, os.Stderr, fmt.Sprintf("Create %s now?", path), true)
	if err != nil || !create {
		return ""
	}

	fmt.Fprintln(os.Stderr)
	if err := SetupFromDist(in, os.Stderr, path); err != nil {
		fmt.Fprintf(os.Stderr, "Setup failed: %v\n", err)
		return ""
	}
	fmt.Fprintln(os.Stderr)
	return path
}

func Get() Config {
	return AppConfig
}
