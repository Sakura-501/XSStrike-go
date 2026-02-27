package compat

import (
	"reflect"
	"testing"

	"github.com/Sakura-501/XSStrike-go/internal/config"
	"github.com/Sakura-501/XSStrike-go/internal/dom"
	"github.com/Sakura-501/XSStrike-go/internal/utils"
)

func TestPythonParityExtractHeaders(t *testing.T) {
	raw := "A: one\\nB: two,\\nC: three"
	got := utils.ExtractHeaders(raw)
	want := map[string]string{"A": "one", "B": "two", "C": "three"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("extract headers mismatch: got=%v want=%v", got, want)
	}
}

func TestPythonParityGetURLAndParseParams(t *testing.T) {
	target := "https://example.com/path?a=1&b=2"
	if got := utils.GetURL(target, true); got != "https://example.com/path" {
		t.Fatalf("get url mismatch: %s", got)
	}
	params := utils.ParseParams(target, "", false)
	want := map[string]string{"a": "1", "b": "2"}
	if !reflect.DeepEqual(params, want) {
		t.Fatalf("params mismatch: got=%v want=%v", params, want)
	}
}

func TestPythonParityHandleAnchor(t *testing.T) {
	base := "https://example.com/app/index"
	cases := map[string]string{
		"//cdn.site/app.js":      "https://cdn.site/app.js",
		"/assets/main.js":        "https://example.com/assets/main.js",
		"scripts/util.js":        "https://example.com/app/scripts/util.js",
		"https://other.com/x.js": "https://other.com/x.js",
	}
	for in, want := range cases {
		if got := utils.HandleAnchor(base, in); got != want {
			t.Fatalf("anchor mismatch for %s: got=%s want=%s", in, got, want)
		}
	}
}

func TestPythonParityDOMAndPayloadRules(t *testing.T) {
	html := `<script>var x = location.search; document.write(x)</script>`
	report := dom.Analyze(html)
	if !report.Potential {
		t.Fatalf("expected dom potential=true")
	}
	if len(config.DefaultPayloads) < 20 || len(config.DefaultFuzzes) < 25 {
		t.Fatalf("expected richer migrated payload rules")
	}
}

func TestPythonParityJSExtractor(t *testing.T) {
	html := `<script src="/a.js"></script><script SRC='//cdn/b.js'></script>`
	got := utils.ExtractJSSources(html)
	want := []string{"/a.js", "//cdn/b.js"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("js sources mismatch: got=%v want=%v", got, want)
	}
}
