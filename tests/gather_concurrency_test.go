package tests

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"gathercode/pkg/gather"
)

func TestCollectConcurrent(t *testing.T) {
	tmp := t.TempDir()
	n := 20
	repos := make([]string, 0, n)

	for i := 0; i < n; i++ {
		rp := filepath.Join(tmp, fmt.Sprintf("repo_%02d", i))
		if err := os.MkdirAll(filepath.Join(rp, "sub"), 0o755); err != nil {
			t.Fatal(err)
		}
		f1 := filepath.Join(rp, "main.go")
		f2 := filepath.Join(rp, "sub", "q.sql")
		if err := os.WriteFile(f1, []byte("package main\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(f2, []byte("SELECT 1;"), 0o644); err != nil {
			t.Fatal(err)
		}
		repos = append(repos, rp)
	}

	exts := []string{".go", ".sql"}
	opts := gather.Options{IncludeHidden: false}
	maxWorkers := runtime.GOMAXPROCS(0)
	ctx := context.Background()

	entries, err := gather.CollectConcurrent(ctx, repos, exts, opts, maxWorkers)
	if err != nil {
		t.Fatal(err)
	}

	expected := n * 2
	if len(entries) != expected {
		t.Fatalf("expected %d entries, got %d", expected, len(entries))
	}

	foundMain := 0
	foundSQL := 0
	for _, e := range entries {
		base := filepath.Base(e.AbsPath)
		if base == "main.go" {
			foundMain++
		}
		if base == "q.sql" {
			foundSQL++
		}
	}

	if foundMain != n || foundSQL != n {
		t.Fatalf("unexpected counts main=%d sql=%d expected=%d", foundMain, foundSQL, n)
	}
}