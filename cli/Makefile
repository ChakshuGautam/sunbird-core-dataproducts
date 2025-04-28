.PHONY: build install clean test

BINARY_NAME=sunbird-cli
INSTALL_PATH=/usr/local/bin

build:
	go build -o $(BINARY_NAME) -v

install: build
	cp $(BINARY_NAME) $(INSTALL_PATH)

clean:
	go clean
	rm -f $(BINARY_NAME)

test:
	go test -v ./...

.DEFAULT_GOAL := build
