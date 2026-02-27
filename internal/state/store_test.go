package state

import "testing"

func TestSetAndGet(t *testing.T) {
	s := New()
	s.Set("timeout", 10)
	value, ok := s.Get("timeout")
	if !ok {
		t.Fatal("expected key to exist")
	}
	if value.(int) != 10 {
		t.Fatalf("unexpected value: %v", value)
	}
}

func TestUpdateAppend(t *testing.T) {
	s := New()
	s.Set("items", []string{"a"})
	if err := s.Update("items", "b", "append"); err != nil {
		t.Fatalf("unexpected append error: %v", err)
	}
	value := s.MustGet("items").([]string)
	if len(value) != 2 || value[1] != "b" {
		t.Fatalf("unexpected append result: %+v", value)
	}
}

func TestUpdateAdd(t *testing.T) {
	s := New()
	s.Set("checked", map[string]struct{}{})
	if err := s.Update("checked", "x.js", "add"); err != nil {
		t.Fatalf("unexpected add error: %v", err)
	}
	set := s.MustGet("checked").(map[string]struct{})
	if _, ok := set["x.js"]; !ok {
		t.Fatalf("expected set to contain x.js: %+v", set)
	}
}

func TestUpdateErrors(t *testing.T) {
	s := New()
	s.Set("notSlice", 1)
	if err := s.Update("notSlice", "x", "append"); err == nil {
		t.Fatal("expected append error")
	}

	s.Set("notSet", []string{})
	if err := s.Update("notSet", "x", "add"); err == nil {
		t.Fatal("expected add error")
	}
}
