package settings

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

var config = viper.New()

func ConfigInit() {
	systemConfigPath, err := os.UserConfigDir()
	if err != nil {
		panic(fmt.Errorf("System Config Directory does not exist: %w", err))
	}

	configPath := fmt.Sprintf("%s%s", systemConfigPath, "/marina")
	createConfigDirIfNotExist(configPath)

	config.SetConfigName("config")
	config.SetConfigType("toml")
	config.AddConfigPath(configPath)

	setDefaults()

	if err := config.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// file not found
			err = config.SafeWriteConfig()
			if err != nil {
				panic(fmt.Errorf("Error creating config file: %w", err))
			}
		} else {
			// error reading config
			panic(fmt.Errorf("Error reading config file: %w", err))
		}
	}
}

func setDefaults() {
	config.SetDefault("UseLinuxCompatibilityVerions", false)
	config.SetDefault("DownloadOSIncompatibleVersions", false)
}

// Getters
func ShouldUseLinuxCompatibilityVersion() bool {
	return config.GetBool("UseLinuxCompatibilityVerions")
}

// Setters
func SetLinuxCompatibility(useCompatibilityVersion bool) {
	config.Set("UseLinuxCompatibilityVerions", useCompatibilityVersion)
	saveChanges()
}

func createConfigDirIfNotExist(path string) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		err := os.Mkdir(path, os.ModePerm)
		if err != nil {
			panic(fmt.Errorf("Error reading config file: %w", err))
		}
	} else if err != nil {
		panic(fmt.Errorf("Error reading config file: %w", err))
	}
}

func saveChanges() {
	err := config.WriteConfig()
	if err != nil {
		panic(fmt.Errorf("Error writing config file: %w", err))
	}
}
