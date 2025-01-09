package goflags

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRuntimeMap(t *testing.T) {
	data := &RuntimeMap{}
	err := data.Set("variable=value")
	require.NoError(t, err, "could not set key-value")

	returned := data.AsMap()["variable"]
	require.Equal(t, "value", returned, "could not get correct return")

	t.Run("file", func(t *testing.T) {
		sb := &strings.Builder{}
		sb.WriteString("variable=value\n")
		sb.WriteString("variable2=value2\n")
		tempFile, err := os.CreateTemp(t.TempDir(), "test")
		require.NoError(t, err, "could not create temp file")
		defer tempFile.Close()
		_, err = tempFile.WriteString(sb.String())
		require.NoError(t, err, "could not write to temp file")
		data2 := &RuntimeMap{}
		err = data2.Set(tempFile.Name())
		require.NoError(t, err, "could not set key-value")
		require.Equal(t, 2, len(data2.AsMap()), "could not get correct number of key-values")
		require.Equal(t, "value", data2.AsMap()["variable"], "could not get correct value")
		require.Equal(t, "value2", data2.AsMap()["variable2"], "could not get correct value")
	})
}
