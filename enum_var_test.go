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

func TestAllowedTypesStringIsSorted(t *testing.T) {
	allowed := AllowdTypes{"yaml": Type1, "json": Type2, "csv": Nil, "sarif": Type1, "xml": Type2}

	assert.Equal(t, "csv, json, sarif, xml, yaml", allowed.String())

	// Ranging over a map gives a different order between calls, so the help
	// text and the "allowed values are" error used to change shape run to run.
	for i := 0; i < 20; i++ {
		assert.Equal(t, "csv, json, sarif, xml, yaml", allowed.String())
	}
}

func TestEnumVarErrorListsAllowedValuesInOrder(t *testing.T) {
	var value string
	enum := &EnumVar{
		allowedTypes: AllowdTypes{"yaml": Type1, "json": Type2, "csv": Nil},
		value:        &value,
	}

	err := enum.Set("toml")

	assert.EqualError(t, err, "allowed values are csv, json, yaml")
}
