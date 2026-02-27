package files

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReadLines(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "payloads.txt")
	if err := os.WriteFile(path, []byte("a\n\n b \n"), 0o644); err != nil {
		t.Fatalf("write file failed: %v", err)
	}

	lines, err := ReadLines(path)
	if err != nil {
		t.Fatalf("unexpected read error: %v", err)
	}
	if len(lines) != 2 || lines[0] != "a" || lines[1] != "b" {
		t.Fatalf("unexpected lines result: %+v", lines)
	}
}
