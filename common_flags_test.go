package goflags

import (
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
	called := make(chan struct{})
	originalSignalSelf := signalSelf
	signalSelf = func() {
		close(called)
	}
	defer func() { signalSelf = originalSignalSelf }()

	flagSet := NewFlagSet()
	flagSet.AddCommonFlags()
	_ = flagSet.Parse("-mt", "100ms")

	select {
	case <-called:
	case <-time.After(500 * time.Millisecond):
		t.Error("Expected signalSelf to be called within 500ms")
	}
}
