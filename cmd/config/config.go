package config

import (
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
	Url string
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
	err = viper.ReadInConfig()
	if err != nil {
		return Config{}, err
	}

	// Set undefined variables
	viper.SetDefault("server.port", "8080")
	viper.SetDefault("api.pharmacy.url", "https://farmanet.minsal.cl/maps/index.php/ws/getLocalesRegion?id_region=6")

	err = viper.Unmarshal(&config)
	if err != nil {
		fmt.Printf("Unable to decode into struct, %v", err)
	}
	return
}
