package tests

import "testing"

import "gathercode/pkg/github"

func TestIsRepoURL(t *testing.T) {
	cases := map[string]bool{
		"https://github.com/owner/repo":    true,
		"https://gitlab.com/owner/repo":    true,
		"https://bitbucket.org/owner/repo": true,
		"/home/user/repo":                  false,
		"file:///tmp/repo":                 false,
		"not a url":                        false,
	}
	for input, want := range cases {
		got := github.IsRepoURL(input)
		if got != want {
			t.Fatalf("IsRepoURL(%q) = %v; want %v", input, got, want)
		}
	}
}
