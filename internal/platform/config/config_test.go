package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadUsesDefaultsWhenConfigFileDoesNotExist(t *testing.T) {
	dir := t.TempDir()

	cfg, err := Load(dir)

	require.NoError(t, err)
	assert.Equal(t, "8080", cfg.Server.Port)
	assert.Equal(t, "https://api.boostr.cl/pharmacies.json", cfg.Provider.RegularURL)
	assert.Equal(t, "https://api.boostr.cl/pharmacies/24h.json", cfg.Provider.DutyURL)
	assert.Equal(t, 15*time.Minute, cfg.Sync.Interval)
	assert.Equal(t, 10*time.Second, cfg.Sync.Timeout)
}

func TestLoadReadsConfigAndEnvironmentOverrides(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, "properties.yml"), []byte(`
server:
  port: 9090
provider:
  regular_url: https://example.com/regular
  duty_url: https://example.com/duty
sync:
  interval: 30m
  timeout: 20s
`), 0o644))

	t.Setenv("PORT", "7070")
	t.Setenv("SYNC_INTERVAL", "45m")

	cfg, err := Load(dir)

	require.NoError(t, err)
	assert.Equal(t, "7070", cfg.Server.Port)
	assert.Equal(t, "https://example.com/regular", cfg.Provider.RegularURL)
	assert.Equal(t, "https://example.com/duty", cfg.Provider.DutyURL)
	assert.Equal(t, 45*time.Minute, cfg.Sync.Interval)
	assert.Equal(t, 20*time.Second, cfg.Sync.Timeout)
}

func TestLoadReturnsParseErrorForInvalidConfigFile(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, "properties.yml"), []byte("server:\n  port: ["), 0o644))

	_, err := Load(dir)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "read config")
}

func TestLoadReturnsValidationErrorForInvalidEnvironmentOverride(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("SYNC_TIMEOUT", "0s")

	_, err := Load(dir)

	require.EqualError(t, err, "sync.timeout must be greater than zero")
}

func TestValidateRejectsInvalidConfig(t *testing.T) {
	testCases := []struct {
		name    string
		cfg     Config
		message string
	}{
		{
			name:    "missing port",
			cfg:     Config{Provider: ProviderConfig{RegularURL: "a", DutyURL: "b"}, Sync: SyncConfig{Interval: time.Second, Timeout: time.Second}},
			message: "server.port is required",
		},
		{
			name:    "missing regular url",
			cfg:     Config{Server: ServerConfig{Port: "8080"}, Provider: ProviderConfig{DutyURL: "b"}, Sync: SyncConfig{Interval: time.Second, Timeout: time.Second}},
			message: "provider.regular_url is required",
		},
		{
			name:    "missing duty url",
			cfg:     Config{Server: ServerConfig{Port: "8080"}, Provider: ProviderConfig{RegularURL: "a"}, Sync: SyncConfig{Interval: time.Second, Timeout: time.Second}},
			message: "provider.duty_url is required",
		},
		{
			name:    "invalid interval",
			cfg:     Config{Server: ServerConfig{Port: "8080"}, Provider: ProviderConfig{RegularURL: "a", DutyURL: "b"}, Sync: SyncConfig{Timeout: time.Second}},
			message: "sync.interval must be greater than zero",
		},
		{
			name:    "invalid timeout",
			cfg:     Config{Server: ServerConfig{Port: "8080"}, Provider: ProviderConfig{RegularURL: "a", DutyURL: "b"}, Sync: SyncConfig{Interval: time.Second}},
			message: "sync.timeout must be greater than zero",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := validate(tc.cfg)

			require.EqualError(t, err, tc.message)
		})
	}
}
