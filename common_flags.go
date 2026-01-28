package goflags

import "time"

// CommonFlags contains common flags shared across ProjectDiscovery tools.
// These flags provide consistent behavior across all tools in the ecosystem.
type CommonFlags struct {
	// MaxTime is the maximum duration for the entire execution.
	// When this duration is reached, an interrupt signal is sent to trigger graceful termination.
	// Example values: "1h", "30m", "1h30m", "2h45m30s"
	MaxTime time.Duration
}

// AddCommonFlags registers common flags to the flagset and returns a CommonFlags struct.
// The handlers are automatically started after Parse() is called.
//
// Usage:
//
//	flagSet := goflags.NewFlagSet()
//	flagSet.AddCommonFlags()
//	flagSet.Parse()
func (flagSet *FlagSet) AddCommonFlags() *CommonFlags {
	cf := &CommonFlags{}

	flagSet.CreateGroup("common", "Common",
		flagSet.DurationVarP(&cf.MaxTime, "max-time", "mt", 0, "maximum time to run before automatic termination (e.g., 1h, 30m)"),
	)

	flagSet.commonFlags = cf
	return cf
}

// startCommonFlagsHandlers is called by Parse() to start handlers.
func (flagSet *FlagSet) startCommonFlagsHandlers() {
	if flagSet.commonFlags != nil {
		flagSet.commonFlags.startMaxTimeHandler()
	}
}

// startMaxTimeHandler starts the max time handler if MaxTime is set.
func (cf *CommonFlags) startMaxTimeHandler() {
	if cf.MaxTime > 0 {
		go func() {
			<-time.After(cf.MaxTime)
			sendInterrupt()
		}()
	}
}
