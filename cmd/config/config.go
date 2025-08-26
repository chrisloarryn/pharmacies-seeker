package config

import (
	"errors"
	"fmt"

	"github.com/spf13/viper"
)

type Server struct {
	Port string `json:"port"`
}

type Api struct {
	Pharmacy Pharmacy
}

type Pharmacy struct {
	Url     string `mapstructure:"url"`
	DutyUrl string `mapstructure:"duty_url"`
}

type Config struct {
	Server Server
	Api    Api
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig(path string) (config Config, err error) {
	viper.SetConfigName("properties")
	viper.AddConfigPath(path)
	viper.SetConfigType("yml")

	// Defaults
	viper.SetDefault("server.port", "8080")
	viper.SetDefault("api.pharmacy.url", "https://api.boostr.cl/pharmacies.json")
	viper.SetDefault("api.pharmacy.duty_url", "https://api.boostr.cl/pharmacies/24h.json")

	// Environment overrides
	viper.AutomaticEnv()
	_ = viper.BindEnv("server.port", "PORT")
	_ = viper.BindEnv("api.pharmacy.url", "API_SERVICE_URL")
	_ = viper.BindEnv("api.pharmacy.duty_url", "API_SERVICE_URL_24H")

	// Read config file (lower precedence than env)
	err = viper.ReadInConfig()
	if err != nil {
		var cfne viper.ConfigFileNotFoundError
		if !errors.As(err, &cfne) { // Only fail for real parsing errors
			return Config{}, err
		}
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		fmt.Printf("Unable to decode into struct, %v", err)
	}
	return
}
