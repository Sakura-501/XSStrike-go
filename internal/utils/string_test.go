package utils

import "testing"

func TestCounter(t *testing.T) {
	if Counter("abc !@# 123") != 3 {
		t.Fatalf("unexpected counter result")
	}
}

func TestClosest(t *testing.T) {
	got := Closest(10, map[int]int{0: 1, 1: 9, 2: 20})
	if got[1] != 9 {
		t.Fatalf("unexpected closest result: %+v", got)
	}
}

func TestFillHoles(t *testing.T) {
	got := FillHoles([]int{1, 3}, []int{1, 2})
	if len(got) != 3 || got[0] != 1 || got[1] != 0 || got[2] != 2 {
		t.Fatalf("unexpected fill holes result: %+v", got)
	}
}

func TestStripper(t *testing.T) {
	if Stripper("a/b/c", "/", "right") != "a/bc" {
		t.Fatalf("unexpected right stripper result")
	}
	if Stripper("a/b/c", "/", "left") != "ab/c" {
		t.Fatalf("unexpected left stripper result")
	}
}

func TestDeJSON(t *testing.T) {
	if DeJSON(`a\\b`) != `a\b` {
		t.Fatalf("unexpected dejson result")
	}
}

func TestIsBadContext(t *testing.T) {
	contexts := []ContextRange{{Start: 2, End: 8, Name: "script"}}
	if IsBadContext(3, contexts) != "script" {
		t.Fatalf("expected bad context")
	}
	if IsBadContext(8, contexts) != "" {
		t.Fatalf("expected empty context on boundary")
	}
}

func TestEqualize(t *testing.T) {
	got := Equalize([]string{"a"}, 2)
	if len(got) != 2 || got[1] != "" {
		t.Fatalf("unexpected equalize result: %+v", got)
	}
}

func TestEscaped(t *testing.T) {
	if !Escaped(2, `a\"`) {
		t.Fatalf("expected escaped char")
	}
	if Escaped(3, `a\\\"`) {
		t.Fatalf("did not expect escaped char")
	}
}
