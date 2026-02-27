package utils

import "testing"

func TestExtractHeaders(t *testing.T) {
	raw := "X-Test: one\\nAccept: text/html,"
	headers := ExtractHeaders(raw)
	if headers["X-Test"] != "one" {
		t.Fatalf("unexpected X-Test: %q", headers["X-Test"])
	}
	if headers["Accept"] != "text/html" {
		t.Fatalf("unexpected Accept: %q", headers["Accept"])
	}
}

func TestGetURL(t *testing.T) {
	in := "https://example.com/path?a=1&b=2"
	out := GetURL(in, true)
	if out != "https://example.com/path" {
		t.Fatalf("unexpected output: %q", out)
	}
}

func TestParseParamsFromQuery(t *testing.T) {
	params := ParseParams("https://example.com/path?a=1&b=2", "", false)
	if params["a"] != "1" || params["b"] != "2" {
		t.Fatalf("unexpected params: %+v", params)
	}
}

func TestParseParamsFromForm(t *testing.T) {
	params := ParseParams("https://example.com/path", "x=9&y=", false)
	if params["x"] != "9" || params["y"] != "" {
		t.Fatalf("unexpected params: %+v", params)
	}
}

func TestParseParamsFromJSON(t *testing.T) {
	params := ParseParams("https://example.com/path", "{\"name\":\"alice\",\"age\":10}", true)
	if params["name"] != "alice" || params["age"] != "10" {
		t.Fatalf("unexpected params: %+v", params)
	}
}
