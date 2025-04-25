package config

import (
    "os"
    "path/filepath"
    "fmt"
    "github.com/spf13/viper"
)
    


func InitializeConfig() error {
    configDir, _ := os.UserConfigDir()
    appConfigDir := filepath.Join(configDir, "gbtest")
    appConfigFile := filepath.Join(appConfigDir, "config.json")
    if _, err := os.Stat(appConfigFile); os.IsNotExist(err) {
		err := os.MkdirAll(appConfigDir, 0755)
		if err != nil {
			return fmt.Errorf("failed to create config directory: %w", err)
		}

		defaultConfig := []byte("{}")
		err = os.WriteFile(appConfigFile, defaultConfig, 0644)
		if err != nil {
			return fmt.Errorf("failed to write default config file: %w", err)
		}
	}

    viper.SetConfigFile(appConfigFile)
    err := viper.ReadInConfig()
    viper.SetDefault("auth.ttl", 86400) 
    viper.SetDefault("auth.authenticated", false)
    viper.SetDefault("auth.lastChecked", 0)


    if err != nil {
            return fmt.Errorf("error reading config file: %w", err)
    }
    return nil
}
