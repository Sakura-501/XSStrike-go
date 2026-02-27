package crawl

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Sakura-501/XSStrike-go/internal/requester"
)

func TestDiscover(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`
			<a href="/a">A</a>
			<a href="/static.jpg">IMG</a>
			<script>var x = location.search; document.write(x)</script>
		`))
	})
	mux.HandleFunc("/a", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`
			<form action="/submit" method="post">
			  <input name="q" value="1">
			</form>
			<a href="/b?x=1">B</a>
		`))
	})
	mux.HandleFunc("/b", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("done"))
	})
	mux.HandleFunc("/submit", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("ok"))
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	crawler := New(requester.New(requester.Config{TimeoutSeconds: 5}), Config{Level: 3, SkipDOM: false})
	report, err := crawler.Discover(server.URL, map[string]string{})
	if err != nil {
		t.Fatalf("discover error: %v", err)
	}

	if report.Processed < 3 {
		t.Fatalf("expected at least 3 processed pages, got %d", report.Processed)
	}
	if len(report.Forms) == 0 {
		t.Fatalf("expected extracted forms")
	}
	if len(report.DOMPages) == 0 {
		t.Fatalf("expected dom reports")
	}
}

func TestExtractFormsAndLinks(t *testing.T) {
	html := `<a href='/x'>x</a><form action='/submit' method='get'><input name='a' value='1'></form>`
	links := ExtractLinks(html)
	if len(links) != 1 || links[0] != "/x" {
		t.Fatalf("unexpected links: %+v", links)
	}

	forms := ExtractForms("https://example.com/p", html)
	if len(forms) != 1 {
		t.Fatalf("unexpected forms count: %d", len(forms))
	}
	if forms[0].Action != "https://example.com/submit" {
		t.Fatalf("unexpected action: %s", forms[0].Action)
	}
}

func TestFormsFromURL(t *testing.T) {
	forms := FormsFromURL("https://example.com/path?a=1&b=2")
	if len(forms) != 1 {
		t.Fatalf("expected one pseudo form")
	}
	if forms[0].Method != "get" {
		t.Fatalf("unexpected method: %s", forms[0].Method)
	}
}

func TestDiscoverDetectsVulnerableJS(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`<script src="/1.6.0/jquery.js"></script>`))
	})
	mux.HandleFunc("/1.6.0/jquery.js", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`/*! jquery */`))
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	crawler := New(requester.New(requester.Config{TimeoutSeconds: 5}), Config{Level: 1, SkipDOM: true})
	report, err := crawler.Discover(server.URL, map[string]string{})
	if err != nil {
		t.Fatalf("discover error: %v", err)
	}
	if len(report.JSFindings) == 0 {
		t.Fatalf("expected js vulnerability findings")
	}
}
