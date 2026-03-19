package crawl

import (
	"strings"

	"github.com/Sakura-501/XSStrike-go/internal/requester"
)

type SeedResult struct {
	Seed      string      `json:"seed"`
	Discovery Report      `json:"discovery"`
	Scan      ScanSummary `json:"scan"`
}

type RunReport struct {
	Seeds           []string     `json:"seeds"`
	Results         []SeedResult `json:"results"`
	TotalProcessed  int          `json:"total_processed"`
	TotalForms      int          `json:"total_forms"`
	TotalFindings   int          `json:"total_findings"`
	TotalJSFindings int          `json:"total_js_findings"`
}

func Run(client *requester.Client, seeds []string, headers map[string]string, cfg Config, blindPayload string) (RunReport, error) {
	normalizedSeeds := normalizeSeeds(seeds)
	runReport := RunReport{Seeds: normalizedSeeds, Results: []SeedResult{}}
	crawler := New(client, cfg)

	for _, seed := range normalizedSeeds {
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
		runReport.TotalJSFindings += len(discovery.JSFindings)
	}

	return runReport, nil
}

func normalizeSeeds(seeds []string) []string {
	seen := make(map[string]struct{}, len(seeds))
	normalized := make([]string, 0, len(seeds))
	for _, seed := range seeds {
		trimmed := strings.TrimSpace(seed)
		if trimmed == "" {
			continue
		}
		if _, ok := seen[trimmed]; ok {
			continue
		}
		seen[trimmed] = struct{}{}
		normalized = append(normalized, trimmed)
	}
	return normalized
}
