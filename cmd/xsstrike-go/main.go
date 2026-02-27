package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Sakura-501/XSStrike-go/internal/options"
	"github.com/Sakura-501/XSStrike-go/internal/ui"
	"github.com/Sakura-501/XSStrike-go/internal/utils"
	"github.com/Sakura-501/XSStrike-go/internal/version"
)

func main() {
	fs := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	fs.SetOutput(os.Stdout)
	fs.Usage = func() {
		fmt.Print(ui.Banner())
		fmt.Println("Usage: xsstrike-go [options]")
		fs.PrintDefaults()
	}

	opts, err := options.Parse(fs, os.Args[1:])
	if err != nil {
		os.Exit(2)
	}

	if opts.Version {
		fmt.Printf("%s %s\n", version.AppName, version.Version)
		return
	}

	if opts.URL == "" {
		fs.Usage()
		return
	}

	fmt.Print(ui.Banner())
	fmt.Printf("Target: %s\n", opts.URL)
	fmt.Printf("Runtime defaults -> timeout=%ds threads=%d delay=%ds\n", opts.Timeout, opts.ThreadCount, opts.Delay)

	headers := utils.ExtractHeaders(opts.HeadersRaw)
	if len(headers) > 0 {
		fmt.Printf("Custom headers parsed: %d\n", len(headers))
	}

	params := utils.ParseParams(opts.URL, opts.Data, opts.JSON)
	if len(params) == 0 {
		fmt.Println("No parameters found")
		return
	}
	fmt.Printf("Parsed parameters: %d\n", len(params))
}
