package goflags

var optionMap map[*StringSlice]Options

func init() {
	optionMap = make(map[*StringSlice]Options)
}

// StringSlice is a slice of strings
type StringSlice struct {
	Value   []string
	Default bool
}

// Set appends a value to the string slice.
func (stringSlice *StringSlice) Set(value string) error {
	if stringSlice.Default {
		stringSlice.Value = []string{}
		stringSlice.Default = false
	}
	option, ok := optionMap[stringSlice]
	if !ok {
		option = StringSliceOptions
	}
	values, err := ToStringSlice(value, option)
	if err != nil {
		return err
	}
	stringSlice.Value = append(stringSlice.Value, values...)
	return nil
}

func (stringSlice StringSlice) String() string {
	return ToString(stringSlice.Value)
}
