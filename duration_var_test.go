package goflags

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDurationVar(t *testing.T) {
	t.Run("day-unit", func(t *testing.T) {
		var tt time.Duration
		flagSet := NewFlagSet()
		flagSet.CreateGroup("Config", "Config",
			flagSet.DurationVarP(&tt, "time-out", "tm", 0, "timeout for the process"),
		)
		os.Args = []string{
			os.Args[0],
			"-time-out", "2d",
		}
		err := flagSet.Parse()
		assert.Nil(t, err)
		assert.Equal(t, 2*24*time.Hour, tt)
		tearDown(t.Name())
	})

	t.Run("without-unit", func(t *testing.T) {
		var tt time.Duration
		flagSet := NewFlagSet()
		flagSet.CreateGroup("Config", "Config",
			flagSet.DurationVarP(&tt, "time-out", "tm", 0, "timeout for the process"),
		)
		os.Args = []string{
			os.Args[0],
			"-time-out", "2",
		}
		err := flagSet.Parse()
		assert.Nil(t, err)
		assert.Equal(t, 2*time.Second, tt)
		tearDown(t.Name())
	})
}
