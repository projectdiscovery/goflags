package goflags

import (
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

// StringSliceRaw is a slice of strings
type StringSliceRaw []string

func (stringSlice *StringSliceRaw) String() string {
	return strings.Join(*stringSlice, " ")
}

// Set appends a value to the string slice.
func (stringSlice *StringSliceRaw) Set(value string) error {
	*stringSlice = append(*stringSlice, value)
	return nil
}

type StringSlice []string

func (stringSlice *StringSlice) String() string {
	return strings.Join(*stringSlice, " ")
}

// SetI appends a value to the string slice
func (StringSlice *StringSlice) Set(value string) error {
	value = strings.ToLower(value)
	slice, err := ToStringSlice(value)
	if err != nil {
		return err
	}

	*StringSlice = append(*StringSlice, slice...)
	return nil
}

type PathSlice []string

func (stringSlice *PathSlice) String() string {
	return strings.Join(*stringSlice, " ")
}

// Set appends a value to the string slice.
func (stringSlice *PathSlice) Set(value string) error {
	slice, err := ToStringSlice(value)
	if err != nil {
		return err
	}

	*stringSlice = append(*stringSlice, slice...)
	return nil
}

var multiValueValidator = regexp.MustCompile("('[^',]+?,.*?')|(\"[^\",]+?,.*?\")|(`[^,]+?,.*?`)")

func ToStringSlice(value string) ([]string, error) {
	if multiValueValidator.FindString(value) != "" {
		return nil, errors.New("Supported values are: value1,value2 etc")
	}

	if strings.Contains(value, ",") {
		var result []string
		slices := strings.Split(value, ",")
		result = make([]string, 0, len(slices))
		for _, slice := range slices {
			result = append(result, strings.TrimSpace(strings.Trim(strings.TrimSpace(slice), "\"'`")))
		}
		return result, nil
	}

	return []string{value}, nil
}
