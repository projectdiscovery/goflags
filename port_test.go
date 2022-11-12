package goflags

import (
	"fmt"
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

func TestSetDefaultPortValue(t *testing.T) {
	var data Port
	flagSet := NewFlagSet()
	flagSet.PortVarP(&data, "port", "p", []string{"1,3"}, "Default value for a test flag example")
	err := flagSet.CommandLine.Parse([]string{"-p", "11"})
	require.Nil(t, err)
	fmt.Println(data)
	require.Equal(t, Port{kv: map[int]struct{}{11: {}}}, data, "could not get correct string slice")

	var data2 Port
	flagSet2 := NewFlagSet()
	flagSet2.PortVarP(&data2, "port", "p", []string{"1,3"}, "Default value for a test flag example")
	err = flagSet2.CommandLine.Parse([]string{"-p", "11,12"})
	require.Nil(t, err)
	fmt.Println(data2)
	require.Equal(t, Port{kv: map[int]struct{}{11: {}, 12: {}}}, data2, "could not get correct string slice")

	var data3 Port
	flagSet3 := NewFlagSet()
	flagSet3.PortVarP(&data3, "port", "p", nil, "Default value for a test flag example")
	err = flagSet3.CommandLine.Parse([]string{"-p", "11,12"})
	fmt.Println(data2)
	require.Nil(t, err)
	require.Equal(t, Port{kv: map[int]struct{}{11: {}, 12: {}}}, data3, "could not get correct string slice")

	tearDown(t.Name())
}
