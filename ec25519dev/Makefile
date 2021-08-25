.PHONY: all

# Go parameters
BINARY=ecdh-dev

all: clean build

build:
	CGO_ENABLED=0  go build -ldflags "-s -w" -o $(BINARY) 

clean:
	go clean
	rm -f $(BINARY)

 
