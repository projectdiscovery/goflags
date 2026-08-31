package goflags

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSuccessfulCallback(t *testing.T) {
	toolName := "tool_1"
	want := `updated successfully!`
	got := &bytes.Buffer{}

	flagSet := NewFlagSet()
	flagSet.CreateGroup("Update", "Update",
		flagSet.CallbackVar(updateCallbackFunc(toolName, got), "update", fmt.Sprintf("update %v to the latest released version", toolName)),
		flagSet.CallbackVarP(func() {}, "disable-update-check", "duc", "disable automatic update check"),
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

func TestConfigFileDoesNotExecuteCallback(t *testing.T) {

	configFile := filepath.Join(t.TempDir(), "config.yaml")
	writeTestConfig(t, configFile, "update: true\n")
	called := false

	flagSet := NewFlagSet()
	flagSet.CallbackVar(func() { called = true }, "update", "update tool")
	flagSet.SetConfigFilePath(configFile)

	assert.NoError(t, flagSet.Parse(""))
	assert.False(t, called)
}

func TestMergeConfigFileDoesNotExecuteCallback(t *testing.T) {
	configFile := filepath.Join(t.TempDir(), "config.yaml")
	writeTestConfig(t, configFile, "update: true\n")
	called := false

	flagSet := NewFlagSet()
	flagSet.CallbackVar(func() { called = true }, "update", "update tool")

	assert.NoError(t, flagSet.MergeConfigFile(configFile))
	assert.False(t, called)
}

func TestFailCallback(t *testing.T) {
	toolName := "tool_1"
	got := &bytes.Buffer{}

	if os.Getenv("IS_SUB_PROCESS") == "1" {
		flagSet := NewFlagSet()
		flagSet.CommandLine.SetOutput(got)
		flagSet.CreateGroup("Update", "Update",
			flagSet.CallbackVar(nil, "update", fmt.Sprintf("update %v to the latest released version", toolName)),
		)
		os.Args = []string{
			os.Args[0],
			"-update",
		}
		_ = flagSet.Parse()
	}
	cmd := exec.Command(os.Args[0], "-test.run=TestFailCallback")
	cmd.Env = append(os.Environ(), "IS_SUB_PROCESS=1")
	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		return
	}
	t.Fatalf("process ran with err %v, want exit error", err)
	tearDown(t.Name())
}

func updateCallbackFunc(_ string, cliOutput io.Writer) func() {
	return func() {
		fmt.Fprintf(cliOutput, "updated successfully!")
	}
}
