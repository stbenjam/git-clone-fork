package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/google/go-github/v28/github"
	"net/url"
	"os"
	"os/exec"
	"strings"
)

func main() {
	remoteName := flag.String("r", "upstream", "Name of the remote to use")
	http := flag.Bool("h", false, "Force HTTP")
	flag.Parse()

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage:\n\tgit clone-fork <user>/<repo>|URI\n\n")
		flag.PrintDefaults()
	}

	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(1)
	}

	// Get owner and repo name
	owner, repo, err := parseOwnerRepo(flag.Arg(0), http)
	if err != nil {
		fmt.Println(err)
		flag.Usage()
		os.Exit(1)
	}

	// Get details from GitHub
	repository, err := fetchRepoDetails(owner, repo)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var originURL, upstreamURL *string
	if *http {
		originURL = repository.CloneURL
		if *repository.Fork && repository.Parent != nil {
			upstreamURL = repository.Parent.CloneURL
		}
	} else {
		originURL = repository.SSHURL
		if *repository.Fork && repository.Parent != nil {
			upstreamURL = repository.Parent.SSHURL
		}
	}

	cloneCmd := exec.Command("git", "clone", *originURL)
	cloneCmd.Stdout = os.Stdout
	cloneCmd.Stderr = os.Stderr
	if err := cloneCmd.Run(); err != nil {
		os.Exit(1)
	}

	if upstreamURL != nil && *upstreamURL != "" {
		fmt.Printf("setting fork remote (%s): %s\n", *remoteName, *upstreamURL)
		cmd := exec.Command("git", "remote", "add", *remoteName, *upstreamURL)
		cmd.Dir = repo
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			os.Exit(1)
		}
	}
}

func parseOwnerRepo(arg string, http *bool) (owner, repo string, err error) {
	errFmtString := "error: expected repo in form <user>/<repo> or a URI. Got: \"%s\""

	// Make SSH URL URI parseable
	arg = strings.Replace(arg, "git@github.com:", "git://github.com/", 1)

	// See if it's a URL, or in the format <user>/<repo>
	url, urlErr := url.ParseRequestURI(arg)
	if urlErr == nil {
		if url.Scheme == "https" || url.Scheme == "http" {
			*http = true
		}

		parts := strings.Split(url.Path, "/")
		if len(parts) < 3 {
			err = fmt.Errorf(errFmtString, arg)
			return
		}

		owner = parts[1]
		repo = strings.Replace(parts[2], ".git", "", 1)
	} else {
		parts := strings.Split(arg, "/")
		if (len(parts) != 2) || (parts[0] == "" || parts[1] == "") {
			err = fmt.Errorf(errFmtString, arg)
			return
		}

		owner = parts[0]
		repo = parts[1]
	}

	return
}

func fetchRepoDetails(owner, repo string) (*github.Repository, error) {
	ctx := context.Background()
	client := github.NewClient(nil)
	repository, _, err := client.Repositories.Get(ctx, owner, repo)

	return repository, err
}
