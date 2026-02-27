package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Sakura-501/XSStrike-go/internal/config"
	"github.com/Sakura-501/XSStrike-go/internal/encoder"
	"github.com/Sakura-501/XSStrike-go/internal/files"
	"github.com/Sakura-501/XSStrike-go/internal/options"
	"github.com/Sakura-501/XSStrike-go/internal/payload"
	"github.com/Sakura-501/XSStrike-go/internal/report"
	"github.com/Sakura-501/XSStrike-go/internal/requester"
	"github.com/Sakura-501/XSStrike-go/internal/scan"
	"github.com/Sakura-501/XSStrike-go/internal/state"
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
	state.Global.Set("options", opts)

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
	headers := mergedHeaders(opts.HeadersRaw)
	state.Global.Set("headers", headers)
	state.Global.Set("checkedScripts", map[string]struct{}{})

	if opts.HeadersRaw != "" {
		fmt.Printf("Custom headers parsed: %d\n", len(utils.ExtractHeaders(opts.HeadersRaw)))
	}

	client := requester.New(requester.Config{TimeoutSeconds: opts.Timeout, DelaySeconds: opts.Delay, Proxy: opts.Proxy})
	runner := scan.NewRunner(client)
	scanReport, err := runner.Run(opts.URL, opts.Data, headers, opts.JSON, opts.Encode)
	if err != nil {
		fmt.Printf("Scan error: %v\n", err)
		os.Exit(1)
	}
	state.Global.Set("scanReport", scanReport)
	printScanReport(scanReport)

	if opts.OutputJSON != "" {
		if err := report.WriteJSON(opts.OutputJSON, scanReport); err != nil {
			fmt.Printf("Write output error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Scan report written: %s\n", opts.OutputJSON)
	}
}

func runFuzzer(opts *options.Options) {
	if opts.PayloadFile != "" {
		list, err := resolvePayloadList(opts.PayloadFile)
		if err != nil {
			fmt.Printf("Payload file error: %v\n", err)
			os.Exit(1)
		}
		state.Global.Set("filePayloads", list)
		printPayloadList("File payloads", list, opts)
		return
	}

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
	state.Global.Set("vectors", vectors)

	limit := opts.Limit
	if limit < 0 {
		limit = 0
	}
	if limit > len(vectors) {
		limit = len(vectors)
	}

	fmt.Printf("Fuzzer payloads: %d\n", len(config.DefaultFuzzes))
	fmt.Printf("Generated context vectors: %d\n", len(vectors))
	if opts.Encode == "base64" {
		fmt.Println("Encoding: base64")
	}

	for i := 0; i < limit; i++ {
		current := encoder.Apply(opts.Encode, vectors[i])
		fmt.Printf("[%d] %s\n", i+1, current)
	}
}

func resolvePayloadList(value string) ([]string, error) {
	if value == "default" {
		return config.DefaultPayloads, nil
	}
	return files.ReadLines(value)
}

func printPayloadList(title string, list []string, opts *options.Options) {
	limit := opts.Limit
	if limit < 0 {
		limit = 0
	}
	if limit > len(list) {
		limit = len(list)
	}
	fmt.Printf("%s: %d\n", title, len(list))
	if opts.Encode == "base64" {
		fmt.Println("Encoding: base64")
	}
	for i := 0; i < limit; i++ {
		fmt.Printf("[%d] %s\n", i+1, encoder.Apply(opts.Encode, list[i]))
	}
}

func mergedHeaders(raw string) map[string]string {
	headers := map[string]string{}
	for key, value := range config.DefaultHeaders {
		headers[key] = value
	}
	for key, value := range utils.ExtractHeaders(raw) {
		headers[key] = value
	}
	return headers
}

func printScanReport(scanReport *scan.Report) {
	if scanReport.NoParams {
		fmt.Println("No parameters to test.")
		return
	}

	fmt.Printf("DOM summary -> sources=%d sinks=%d potential=%t\n", scanReport.DOM.Sources, scanReport.DOM.Sinks, scanReport.DOM.Potential)
	fmt.Printf("Scan summary -> method=%s tested=%d reflected=%d\n", scanReport.Method, scanReport.Tested, scanReport.Reflected)
	for _, item := range scanReport.Findings {
		status := "not-reflected"
		if item.Reflected {
			status = "reflected"
		}
		if item.Error != "" {
			fmt.Printf("- %s: error=%s\n", item.Name, item.Error)
			continue
		}
		fmt.Printf("- %s: reflections=%d status=%s\n", item.Name, item.Reflections, status)
	}
}
