# Makefile for warp

all: warp

warp: *.go
	godep go build -o warp

fmt:
	go fmt ./...

.PHONY: clean
clean:
	rm -f warp
	go clean -r

test:
	godep go test

.PHONY: secret
secret:
	@openssl rand -base64 64 | sed ':a;N;$$!ba;s/\n//g'