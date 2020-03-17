GO=go
GOTEST=go test -v -cover
INSTALL=go install

all: build clean

build:
	$(INSTALL) github.com/harrybrwn/yt

test:
	$(GOTEST) ./...

clean:
	$(GO) clean -testcache

.PHONEY: all build test clean
