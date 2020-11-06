.PHONY: build
build:
	go build -v -o start

.PHONY: test
test:
	go test -v -timeout 30s ./...

.DEFAULT_GOAL := build