APP_NAME := folio
VERSION := 0.1.0
BUILD_TIME := $(shell date -u '+%Y-%m-%dT%H:%M:%SZ')

# Go flags
GOFLAGS := -ldflags="-s -w -X main.version=$(VERSION) -X main.buildTime=$(BUILD_TIME)"

# Targets
.PHONY: all build clean dev lint test docker docker-arm64

all: build

build:
	CGO_ENABLED=0 go build $(GOFLAGS) -o $(APP_NAME) ./cmd/folio

dev:
	go run ./cmd/folio

clean:
	rm -f $(APP_NAME)
	rm -f *.db *.db-wal *.db-shm

lint:
	go vet ./...
	gofmt -l -s .

test:
	go test -v ./...

# Cross-compilation for RPi
build-arm64:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build $(GOFLAGS) -o $(APP_NAME)-linux-arm64 ./cmd/folio

build-armv7:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7 go build $(GOFLAGS) -o $(APP_NAME)-linux-armv7 ./cmd/folio

build-amd64:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build $(GOFLAGS) -o $(APP_NAME)-linux-amd64 ./cmd/folio

build-all: build-amd64 build-arm64 build-armv7

docker:
	docker build -t $(APP_NAME):$(VERSION) -t $(APP_NAME):latest .

docker-arm64:
	docker buildx build --platform linux/arm64 -t $(APP_NAME):$(VERSION)-arm64 -t $(APP_NAME):latest-arm64 .

tidy:
	go mod tidy
