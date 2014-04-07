# Makefile for warp

all: warp

warp: *.go
	godep go build -o warp

fmt:
	go fmt ./...

clean:
	rm -f warp
	go clean -r

test:
	godep go test

dist:
	./scripts/dist.sh

.PHONY: clean dist
