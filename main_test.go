package main

import "testing"

func Test_parseOwnerRepo(t *testing.T) {
	valid := []string{
		"stbenjam/git-clone-fork",
		"git@github.com:stbenjam/git-clone-fork.git",
		"https://github.com/stbenjam/git-clone-fork.git",
		"https://github.com/stbenjam/git-clone-fork/",
	}

	var http bool
	for _, repo := range valid {
		owner, repo, err := parseOwnerRepo(repo, &http)
		if err != nil {
			t.Error(err)
			return
		}

		if owner != "stbenjam" {
			t.Errorf("owner incorrect, got \"%s\", expected stbenjam", owner)
		}

		if repo != "git-clone-fork" {
			t.Errorf("repo incorrect, got \"%s\", expected git-clone-fork", repo)
		}
	}

	invalid := []string{
		"invalid",
		"user/repo/something-else",
		"justuser/",
		"/justrepo",
	}

	for _, repo := range invalid {
		owner, repo, err := parseOwnerRepo(repo, &http)
		if err == nil {
			t.Errorf("expected error, but parse succeeded. user: %s, repo: %s", owner, repo)
		}

	}
}
