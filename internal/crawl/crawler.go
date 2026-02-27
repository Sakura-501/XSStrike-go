package crawl

import (
	"net/url"
	"sort"
	"strings"

	"github.com/Sakura-501/XSStrike-go/internal/dom"
	"github.com/Sakura-501/XSStrike-go/internal/requester"
	"github.com/Sakura-501/XSStrike-go/internal/retirejs"
	"github.com/Sakura-501/XSStrike-go/internal/utils"
)

type Input struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Form struct {
	PageURL string  `json:"page_url"`
	Action  string  `json:"action"`
	Method  string  `json:"method"`
	Inputs  []Input `json:"inputs"`
}

type DOMPage struct {
	URL    string     `json:"url"`
	Report dom.Report `json:"report"`
}

type Report struct {
	Seed       string             `json:"seed"`
	Visited    []string           `json:"visited"`
	Forms      []Form             `json:"forms"`
	DOMPages   []DOMPage          `json:"dom_pages"`
	JSFindings []retirejs.Finding `json:"js_findings"`
	Processed  int                `json:"processed"`
}

type Config struct {
	Level   int
	SkipDOM bool
}

type Crawler struct {
	Client  *requester.Client
	Cfg     Config
	Scanner *retirejs.Scanner
}

func New(client *requester.Client, cfg Config) *Crawler {
	if cfg.Level <= 0 {
		cfg.Level = 2
	}
	scanner, _ := retirejs.NewDefault()
	return &Crawler{Client: client, Cfg: cfg, Scanner: scanner}
}

func (c *Crawler) Discover(seedURL string, headers map[string]string) (Report, error) {
	report := Report{Seed: seedURL, Forms: []Form{}, Visited: []string{}, DOMPages: []DOMPage{}, JSFindings: []retirejs.Finding{}}
	if c.Client == nil {
		return report, nil
	}

	normalized := normalizeSeed(seedURL)
	visited := map[string]struct{}{}
	current := []string{normalized}
	main := rootURL(normalized)

	for depth := 0; depth < c.Cfg.Level && len(current) > 0; depth++ {
		next := []string{}
		for _, pageURL := range current {
			if _, ok := visited[pageURL]; ok {
				continue
			}
			visited[pageURL] = struct{}{}
			report.Processed++

			resp, err := c.Client.DoGet(pageURL, map[string]string{}, headers)
			if err != nil {
				continue
			}

			if !c.Cfg.SkipDOM {
				report.DOMPages = append(report.DOMPages, DOMPage{URL: pageURL, Report: dom.Analyze(resp.Body)})
			}

			report.Forms = append(report.Forms, ExtractForms(pageURL, resp.Body)...)
			report.Forms = append(report.Forms, FormsFromURL(pageURL)...)
			if c.Scanner != nil {
				report.JSFindings = append(report.JSFindings, c.Scanner.ScanPage(c.Client, pageURL, resp.Body, headers)...)
			}

			for _, href := range ExtractLinks(resp.Body) {
				full := utils.HandleAnchor(pageURL, href)
				if !sameHost(main, full) {
					continue
				}
				if isStaticResource(full) {
					continue
				}
				if _, ok := visited[full]; ok {
					continue
				}
				next = append(next, full)
			}
		}
		current = uniqueSorted(next)
	}

	report.Visited = mapKeysSorted(visited)
	report.Forms = dedupForms(report.Forms)
	return report, nil
}

func normalizeSeed(seed string) string {
	trimmed := strings.TrimSpace(seed)
	if strings.HasPrefix(trimmed, "http://") || strings.HasPrefix(trimmed, "https://") {
		return trimmed
	}
	return "http://" + trimmed
}

func rootURL(raw string) string {
	u, err := url.Parse(raw)
	if err != nil {
		return raw
	}
	return u.Scheme + "://" + u.Host
}

func sameHost(root, candidate string) bool {
	uRoot, err1 := url.Parse(root)
	uCand, err2 := url.Parse(candidate)
	if err1 != nil || err2 != nil {
		return false
	}
	return strings.EqualFold(uRoot.Host, uCand.Host)
}

func isStaticResource(link string) bool {
	lower := strings.ToLower(link)
	blocked := []string{".pdf", ".png", ".jpg", ".jpeg", ".xls", ".xml", ".docx", ".doc", ".css", ".js", ".svg"}
	for _, ext := range blocked {
		if strings.HasSuffix(lower, ext) {
			return true
		}
	}
	return false
}

func mapKeysSorted(in map[string]struct{}) []string {
	out := make([]string, 0, len(in))
	for key := range in {
		out = append(out, key)
	}
	sort.Strings(out)
	return out
}

func uniqueSorted(in []string) []string {
	seen := map[string]struct{}{}
	out := make([]string, 0, len(in))
	for _, v := range in {
		if _, ok := seen[v]; ok {
			continue
		}
		seen[v] = struct{}{}
		out = append(out, v)
	}
	sort.Strings(out)
	return out
}

func dedupForms(in []Form) []Form {
	seen := map[string]struct{}{}
	out := make([]Form, 0, len(in))
	for _, form := range in {
		key := form.PageURL + "|" + form.Action + "|" + form.Method + "|" + inputsKey(form.Inputs)
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		out = append(out, form)
	}
	return out
}

func inputsKey(inputs []Input) string {
	parts := make([]string, 0, len(inputs))
	for _, input := range inputs {
		parts = append(parts, input.Name+"="+input.Value)
	}
	sort.Strings(parts)
	return strings.Join(parts, "&")
}
