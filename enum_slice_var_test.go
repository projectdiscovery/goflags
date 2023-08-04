package goflags

import (
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

var enumSliceData []string

func TestEnumSliceVar(t *testing.T) {
	t.Run("Test with single value", func(t *testing.T) {
		flagSet := NewFlagSet()
		flagSet.EnumSliceVar(&enumSliceData, "enum", []EnumVariable{Type1}, "enum", AllowdTypes{"type1": Type1, "type2": Type2})
		os.Args = []string{
			os.Args[0],
			"--enum", "type1",
		}
		err := flagSet.Parse()
		assert.Nil(t, err)
		assert.Equal(t, []string{"type1"}, enumSliceData)
		tearDown(t.Name())
	})

	t.Run("Test with multiple value", func(t *testing.T) {
		flagSet := NewFlagSet()
		flagSet.EnumSliceVar(&enumSliceData, "enum", []EnumVariable{Type1}, "enum", AllowdTypes{"type1": Type1, "type2": Type2})
		os.Args = []string{
			os.Args[0],
			"--enum", "type1,type2",
		}
		err := flagSet.Parse()
		assert.Nil(t, err)
		assert.Equal(t, []string{"type1", "type2"}, enumSliceData)
		tearDown(t.Name())
	})

	t.Run("Test with invalid value", func(t *testing.T) {
		if os.Getenv("IS_SUB_PROCESS") == "1" {
			flagSet := NewFlagSet()

			flagSet.EnumSliceVar(&enumSliceData, "enum", []EnumVariable{Nil}, "enum", AllowdTypes{"type1": Type1, "type2": Type2})
			os.Args = []string{
				os.Args[0],
				"--enum", "type3",
			}
			_ = flagSet.Parse()
			return
		}
		cmd := exec.Command(os.Args[0], "-test.run=TestFailEnumVar")
		cmd.Env = append(os.Environ(), "IS_SUB_PROCESS=1")
		err := cmd.Run()
		if e, ok := err.(*exec.ExitError); ok && !e.Success() {
			return
		}
		t.Fatalf("process ran with err %v, want exit error", err)
		tearDown(t.Name())
	})
}
