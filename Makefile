.PHONY: build install clean test fmt vet

BINARY_NAME=sact
VERSION?=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS=-ldflags "-X main.version=${VERSION}"

build:
	go build ${LDFLAGS} -o ${BINARY_NAME} ./cmd/sact

install:
	go install ${LDFLAGS} ./cmd/sact

clean:
	go clean
	rm -f ${BINARY_NAME}

test:
	go test -v ./...

fmt:
	go fmt ./...

vet:
	go vet ./...

run: build
	./${BINARY_NAME}

all: fmt vet build
