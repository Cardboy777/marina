package settings

import (
	"fmt"
	"marina/constants"
	"os"
	"path/filepath"
	"strings"

	"github.com/adrg/xdg"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var config = viper.New()

var restart = false

var configDir *string

func Init() {
	configDir = pflag.StringP("config-dir", "c", getDefaultConfigDir(), "Directory that contains \"config.toml\"")
	pflag.Parse()

	createConfigDirIfNotExist(*configDir)

	config.SetConfigName("config")
	config.SetConfigType("toml")
	config.AddConfigPath(*configDir)

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

func getDefaultConfigDir() string {
	path, err := os.UserConfigDir()
	if err != nil {
		panic(fmt.Errorf("System Config Directory does not exist: %w", err))
	}
	return filepath.Join(path, strings.ToLower(constants.AppName))
}

func setDefaults() {
	config.SetDefault("InstallDir", GetDefaultInstallDir())
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

func GetDefaultInstallDir() string {
	return filepath.Join(xdg.DataHome, strings.ToLower(constants.AppName))
}

func SetInstallDir(dir string) {
	if len(dir) == 0 {
		dir = GetDefaultInstallDir()
	}
	if dir != GetInstallDirName() {
		config.Set("InstallDir", dir)
		saveChanges()
		restart = true
	}
}

func saveChanges() {
	err := config.WriteConfig()
	if err != nil {
		panic(fmt.Errorf("Error writing config file: %w", err))
	}
}

func GetInstallDirName() string {
	return filepath.Join(config.GetString("InstallDir"))
}

func ShouldRestart() bool {
	if restart {
		restart = false
		return true
	}
	return false
}
