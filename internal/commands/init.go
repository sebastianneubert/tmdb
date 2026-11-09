package commands

import (
	"bufio"
	"errors"
	"fmt"
	"os"

	"github.com/sebastianneubert/tmdb/internal/config"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Create the tmdb config file in your config folder.",
	Long: `Creates the tmdb config file interactively, based on the dist file .env.example.

The file is stored in ~/.config/tmdb/.env (or $XDG_CONFIG_HOME/tmdb/.env), so tmdb works
from every folder. A .env in the current folder still takes precedence if it starts with
the header "# tmdb version x.y.z".

The configuration is looked up in this order:
  1. .env in the current directory
  2. .env in the users home/.config folder (~/.config/tmdb), generate it with the command "tmdb init"
  3. Environment vars and default values
`,
	Args: cobra.NoArgs,
	RunE: runInit,
}

func runInit(cmd *cobra.Command, args []string) error {
	cmd.SilenceUsage = true
	cmd.SilenceErrors = true // Execute() prints the error

	if !config.IsInteractive() {
		return errors.New("tmdb init needs an interactive terminal")
	}

	path, err := config.ConfigFilePath()
	if err != nil {
		return fmt.Errorf("cannot determine the config folder: %w", err)
	}

	in := bufio.NewReader(cmd.InOrStdin())
	out := cmd.OutOrStdout()

	if _, err := os.Stat(path); err == nil {
		overwrite, err := config.Confirm(in, out, fmt.Sprintf("%s already exists. Overwrite it?", path), false)
		if err != nil {
			return err
		}
		if !overwrite {
			fmt.Fprintln(out, "Nothing changed.")
			return nil
		}
	}

	return config.SetupFromDist(in, out, path)
}
