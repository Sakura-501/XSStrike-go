package crawl

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Sakura-501/XSStrike-go/internal/requester"
)

func TestRun(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`<a href="/a">go</a>`))
	})
	mux.HandleFunc("/a", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`<form action="/submit" method="get"><input name="q" value="1"></form>`))
	})
	mux.HandleFunc("/submit", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(r.URL.Query().Get("q")))
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := requester.New(requester.Config{TimeoutSeconds: 5})
	runReport, err := Run(client, []string{server.URL}, map[string]string{}, Config{Level: 3}, "")
	if err != nil {
		t.Fatalf("run error: %v", err)
	}
	if len(runReport.Results) != 1 {
		t.Fatalf("expected one result")
	}
	if runReport.TotalForms == 0 {
		t.Fatalf("expected discovered forms")
	}
	if runReport.TotalProcessed == 0 {
		t.Fatalf("expected processed pages")
	}
}
