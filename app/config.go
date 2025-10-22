package app

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"

	"sqlcmder/drivers"
	"sqlcmder/models"
)

type Config struct {
	ConfigFile  string
	AppConfig   *models.AppConfig   `toml:"application"`
	Connections []models.Connection `toml:"database"`
}

func defaultConfig() *Config {
	return &Config{
		AppConfig: &models.AppConfig{
			DefaultPageSize:              300,
			SidebarOverlay:               false,
			MaxQueryHistoryPerConnection: 100,
		},
	}
}

func GetConfigPath() (string, error) {
	configDir := os.Getenv("XDG_CONFIG_HOME")
	if configDir == "" {
		dir, err := os.UserConfigDir()
		if err != nil {
			return "", err
		}
		configDir = dir
	}
	return configDir, nil
}

func DefaultConfigFile() (string, error) {
	// Get the executable path
	exePath, err := os.Executable()
	if err != nil {
		return "", err
	}
	// Get the directory containing the executable
	exeDir := filepath.Dir(exePath)
	// Return config file path in the same directory as the executable
	return filepath.Join(exeDir, "config.toml"), nil
}

func LoadConfig(configFile string) error {
	file, err := os.ReadFile(configFile)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	err = toml.Unmarshal(file, App.config)
	if err != nil {
		return err
	}

	for i, conn := range App.config.Connections {
		App.config.Connections[i].URL = parseConfigURL(&conn)
	}

	return nil
}

func (c *Config) SaveConnections(connections []models.Connection) error {
	c.Connections = connections

	if err := os.MkdirAll(filepath.Dir(c.ConfigFile), 0o755); err != nil {
		return err
	}

	file, err := os.Create(c.ConfigFile)
	if err != nil {
		return err
	}
	defer file.Close()

	return toml.NewEncoder(file).Encode(c)
}

// parseConfigURL automatically generates the URL from the connection struct
// if the URL is empty. It is useful for handling usernames and passwords with
// special characters. NOTE: Only MSSQL is supported for now!
func parseConfigURL(conn *models.Connection) string {
	if conn.URL != "" {
		return conn.URL
	}

	// Only MSSQL is supported for now.
	if conn.Provider != drivers.DriverMSSQL {
		return conn.URL
	}

	user := url.QueryEscape(conn.Username)
	pass := url.QueryEscape(conn.Password)

	return fmt.Sprintf(
		"%s://%s:%s@%s:%s?database=%s%s",
		conn.Provider,
		user,
		pass,
		conn.Hostname,
		conn.Port,
		conn.DBName,
		conn.URLParams,
	)
}
