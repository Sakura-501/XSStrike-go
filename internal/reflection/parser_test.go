package reflection

import (
	"testing"

	"github.com/Sakura-501/XSStrike-go/internal/config"
)

func TestParseFindsScriptContext(t *testing.T) {
	html := `<script>var a = "` + config.XSSChecker + `";</script>`
	occ := Parse(html, "")
	if occ.Count() == 0 {
		t.Fatalf("expected occurrences")
	}
	foundScript := false
	for _, item := range occ {
		if item.Context == "script" {
			foundScript = true
			if item.Details.Quote == "" {
				t.Fatalf("expected script quote detail")
			}
		}
	}
	if !foundScript {
		t.Fatalf("expected script context")
	}
}

func TestParseFindsAttributeAndComment(t *testing.T) {
	html := `<a href="` + config.XSSChecker + `">x</a><!-- ` + config.XSSChecker + ` -->`
	occ := Parse(html, "")
	attr, comment := 0, 0
	for _, item := range occ {
		if item.Context == "attribute" {
			attr++
		}
		if item.Context == "comment" {
			comment++
		}
	}
	if attr == 0 {
		t.Fatalf("expected attribute context")
	}
	if comment == 0 {
		t.Fatalf("expected comment context")
	}
}

func TestParseBadTag(t *testing.T) {
	html := `<style>.x{content:"` + config.XSSChecker + `"}</style>`
	occ := Parse(html, "")
	for _, item := range occ {
		if item.Details.BadTag == "style" {
			return
		}
	}
	t.Fatalf("expected bad tag style")
}

func TestJSContexter(t *testing.T) {
	script := `if(window.x){` + config.XSSChecker
	breaker := JSContexter(script)
	if breaker == "" {
		t.Fatalf("expected non-empty breaker")
	}
}
