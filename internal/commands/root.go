package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/sebastianneubert/tmdb/internal/config"
)

var rootCmd = &cobra.Command{
	Use:   "tmdb",
	Short: "A CLI to find movies and TV shows on your streaming services.",
	Long: `tmdb ist ein Command Line Interface zur Suche von Inhalten auf Streaming-Diensten.

Die Konfiguration erfolgt über die .env Datei oder Umgebungsvariablen.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	cobra.OnInitialize(config.Init)
	rootCmd.AddCommand(topCmd)
	rootCmd.AddCommand(actorCmd)
	rootCmd.AddCommand(showsCmd)
	rootCmd.AddCommand(searchCmd)
	rootCmd.AddCommand(genresCmd)
}

func Execute() {
  bindCommandFlags(topCmd)
	bindCommandFlags(actorCmd)

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