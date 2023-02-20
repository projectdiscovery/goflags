package goflags

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUpdateGroupPositive(t *testing.T) {
	toolName := "tool_1"
	want := `updated successfully!`
	got := &bytes.Buffer{}

	flagSet := NewFlagSet()
	flagSet.NewUpdateGroup(toolName, updateCheckCallbackFunc(toolName, "v0.0.1", got), updateCallbackFunc("tool_1", got))
	os.Args = []string{
		os.Args[0],
		"-up",
	}
	err := flagSet.Parse()
	assert.Nil(t, err)
	assert.Equal(t, want, got.String())
	tearDown(t.Name())
}

func TestUpdateGroupNotifyUpdate(t *testing.T) {
	toolName := "tool_1"
	want := fmt.Sprintf("update available on %v!", toolName)
	got := &bytes.Buffer{}

	flagSet := NewFlagSet()
	flagSet.NewUpdateGroup(toolName, updateCheckCallbackFunc(toolName, "v0.0.1", got), updateCallbackFunc("tool_1", got))
	os.Args = []string{
		os.Args[0],
	}
	err := flagSet.Parse()
	assert.Nil(t, err)
	assert.Equal(t, want, got.String())
	tearDown(t.Name())
}

func TestUpdateGroupNegative(t *testing.T) {
	toolName := "tool_1"
	want := ""
	got := &bytes.Buffer{}

	flagSet := NewFlagSet()
	flagSet.NewUpdateGroup(toolName, updateCheckCallbackFunc(toolName, "v0.0.1", got), updateCallbackFunc("tool_1", got))
	os.Args = []string{
		os.Args[0],
		"-up",
		"-duc",
	}
	err := flagSet.Parse()
	assert.Nil(t, err)
	assert.Equal(t, want, got.String())
	tearDown(t.Name())
}

func updateCheckCallbackFunc(toolName, toolVersion string, cliOutput io.Writer) func() {
	return func() {
		if toolVersion < "v0.0.2" {
			fmt.Fprintf(cliOutput, "update available on %v!", toolName)
		}
	}
}

func updateCallbackFunc(toolName string, cliOutput io.Writer) func() {
	return func() {
		fmt.Fprintf(cliOutput, "updated successfully!")
	}
}
