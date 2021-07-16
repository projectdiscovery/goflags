package goflags

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_toStringSlice(t *testing.T) {
	expected := []string{"aa", "bb", "cc"}
	values := map[string]interface{}{
		"aa,bb,cc":                       expected,
		"  aa, bb,  cc   ":               expected,
		"  `aa`, 'bb',  \"cc\"   ":       expected,
		"  `aa`, bb,  \"cc\"   ":         expected,
		"  `aa, bb,  cc\"   ":            expected,
		"  \"aa\", bb,  cc\"   ":         expected,
		"\n  aa, \tbb,  cc\r   ":         expected,
		"\n  \"aa', \t`bb',  \"cc\r`   ": nil,
		"\"\n  aa', `\tbb',  \"cc\r`   ": nil,
		"'aa,bb,cc'":                     nil,
		"`aa,bb,cc`":                     nil,
		"\"aa,bb,cc\"":                   nil,
	}
	for input, expectedValue := range values {
		t.Run(input, func(t *testing.T) {
			slice, err := ToStringSlice(input)
			if expectedValue == nil {
				if err == nil {
					t.Fail()
				}
			} else {
				assert.Equal(t, expectedValue, slice)
			}
		})
	}
}
