package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Sakura-501/XSStrike-go/internal/options"
	"github.com/Sakura-501/XSStrike-go/internal/ui"
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
}
