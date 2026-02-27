package payload

import (
	"math/rand"
	"strings"
	"testing"
)

func TestGenerateVectorsCount(t *testing.T) {
	in := GeneratorInput{
		Fillings:  []string{" "},
		EFillings: []string{""},
		LFillings: []string{""},
		EventHandlers: map[string][]string{
			"onmouseover": {"a", "html"},
		},
		Tags:      []string{"a", "html"},
		Functions: []string{"confirm()"},
		Ends:      []string{">"},
		Bait:      "v3dm0s",
	}

	vectors := GenerateVectors(in, rand.New(rand.NewSource(7)))
	if len(vectors) != 2 {
		t.Fatalf("unexpected vector count: got %d want 2", len(vectors))
	}

	if !strings.Contains(vectors[0], "confirm()") || !strings.Contains(vectors[1], "confirm()") {
		t.Fatalf("generated vectors are missing function call: %+v", vectors)
	}
}

func TestRandomUpperPreservesCharacters(t *testing.T) {
	input := "onmouseover"
	out := RandomUpper(input, rand.New(rand.NewSource(9)))
	if strings.ToLower(out) != input {
		t.Fatalf("random upper changed content: %q", out)
	}
}
