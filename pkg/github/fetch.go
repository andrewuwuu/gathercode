package github

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"gathercode/pkg/utils"
)

func IsRepoURL(s string) bool {
	u, err := url.Parse(s)
	if err != nil {
		return false
	}
	host := strings.ToLower(u.Host)
	if host == "" {
		return false
	}
	return strings.Contains(host, "github.com") || strings.Contains(host, "gitlab.com") || strings.Contains(host, "bitbucket.org")
}

func FetchRepo(repoURL, branch string) (string, string, error) {
	tempDir, err := os.MkdirTemp("", "gath_repo_*")
	if err != nil {
		return "", "", err
	}
	if err := tryGitClone(repoURL, tempDir, branch); err == nil {
		repoName := filepath.Base(tempDir)
		return tempDir, repoName, nil
	}
	if err := tryDownloadZip(repoURL, tempDir, branch); err == nil {
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
	_ = os.RemoveAll(tempDir)
	return "", "", errors.New("failed to fetch repo")
}

func tryGitClone(repoURL, destDir, branch string) error {
	gitPath, err := exec.LookPath("git")
	if err != nil {
		return fmt.Errorf("git not found: %w", err)
	}
	args := []string{"clone", "--depth", "1"}
	if branch != "" {
		args = append(args, "--branch", branch)
	}
	args = append(args, repoURL, destDir)
	cmd := exec.Command(gitPath, args...)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("git clone failed: %w", err)
	}
	return nil
}

func tryDownloadZip(repoURL, destDir, branch string) error {
	u, err := url.Parse(repoURL)
	if err != nil {
		return err
	}
	parts := strings.Split(strings.Trim(u.Path, "/"), "/")
	if len(parts) < 2 {
		return fmt.Errorf("could not parse owner/repo from %s", repoURL)
	}
	owner, repo := parts[0], parts[1]
	candidates := []string{}
	if branch != "" {
		candidates = append(candidates, branch)
	} else {
		candidates = append(candidates, "main", "master")
	}
	var lastErr error
	for _, br := range candidates {
		zipURL := "https://github.com/" + owner + "/" + repo + "/archive/refs/heads/" + br + ".zip"
		resp, err := http.Get(zipURL)
		if err != nil {
			lastErr = err
			continue
		}
		if resp.StatusCode != 200 {
			lastErr = fmt.Errorf("http status %d", resp.StatusCode)
			resp.Body.Close()
			continue
		}
		if err := utils.UnzipReaderToDir(resp.Body, destDir); err != nil {
			resp.Body.Close()
			lastErr = err
			continue
		}
		resp.Body.Close()
		return nil
	}
	return lastErr
}