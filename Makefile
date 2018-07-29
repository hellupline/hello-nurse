# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get


BINARY_NAME=hello-nurse


all: build

build:
	CGO_ENABLED=0 $(GOBUILD) -o $(BINARY_NAME) -v ./...
run: build
	./$(BINARY_NAME)
test:
	CGO_ENABLED=0 $(GOTEST) -v ./...
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)


# Cross compilation
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_NAME)_linux -v ./...

build-darwin:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(BINARY_NAME)_darwin -v ./...

build-windows:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(BINARY_NAME)_windows.exe -v ./...

build-docker:
	docker run --rm -it -v "${PWD}:/go/src/github.com/hellupline/hello-nurse" -w "/go/src/github.com/hellupline/hello-nurse" -e "CGO_ENABLED=0" golang make build
