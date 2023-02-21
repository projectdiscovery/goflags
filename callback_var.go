package goflags

type (
	CallbackFunc func()

	CallbackVar struct {
		option *bool
		action CallbackFunc
	}
)

// CallbackVar adds a Callback flag with a longname
func (flagSet *FlagSet) CallbackVar(field *bool, long string, defaultValue CallbackFunc, usage string) *FlagData {
	return flagSet.CallbackVarP(field, long, "", defaultValue, usage)
}

// CallbackVarP adds a Callback flag with a shortname and longname
func (flagSet *FlagSet) CallbackVarP(field *bool, long, short string, defaultValue CallbackFunc, usage string) *FlagData {
	if defaultValue == nil {
		return &FlagData{}
	}
	flagSet.callbacks = append(flagSet.callbacks, CallbackVar{option: field, action: defaultValue})
	return flagSet.BoolVarP(field, long, short, *field, usage)
}
