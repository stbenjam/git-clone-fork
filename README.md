[![Build Status](https://travis-ci.org/stbenjam/git-clone-fork.svg?branch=master)](https://travis-ci.org/stbenjam/git-clone-fork) [![License](https://img.shields.io/badge/license-MIT-blue)](https://opensource.org/licenses/Apache-2.0) ![GitHub release](https://img.shields.io/github/release/stbenjam/git-clone-fork.svg)


# git clone-fork

Clones a repository, and if it is a fork of a GitHub repository, sets
that repo as a remote named 'upstream'.

## Motivation

Often, I find myself repeating the same commands over and over:

  1. git clone http://github.com/stbenjam/foo.git
  2. cd foo
  3. git remote add upstream http://github.com/bigcorp/foo.git

This automates that process, by querying the GitHub API to find out if
the repo is a fork, and if so, creating a remote for it.

## Installation

1. Install [go](https://golang.org/doc/install)
2. Run `go get -u github.com/stbenjam/git-clone-fork`
3. Call as `git clone-fork` or `git-clone-fork`

## Usage

```
  git clone-fork user/repo|<clone url>

  -h Forces use of HTTP
  -r [name] changes the name of 'upstream' to something else
```

## Examples

### Standard

```
$ git clone-fork stbenjam/installer
Cloning into 'installer'...
remote: Enumerating objects: 8, done.
remote: Counting objects: 100% (8/8), done.
remote: Compressing objects: 100% (8/8), done.
remote: Total 88745 (delta 0), reused 2 (delta 0), pack-reused 88737
Receiving objects: 100% (88745/88745), 68.22 MiB | 19.48 MiB/s, done.
Resolving deltas: 100% (53901/53901), done.
setting fork remote (upstream): git@github.com:openshift/installer.git
$ cd installer/
$ git remote -v
origin	git@github.com:stbenjam/installer.git (fetch)
origin	git@github.com:stbenjam/installer.git (push)
upstream	git@github.com:openshift/installer.git (fetch)
upstream	git@github.com:openshift/installer.git (push)
```

### Different remote

```
$ git clone-fork -r potato stbenjam/installer
Cloning into 'installer'...
remote: Enumerating objects: 8, done.
remote: Counting objects: 100% (8/8), done.
remote: Compressing objects: 100% (8/8), done.
remote: Total 88745 (delta 0), reused 2 (delta 0), pack-reused 88737
Receiving objects: 100% (88745/88745), 68.22 MiB | 17.68 MiB/s, done.
Resolving deltas: 100% (53901/53901), done.
Checking out files: 100% (12457/12457), done.
setting fork remote (potato): git@github.com:openshift/installer.git
$ cd installer
$ git remote -v
origin	git@github.com:stbenjam/installer.git (fetch)
origin	git@github.com:stbenjam/installer.git (push)
potato	git@github.com:openshift/installer.git (fetch)
potato	git@github.com:openshift/installer.git (push)
```

### Force HTTP

```
$ git clone-fork -h stbenjam/installer
Cloning into 'installer'...
remote: Enumerating objects: 8, done.
remote: Counting objects: 100% (8/8), done.
remote: Compressing objects: 100% (8/8), done.
remote: Total 88745 (delta 0), reused 2 (delta 0), pack-reused 88737
Receiving objects: 100% (88745/88745), 68.22 MiB | 18.22 MiB/s, done.
Resolving deltas: 100% (53901/53901), done.
Checking out files: 100% (12457/12457), done.
setting fork remote (upstream): https://github.com/openshift/installer.git
$ cd installer
$ git remote -v
origin	https://github.com/stbenjam/installer.git (fetch)
origin	https://github.com/stbenjam/installer.git (push)
upstream	https://github.com/openshift/installer.git (fetch)
upstream	https://github.com/openshift/installer.git (push)
```

## License

See LICENSE
