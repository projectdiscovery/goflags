package goflags

import (
	"os"
	"os/signal"
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

func TestMaxTimeGracefulShutdown(t *testing.T) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	defer signal.Stop(sigChan)

	flagSet := NewFlagSet()
	flagSet.AddCommonFlags()
	_ = flagSet.Parse("-mt", "100ms")

	select {
	case <-sigChan:
	case <-time.After(500 * time.Millisecond):
		t.Error("Expected SIGINT within 500ms")
	}
}
