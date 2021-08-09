package goflags

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNormalizedStringSlicePositive(t *testing.T) {
	expectedABC := []string{"aa", "bb", "cc"}
	expectedFilePath := []string{"/root/home/file0"}

	slices := map[string][]string{
		"aa,bb,cc":                 expectedABC,
		"  aa, bb,  cc   ":         expectedABC,
		"  `aa`, 'bb',  \"cc\"   ": expectedABC,
		"  `aa`, bb,  \"cc\"   ":   expectedABC,
		"  `aa, bb,  cc\"   ":      expectedABC,
		"  \"aa\", bb,  cc\"   ":   expectedABC,
		"\n  aa, \tbb,  cc\r   ":   expectedABC,

		"\"value1\",value,'value3'": {"value1", "value", "value3"},
		"\"value1\",VALUE,'value3'": {"value1", "value", "value3"},

		"\"/root/home/file0\"":       expectedFilePath,
		"'/root/home/file0'":         expectedFilePath,
		"`/root/home/file0`":         expectedFilePath,
		"\"/root/home/file0\",":      expectedFilePath,
		",\"/root/home/file0\",":     expectedFilePath,
		",\"/root/home/file0\"":      expectedFilePath,
		",,\"/root/home/file0\"":     expectedFilePath,
		"\"\",,\"/root/home/file0\"": expectedFilePath,
		"\" \",\"/root/home/file0\"": expectedFilePath,
		"\"/root/home/file0\",\"\"":  expectedFilePath,
		"/root/home/file0":           expectedFilePath,

		"\"/root/home/file2\",\"/root/home/file3\"":             {"/root/home/file2", "/root/home/file3"},
		"/root/home/file4,/root/home/file5":                     {"/root/home/file4", "/root/home/file5"},
		"\"/root/home/file4,/root/home/file5\"":                 {"/root/home/file4,/root/home/file5"},
		"\"/root/home/file6\",/root/home/file7":                 {"/root/home/file6", "/root/home/file7"},
		"\"c:\\my files\\bug,bounty\"":                          {"c:\\my files\\bug,bounty"},
		"\"c:\\my files\\bug,bounty\",c:\\my_files\\bug bounty": {"c:\\my files\\bug,bounty", "c:\\my_files\\bug bounty"},
	}

	for value, expected := range slices {
		result, err := ToNormalizedStringSlice(value)
		fmt.Println(result)
		assert.Nil(t, err)
		assert.Equal(t, result, expected)
	}
}

func TestNormalizedStringSliceNegative(t *testing.T) {
	slices := []string{
		"\"/root/home/file0",
		"'/root/home/file0",
		"`/root/home/file0",
		"\"/root/home/file0'",
		"\"/root/home/file0`",
	}

	for _, value := range slices {
		result, err := ToNormalizedStringSlice(value)
		assert.Nil(t, result)
		assert.NotNil(t, err)
	}
}
