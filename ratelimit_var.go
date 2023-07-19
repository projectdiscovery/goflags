package goflags

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	stringsutil "github.com/projectdiscovery/utils/strings"
	timeutil "github.com/projectdiscovery/utils/time"
)

type RateLimit struct {
	MaxCount uint
	Duration time.Duration
}

type RateLimitMap struct {
	kv map[string]RateLimit
}

// Set inserts a value to the map. Format: key=value
func (rateLimitMap *RateLimitMap) Set(value string) error {
	if rateLimitMap.kv == nil {
		rateLimitMap.kv = make(map[string]RateLimit)
	}
	var k, v string
	if idxSep := strings.Index(value, kvSep); idxSep > 0 {
		k = value[:idxSep]
		v = value[idxSep+1:]
	}
	// note:
	// - inserting multiple times the same key will override the previous value
	// - empty string is legitimate value

	if k != "" {
		rateLimit, err := parseRateLimit(v)
		if err != nil {
			return err
		}
		rateLimitMap.kv[k] = rateLimit
	}
	return nil
}

// Del removes the specified key
func (rateLimitMap *RateLimitMap) Del(key string) error {
	if rateLimitMap.kv == nil {
		return errors.New("empty runtime map")
	}
	delete(rateLimitMap.kv, key)
	return nil
}

// IsEmpty specifies if the underlying map is empty
func (rateLimitMap *RateLimitMap) IsEmpty() bool {
	return rateLimitMap.kv == nil || len(rateLimitMap.kv) == 0
}

// AsMap returns the internal map as reference - changes are allowed
func (rateLimitMap *RateLimitMap) AsMap() map[string]RateLimit {
	return rateLimitMap.kv
}

func (rateLimitMap RateLimitMap) String() string {
	defaultBuilder := &strings.Builder{}
	defaultBuilder.WriteString("{")

	var items string
	for k, v := range rateLimitMap.kv {
		items += fmt.Sprintf("\"%s\"=\"%s\"%s", k, v.Duration.String(), kvSep)
	}
	defaultBuilder.WriteString(stringsutil.TrimSuffixAny(items, ",", "="))
	defaultBuilder.WriteString("}")
	return defaultBuilder.String()
}

// RateLimitMapVar adds a ratelimit flag with a longname
func (flagSet *FlagSet) RateLimitMapVar(field *RateLimitMap, long string, defaultValue []string, usage string) *FlagData {
	return flagSet.RateLimitMapVarP(field, long, "", defaultValue, usage)
}

// RateLimitMapVarP adds a ratelimit flag with a short name and long name.
// It is equivalent to RateLimitMapVar, and also allows specifying ratelimits in days (e.g., "hackertarget=2/d" 2 requests per day, which is equivalent to 24h).
func (flagSet *FlagSet) RateLimitMapVarP(field *RateLimitMap, long, short string, defaultValue []string, usage string) *FlagData {
	if field == nil {
		panic(fmt.Errorf("field cannot be nil for flag -%v", long))
	}

	for _, item := range defaultValue {
		if err := field.Set(item); err != nil {
			panic(fmt.Errorf("failed to set default value for flag -%v: %v", long, err))
		}
	}

	flagData := &FlagData{
		usage:        usage,
		long:         long,
		defaultValue: defaultValue,
		skipMarshal:  true,
	}
	if short != "" {
		flagData.short = short
		flagSet.CommandLine.Var(field, short, usage)
		flagSet.flagKeys.Set(short, flagData)
	}
	flagSet.CommandLine.Var(field, long, usage)
	flagSet.flagKeys.Set(long, flagData)
	return flagData
}

func parseRateLimit(s string) (RateLimit, error) {
	sArr := strings.Split(s, "/")

	if len(sArr) < 2 {
		return RateLimit{}, errors.New("parse error")
	}

	maxCount, err := strconv.ParseUint(sArr[0], 10, 64)
	if err != nil {
		return RateLimit{}, errors.New("parse error: " + err.Error())
	}
	duration, err := timeutil.ParseDuration("1" + sArr[1])
	if err != nil {
		return RateLimit{}, errors.New("parse error: " + err.Error())
	}
	return RateLimit{MaxCount: uint(maxCount), Duration: duration}, nil
}
