package goflags

import (
	"fmt"
	"sort"
	"strings"
)

type EnumVariable int8

func (e *EnumVariable) String() string {
	return fmt.Sprintf("%v", *e)
}

type AllowdTypes map[string]EnumVariable

func (a AllowdTypes) String() string {
	// Sorted, because ranging over the map gives a different order from one
	// call to the next, which makes both the help output and the "allowed
	// values are" error text unstable.
	keys := make([]string, 0, len(a))
	for k := range a {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return strings.Join(keys, ", ")
}

type EnumVar struct {
	allowedTypes AllowdTypes
	value        *string
}

func (e *EnumVar) String() string {
	if e.value != nil {
		return *e.value
	}
	return ""
}

func (e *EnumVar) Set(value string) error {
	_, ok := e.allowedTypes[value]
	if !ok {
		return fmt.Errorf("allowed values are %v", e.allowedTypes.String())
	}
	*e.value = value
	return nil
}
