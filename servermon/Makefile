.PHONY: all

# Go parameters
BINARY=servermon

all: clean build

build:
	CGO_ENABLED=0  go build -ldflags "-s -w" -o $(BINARY) 

clean:
	go fmt ./...
	rm -f $(BINARY)

 
