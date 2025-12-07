package tests

import (
	"os"
	"path/filepath"
	"testing"

	"gathercode/pkg/gather"
)

func TestCollectLocalPaths(t *testing.T) {
	root := t.TempDir()
	repo := filepath.Join(root, "repo")
	if err := os.MkdirAll(filepath.Join(repo, "pkg"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(repo, "main.go"), []byte("package main\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(repo, "pkg", "migrate.sql"), []byte("CREATE TABLE t();"), 0o644); err != nil {
		t.Fatal(err)
	}
	inputs := []string{repo}
	exts := []string{".go", ".sql"}
	opts := gather.Options{IncludeHidden: false}
	entries, err := gather.Collect(inputs, exts, opts)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	found := map[string]bool{}
	for _, e := range entries {
		switch filepath.Base(e.AbsPath) {
		case "main.go":
			found["main.go"] = true
		case "migrate.sql":
			found["migrate.sql"] = true
		default:
			t.Fatalf("unexpected file %s", e.AbsPath)
		}
		if filepath.Base(filepath.Dir(e.AbsPath)) == ".git" {
			t.Fatalf(".git should not be included")
		}
	}
	if !found["main.go"] || !found["migrate.sql"] {
		t.Fatalf("missing expected files: %+v", found)
	}
}
