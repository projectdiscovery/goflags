package goflags

import (
	"errors"
	"strconv"
)

var errParse = errors.New("parse error")

type (
	CallbackFunc func()

	CallbackVar struct {
		option  *bool
		visited bool
		action  CallbackFunc
	}
)

func newCallbackVar(option *bool, action CallbackFunc) *CallbackVar {
	return &CallbackVar{option: option, action: action}
}

func (c *CallbackVar) Set(s string) error {
	v, err := strconv.ParseBool(s)
	if err != nil {
		err = errParse
	}
	c.option = &v
	return err
}

func (c *CallbackVar) IsBoolFlag() bool {
	return true
}

func (c *CallbackVar) String() string {
	if c.option != nil {
		return strconv.FormatBool(bool(*c.option))
	}
	return "false"
}

// CallbackVar adds a Callback flag with a longname
func (flagSet *FlagSet) CallbackVar(field *bool, long string, defaultValue CallbackFunc, usage string) *FlagData {
	return flagSet.CallbackVarP(field, long, "", defaultValue, usage)
}

// CallbackVarP adds a Callback flag with a shortname and longname
func (flagSet *FlagSet) CallbackVarP(field *bool, long, short string, defaultValue CallbackFunc, usage string) *FlagData {
	if defaultValue == nil {
		return &FlagData{}
	}

	flagData := &FlagData{
		usage:        usage,
		long:         long,
		defaultValue: strconv.FormatBool(*field),
		field:        newCallbackVar(field, defaultValue),
	}
	if short != "" {
		flagData.short = short
		flagSet.CommandLine.Var(flagData.field, short, usage)
		flagSet.flagKeys.Set(short, flagData)
	}
	flagSet.CommandLine.Var(flagData.field, long, usage)
	flagSet.flagKeys.Set(long, flagData)
	return flagData
}
