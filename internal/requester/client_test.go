package requester

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestDoGet(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("expected GET, got %s", r.Method)
		}
		if r.URL.Query().Get("q") != "test" {
			t.Fatalf("missing q parameter")
		}
		if r.Header.Get("X-Test") != "yes" {
			t.Fatalf("missing custom header")
		}
		if r.Header.Get("User-Agent") == "" || r.Header.Get("User-Agent") == "$" {
			t.Fatalf("invalid user-agent header: %q", r.Header.Get("User-Agent"))
		}
		_, _ = w.Write([]byte("ok"))
	}))
	defer server.Close()

	client := New(Config{TimeoutSeconds: 5})
	resp, err := client.DoGet(server.URL, map[string]string{"q": "test"}, map[string]string{"X-Test": "yes", "User-Agent": "$"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusOK || resp.Body != "ok" {
		t.Fatalf("unexpected response: %+v", resp)
	}
}

func TestDoPostForm(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", r.Method)
		}
		if r.Header.Get("Content-Type") == "" {
			t.Fatalf("missing content-type")
		}
		raw, _ := io.ReadAll(r.Body)
		decoded, _ := url.ParseQuery(string(raw))
		if decoded.Get("a") != "1" {
			t.Fatalf("missing form value")
		}
		_, _ = w.Write([]byte("form-ok"))
	}))
	defer server.Close()

	client := New(Config{TimeoutSeconds: 5})
	resp, err := client.DoPost(server.URL, map[string]string{"a": "1"}, map[string]string{}, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Body != "form-ok" {
		t.Fatalf("unexpected response body: %q", resp.Body)
	}
}

func TestDoPostJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/json" {
			t.Fatalf("expected application/json")
		}
		var decoded map[string]string
		if err := json.NewDecoder(r.Body).Decode(&decoded); err != nil {
			t.Fatalf("invalid json body: %v", err)
		}
		if decoded["name"] != "alice" {
			t.Fatalf("invalid json payload")
		}
		_, _ = w.Write([]byte("json-ok"))
	}))
	defer server.Close()

	client := New(Config{TimeoutSeconds: 5})
	resp, err := client.DoPost(server.URL, map[string]string{"name": "alice"}, map[string]string{}, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Body != "json-ok" {
		t.Fatalf("unexpected response body: %q", resp.Body)
	}
}
