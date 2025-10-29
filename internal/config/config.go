package config

import (
	"os"
	"path/filepath"

	"github.com/spf13/viper"
	"github.com/subosito/gotenv"
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

// LoadEnvFile loads .env file from project root if it exists
func LoadEnvFile() {
	// Try to find project root (look for go.mod or .git)
	wd, err := os.Getwd()
	if err != nil {
		return // Can't determine working directory, skip
	}

	// Walk up the directory tree to find project root
	dir := wd
	for {
		envPath := filepath.Join(dir, ".env")
		if _, err := os.Stat(envPath); err == nil {
			// Found .env file, load it
			if err := gotenv.Load(envPath); err == nil {
				// Environment variables are now loaded
			}
			break
		}

		// Check if we're at root (has go.mod or .git)
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			break // Found project root
		}
		if _, err := os.Stat(filepath.Join(dir, ".git")); err == nil {
			break // Found git root
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break // Reached filesystem root
		}
		dir = parent
	}
}

func Load() Config {
	// Load .env file first (if present in project root)
	LoadEnvFile()

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