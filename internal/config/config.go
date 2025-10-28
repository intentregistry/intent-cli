package config

import (
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	APIURL    string
	Token     string
	Telemetry bool
}

func configDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".intent")
}

func EnsureDir() error {
	return os.MkdirAll(configDir(), 0o755)
}

func Load() Config {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(configDir())
	_ = v.ReadInConfig()

	api := v.GetString("api_url")
	if env := os.Getenv("INTENT_API_URL"); env != "" { api = env }
	if api == "" { api = "https://api.intentregistry.com" } // por defecto

	tok := v.GetString("token")
	if env := os.Getenv("INTENT_TOKEN"); env != "" { tok = env }

	// Telemetry: check environment variable first, then config file
	telemetry := false
	if env := os.Getenv("INTENT_TELEMETRY"); env != "" {
		telemetry = env == "1" || env == "true"
	} else {
		telemetry = v.GetBool("telemetry")
	}

	return Config{APIURL: api, Token: tok, Telemetry: telemetry}
}

func SaveToken(token string) error {
	if err := EnsureDir(); err != nil { return err }
	v := viper.New()
	v.Set("token", token)
	return v.WriteConfigAs(filepath.Join(configDir(), "config.yaml"))
}

func SaveConfig(apiURL, token string) error {
	if err := EnsureDir(); err != nil { return err }
	v := viper.New()
	v.Set("api_url", apiURL)
	v.Set("token", token)
	return v.WriteConfigAs(filepath.Join(configDir(), "config.yaml"))
}