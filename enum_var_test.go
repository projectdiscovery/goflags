package goflags

import (
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

var enumString string

const (
	Nil EnumVariable = iota
	Type1
	Type2
)

func TestSuccessfulEnumVar(t *testing.T) {
	flagSet := NewFlagSet()
	flagSet.EnumVar(&enumString, "enum", Type1, "enum", AllowdTypes{"type1": Type1, "type2": Type2})
	os.Args = []string{
		os.Args[0],
		"--enum", "type1",
	}
	err := flagSet.Parse()
	assert.Nil(t, err)
	assert.Equal(t, "type1", enumString)
	tearDown(t.Name())
}

func TestFailEnumVar(t *testing.T) {
	if os.Getenv("IS_SUB_PROCESS") == "1" {
		flagSet := NewFlagSet()

		flagSet.EnumVar(&enumString, "enum", Nil, "enum", AllowdTypes{"type1": Type1, "type2": Type2})
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
}
