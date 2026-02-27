package utils

import (
	"strings"
	"testing"
)

func TestURLPathToMap(t *testing.T) {
	result := URLPathToMap("https://example.com/a/b")
	if result["a"] != "a" || result["b"] != "b" {
		t.Fatalf("unexpected path map: %+v", result)
	}
}

func TestMapToURLPath(t *testing.T) {
	raw := MapToURLPath("https://example.com/old/path", map[string]string{
		"a": "one",
		"b": "two",
	})
	if raw != "https://example.com/one/two" {
		t.Fatalf("unexpected built path: %q", raw)
	}
}

func TestJSONToMap(t *testing.T) {
	result, err := JSONToMap(`{"name":"alice","age":10}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["name"] != "alice" || result["age"] != "10" {
		t.Fatalf("unexpected map result: %+v", result)
	}
}

func TestMapToJSON(t *testing.T) {
	raw, err := MapToJSON(map[string]string{"a": "1", "b": "2"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(raw, `"a":"1"`) || !strings.Contains(raw, `"b":"2"`) {
		t.Fatalf("unexpected json output: %q", raw)
	}
}

func TestFlattenParams(t *testing.T) {
	query := FlattenParams("q", map[string]string{"q": "test", "page": "1"}, "PAYLOAD")
	if query != "?page=1&q=PAYLOAD" {
		t.Fatalf("unexpected flattened query: %q", query)
	}
}
