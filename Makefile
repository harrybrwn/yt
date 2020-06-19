GO=go
GOTEST=go test -v -cover
INSTALL=go install

all: clean test

install:
	# $(INSTALL) github.com/harrybrwn/yt
	go install

test:
	go test -v ./... -coverprofile=coverage.txt -covermode=atomic

clean:
	go clean -testcache
	go clean -i
	$(RM) coverage.txt *.a dist -r

build:
	go build -o youtube.a ./youtube
	go build -o cmd.a ./cmd

dist:
	goreleaser release \
		--rm-dist \
		--parallelism 8

snapshot:
	goreleaser release --skip-publish --snapshot

.PHONY: all build test clean
