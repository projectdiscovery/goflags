package goflags

import "fmt"

type (
	UpdateCheckCallback func()
	UpdateCallback      func()

	UpdateGroup struct {
		Options options
		Actions actions
	}

	options struct {
		IsUpdate, DisableUpdate bool
	}

	actions struct {
		UpdateCallback
		UpdateCheckCallback
	}
)

// NewUpdateGroup creates self-update options
func (flagset *FlagSet) NewUpdateGroup(toolName string, uc UpdateCheckCallback, u UpdateCallback) {
	ug := UpdateGroup{
		Options: options{IsUpdate: false, DisableUpdate: false},
		Actions: actions{UpdateCallback: u, UpdateCheckCallback: uc},
	}
	flagset.CreateGroup("update", "Update",
		flagset.BoolVarP(&ug.Options.IsUpdate, "update", "up", false, fmt.Sprintf("update %v to the latest released version", toolName)),
		flagset.BoolVarP(&ug.Options.DisableUpdate, "disable-update-check", "duc", false, "disable automatic update check"),
	)
	flagset.updateGroup = &ug
}
