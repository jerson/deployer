BUILD?=go build -ldflags="-w -s"
ENV?=dev

default: test

deps:
	go mod download

test:
	go test ./...

format:
	go fmt ./...

vet:
	go vet ./...

lint:
	golint -set_exit_status -min_confidence 0.3 ./...