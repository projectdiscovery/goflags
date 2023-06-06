package goflags

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSizeVar(t *testing.T) {
	t.Run("valid-size", func(t *testing.T) {
		var fileSize Size
		flagSet := NewFlagSet()
		flagSet.CreateGroup("Config", "Config",
			flagSet.SizeVarP(&fileSize, "max-size", "ms", "", "max size of the file"),
		)
		os.Args = []string{
			os.Args[0],
			"-max-size", "2kb",
		}
		err := flagSet.Parse()
		assert.Nil(t, err)
		assert.Equal(t, Size(2048), fileSize)
		tearDown(t.Name())
	})

	t.Run("default-value", func(t *testing.T) {
		var fileSize Size
		flagSet := NewFlagSet()
		flagSet.CreateGroup("Config", "Config",
			flagSet.SizeVarP(&fileSize, "max-size", "ms", "2kb", "max size of the file"),
		)
		os.Args = []string{
			os.Args[0],
		}
		err := flagSet.Parse()
		assert.Nil(t, err)
		assert.Equal(t, Size(2048), fileSize)
		tearDown(t.Name())
	})

	t.Run("without-unit", func(t *testing.T) {
		var fileSize Size
		err := fileSize.Set("2")
		assert.Nil(t, err)
		assert.Equal(t, Size(2097152), fileSize)
		tearDown(t.Name())
	})

	t.Run("invalid-size-unit", func(t *testing.T) {
		var fileSize Size
		err := fileSize.Set("2kilobytes")
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "parse error")
		tearDown(t.Name())
	})
}
