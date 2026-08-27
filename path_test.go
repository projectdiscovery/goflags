package goflags

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFlagSet_SetConfigFilePath(t *testing.T) {
	configFilePath := "/tmp/config.yaml"
	flagSet := NewFlagSet()

	var stringSlice StringSlice
	flagSet.StringSliceVarP(&stringSlice, "header", "H", []string{}, "Header values. Expected usage: -H \"header1\":\"value1\" -H \"header2\":\"value2\"", StringSliceOptions)
	os.Args = []string{
		os.Args[0],
	}
	flagSet.SetConfigFilePath(configFilePath)

	err := flagSet.Parse()
	assert.Nil(t, err)
	gotFilePath, err := flagSet.GetConfigFilePath()
	assert.Nil(t, err)
	assert.Equal(t, configFilePath, gotFilePath)
	tearDown(t.Name())
}

func TestFlagSet_SetConfigFilePaths(t *testing.T) {
	paths := []string{"/etc/tool/config.yaml", "/tmp/user/config.yaml"}
	flagSet := NewFlagSet()
	flagSet.SetConfigFilePaths(paths...)
	paths[1] = "/tmp/changed.yaml"

	gotFilePath, err := flagSet.GetConfigFilePath()
	assert.NoError(t, err)
	assert.Equal(t, "/tmp/user/config.yaml", gotFilePath)
}

func TestFlagSet_SetConfigFilePathsIgnoresEmptyPaths(t *testing.T) {
	flagSet := NewFlagSet()
	defaultPath, err := flagSet.GetConfigFilePath()
	assert.NoError(t, err)

	flagSet.SetConfigFilePaths("", "/etc/tool/config.yaml", "", "/tmp/user/config.yaml", "")
	gotFilePath, err := flagSet.GetConfigFilePath()
	assert.NoError(t, err)
	assert.Equal(t, "/tmp/user/config.yaml", gotFilePath)

	flagSet.SetConfigFilePaths("")
	gotFilePath, err = flagSet.GetConfigFilePath()
	assert.NoError(t, err)
	assert.Equal(t, defaultPath, gotFilePath)
}

func TestFlagSet_SetConfigFilePathEmptyRestoresDefault(t *testing.T) {
	flagSet := NewFlagSet()
	wantFilePath, err := flagSet.GetConfigFilePath()
	assert.NoError(t, err)

	flagSet.SetConfigFilePath("/tmp/config.yaml")
	flagSet.SetConfigFilePath("")

	gotFilePath, err := flagSet.GetConfigFilePath()
	assert.NoError(t, err)
	assert.Equal(t, wantFilePath, gotFilePath)
}
