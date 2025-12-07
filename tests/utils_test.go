package tests

import (
	"archive/zip"
	"bytes"
	"io"
	"os"
	"path/filepath"
	"testing"

	"gathercode/pkg/utils"
)

func TestGatherFilesFromRoot(t *testing.T) {
	root := t.TempDir()
	base := filepath.Join(root, "myrepo")
	if err := os.MkdirAll(filepath.Join(base, "sub"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(base, "main.go"), []byte("package main\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(base, "sub", "query.sql"), []byte("SELECT 1;"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(base, ".hidden.go"), []byte("package hidden\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	extSet := map[string]bool{".go": true, ".sql": true}
	entries, err := utils.GatherFilesFromRoot(base, extSet, false)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	foundMain := false
	foundSQL := false
	for _, e := range entries {
		if filepath.Base(e.AbsPath) == "main.go" {
			foundMain = true
		}
		if filepath.Base(e.AbsPath) == "query.sql" {
			foundSQL = true
		}
		if len(e.DisplayPath) == 0 {
			t.Fatalf("empty display path")
		}
	}
	if !foundMain || !foundSQL {
		t.Fatalf("missing expected files: main.go=%v query.sql=%v", foundMain, foundSQL)
	}
}

func TestUnzipReaderToDir(t *testing.T) {
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	f1, err := zw.Create("folder/a.txt")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := io.WriteString(f1, "hello"); err != nil {
		t.Fatal(err)
	}
	if err := zw.Close(); err != nil {
		t.Fatal(err)
	}
	dest := t.TempDir()
	if err := utils.UnzipReaderToDir(bytes.NewReader(buf.Bytes()), dest); err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(filepath.Join(dest, "folder", "a.txt"))
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "hello" {
		t.Fatalf("unexpected content: %s", string(got))
	}
}
