# yt

[![Build Status](https://travis-ci.com/harrybrwn/yt.svg?branch=master)](https://travis-ci.com/harrybrwn/yt)
[![codecov](https://codecov.io/gh/harrybrwn/yt/branch/master/graph/badge.svg)](https://codecov.io/gh/harrybrwn/yt)
[![Go Report Card](https://goreportcard.com/badge/github.com/harrybrwn/yt)](https://goreportcard.com/report/github.com/harrybrwn/yt)
[![GoDoc](https://godoc.org/github.com/github.com/harrybrwn/yt?status.svg)](https://godoc.org/github.com/harrybrwn/yt)

A cli for downloading youtube videos


### Installation
#### Homebrew
```
brew install harrybrwn/tap/yt
```
#### Snap
```
snap install go-yt
alias yt=go-yt
```
#### With Go
```
go install github.com/harrybrwn/yt
```

### Examples
```sh
yt video https://www.youtube.com/watch?v=1234
yt video 1234 # same result with the same video id
```

### Completion
#### zsh
```
source <(yt completion zsh)
compdef _yt yt
```
#### bash
```
source <(yt completion bash)
```
