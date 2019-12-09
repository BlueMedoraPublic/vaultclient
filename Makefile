all: fmt lint test

test:
	go test ./...

lint:
	golint ./...

fmt:
	go fmt ./...
