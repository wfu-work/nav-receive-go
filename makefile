# Makefile for the navReceiveApp project

GOCMD=CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOMOD=$(GOCMD) mod
BINARY_NAME=navReceiveApp
PACKAGE=./

# Default target executed when no arguments are given to make.
all: build

# Build the binary
build:
	$(GOBUILD) -o ./docker/$(BINARY_NAME) $(PACKAGE)

# Clean the build files
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

# Run the application
run: build
	./$(BINARY_NAME)

# Tidy up the dependencies
tidy:
	$(GOMOD) tidy

.PHONY: all build clean run tidy