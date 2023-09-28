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
	rls      goflags.RateLimitMap
	severity []string
	// Dynamic
	titleSize int
	target    string
	hashes    []string
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
		flagSet.RateLimitMapVarP(&testOptions.rls, "rate-limits", "rls", nil, "rate limits in format k=v/d i.e hackertarget=10/s", goflags.CommaSeparatedStringSliceOptions),
	)
	flagSet.CreateGroup("additional", "Additional",
		flagSet.StringVarP(&testOptions.Phone, "phone", "ph", "", "phone of the user"),
		flagSet.StringSliceVarP(&testOptions.Address, "address", "add", nil, "address of the user", goflags.StringSliceOptions),
		flagSet.CallbackVarP(CheckUpdate, "update", "ut", "update this tool to latest version"),
		flagSet.SizeVarP(&testOptions.fileSize, "max-size", "ms", "", "max file size"),
		flagSet.DurationVar(&testOptions.duration, "timeout", time.Hour, "timeout"),
		flagSet.EnumSliceVarP(&testOptions.severity, "severity", "s", []goflags.EnumVariable{2}, "severity of the scan", goflags.AllowdTypes{"low": goflags.EnumVariable(0), "medium": goflags.EnumVariable(1), "high": goflags.EnumVariable(2)}),
	)
	flagSet.CreateGroup("Dynmaic", "Dynamic",
		flagSet.DynamicVarP(&testOptions.titleSize, "title", "t", 50, "first N characters of the title"),
		flagSet.DynamicVarP(&testOptions.target, "target", "u", "https://example.com", "target url"),
		flagSet.DynamicVarP(&testOptions.hashes, "hashes", "hs", []string{"md5", "sha1"}, "supported hashes"),
	)
	flagSet.SetCustomHelpText("EXAMPLE USAGE:\ngo run ./examples/basic [OPTIONS]")

	if err := flagSet.Parse(); err != nil {
		log.Fatal(err)
	}

	// ratelimits value is
	if len(testOptions.rls.AsMap()) > 0 {
		fmt.Printf("Got RateLimits: %+v\n", testOptions.rls)
	}

	if len(testOptions.severity) > 0 {
		fmt.Printf("Got Severity: %+v\n", testOptions.severity)
	}

	fmt.Println("Dynamic Values Output")
	fmt.Println("title size:", testOptions.titleSize)
	fmt.Println("target:", testOptions.target)
	fmt.Println("hashes:", testOptions.hashes)
}
