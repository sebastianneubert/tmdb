package commands

import (
	"fmt"
	"os"

	"github.com/sebastianneubert/tmdb/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "tmdb",
	Short: "A CLI to find movies and TV shows on your streaming services.",
	Long: `tmdb is a command line interface to find movies and TV shows on your
	       subscribed streaming services with params and filters that match your criteria.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.AddCommand(topCmd)
	rootCmd.AddCommand(popularCmd)
	rootCmd.AddCommand(actorCmd)
	rootCmd.AddCommand(showsCmd)
	rootCmd.AddCommand(searchCmd)
	rootCmd.AddCommand(genresCmd)
	rootCmd.AddCommand(providersCmd)
	rootCmd.AddCommand(initCmd)
}

// commandsWithoutConfig never need a config file. For them the config is loaded silently and
// without questions: output or prompts would break shell completion and help.
var commandsWithoutConfig = map[string]bool{
	"help":             true,
	"completion":       true,
	"bash":             true,
	"zsh":              true,
	"fish":             true,
	"powershell":       true,
	"__complete":       true,
	"__completeNoDesc": true,
	"init":             true,
}

// initConfig loads the config before a command runs.
func initConfig() {
	args := os.Args[1:]

	quiet := false
	for _, arg := range args {
		if arg == "--" {
			break
		}
		if arg == "-h" || arg == "--help" || arg == "--version" {
			quiet = true
		}
	}

	cmd, _, err := rootCmd.Find(args)
	if err != nil || cmd == rootCmd {
		quiet = true
	}
	for c := cmd; c != nil; c = c.Parent() {
		if commandsWithoutConfig[c.Name()] {
			quiet = true
		}
	}

	config.InitWith(config.InitOptions{Quiet: quiet, AllowSetup: !quiet})
}

func activeConfigHint() string {
	if path := config.ActiveFile(); path != "" {
		return "Active config: " + path
	}
	return "Active config: none (Environment variables and defaults)"
}

func Execute() {
	rootCmd.Version = config.AppVersion()
	rootCmd.Long += "\n" + activeConfigHint()

	bindCommandFlags(topCmd)
	bindCommandFlags(popularCmd)
	bindCommandFlags(actorCmd)
	bindCommandFlags(searchCmd)
	bindCommandFlags(showsCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func bindCommandFlags(cmd *cobra.Command) {
	viper.BindPFlag("PROVIDERS", cmd.Flags().Lookup("providers"))
	viper.BindPFlag("REGION", cmd.Flags().Lookup("region"))
	viper.BindPFlag("MIN_RATING", cmd.Flags().Lookup("min-rating"))
	viper.BindPFlag("MIN_VOTES", cmd.Flags().Lookup("min-votes"))
	viper.BindPFlag("API_TIMEOUT_SECONDS", cmd.Flags().Lookup("timeout"))
}
