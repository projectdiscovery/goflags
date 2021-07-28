package goflags

import (
	"strings"
)

// StringSlice is a slice of strings
type StringSlice []string

func (stringSlice *StringSlice) String() string {
	return strings.Join(*stringSlice, " ")
}

// Set appends a value to the string slice.
func (stringSlice *StringSlice) Set(value string) error {
	slice := ToStringSlice(value)

	*stringSlice = append(*stringSlice, slice...)
	return nil
}

func ToStringSlice(value string) []string {
	var result []string
	if strings.Contains(value, ",") {
		slices := strings.Split(value, ",")
		result = make([]string, 0, len(slices))
		for _, slice := range slices {
			result = append(result, slice)
		}
	} else {
		result = []string{value}
	}
	return result
}

func (stringSlice *StringSlice) createStringArrayDefaultValue() string {
	defaultBuilder := &strings.Builder{}
	defaultBuilder.WriteString("[")
	for i, k := range *stringSlice {
		defaultBuilder.WriteString("\"")
		defaultBuilder.WriteString(k)
		defaultBuilder.WriteString("\"")
		if i != len(*stringSlice)-1 {
			defaultBuilder.WriteString(", ")
		}
	}
	defaultBuilder.WriteString("]")
	return defaultBuilder.String()
}
