package config

import "testing"

func TestDefaultRuleSetsAreRich(t *testing.T) {
	if len(DefaultFunctions) < 6 {
		t.Fatalf("expected at least 6 default functions, got %d", len(DefaultFunctions))
	}
	if len(DefaultPayloads) < 20 {
		t.Fatalf("expected at least 20 default payloads, got %d", len(DefaultPayloads))
	}
	if len(DefaultFuzzes) < 25 {
		t.Fatalf("expected at least 25 default fuzzes, got %d", len(DefaultFuzzes))
	}
}
