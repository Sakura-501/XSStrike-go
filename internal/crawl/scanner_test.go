package crawl

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Sakura-501/XSStrike-go/internal/requester"
)

func TestScanFormsGET(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(fmt.Sprintf("q=%s", r.URL.Query().Get("q"))))
	}))
	defer server.Close()

	forms := []Form{{
		PageURL: server.URL,
		Action:  server.URL,
		Method:  "get",
		Inputs:  []Input{{Name: "q", Value: "1"}},
	}}

	summary := ScanForms(requester.New(requester.Config{TimeoutSeconds: 5}), forms, map[string]string{}, "")
	if summary.Tested != 1 {
		t.Fatalf("expected tested=1, got %d", summary.Tested)
	}
	if len(summary.Findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(summary.Findings))
	}
	if summary.Findings[0].Reflections == 0 {
		t.Fatalf("expected reflection")
	}
}

func TestScanFormsPOSTAndBlind(t *testing.T) {
	blindSeen := false
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = r.ParseForm()
		if r.FormValue("token") == "blind-1" {
			blindSeen = true
		}
		_, _ = w.Write([]byte(r.FormValue("token")))
	}))
	defer server.Close()

	forms := []Form{{
		PageURL: server.URL,
		Action:  server.URL,
		Method:  "post",
		Inputs:  []Input{{Name: "token", Value: "abc"}},
	}}

	summary := ScanForms(requester.New(requester.Config{TimeoutSeconds: 5}), forms, map[string]string{}, "blind-1")
	if summary.Tested != 1 {
		t.Fatalf("expected tested=1")
	}
	if len(summary.Findings) != 1 {
		t.Fatalf("expected finding")
	}
	if !summary.Findings[0].BlindSent {
		t.Fatalf("expected blind payload send flag")
	}
	if !blindSeen {
		t.Fatalf("expected blind payload request to be sent")
	}
}
