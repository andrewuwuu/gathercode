package utils

import (
	"archive/zip"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type FileEntry struct {
	DisplayPath string
	AbsPath     string
}

func PathExists(p string) bool {
	_, err := os.Stat(p)
	return err == nil
}

func GatherFilesFromRoot(rootPath string, extSet map[string]bool, includeHidden bool) ([]FileEntry, error) {
	var out []FileEntry
	rootAbs, err := filepath.Abs(rootPath)
	if err != nil {
		return nil, err
	}
	basePrefix := filepath.Base(rootAbs)
	err = filepath.WalkDir(rootAbs, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return nil
		}
		rel, _ := filepath.Rel(rootAbs, path)
		parts := strings.Split(rel, string(os.PathSeparator))
		for _, p := range parts {
			if p == "." || p == "" {
				continue
			}
			if strings.HasPrefix(p, ".") && !includeHidden {
				if d.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}
		}
		if d.IsDir() {
			return nil
		}
		ext := strings.ToLower(filepath.Ext(d.Name()))
		if !extSet[ext] {
			return nil
		}
		display := filepath.ToSlash(filepath.Join(basePrefix, rel))
		out = append(out, FileEntry{DisplayPath: display, AbsPath: path})
		return nil
	})
	if err != nil {
		return nil, err
	}
	return out, nil
}

func UnzipReaderToDir(r io.Reader, destDir string) error {
	tmp, err := os.CreateTemp("", "repozip_*.zip")
	if err != nil {
		return err
	}
	tmpPath := tmp.Name()
	_, err = io.Copy(tmp, r)
	tmp.Close()
	if err != nil {
		_ = os.Remove(tmpPath)
		return err
	}
	zr, err := zip.OpenReader(tmpPath)
	if err != nil {
		_ = os.Remove(tmpPath)
		return err
	}
	defer func() {
		_ = zr.Close()
		_ = os.Remove(tmpPath)
	}()
	for _, f := range zr.File {
		target := filepath.Join(destDir, f.Name)
		cleanTarget := filepath.Clean(target)
		cleanDest := filepath.Clean(destDir)
		if !strings.HasPrefix(cleanTarget, cleanDest+string(os.PathSeparator)) && cleanTarget != cleanDest {
			return fmt.Errorf("illegal file path: %s", target)
		}
		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(target, 0o755); err != nil {
				return err
			}
			continue
		}
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			return err
		}
		outf, err := os.OpenFile(target, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}
		rc, err := f.Open()
		if err != nil {
			outf.Close()
			return err
		}
		_, err = io.Copy(outf, rc)
		outf.Close()
		rc.Close()
		if err != nil {
			return err
		}
	}
	return nil
}