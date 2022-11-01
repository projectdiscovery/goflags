package main

import (
	"log"

	"github.com/projectdiscovery/goflags"
)

type Options struct {
	name    string
	Email   goflags.StringSlice
	Phone   string
	Address goflags.StringSlice
}

func main() {
	testOptions := &Options{}
	flagSet := goflags.NewFlagSet()
	flagSet.CreateGroup("info", "Info",
		flagSet.StringVarP(&testOptions.name, "name", "n", "", "name of the user"),
		flagSet.StringSliceVarP(&testOptions.Email, "email", "e", nil, "email of the user", goflags.CommaSeparatedStringSliceOptions),
	)
	flagSet.CreateGroup("additional", "Additional",
		flagSet.StringVarP(&testOptions.Phone, "phone", "ph", "", "phone of the user"),
		flagSet.StringSliceVarP(&testOptions.Address, "address", "add", nil, "address of the user", goflags.StringSliceOptions),
	)

	if err := flagSet.Parse(); err != nil {
		log.Fatal(err)
	}
}
