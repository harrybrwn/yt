COVER_FILE=test-coverage
COVER=go tool cover


all: test build clean

build:
	go install yt

test:
	go test -v ./... -coverprofile=$(COVER_FILE)
	$(COVER) -func=$(COVER_FILE)
	$(COVER) -html=$(COVER_FILE) -o coverage.html

clean:
	rm $(COVER_FILE)