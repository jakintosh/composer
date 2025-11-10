.PHONY: build build-daemon

build: build-daemon
	go build -o bin/composer ./cmd/composer

build-daemon:
	go build -o bin/composerd ./cmd/composerd

test:
	go test ./...

rund: build-daemon
	./bin/composerd
