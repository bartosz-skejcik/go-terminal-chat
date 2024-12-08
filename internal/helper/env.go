package helper

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Config represents the entire application configuration
type Config struct {
	Port    int    `mapstructure:"port"`
	Twitch  Twitch `mapstructure:"twitch"`
	Runner  Runner `mapstructure:"runner"`
	AppName string `mapstructure:"app_name"`
}

// Twitch holds Twitch-related configuration
type Twitch struct {
	ClientID     string `mapstructure:"client_id"`
	ClientSecret string `mapstructure:"client_secret"`
	Channel      string `mapstructure:"channel"`
	AuthToken    string `mapstructure:"-"` // Excluded from configuration
}

// Runner holds runner-specific configuration
type Runner struct {
	ShowTimestamps bool `mapstructure:"timestamps"`
	LogMessages    bool `mapstructure:"log_messages"`
}

// NewConfigManager creates and initializes a new configuration manager
func NewConfigManager() (*Config, error) {
	// Get user home directory
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("error getting user home directory: %v", err)
	}

	// Set up configuration paths
	configDir := filepath.Join(home, ".config", "gtc")
	configName := "config"
	configType := "yaml" // Changed to YAML for better readability

	// Ensure config directory exists
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create config directory: %v", err)
	}

	// Configure Viper
	v := viper.New()
	v.SetConfigName(configName)
	v.SetConfigType(configType)
	v.AddConfigPath(configDir)

	// Set default values
	setDefaultConfiguration(v)

	// Create config file if it doesn't exist
	configPath := filepath.Join(configDir, configName+"."+configType)
	if err := createDefaultConfigIfNotExists(v, configPath); err != nil {
		return nil, err
	}

	// Read configuration
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config: %v", err)
	}

	// Unmarshal configuration
	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("unable to decode configuration: %v", err)
	}

	return &config, nil
}

// setDefaultConfiguration sets up default values for the configuration
func setDefaultConfiguration(v *viper.Viper) {
	v.SetDefault("port", 8080)
	v.SetDefault("app_name", "gtc")

	v.SetDefault("twitch.client_id", "")
	v.SetDefault("twitch.client_secret", "")
	v.SetDefault("twitch.channel", "")

	v.SetDefault("runner.timestamps", false)
	v.SetDefault("runner.log_messages", false)
}

// createDefaultConfigIfNotExists creates a default config file if it doesn't exist
func createDefaultConfigIfNotExists(v *viper.Viper, configPath string) error {
	// Check if config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Write default configuration
		if err := v.WriteConfigAs(configPath); err != nil {
			return fmt.Errorf("failed to write default config file: %v", err)
		}
		log.Printf("Created default configuration file at %s", configPath)
	}
	return nil
}

// Print prints the current configuration
func (c *Config) Print() {
	fmt.Println("Application Configuration:")
	fmt.Printf("App Name: %s\n", c.AppName)
	fmt.Printf("Port: %d\n", c.Port)
	fmt.Printf("Twitch Client ID: %s\n", c.Twitch.ClientID)
	fmt.Printf("Twitch Channel: %s\n", c.Twitch.Channel)
	fmt.Printf("Show Timestamps: %t\n", c.Runner.ShowTimestamps)
	fmt.Printf("Log Messages: %t\n", c.Runner.LogMessages)
}

// Update updates a specific configuration value
func (c *Config) Update(key string, value interface{}) error {
	// Create a new Viper instance to avoid modifying the global configuration
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(filepath.Join(os.Getenv("HOME"), ".config", "gtc"))

	// Read existing configuration
	if err := v.ReadInConfig(); err != nil {
		return fmt.Errorf("error reading config: %v", err)
	}

	// Set the new value
	v.Set(key, value)

	// Write the updated configuration
	configPath := v.ConfigFileUsed()
	if err := v.WriteConfigAs(configPath); err != nil {
		return fmt.Errorf("failed to update configuration: %v", err)
	}

	// Reload the configuration
	updatedConfig, err := NewConfigManager()
	if err != nil {
		return fmt.Errorf("failed to reload configuration: %v", err)
	}

	// Update the current config
	*c = *updatedConfig

	return nil
}
