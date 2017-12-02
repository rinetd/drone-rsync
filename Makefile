.PHONY: all clean deps fmt vet test docker

EXECUTABLE ?= drone-rsync
IMAGE ?= rinetd/$(EXECUTABLE)
COMMIT ?= $(shell git rev-parse --short HEAD)

LDFLAGS = -X "main.buildCommit=$(COMMIT)"
PACKAGES = $(shell go list ./... | grep -v /vendor/)

all: deps build docker

clean:
	go clean -i ./...

deps:
	go get -t ./...

fmt:
	go fmt $(PACKAGES)

vet:
	go vet $(PACKAGES)

test:
	@for PKG in $(PACKAGES); do go test -cover -coverprofile $$GOPATH/src/$$PKG/coverage.out $$PKG || exit 1; done;

docker:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags '-s -w'
	upx $(EXECUTABLE)
	docker build --rm -t $(IMAGE) .
	docker push $(IMAGE)
$(EXECUTABLE): $(wildcard *.go)
	go build -ldflags '-s -w $(LDFLAGS)'

build: $(EXECUTABLE)
