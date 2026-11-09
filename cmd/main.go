package main

import (
	"github.com/sebastianneubert/tmdb/internal/commands"
	"github.com/sebastianneubert/tmdb/internal/config"
)

// version is set at build time: -ldflags "-X main.version=0.7.0" (see the release workflow).
var version = "dev"

func main() {
	config.SetAppVersion(version)
	commands.Execute()
}
