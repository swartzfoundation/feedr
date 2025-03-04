.PHONY: dev build clean run test build-linux build-windows build-mac build-arm build-arm64 build-frontend build-all

now=$(shell date +%s)
GOCMD=go
LDFLAGS=-ldflags "-X 'main.Version=v0.0.1' -X 'main.BuildTime=$(now)' -s -w"
GOBUILD=$(GOCMD) build -v $(LDFLAGS)
MAIN_FILE=main.go

BINARY_NAME=feedr

dev:
	$(GOCMD) run $(MAIN_FILE)

build:
	$(GOBUILD) -o bin/$(BINARY_NAME)

clean:
	rm -f bin/$(BINARY_NAME)

run:
	$(GOBUILD) -o bin/$(BINARY_NAME)
	./bin/$(BINARY_NAME)

test:
	$(GOCMD) test -v ./...

build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o bin/$(BINARY_NAME)

build-windows:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GOBUILD) -o bin/$(BINARY_NAME).exe

build-mac:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GOBUILD) -o bin/$(BINARY_NAME)

build-arm:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm $(GOBUILD) -o bin/$(BINARY_NAME)

build-arm64:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 $(GOBUILD) -o bin/$(BINARY_NAME)


build-frontend:
	cd frontend && pnpm build

build-all: build-linux build-windows build-mac build-rpi build-rpi2 build-rpi3 build-rpi4 build-arm build-arm64

help:
	@echo "Please use 'make <target>' where <target> is one of"
	@echo "  build         to build the binary"
	@echo "  clean         to remove the binary"
	@echo "  run           to build and run the binary"
	@echo "  test          to run tests"
	@echo "  build-linux   to build the binary for Linux"
	@echo "  build-windows to build the binary for Windows"
	@echo "  build-mac     to build the binary for Mac"
	@echo "  build-arm     to build the binary for ARM"
	@echo "  build-arm64   to build the binary for ARM64"
	@echo "  build-frontend to build the frontend"
	@echo "  build-all     to build the binary for all platforms"
	@echo "  help          to show this help message"

# Default target
.DEFAULT_GOAL := help