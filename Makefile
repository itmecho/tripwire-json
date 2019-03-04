.PHONY: build

VERSION = 0.1.0

build:
	go build -ldflags "-X main.version=$(VERSION) -X main.commitHash=$$(git rev-parse --short HEAD)"