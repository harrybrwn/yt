GO=go
GOTEST=go test -v -cover
INSTALL=go install

all: clean test

install:
	$(INSTALL) github.com/harrybrwn/yt

test:
	go test -v ./... -coverprofile=coverage.txt -covermode=atomic

clean:
	go clean -testcache
	go clean -i
	$(RM) coverage.txt *.a

build:
	go build -o youtube.a ./youtube
	go build -o cmd.a ./cmd

.PHONY: all build test clean
