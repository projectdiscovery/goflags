package goflags

import (
	"fmt"
	"testing"
)

func TestGenerateDefaultConfig(t *testing.T) {
	f := New()

	var data string
	var data2 StringSlice
	f.StringVar(&data, "test", "test-default-value", "Default value for a test flag example")
	f.StringSliceVar(&data2, "slice", []string{"item1", "item2"}, "String slice flag example value")
	f.StringVarP(&data, "s", "s-long", "shortlong", "String shortlong flag example value")
	fmt.Printf("%s\n", string(f.generateDefaultConfig()))
}
