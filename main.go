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

	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(1)
	}

	// Get owner and repo name
	owner, repo := parseOwnerRepo(os.Args[1], http)

	// Get details from GitHub
	repository := *fetchRepoDetails(owner, repo)

	var originURL, upstreamURL string
	if *http {
		originURL = *repository.CloneURL
		if *repository.Fork && repository.Parent != nil {
			upstreamURL = *repository.Parent.CloneURL
		}
	} else {
		originURL = *repository.SSHURL
		if *repository.Fork && repository.Parent != nil {
			upstreamURL = *repository.Parent.SSHURL
		}
	}

	cloneCmd := exec.Command("git", "clone", originURL)
	cloneCmd.Stdout = os.Stdout
	cloneCmd.Stderr = os.Stderr
	if err := cloneCmd.Run(); err != nil {
		exitOnError(err)
	}

	if *repository.Fork {
		fmt.Printf("setting fork remote (%s): %s\n", *remoteName, upstreamURL)
		cmd := exec.Command("git", "remote", "add", *remoteName, upstreamURL)
		cmd.Dir = repo
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			exitOnError(err)
		}
	}
}

func parseOwnerRepo(arg string, http *bool) (owner, repo string) {
	// See if it's a URL
	url, err := url.ParseRequestURI(arg)
	if err == nil {
		if url.Scheme == "https" || url.Scheme == "http" {
			*http = true
		}

		parts := strings.Split(url.Path, "/")
		owner = parts[1]
		repo = strings.Replace(parts[2], ".git", "", 1)
	} else {
		parts := strings.Split(arg, "/")
		if len(parts) != 2 {
			flag.Usage()
			os.Exit(1)
		}

		owner = parts[0]
		repo = parts[1]
	}

	return
}

func fetchRepoDetails(owner, repo string) *github.Repository {
	ctx := context.Background()
	client := github.NewClient(nil)
	repository, _, err := client.Repositories.Get(ctx, owner, repo)
	if err != nil {
		fmt.Printf("Could not fetch %s/%s: %+v", owner, repo, err)
		os.Exit(1)
	}

	return repository
}

func exitOnError(err error) {
	if err != nil {
		if exit, ok := err.(*exec.ExitError); ok {
			os.Exit(exit.ExitCode())
		} else {
			os.Exit(1)
		}
	}
}
