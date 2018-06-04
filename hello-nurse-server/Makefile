GOCMD=vgo
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GORUN=$(GOCMD) run
GOCLEAN=$(GOCMD) clean
BINARY_NAME=hello-nurse
BINARY_UNIX=$(BINARY_NAME)_unix

all: build

build:
	$(GOBUILD) -o $(BINARY_NAME) -v ./...
run: build
	./$(BINARY_NAME)

test:
	$(GOTEST) -v ./...
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)


# Cross compilation
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v
docker-build:
	docker run --rm -it -v "$(GOPATH)":/go -w /go/src/github.com/hellupline/hello-nurse/server golang:latest vgo build -o "$(BINARY_UNIX)" -v
