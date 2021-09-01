.PHONY: all

# Go parameters
BINARY=fclient

all: clean build

build:
	go mod download
	CGO_ENABLED=1  go build -ldflags "-s -w" -o $(BINARY) 

clean:
	go clean
	go fmt ./...
	rm -f $(BINARY)

 
