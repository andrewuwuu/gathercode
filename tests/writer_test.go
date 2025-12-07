package tests

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gathercode/pkg/gather"
	"gathercode/pkg/writer"
)

func TestWriteAggregatedMarkdown(t *testing.T) {
	tmp := t.TempDir()
	f1 := filepath.Join(tmp, "a.go")
	f2 := filepath.Join(tmp, "sub", "b.sql")
	if err := os.MkdirAll(filepath.Dir(f2), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(f1, []byte("package a\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(f2, []byte("SELECT 42;"), 0o644); err != nil {
		t.Fatal(err)
	}
	entries := []gather.Entry{
		{DisplayPath: "repo/a.go", AbsPath: f1},
		{DisplayPath: "repo/sub/b.sql", AbsPath: f2},
	}
	out := filepath.Join(tmp, "out.md")
	if err := writer.WriteAggregatedMarkdown(entries, out); err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(out)
	if err != nil {
		t.Fatal(err)
	}
	s := string(got)
	if !strings.Contains(s, "-- repo/a.go") {
		t.Fatalf("missing header for a.go")
	}
	if !strings.Contains(s, "package a") {
		t.Fatalf("missing content for a.go")
	}
	if !strings.Contains(s, "-- repo/sub/b.sql") {
		t.Fatalf("missing header for b.sql")
	}
	if !strings.Contains(s, "SELECT 42;") {
		t.Fatalf("missing content for b.sql")
	}
	if !strings.Contains(s, "---------------") {
		t.Fatalf("missing separator")
	}
}