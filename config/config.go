package config

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

func DefaultConfig() *Config {
	return &Config{
		AppConfig: &models.AppConfig{
			DefaultPageSize:              300,
			SidebarOverlay:               false,
			MaxQueryHistoryPerConnection: 100,
			Theme:                        models.ThemeDark, // Default to dark theme
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

func LoadConfig(configFile string, config *Config) error {
	file, err := os.ReadFile(configFile)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	err = toml.Unmarshal(file, config)
	if err != nil {
		return err
	}

	for i, conn := range config.Connections {
		config.Connections[i].DSN = parseConfigDSN(&conn)
		config.Connections[i].SetDSNValue() // Ensure DsnValue is set
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

// parseConfigDSN automatically generates the DSN from the connection struct
// if the DSN is empty. It is useful for handling usernames and passwords with
// special characters. NOTE: Only MSSQL is supported for now!
func parseConfigDSN(conn *models.Connection) string {
	// Handle new DSN structure with priority: custom > auto-generated
	if conn.DsnCustom != "" {
		conn.DsnValue = conn.DsnCustom
		return conn.DsnCustom
	}
	
	if conn.DsnAuto != "" {
		conn.DsnValue = conn.DsnAuto
		return conn.DsnAuto
	}
	
	// Fallback to legacy DSN field for backward compatibility
	if conn.DSN != "" {
		conn.DsnValue = conn.DSN
		return conn.DSN
	}

	// Only MSSQL is supported for now.
	if conn.Driver != drivers.DriverMSSQL {
		return conn.DSN
	}

	user := url.QueryEscape(conn.Username)
	pass := url.QueryEscape(conn.Password)

	autoDSN := fmt.Sprintf(
		"%s://%s:%s@%s:%s?database=%s%s",
		conn.Driver,
		user,
		pass,
		conn.Hostname,
		conn.Port,
		conn.DBName,
		conn.DSNParams,
	)
	
	// Update the new DSN fields
	conn.DsnAuto = autoDSN
	conn.DsnValue = autoDSN
	
	return autoDSN
}
