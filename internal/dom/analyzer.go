package dom

import (
	"regexp"
	"strings"
)

var (
	sourcePattern = regexp.MustCompile(`\b(?:document\.(?:URL|documentURI|URLUnencoded|baseURI|cookie|referrer)|location\.(?:href|search|hash|pathname)|window\.name|history\.(?:pushState|replaceState)|(?:local|session)Storage)\b`)
	sinkPattern   = regexp.MustCompile(`\b(?:eval|Function|set(?:Timeout|Interval|Immediate)|document\.(?:write|writeln)|(?:[a-zA-Z0-9_$]+\.)?innerHTML|(?:document|window)\.location|assign|navigate)\b`)
	scriptPattern = regexp.MustCompile(`(?is)<script[^>]*>(.*?)</script>`)
)

type Finding struct {
	ScriptIndex int    `json:"script_index"`
	Line        int    `json:"line"`
	Kind        string `json:"kind"`
	Match       string `json:"match"`
	Snippet     string `json:"snippet"`
}

type Report struct {
	Checked   bool      `json:"checked"`
	Sources   int       `json:"sources"`
	Sinks     int       `json:"sinks"`
	Potential bool      `json:"potential"`
	Findings  []Finding `json:"findings"`
}

// Analyze scans inline scripts and reports DOM source/sink findings.
func Analyze(response string) Report {
	report := Report{Checked: true, Findings: []Finding{}}
	scripts := scriptPattern.FindAllStringSubmatch(response, -1)

	for scriptIndex, scriptMatch := range scripts {
		if len(scriptMatch) < 2 {
			continue
		}
		lines := strings.Split(scriptMatch[1], "\n")
		for lineNo, rawLine := range lines {
			line := strings.TrimSpace(rawLine)
			if line == "" {
				continue
			}

			for _, match := range sourcePattern.FindAllString(line, -1) {
				report.Sources++
				report.Findings = append(report.Findings, Finding{
					ScriptIndex: scriptIndex + 1,
					Line:        lineNo + 1,
					Kind:        "source",
					Match:       match,
					Snippet:     line,
				})
			}

			for _, match := range sinkPattern.FindAllString(line, -1) {
				report.Sinks++
				report.Findings = append(report.Findings, Finding{
					ScriptIndex: scriptIndex + 1,
					Line:        lineNo + 1,
					Kind:        "sink",
					Match:       match,
					Snippet:     line,
				})
			}
		}
	}

	report.Potential = report.Sources > 0 && report.Sinks > 0
	return report
}
