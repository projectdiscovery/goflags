package goflags

import (
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

// StringSlice is a slice of strings
type StringSlice []string

func (stringSlice *StringSlice) String() string {
	return strings.Join(*stringSlice, " ")
}

// Set appends a value to the string slice.
func (stringSlice *StringSlice) Set(value string) error {
	if slice, err := ToStringSlice(value); err != nil {
		return err
	} else {
		*stringSlice = append(*stringSlice, slice...)
		return nil
	}
}

// SetI appends a value to the string slice converting it to lowercase
func (stringSlice *StringSlice) SetToLower(value string) error {
	if slice, err := ToStringSlice(value); err != nil {
		return err
	} else {
		slice = sliceToLower(slice)
		*stringSlice = append(*stringSlice, slice...)
		return nil
	}
}

var multiValueValidator = regexp.MustCompile("('[^',]+?,.*?')|(\"[^\",]+?,.*?\")|(`[^,]+?,.*?`)")

func ToStringSlice(value string) ([]string, error) {
	if multiValueValidator.FindString(value) != "" {
		return nil, errors.New("Supported values are: value1,value2 etc")
	}

	var result []string
	if strings.Contains(value, ",") {
		slices := strings.Split(value, ",")
		result = make([]string, 0, len(slices))
		for _, slice := range slices {
			result = append(result, strings.TrimSpace(strings.Trim(strings.TrimSpace(slice), "\"'`")))
		}
	} else {
		result = []string{value}
	}
	return result, nil
}

func sliceToLower(values []string) []string {
	for i := range values {
		values[i] = strings.ToLower(values[i])
	}
	return values
}
