package tests

import (
	"archive/zip"
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"gathercode/pkg/github"
)

func TestFetchRepoFromZipURL(t *testing.T) {
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	f1, err := zw.Create("repo-main/main.go")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := io.WriteString(f1, "package main\n"); err != nil {
		t.Fatal(err)
	}
	if err := zw.Close(); err != nil {
		t.Fatal(err)
	}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/zip")
		_, _ = w.Write(buf.Bytes())
	}))
	defer ts.Close()
	repoRoot, repoName, err := github.FetchRepoFromZipURL(ts.URL)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(repoRoot)
	mainPath := filepath.Join(repoRoot, "main.go")
	b, err := os.ReadFile(mainPath)
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != "package main\n" {
		t.Fatalf("unexpected content: %s", string(b))
	}
	if repoName == "" {
		t.Fatalf("expected repo name")
	}
}