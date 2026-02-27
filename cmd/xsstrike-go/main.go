package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Sakura-501/XSStrike-go/internal/bruteforce"
	"github.com/Sakura-501/XSStrike-go/internal/config"
	"github.com/Sakura-501/XSStrike-go/internal/crawl"
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

	if opts.URL == "" && !opts.Fuzzer && !opts.Crawl && opts.SeedsFile == "" {
		fs.Usage()
		return
	}

	fmt.Print(ui.Banner())
	fmt.Printf("Runtime defaults -> timeout=%ds threads=%d delay=%ds\n", opts.Timeout, opts.ThreadCount, opts.Delay)

	if opts.Fuzzer {
		runFuzzer(opts)
		return
	}

	headers := mergedHeaders(opts.HeadersRaw)
	client := requester.New(requester.Config{TimeoutSeconds: opts.Timeout, DelaySeconds: opts.Delay, Proxy: opts.Proxy})

	if opts.Crawl || opts.SeedsFile != "" {
		runCrawl(opts, headers, client)
		return
	}

	runSingleScan(opts, headers, client)
}

func runSingleScan(opts *options.Options, headers map[string]string, client *requester.Client) {
	fmt.Printf("Target: %s\n", opts.URL)
	state.Global.Set("headers", headers)
	state.Global.Set("checkedScripts", map[string]struct{}{})

	if opts.HeadersRaw != "" {
		fmt.Printf("Custom headers parsed: %d\n", len(utils.ExtractHeaders(opts.HeadersRaw)))
	}

	if opts.PayloadFile != "" {
		payloadList, err := resolvePayloadList(opts.PayloadFile)
		if err != nil {
			fmt.Printf("Payload file error: %v\n", err)
			os.Exit(1)
		}
		bfReport, err := bruteforce.Run(client, opts.URL, opts.Data, opts.JSON, opts.Path, headers, payloadList, opts.Encode)
		if err != nil {
			fmt.Printf("Bruteforce error: %v\n", err)
			os.Exit(1)
		}
		state.Global.Set("bruteforceReport", bfReport)
		printBruteforceReport(bfReport)
		if opts.OutputJSON != "" {
			if err := report.WriteJSON(opts.OutputJSON, bfReport); err != nil {
				fmt.Printf("Write output error: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("Bruteforce report written: %s\n", opts.OutputJSON)
		}
		return
	}

	runner := scan.NewRunner(client)
	scanReport, err := runner.Run(opts.URL, opts.Data, headers, opts.JSON, opts.Path, opts.Encode)
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

func runCrawl(opts *options.Options, headers map[string]string, client *requester.Client) {
	seeds := []string{}
	if opts.URL != "" {
		seeds = append(seeds, opts.URL)
	}
	if opts.SeedsFile != "" {
		items, err := files.ReadLines(opts.SeedsFile)
		if err != nil {
			fmt.Printf("Seed file error: %v\n", err)
			os.Exit(1)
		}
		seeds = append(seeds, items...)
	}
	if len(seeds) == 0 {
		fmt.Println("No crawl seeds provided.")
		os.Exit(1)
	}

	blindPayload := ""
	if opts.Blind {
		blindPayload = opts.BlindPayload
	}

	runReport, err := crawl.Run(client, seeds, headers, crawl.Config{Level: opts.Level, SkipDOM: opts.SkipDOM}, blindPayload)
	if err != nil {
		fmt.Printf("Crawl error: %v\n", err)
		os.Exit(1)
	}

	state.Global.Set("crawlReport", runReport)
	printCrawlReport(runReport)

	if opts.OutputJSON != "" {
		if err := report.WriteJSON(opts.OutputJSON, runReport); err != nil {
			fmt.Printf("Write output error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Crawl report written: %s\n", opts.OutputJSON)
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
	if scanReport.WAF.Detected {
		fmt.Printf("WAF detected -> %s (score=%d)\n", scanReport.WAF.Name, scanReport.WAF.Score)
	} else {
		fmt.Println("WAF detected -> no")
	}
	fmt.Printf("Scan summary -> method=%s tested=%d reflected=%d generated_candidates=%d\n", scanReport.Method, scanReport.Tested, scanReport.Reflected, scanReport.GeneratedCandidates)
	for _, item := range scanReport.Findings {
		status := "not-reflected"
		if item.Reflected {
			status = "reflected"
		}
		if item.Error != "" {
			fmt.Printf("- %s: error=%s\n", item.Name, item.Error)
			continue
		}
		fmt.Printf("- %s: reflections=%d occurrences=%d candidates=%d top_conf=%d status=%s\n", item.Name, item.Reflections, item.Occurrences, item.Candidates, item.TopConfidence, status)
	}
}

func printBruteforceReport(bruteforceReport bruteforce.Report) {
	if bruteforceReport.NoParams {
		fmt.Println("Bruteforce -> no parameters to test.")
		return
	}
	fmt.Printf("Bruteforce summary -> tested=%d hits=%d\n", bruteforceReport.Tested, len(bruteforceReport.Hits))
	for _, hit := range bruteforceReport.Hits {
		fmt.Printf("- param=%s reflections=%d payload=%s\n", hit.Param, hit.Reflections, hit.Payload)
	}
}

func printCrawlReport(runReport crawl.RunReport) {
	fmt.Printf("Crawl summary -> seeds=%d processed=%d forms=%d findings=%d js_findings=%d\n", len(runReport.Seeds), runReport.TotalProcessed, runReport.TotalForms, runReport.TotalFindings, runReport.TotalJSFindings)
	for _, result := range runReport.Results {
		domPotential := 0
		for _, page := range result.Discovery.DOMPages {
			if page.Report.Potential {
				domPotential++
			}
		}
		fmt.Printf("- seed=%s visited=%d forms=%d dom_potential=%d findings=%d js_findings=%d\n", result.Seed, len(result.Discovery.Visited), len(result.Discovery.Forms), domPotential, len(result.Scan.Findings), len(result.Discovery.JSFindings))
	}
}
