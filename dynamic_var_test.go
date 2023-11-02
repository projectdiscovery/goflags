package goflags

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDynamicVar(t *testing.T) {
	t.Run("with bool as type", func(t *testing.T) {
		var b bool
		flagSet := NewFlagSet()
		flagSet.CreateGroup("Option", "option",
			flagSet.DynamicVar(&b, "kev", false, "kev with or without value"),
		)
		os.Args = []string{
			os.Args[0],
			"-kev=false",
		}
		err := flagSet.Parse()
		assert.Nil(t, err)
		assert.Equal(t, false, b)
		tearDown(t.Name())
	})

	t.Run("without value for int as type", func(t *testing.T) {
		var i int
		flagSet := NewFlagSet()
		flagSet.CreateGroup("Option", "option",
			flagSet.DynamicVarP(&i, "concurrency", "c", 25, "concurrency for the process"),
		)
		os.Args = []string{
			os.Args[0],
			"-c",
		}
		err := flagSet.Parse()
		assert.Nil(t, err)
		assert.Equal(t, 25, i)
		tearDown(t.Name())
	})

	t.Run("with value for int as type", func(t *testing.T) {
		var i int
		flagSet := NewFlagSet()
		flagSet.CreateGroup("Option", "option",
			flagSet.DynamicVarP(&i, "concurrency", "c", 25, "concurrency for the process"),
		)
		os.Args = []string{
			os.Args[0],
			"-c=100",
		}
		err := flagSet.Parse()
		assert.Nil(t, err)
		assert.Equal(t, 100, i)
		tearDown(t.Name())
	})

	t.Run("with value for float64 as type", func(t *testing.T) {
		var f float64
		flagSet := NewFlagSet()
		flagSet.CreateGroup("Option", "option",
			flagSet.DynamicVarP(&f, "percentage", "p", 0.0, "percentage for the process"),
		)
		os.Args = []string{
			os.Args[0],
			"-p=100.0",
		}
		err := flagSet.Parse()
		assert.Nil(t, err)
		assert.Equal(t, 100.0, f)
		tearDown(t.Name())
	})

	t.Run("with value for string as type", func(t *testing.T) {
		var s string
		flagSet := NewFlagSet()
		flagSet.CreateGroup("Option", "option",
			flagSet.DynamicVarP(&s, "name", "n", "", "name of the user"),
		)
		os.Args = []string{
			os.Args[0],
			"-n=test",
		}
		err := flagSet.Parse()
		assert.Nil(t, err)
		assert.Equal(t, "test", s)
		tearDown(t.Name())
	})

	t.Run("with value for string slice as type", func(t *testing.T) {
		var s []string
		flagSet := NewFlagSet()
		flagSet.CreateGroup("Option", "option",
			flagSet.DynamicVarP(&s, "name", "n", []string{}, "name of the user"),
		)
		os.Args = []string{
			os.Args[0],
			"-n=test1,test2",
		}
		err := flagSet.Parse()
		assert.Nil(t, err)
		assert.Equal(t, []string{"test1", "test2"}, s)
		tearDown(t.Name())
	})

}
