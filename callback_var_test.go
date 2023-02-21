package goflags

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCallbackVarPositive(t *testing.T) {
	toolName := "tool_1"
	want := `updated successfully!`
	got := &bytes.Buffer{}

	flagSet := NewFlagSet()
	var update, disableUpdate bool
	flagSet.CreateGroup("Update", "Update",
		flagSet.CallbackVar(&update, "update", updateCallbackFunc(toolName, got), fmt.Sprintf("update %v to the latest released version", toolName)),
		flagSet.CallbackVarP(&disableUpdate, "disable-update-check", "duc", func() {}, "disable automatic update check"),
	)
	os.Args = []string{
		os.Args[0],
		"-update",
	}
	err := flagSet.Parse()
	assert.Nil(t, err)
	assert.Equal(t, want, got.String())
	tearDown(t.Name())
}

func TestCallbackVarNegative(t *testing.T) {
	toolName := "tool_1"
	got := &bytes.Buffer{}

	if os.Getenv("IS_SUB_PROCESS") == "1" {
		flagSet := NewFlagSet()
		flagSet.CommandLine.SetOutput(got)
		var update bool
		flagSet.CreateGroup("Update", "Update",
			flagSet.CallbackVar(&update, "update", nil, fmt.Sprintf("update %v to the latest released version", toolName)),
		)
		os.Args = []string{
			os.Args[0],
			"-update",
		}
		_ = flagSet.Parse()
	}
	cmd := exec.Command(os.Args[0], "-test.run=TestCallbackVarNegative")
	cmd.Env = append(os.Environ(), "IS_SUB_PROCESS=1")
	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		return
	}
	t.Fatalf("process ran with err %v, want exit error", err)
	tearDown(t.Name())
}

func updateCallbackFunc(toolName string, cliOutput io.Writer) func() {
	return func() {
		fmt.Fprintf(cliOutput, "updated successfully!")
	}
}
