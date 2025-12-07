BINARY_NAME := gather
CMD_PATH := ./cmd/gather
OUT := $(BINARY_NAME)

.PHONY: all build run test clean tidy fmt vet

all: build

build:
	go build -o $(OUT) $(CMD_PATH)

run:
	go run $(CMD_PATH)

test:
	go test ./... -v

tidy:
	go mod tidy

fmt:
	go fmt ./...

vet:
	go vet ./...

clean:
	rm -f $(OUT)