package goflags

import (
	"os"
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

func TestPortSetMixedKeywordLists(t *testing.T) {
	t.Run("top-100,31337", func(t *testing.T) {
		keywordOnly := &Port{}
		_ = keywordOnly.Set("top-100")
		require.NotContains(t, keywordOnly.AsPorts(), 31337, "31337 must be outside top-100")
		require.Contains(t, keywordOnly.AsPorts(), 80)

		port := &Port{}
		_ = port.Set("top-100,31337")
		require.Len(t, port.AsPorts(), 101, "could not get correct ports")
		require.Contains(t, port.AsPorts(), 31337)
		require.Contains(t, port.AsPorts(), 80)
	})
	t.Run("80,full", func(t *testing.T) {
		port := &Port{}
		_ = port.Set("80,full")
		require.Len(t, port.AsPorts(), 65535, "could not get correct ports")
		require.Contains(t, port.AsPorts(), 80)
	})
	t.Run("22,*", func(t *testing.T) {
		port := &Port{}
		_ = port.Set("22,*")
		require.Len(t, port.AsPorts(), 65535, "could not get correct ports")
		require.Contains(t, port.AsPorts(), 22)
	})
	t.Run("80,-", func(t *testing.T) {
		port := &Port{}
		_ = port.Set("80,-")
		require.Len(t, port.AsPorts(), 65535, "could not get correct ports")
		require.Contains(t, port.AsPorts(), 80)
	})
	t.Run("top-1000,54321", func(t *testing.T) {
		keywordOnly := &Port{}
		_ = keywordOnly.Set("top-1000")
		require.NotContains(t, keywordOnly.AsPorts(), 54321, "54321 must be outside top-1000")

		port := &Port{}
		_ = port.Set("top-1000,54321")
		require.Len(t, port.AsPorts(), 1001, "could not get correct ports")
		require.Contains(t, port.AsPorts(), 54321)
	})
}

func TestPortKeywordInListViaFlagSet(t *testing.T) {
	flagSet := NewFlagSet()
	var p Port
	flagSet.PortVarP(&p, "ports", "p", nil, "ports to scan")

	os.Args = []string{
		os.Args[0],
		"-p", "top-100,31337",
	}
	err := flagSet.Parse()
	require.Nil(t, err)
	require.Len(t, p.AsPorts(), 101, "could not get correct ports")
	require.Contains(t, p.AsPorts(), 31337)
	tearDown(t.Name())
}

func BenchmarkPortStringFull(b *testing.B) {
	port := &Port{}
	require.NoError(b, port.Set("full"))

	for b.Loop() {
		if result := port.String(); len(result) == 0 {
			b.Fatal("Port.String returned an empty result")
		}
	}
}
