COVER_FILE=test-coverage
COVER=go tool cover

PKGS=./cmd ./youtube
GOTEST=go test -v -cover


all: build clean

build: test
	go install yt

test:
	@for pkg in $(PKGS); do \
		$(GOTEST) $$pkg -c; \
	done

	@for pkg in $(PKGS); do \
		./"$$pkg".test; \
	done

clean:
	# rm $(COVER_FILE)
	@for pkg in $(PKGS); do \
		rm "$$pkg".test; \
	done