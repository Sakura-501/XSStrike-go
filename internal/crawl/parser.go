package crawl

import (
	"net/url"
	"regexp"
	"strings"

	"github.com/Sakura-501/XSStrike-go/internal/utils"
)

var (
	reLink   = regexp.MustCompile(`(?is)<a[^>]+href=["']?([^"'\s>]+)`)
	reForm   = regexp.MustCompile(`(?is)<form([^>]*)>(.*?)</form>`)
	reInput  = regexp.MustCompile(`(?is)<input([^>]*)>`)
	reAttr   = regexp.MustCompile(`(?is)([a-zA-Z_:][-a-zA-Z0-9_:.]*)\s*=\s*["']([^"']*)["']`)
	reMethod = regexp.MustCompile(`(?is)method\s*=\s*["']?([a-zA-Z]+)`)
	reAction = regexp.MustCompile(`(?is)action\s*=\s*["']?([^"'\s>]+)`)
)

func ExtractLinks(html string) []string {
	matches := reLink.FindAllStringSubmatch(html, -1)
	out := make([]string, 0, len(matches))
	for _, match := range matches {
		if len(match) < 2 {
			continue
		}
		href := strings.TrimSpace(match[1])
		if href == "" || strings.HasPrefix(href, "#") || strings.HasPrefix(strings.ToLower(href), "javascript:") {
			continue
		}
		out = append(out, href)
	}
	return out
}

func ExtractForms(pageURL, html string) []Form {
	forms := []Form{}
	for _, match := range reForm.FindAllStringSubmatch(html, -1) {
		if len(match) < 3 {
			continue
		}
		attrs := match[1]
		body := match[2]

		method := "get"
		if mm := reMethod.FindStringSubmatch(attrs); len(mm) > 1 {
			method = strings.ToLower(strings.TrimSpace(mm[1]))
		}

		action := pageURL
		if am := reAction.FindStringSubmatch(attrs); len(am) > 1 {
			action = utils.HandleAnchor(pageURL, strings.TrimSpace(am[1]))
		}

		inputs := []Input{}
		for _, inputMatch := range reInput.FindAllStringSubmatch(body, -1) {
			if len(inputMatch) < 2 {
				continue
			}
			attrMap := parseAttrs(inputMatch[1])
			name := attrMap["name"]
			if strings.TrimSpace(name) == "" {
				continue
			}
			inputs = append(inputs, Input{Name: name, Value: attrMap["value"]})
		}

		forms = append(forms, Form{PageURL: pageURL, Action: action, Method: method, Inputs: inputs})
	}
	return forms
}

func FormsFromURL(rawURL string) []Form {
	u, err := url.Parse(rawURL)
	if err != nil || u.RawQuery == "" {
		return nil
	}
	query := u.Query()
	inputs := []Input{}
	for key, values := range query {
		value := ""
		if len(values) > 0 {
			value = values[0]
		}
		inputs = append(inputs, Input{Name: key, Value: value})
	}
	u.RawQuery = ""
	return []Form{{
		PageURL: rawURL,
		Action:  u.String(),
		Method:  "get",
		Inputs:  inputs,
	}}
}

func parseAttrs(raw string) map[string]string {
	out := map[string]string{}
	for _, match := range reAttr.FindAllStringSubmatch(raw, -1) {
		if len(match) < 3 {
			continue
		}
		key := strings.ToLower(strings.TrimSpace(match[1]))
		out[key] = strings.TrimSpace(match[2])
	}
	return out
}
