package goflags

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRateLimitMapVar(t *testing.T) {

	t.Run("default-value", func(t *testing.T) {
		var rateLimitMap RateLimitMap
		flagSet := NewFlagSet()
		flagSet.CreateGroup("Config", "Config",
			flagSet.RateLimitMapVarP(&rateLimitMap, "rate-limits", "rls", []string{"hackertarget=1/ms"}, "rate limits", CommaSeparatedStringSliceOptions),
		)
		os.Args = []string{
			os.Args[0],
		}
		err := flagSet.Parse()
		assert.Nil(t, err)
		assert.Equal(t, RateLimit{MaxCount: 1, Duration: time.Millisecond}, rateLimitMap.AsMap()["hackertarget"])
		tearDown(t.Name())
	})

	t.Run("multiple-default-value", func(t *testing.T) {
		var rateLimitMap RateLimitMap
		flagSet := NewFlagSet()
		flagSet.CreateGroup("Config", "Config",
			flagSet.RateLimitMapVarP(&rateLimitMap, "rate-limits", "rls", []string{"hackertarget=1/s,github=1/ms"}, "rate limits", CommaSeparatedStringSliceOptions),
		)
		os.Args = []string{
			os.Args[0],
		}
		err := flagSet.Parse()
		assert.Nil(t, err)
		assert.Equal(t, RateLimit{MaxCount: 1, Duration: time.Second}, rateLimitMap.AsMap()["hackertarget"])
		assert.Equal(t, RateLimit{MaxCount: 1, Duration: time.Millisecond}, rateLimitMap.AsMap()["github"])
		tearDown(t.Name())
	})

	t.Run("valid-rate-limit", func(t *testing.T) {
		var rateLimitMap RateLimitMap
		flagSet := NewFlagSet()
		flagSet.CreateGroup("Config", "Config",
			flagSet.RateLimitMapVarP(&rateLimitMap, "rate-limits", "rls", nil, "rate limits", CommaSeparatedStringSliceOptions),
		)
		os.Args = []string{
			os.Args[0],
			"-rls", "hackertarget=10/m",
		}
		err := flagSet.Parse()
		assert.Nil(t, err)
		assert.Equal(t, RateLimit{MaxCount: 10, Duration: time.Minute}, rateLimitMap.AsMap()["hackertarget"])

		tearDown(t.Name())
	})

	t.Run("valid-rate-limits", func(t *testing.T) {
		var rateLimitMap RateLimitMap
		flagSet := NewFlagSet()
		flagSet.CreateGroup("Config", "Config",
			flagSet.RateLimitMapVarP(&rateLimitMap, "rate-limits", "rls", nil, "rate limits", CommaSeparatedStringSliceOptions),
		)
		os.Args = []string{
			os.Args[0],
			"-rls", "hackertarget=1/s,github=1/ms",
		}
		err := flagSet.Parse()
		assert.Nil(t, err)
		assert.Equal(t, RateLimit{MaxCount: 1, Duration: time.Second}, rateLimitMap.AsMap()["hackertarget"])
		assert.Equal(t, RateLimit{MaxCount: 1, Duration: time.Millisecond}, rateLimitMap.AsMap()["github"])
		tearDown(t.Name())
	})

	t.Run("without-unit", func(t *testing.T) {
		var rateLimitMap RateLimitMap
		err := rateLimitMap.Set("hackertarget=1")
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "parse error")
		tearDown(t.Name())
	})

	t.Run("invalid-unit", func(t *testing.T) {
		var rateLimitMap RateLimitMap
		err := rateLimitMap.Set("hackertarget=1/x")
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "parse error")
		tearDown(t.Name())
	})
}
