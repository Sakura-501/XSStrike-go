package utils

import "testing"

func TestExtractReflectedScripts(t *testing.T) {
	html := `<script>var a = "v3dm0s"</script><script>alert(1)</script>`
	found := ExtractReflectedScripts(html, "v3dm0s")
	if len(found) != 1 {
		t.Fatalf("unexpected matched count: %d", len(found))
	}
}

func TestExtractJSSources(t *testing.T) {
	html := `<script src="/app.js"></script><script SRC='https://cdn/x.js'></script>`
	sources := ExtractJSSources(html)
	if len(sources) != 2 {
		t.Fatalf("unexpected source count: %d", len(sources))
	}
	if sources[0] != "/app.js" || sources[1] != "https://cdn/x.js" {
		t.Fatalf("unexpected sources: %+v", sources)
	}
}

func TestHandleAnchor(t *testing.T) {
	cases := []struct {
		parent string
		anchor string
		want   string
	}{
		{"https://example.com/base/page", "//cdn.site/x.js", "https://cdn.site/x.js"},
		{"https://example.com/base/page", "/static/app.js", "https://example.com/static/app.js"},
		{"https://example.com/base/", "next", "https://example.com/base/next"},
		{"https://example.com/base", "next", "https://example.com/next"},
	}
	for _, tc := range cases {
		got := HandleAnchor(tc.parent, tc.anchor)
		if got != tc.want {
			t.Fatalf("HandleAnchor(%q, %q): got %q want %q", tc.parent, tc.anchor, got, tc.want)
		}
	}
}
