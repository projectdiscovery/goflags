package goflags

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPortType(t *testing.T) {
	port := &Port{}
	_ = port.Set("21-25,80,TCP:443")
	require.ElementsMatch(t, port.AsPorts(), []int{21, 22, 23, 24, 25, 80, 443}, "could not get correct ports")

	t.Run("comma-separated", func(t *testing.T) {
		port := &Port{}
		_ = port.Set("80,443")
		require.ElementsMatch(t, port.AsPorts(), []int{80, 443}, "could not get correct ports")
	})
	t.Run("dash", func(t *testing.T) {
		port := &Port{}
		_ = port.Set("21-25")
		require.ElementsMatch(t, port.AsPorts(), []int{21, 22, 23, 24, 25}, "could not get correct ports")
	})
	t.Run("dash-suffix", func(t *testing.T) {
		port := &Port{}
		_ = port.Set("1-")
		require.Len(t, port.AsPorts(), 65535, "could not get correct ports")
	})
	t.Run("full", func(t *testing.T) {
		port := &Port{}
		_ = port.Set("full")
		require.Len(t, port.AsPorts(), 65535, "could not get correct ports")

		port = &Port{}
		_ = port.Set("*")
		require.Len(t, port.AsPorts(), 65535, "could not get correct ports")
	})
	t.Run("top-xxx", func(t *testing.T) {
		port := &Port{}
		_ = port.Set("top-100")
		require.Len(t, port.AsPorts(), 100, "could not get correct ports")

		port = &Port{}
		_ = port.Set("top-1000")
		require.Len(t, port.AsPorts(), 1000, "could not get correct ports")
	})
	t.Run("services", func(t *testing.T) {
		port := &Port{}
		_ = port.Set("http,ftp")
		require.ElementsMatch(t, port.AsPorts(), []int{80, 8008, 21}, "could not get correct ports")
	})
	t.Run("services-wildcard", func(t *testing.T) {
		port := &Port{}
		_ = port.Set("ftp*")
		require.ElementsMatch(t, port.AsPorts(), []int{21, 20, 989, 990, 574, 8021}, "could not get correct ports")
	})
	t.Run("colon", func(t *testing.T) {
		port := &Port{}
		_ = port.Set("TCP:443,UDP:53")
		require.ElementsMatch(t, port.AsPorts(), []int{443, 53}, "could not get correct ports")
	})
}
