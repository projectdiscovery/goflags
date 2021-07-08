package goflags

import (
	"fmt"
	"github.com/pkg/errors"
	"strings"
)

// StringSlice is a slice of strings
type StringSlice []string

func (stringSlice *StringSlice) String() string {
	return strings.Join(*stringSlice, " ")
}

// Set appends a value to the string slice.
func (stringSlice *StringSlice) Set(value string) error {
	*stringSlice = append(*stringSlice, value)
	return nil
}

type Severities []Severity

func (severities *Severities) String() string {
	return strings.Join(severities.ToStringArray(), " ")
}

func (severities *Severities) Set(value string) error {
	computedSeverity, err := toSeverity(value)
	if err != nil {
		return errors.New(fmt.Sprintf("'%s' is not a valid severity!", value))
	}
	*severities = append(*severities, computedSeverity)
	return nil
}

func (severities *Severities) ToStringArray() []string {
	var result []string
	for _, severity := range *severities {
		result = append(result, severity.String())
	}
	return result
}
