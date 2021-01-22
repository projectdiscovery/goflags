package goflags

import "strings"

// -- string Value
type stringValue string

func newStringValue(val string, p *string) *stringValue {
	*p = val
	return (*stringValue)(p)
}

func (s *stringValue) Set(val string) error {
	*s = stringValue(val)
	return nil
}

func (s *stringValue) Get() interface{} { return string(*s) }

func (s *stringValue) String() string { return string(*s) }

// StringSlice is a slice of strings
type StringSlice []string

func (i *StringSlice) String() string {
	return strings.Join(*i, " ")
}

// Set appends a value to the string slice.
func (i *StringSlice) Set(value string) error {
	*i = append(*i, value)
	return nil
}
