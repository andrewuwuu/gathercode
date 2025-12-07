package github

import (
	"errors"
	"net/http"
	"os"
	"path/filepath"

	"gathercode/pkg/utils"
)

func FetchRepoFromZipURL(zipURL string) (string, string, error) {
	tempDir, err := os.MkdirTemp("", "gath_repo_*")
	if err != nil {
		return "", "", err
	}
	resp, err := http.Get(zipURL)
	if err != nil {
		_ = os.RemoveAll(tempDir)
		return "", "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		_ = os.RemoveAll(tempDir)
		return "", "", errors.New("non-200 from zip url")
	}
	if err := utils.UnzipReaderToDir(resp.Body, tempDir); err != nil {
		_ = os.RemoveAll(tempDir)
		return "", "", err
	}
	children, err := os.ReadDir(tempDir)
	if err != nil {
		return tempDir, filepath.Base(tempDir), nil
	}
	for _, c := range children {
		if c.IsDir() {
			return filepath.Join(tempDir, c.Name()), c.Name(), nil
		}
	}
	return tempDir, filepath.Base(tempDir), nil
}