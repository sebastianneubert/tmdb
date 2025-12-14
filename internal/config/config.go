package config

import (
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
}

var AppConfig Config

func Init() {
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	viper.SetDefault("REGION", DefaultRegion)
	viper.SetDefault("PROVIDERS", DefaultProviders)
	viper.SetDefault("MIN_RATING", DefaultMinRating)
	viper.SetDefault("MIN_VOTES", DefaultMinVotes)
	viper.SetDefault("API_TIMEOUT_SECONDS", DefaultTimeout)
	viper.SetDefault("DEBUG", DefaultDebug)
	viper.SetDefault("TMDB_API_KEY", "")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("WARNUNG: Keine .env-Datei gefunden. Verwende Umgebungsvariablen und Defaults.")
		} else {
			fmt.Println("FEHLER: Konnte Config nicht lesen:", err)
		}
	} else {
		// fmt.Println("Config erfolgreich geladen aus:", viper.ConfigFileUsed())
	}

	if err := viper.Unmarshal(&AppConfig); err != nil {
		fmt.Printf("FATAL: Konvertierungsfehler in das Config-Struct: %v\n", err)
		os.Exit(1) // Das ist ein kritischer Fehler, Programm sollte stoppen
	}
}

func Get() Config {
	return AppConfig
}
