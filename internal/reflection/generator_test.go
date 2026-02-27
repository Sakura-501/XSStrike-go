package reflection

import (
	"testing"

	"github.com/Sakura-501/XSStrike-go/internal/config"
)

func TestGenerateCandidatesHTML(t *testing.T) {
	occ := Occurrences{
		1: &Occurrence{
			Position: 1,
			Context:  "html",
			Details:  Details{},
			Score:    map[string]int{"<": 100, ">": 100},
		},
	}
	vectors := GenerateCandidates(occ, "<html></html>")
	if len(vectors[10]) == 0 {
		t.Fatalf("expected confidence 10 vectors")
	}
}

func TestGenerateCandidatesScript(t *testing.T) {
	html := `<script>if(x){` + config.XSSChecker + `}</script>`
	occ := Occurrences{
		2: &Occurrence{
			Position: 2,
			Context:  "script",
			Details:  Details{Quote: "\""},
			Score: map[string]int{
				"</scRipT/>": 100,
				">":          100,
				"\"":         100,
			},
		},
	}
	vectors := GenerateCandidates(occ, html)
	total := 0
	for _, list := range vectors {
		total += len(list)
	}
	if total == 0 {
		t.Fatalf("expected generated vectors")
	}
}
