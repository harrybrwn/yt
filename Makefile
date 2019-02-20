PKGS=./cmd \
     ./youtube

GOTEST=go test -v -cover


all: build clean

build: test
	go install github.com/harrybrwn/yt

test:
	@for pkg in $(PKGS); do \
		$(GOTEST) $$pkg -c; \
		chmod +xw "$$pkg".test; \
	done

	@for pkg in $(PKGS); do \
		./"$$pkg".test; \
	done

clean:
	@for pkg in $(PKGS); do \
		rm "$$pkg".test; \
	done
