package retirejs

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Sakura-501/XSStrike-go/internal/requester"
)

func TestScanScriptCustomDefinition(t *testing.T) {
	defs := map[string]componentDef{
		"jquery": {
			Vulnerabilities: []vulnerabilityDef{{Below: "1.6.3", Severity: "medium", Info: []string{"xss"}}},
			Extractors: map[string]interface{}{
				"filename": []interface{}{"jquery-([0-9][0-9.a-z_-]+)(\\.min)?\\.js"},
			},
		},
	}
	scanner := New(defs)
	findings := scanner.ScanScript("https://cdn/jquery-1.6.0.min.js", "")
	if len(findings) != 1 {
		t.Fatalf("expected one finding, got %d", len(findings))
	}
	if findings[0].Component != "jquery" {
		t.Fatalf("unexpected component: %s", findings[0].Component)
	}
}

func TestScanPageWithFetcher(t *testing.T) {
	defs := map[string]componentDef{
		"lib": {
			Vulnerabilities: []vulnerabilityDef{{Below: "2.0.0", Severity: "low", Info: []string{"info"}}},
			Extractors: map[string]interface{}{
				"filename": []interface{}{"lib-([0-9][0-9.a-z_-]+)\\.js"},
			},
		},
	}
	scanner := New(defs)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/app" {
			_, _ = w.Write([]byte(`<script src="/js/lib-1.0.0.js"></script>`))
			return
		}
		_, _ = w.Write([]byte("console.log('x')"))
	}))
	defer server.Close()

	client := requester.New(requester.Config{TimeoutSeconds: 5})
	pageResp, err := client.DoGet(server.URL+"/app", map[string]string{}, map[string]string{})
	if err != nil {
		t.Fatalf("page fetch error: %v", err)
	}
	findings := scanner.ScanPage(client, server.URL+"/app", pageResp.Body, map[string]string{})
	if len(findings) != 1 {
		t.Fatalf("expected one finding from page scripts, got %d", len(findings))
	}
}

func TestNewDefault(t *testing.T) {
	scanner, err := NewDefault()
	if err != nil {
		t.Fatalf("expected default scanner: %v", err)
	}
	if len(scanner.Definitions) == 0 {
		t.Fatalf("expected definitions loaded")
	}
}
