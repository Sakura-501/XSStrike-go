package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Sakura-501/XSStrike-go/internal/config"
	"github.com/Sakura-501/XSStrike-go/internal/options"
	"github.com/Sakura-501/XSStrike-go/internal/payload"
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

	if opts.URL == "" && !opts.Fuzzer {
		fs.Usage()
		return
	}

	fmt.Print(ui.Banner())
	fmt.Printf("Runtime defaults -> timeout=%ds threads=%d delay=%ds\n", opts.Timeout, opts.ThreadCount, opts.Delay)

	if opts.Fuzzer {
		runFuzzer(opts)
		return
	}

	fmt.Printf("Target: %s\n", opts.URL)

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

func runFuzzer(opts *options.Options) {
	vectors := payload.GenerateVectors(payload.GeneratorInput{
		Fillings:      config.DefaultFillings,
		EFillings:     config.DefaultEFillings,
		LFillings:     config.DefaultLFillings,
		EventHandlers: config.DefaultEventHandlers,
		Tags:          config.DefaultTags,
		Functions:     config.DefaultFunctions,
		Ends:          config.DefaultEnds,
		Bait:          config.XSSChecker,
	}, nil)

	limit := opts.Limit
	if limit < 0 {
		limit = 0
	}
	if limit > len(vectors) {
		limit = len(vectors)
	}

	fmt.Printf("Fuzzer payloads: %d\n", len(config.DefaultFuzzes))
	fmt.Printf("Generated context vectors: %d\n", len(vectors))
	for i := 0; i < limit; i++ {
		fmt.Printf("[%d] %s\n", i+1, vectors[i])
	}
}
