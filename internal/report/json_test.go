package report

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestWriteJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "result.json")
	err := WriteJSON(path, map[string]interface{}{"ok": true, "count": 2})
	if err != nil {
		t.Fatalf("unexpected write error: %v", err)
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("unexpected read error: %v", err)
	}
	text := string(raw)
	if !strings.Contains(text, `"ok": true`) || !strings.Contains(text, `"count": 2`) {
		t.Fatalf("unexpected json output: %s", text)
	}
}
