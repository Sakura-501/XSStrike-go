package dom

import "testing"

func TestAnalyzePotential(t *testing.T) {
	html := `<html><script>
var user = location.search;
document.write(user);
</script></html>`
	report := Analyze(html)

	if !report.Checked {
		t.Fatalf("expected checked report")
	}
	if report.Sources == 0 {
		t.Fatalf("expected source findings")
	}
	if report.Sinks == 0 {
		t.Fatalf("expected sink findings")
	}
	if !report.Potential {
		t.Fatalf("expected potential source-to-sink result")
	}
}

func TestAnalyzeSinkOnly(t *testing.T) {
	html := `<script>document.write("safe")</script>`
	report := Analyze(html)
	if report.Sources != 0 {
		t.Fatalf("expected 0 sources")
	}
	if report.Sinks != 1 {
		t.Fatalf("expected 1 sink, got %d", report.Sinks)
	}
	if report.Potential {
		t.Fatalf("did not expect potential result")
	}
}

func TestAnalyzeNoScript(t *testing.T) {
	html := `<html><body>no scripts</body></html>`
	report := Analyze(html)
	if report.Sources != 0 || report.Sinks != 0 {
		t.Fatalf("expected no findings")
	}
	if len(report.Findings) != 0 {
		t.Fatalf("expected empty findings")
	}
}
