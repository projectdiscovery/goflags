package goflags

import (
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"testing"
	"time"
)

func TestAddCommonFlags(t *testing.T) {
	flagSet := NewFlagSet()
	commonFlags := flagSet.AddCommonFlags()

	if commonFlags == nil {
		t.Fatal("AddCommonFlags returned nil")
	}

	if commonFlags.MaxTime != 0 {
		t.Errorf("Expected default MaxTime to be 0, got %v", commonFlags.MaxTime)
	}

	err := flagSet.Parse("-max-time", "1h30m")
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	expectedDuration := 90 * time.Minute
	if commonFlags.MaxTime != expectedDuration {
		t.Errorf("Expected MaxTime to be %v, got %v", expectedDuration, commonFlags.MaxTime)
	}
}

func TestAddCommonFlagsShortFlag(t *testing.T) {
	flagSet := NewFlagSet()
	commonFlags := flagSet.AddCommonFlags()

	err := flagSet.Parse("-mt", "45m")
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	expectedDuration := 45 * time.Minute
	if commonFlags.MaxTime != expectedDuration {
		t.Errorf("Expected MaxTime to be %v, got %v", expectedDuration, commonFlags.MaxTime)
	}
}

func TestMaxTimeInterrupt(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("skipping on Windows: GenerateConsoleCtrlEvent sends to all console processes including test runner")
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	defer signal.Stop(sigChan)

	flagSet := NewFlagSet()
	flagSet.AddCommonFlags()
	_ = flagSet.Parse("-mt", "100ms")

	select {
	case <-sigChan:
		// Success - received interrupt
	case <-time.After(500 * time.Millisecond):
		t.Error("Expected interrupt signal within 500ms")
	}
}

func TestMaxTimeHandlerStartsAfterConfigMerge(t *testing.T) {
	dir := t.TempDir()
	primaryConfig := filepath.Join(dir, "primary.yaml")
	explicitConfig := filepath.Join(dir, "explicit.yaml")
	writeTestConfig(t, primaryConfig, "")
	writeTestConfig(t, explicitConfig, "max-time: 1h\n")

	flagSet := NewFlagSet()
	commonFlags := flagSet.AddCommonFlags()
	stopMaxTimeHandler(t, commonFlags)
	flagSet.SetConfigFilePath(primaryConfig)

	if err := flagSet.Parse(""); err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
	if commonFlags.maxTimeTimer != nil {
		t.Fatal("max-time timer started with a zero duration")
	}

	if err := flagSet.MergeConfigFile(explicitConfig); err != nil {
		t.Fatalf("MergeConfigFile failed: %v", err)
	}
	if commonFlags.maxTimeTimer == nil {
		t.Fatal("max-time timer was not started after config merge")
	}
	if commonFlags.activeMaxTime != time.Hour {
		t.Fatalf("active max-time = %v, want %v", commonFlags.activeMaxTime, time.Hour)
	}
}

func TestMaxTimeHandlerReplacesChangedConfig(t *testing.T) {
	dir := t.TempDir()
	primaryConfig := filepath.Join(dir, "primary.yaml")
	explicitConfig := filepath.Join(dir, "explicit.yaml")
	writeTestConfig(t, primaryConfig, "max-time: 1h\n")
	writeTestConfig(t, explicitConfig, "max-time: 2h\n")

	flagSet := NewFlagSet()
	commonFlags := flagSet.AddCommonFlags()
	stopMaxTimeHandler(t, commonFlags)
	flagSet.SetConfigFilePath(primaryConfig)

	if err := flagSet.Parse(""); err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
	oldTimer := commonFlags.maxTimeTimer
	if oldTimer == nil {
		t.Fatal("max-time timer was not started from primary config")
	}

	if err := flagSet.MergeConfigFile(explicitConfig); err != nil {
		t.Fatalf("MergeConfigFile failed: %v", err)
	}
	if commonFlags.maxTimeTimer == oldTimer {
		t.Fatal("max-time timer was not replaced after duration changed")
	}
	if oldTimer.Stop() {
		t.Fatal("replaced max-time timer was still active")
	}
	if commonFlags.activeMaxTime != 2*time.Hour {
		t.Fatalf("active max-time = %v, want %v", commonFlags.activeMaxTime, 2*time.Hour)
	}
}

func TestMaxTimeHandlerKeepsTimerWhenConfigDoesNotChangeDuration(t *testing.T) {
	dir := t.TempDir()
	primaryConfig := filepath.Join(dir, "primary.yaml")
	explicitConfig := filepath.Join(dir, "explicit.yaml")
	writeTestConfig(t, primaryConfig, "max-time: 1h\n")
	writeTestConfig(t, explicitConfig, "value: changed\n")

	flagSet := NewFlagSet()
	commonFlags := flagSet.AddCommonFlags()
	stopMaxTimeHandler(t, commonFlags)
	flagSet.SetConfigFilePath(primaryConfig)
	var value string
	flagSet.StringVar(&value, "value", "built-in", "test value")

	if err := flagSet.Parse(""); err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
	oldTimer := commonFlags.maxTimeTimer
	if oldTimer == nil {
		t.Fatal("max-time timer was not started from primary config")
	}

	if err := flagSet.MergeConfigFile(explicitConfig); err != nil {
		t.Fatalf("MergeConfigFile failed: %v", err)
	}
	if commonFlags.maxTimeTimer != oldTimer {
		t.Fatal("unrelated config merge replaced the max-time timer")
	}
	if value != "changed" {
		t.Fatalf("value = %q, want %q", value, "changed")
	}
}

func TestMaxTimeHandlerStopsWhenConfigDisablesIt(t *testing.T) {
	dir := t.TempDir()
	primaryConfig := filepath.Join(dir, "primary.yaml")
	explicitConfig := filepath.Join(dir, "explicit.yaml")
	writeTestConfig(t, primaryConfig, "max-time: 1h\n")
	writeTestConfig(t, explicitConfig, "max-time: 0s\n")

	flagSet := NewFlagSet()
	commonFlags := flagSet.AddCommonFlags()
	stopMaxTimeHandler(t, commonFlags)
	flagSet.SetConfigFilePath(primaryConfig)

	if err := flagSet.Parse(""); err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
	oldTimer := commonFlags.maxTimeTimer
	if oldTimer == nil {
		t.Fatal("max-time timer was not started from primary config")
	}

	if err := flagSet.MergeConfigFile(explicitConfig); err != nil {
		t.Fatalf("MergeConfigFile failed: %v", err)
	}
	if commonFlags.maxTimeTimer != nil {
		t.Fatal("max-time timer remained active after max-time was disabled")
	}
	if oldTimer.Stop() {
		t.Fatal("disabled max-time timer was still active")
	}
}

func TestMaxTimeHandlerStartsAfterConfigError(t *testing.T) {
	configFile := filepath.Join(t.TempDir(), "config.yaml")
	writeTestConfig(t, configFile, "max-time: [\n")

	flagSet := NewFlagSet()
	commonFlags := flagSet.AddCommonFlags()
	stopMaxTimeHandler(t, commonFlags)
	flagSet.SetConfigFilePath(configFile)

	if err := flagSet.Parse("-max-time", "1h"); err == nil {
		t.Fatal("Parse succeeded with malformed config")
	}
	if commonFlags.maxTimeTimer == nil {
		t.Fatal("CLI max-time timer was not started after config error")
	}
	if commonFlags.activeMaxTime != time.Hour {
		t.Fatalf("active max-time = %v, want %v", commonFlags.activeMaxTime, time.Hour)
	}
}

func TestMaxTimeHandlerStartsFromMergeAfterConfigError(t *testing.T) {
	dir := t.TempDir()
	primaryConfig := filepath.Join(dir, "primary.yaml")
	explicitConfig := filepath.Join(dir, "explicit.yaml")
	writeTestConfig(t, primaryConfig, "max-time: [\n")
	writeTestConfig(t, explicitConfig, "max-time: 1h\n")

	flagSet := NewFlagSet()
	commonFlags := flagSet.AddCommonFlags()
	stopMaxTimeHandler(t, commonFlags)
	flagSet.SetConfigFilePath(primaryConfig)

	if err := flagSet.Parse(""); err == nil {
		t.Fatal("Parse succeeded with malformed config")
	}
	if commonFlags.maxTimeTimer != nil {
		t.Fatal("max-time timer started with a zero duration")
	}

	if err := flagSet.MergeConfigFile(explicitConfig); err != nil {
		t.Fatalf("MergeConfigFile failed: %v", err)
	}
	if commonFlags.maxTimeTimer == nil {
		t.Fatal("max-time timer was not started by config merge after Parse error")
	}
	if commonFlags.activeMaxTime != time.Hour {
		t.Fatalf("active max-time = %v, want %v", commonFlags.activeMaxTime, time.Hour)
	}
}

func TestConfigMergeBeforeParseDoesNotStartMaxTimeHandler(t *testing.T) {
	configFile := filepath.Join(t.TempDir(), "config.yaml")
	writeTestConfig(t, configFile, "max-time: 1h\n")

	flagSet := NewFlagSet()
	commonFlags := flagSet.AddCommonFlags()
	stopMaxTimeHandler(t, commonFlags)

	if err := flagSet.MergeConfigFile(configFile); err != nil {
		t.Fatalf("MergeConfigFile failed: %v", err)
	}
	if commonFlags.maxTimeTimer != nil {
		t.Fatal("config merge started max-time handler before Parse")
	}
}

func stopMaxTimeHandler(t *testing.T, commonFlags *CommonFlags) {
	t.Helper()
	t.Cleanup(func() {
		if commonFlags.maxTimeTimer != nil {
			commonFlags.maxTimeTimer.Stop()
		}
	})
}
