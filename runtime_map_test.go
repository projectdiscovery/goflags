package goflags

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRuntimeMap(t *testing.T) {
	data := &RuntimeMap{}
	err := data.Set("variable=value")
	require.NoError(t, err, "could not set key-value")

	returned := data.AsMap()["variable"]
	require.Equal(t, "value", returned, "could not get correct return")
}
