package config

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	Provider ProviderConfig
	Sync     SyncConfig
}

type ServerConfig struct {
	Port string
}

type ProviderConfig struct {
	RegularURL string
	DutyURL    string
}

type SyncConfig struct {
	Interval time.Duration
	Timeout  time.Duration
}

func Load(path string) (Config, error) {
	v := viper.New()
	v.SetConfigName("properties")
	v.SetConfigType("yml")
	v.AddConfigPath(path)
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	v.SetDefault("server.port", "8080")
	v.SetDefault("provider.regular_url", "https://api.boostr.cl/pharmacies.json")
	v.SetDefault("provider.duty_url", "https://api.boostr.cl/pharmacies/24h.json")
	v.SetDefault("sync.interval", "15m")
	v.SetDefault("sync.timeout", "10s")

	_ = v.BindEnv("server.port", "PORT")
	_ = v.BindEnv("provider.regular_url", "PROVIDER_REGULAR_URL", "API_SERVICE_URL")
	_ = v.BindEnv("provider.duty_url", "PROVIDER_DUTY_URL", "API_SERVICE_URL_24H")
	_ = v.BindEnv("sync.interval", "SYNC_INTERVAL")
	_ = v.BindEnv("sync.timeout", "SYNC_TIMEOUT")

	if err := v.ReadInConfig(); err != nil {
		var configNotFound viper.ConfigFileNotFoundError
		if !errors.As(err, &configNotFound) {
			return Config{}, fmt.Errorf("read config: %w", err)
		}
	}

	cfg := Config{
		Server: ServerConfig{
			Port: v.GetString("server.port"),
		},
		Provider: ProviderConfig{
			RegularURL: v.GetString("provider.regular_url"),
			DutyURL:    v.GetString("provider.duty_url"),
		},
		Sync: SyncConfig{
			Interval: v.GetDuration("sync.interval"),
			Timeout:  v.GetDuration("sync.timeout"),
		},
	}

	if err := validate(cfg); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func validate(cfg Config) error {
	if cfg.Server.Port == "" {
		return errors.New("server.port is required")
	}
	if cfg.Provider.RegularURL == "" {
		return errors.New("provider.regular_url is required")
	}
	if cfg.Provider.DutyURL == "" {
		return errors.New("provider.duty_url is required")
	}
	if cfg.Sync.Interval <= 0 {
		return errors.New("sync.interval must be greater than zero")
	}
	if cfg.Sync.Timeout <= 0 {
		return errors.New("sync.timeout must be greater than zero")
	}
	return nil
}
