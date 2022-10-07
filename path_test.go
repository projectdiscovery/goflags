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
