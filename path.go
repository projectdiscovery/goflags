package goflags

import (
	"os"
	"path/filepath"
	"strings"

	fileutil "github.com/projectdiscovery/utils/file"
	folderutil "github.com/projectdiscovery/utils/folder"
)

// GetConfigFilePath returns the config file path
func (flagSet *FlagSet) GetConfigFilePath() (string, error) {
	// return configFilePath if already set
	if flagSet.configFilePath != "" {
		return flagSet.configFilePath, nil
	}

	err := migrateConfigFiles()
	if err != nil {
		return "", err
	}

	return buildConfigFilePath(), nil
}

// SetConfigFilePath sets custom config file path
func (flagSet *FlagSet) SetConfigFilePath(filePath string) {
	flagSet.configFilePath = filePath
}

// Deprecated: Use FlagSet.GetConfigFilePath instead.
// GetConfigFilePath returns the default config file path
func GetConfigFilePath() (string, error) {
	err := migrateConfigFiles()
	if err != nil {
		return "", err
	}

	return buildConfigFilePath(), nil
}

func buildConfigFilePath() string {
	appConfigDir := buildAppConfigDirPath()
	return filepath.Join(appConfigDir, "config.yaml")
}

func buildAppConfigDirPath() string {
	appName := buildAppName()
	return folderutil.AppConfigDirOrDefault(".", appName)
}

func buildAppName() string {
	appName := filepath.Base(os.Args[0])
	return strings.TrimSuffix(appName, filepath.Ext(appName))
}

// Note: This is a temporary function to migrate config files from old os-specific config path to os-agnostic config path
func migrateConfigFiles() error {
	appName := buildAppName()
	homePath, _ := os.UserHomeDir()
	sourceDir := filepath.Join(homePath, ".config", appName)
	destinationDir := buildAppConfigDirPath()

	ok := fileutil.FolderExists(sourceDir)
	if !ok {
		return nil
	}
	return folderutil.MigrateDir(sourceDir, destinationDir)
}
