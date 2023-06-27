package main

import (
	"fmt"
	"log"
	"time"

	"github.com/projectdiscovery/goflags"
)

type Options struct {
	name     string
	Email    goflags.StringSlice
	Phone    string
	Address  goflags.StringSlice
	fileSize goflags.Size
	duration time.Duration
}

func main() {
	testOptions := &Options{}
	CheckUpdate := func() {
		fmt.Println("checking if new version is available")
		fmt.Println("updating tool....")
	}

	flagSet := goflags.NewFlagSet()
	flagSet.CreateGroup("info", "Info",
		flagSet.StringVarP(&testOptions.name, "name", "n", "", "name of the user"),
		flagSet.StringSliceVarP(&testOptions.Email, "email", "e", nil, "email of the user", goflags.CommaSeparatedStringSliceOptions),
	)
	flagSet.CreateGroup("additional", "Additional",
		flagSet.StringVarP(&testOptions.Phone, "phone", "ph", "", "phone of the user"),
		flagSet.StringSliceVarP(&testOptions.Address, "address", "add", nil, "address of the user", goflags.StringSliceOptions),
		flagSet.CallbackVarP(CheckUpdate, "update", "ut", "update this tool to latest version"),
		flagSet.SizeVarP(&testOptions.fileSize, "max-size", "ms", "", "max file size"),
		flagSet.DurationVar(&testOptions.duration, "timeout", time.Hour, "timeout"),
	)

	if err := flagSet.Parse(); err != nil {
		log.Fatal(err)
	}
}
