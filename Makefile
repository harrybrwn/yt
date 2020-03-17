GOTEST=go test -v -cover


all: build clean

build:
	go install github.com/harrybrwn/yt

test:
	$(GOTEST) ./...
