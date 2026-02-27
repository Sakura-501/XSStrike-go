package crawl

import (
	"github.com/Sakura-501/XSStrike-go/internal/requester"
)

type SeedResult struct {
	Seed      string      `json:"seed"`
	Discovery Report      `json:"discovery"`
	Scan      ScanSummary `json:"scan"`
}

type RunReport struct {
	Seeds          []string     `json:"seeds"`
	Results        []SeedResult `json:"results"`
	TotalProcessed int          `json:"total_processed"`
	TotalForms     int          `json:"total_forms"`
	TotalFindings  int          `json:"total_findings"`
}

func Run(client *requester.Client, seeds []string, headers map[string]string, cfg Config, blindPayload string) (RunReport, error) {
	runReport := RunReport{Seeds: seeds, Results: []SeedResult{}}
	crawler := New(client, cfg)

	for _, seed := range seeds {
		discovery, err := crawler.Discover(seed, headers)
		if err != nil {
			return runReport, err
		}
		scanSummary := ScanForms(client, discovery.Forms, headers, blindPayload)
		runReport.Results = append(runReport.Results, SeedResult{
			Seed:      seed,
			Discovery: discovery,
			Scan:      scanSummary,
		})
		runReport.TotalProcessed += discovery.Processed
		runReport.TotalForms += len(discovery.Forms)
		runReport.TotalFindings += len(scanSummary.Findings)
	}

	return runReport, nil
}
