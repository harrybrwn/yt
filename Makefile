GO=go

all: clean test

install:
	go install

test:
	go test -v ./... -coverprofile=coverage.txt -covermode=atomic

clean:
	go clean -testcache
	go clean -i
	$(RM) gocoverage.txt coverage.txt *.a dist -r

build:
	go build -o youtube.a ./youtube
	go build -o cmd.a ./cmd

dist:
	if [ ! -f release-notes.md ]; then echo 'no release notes'; exit 1; fi
	goreleaser release \
		--rm-dist \
		--parallelism 8 \
		--release-notes release-notes.md

docs:
	go run docs/gen.go docs

snapshot:
	goreleaser release --skip-publish --snapshot

.PHONY: all build test clean docs
