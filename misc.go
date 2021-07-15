package goflags

import (
	"fmt"
	"github.com/pkg/errors"
	"regexp"
	"strings"
)

// StringSlice is a slice of strings
type StringSlice []string

func (stringSlice *StringSlice) String() string {
	return strings.Join(*stringSlice, " ")
}

// Set appends a value to the string slice.
func (stringSlice *StringSlice) Set(value string) error {
	if slice, err := toStringSlice(value); err != nil {
		return err
	} else {
		*stringSlice = append(*stringSlice, slice...)
		return nil
	}
}

type Severities []Severity

func (severities Severities) String() string {
	return strings.Join(severities.ToStringArray(), ", ")
}

func (severities *Severities) Set(value string) error {
	if inputSeverities, err := toStringSlice(value); err != nil {
		return err
	} else {
		for _, inputSeverity := range inputSeverities {
			if err := setSeverity(severities, inputSeverity); err != nil {
				return err
			}
		}
		return nil
	}
}

func setSeverity(severities *Severities, value string) error {
	computedSeverity, err := toSeverity(value)
	if err != nil {
		return errors.New(fmt.Sprintf("'%s' is not a valid severity!", value))
	}
	// TODO change the Severities type to map[Severity]interface{}, where the values are struct{}{}, to "simulates" a "set" data structure
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

var multiValueValidator = regexp.MustCompile("('[^',]+?,.*?')|(\"[^\",]+?,.*?\")|(`[^,]+?,.*?`)")

func toStringSlice(value string) ([]string, error) {
	if multiValueValidator.FindString(value) != "" {
		return nil, errors.New("Supported values are: value1,value2 etc")
	}

	value = strings.ToLower(value)

	var result []string
	if strings.Contains(value, ",") {
		slices := strings.Split(value, ",")
		result = make([]string, len(slices))
		for _, slice := range slices {
			result = append(result, strings.TrimSpace(strings.Trim(strings.TrimSpace(slice), "\"'`")))
		}
	} else {
		result = []string{value}
	}
	return result, nil
}
