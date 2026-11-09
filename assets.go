// Package tmdb only exists to embed files from the repository root into the binary.
package tmdb

import _ "embed"

// EnvExample is the dist file .env.example. It is the template for the config file
// that `tmdb init` creates in the config folder of the user.
//
//go:embed .env.example
var EnvExample string
