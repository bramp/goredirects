.PHONY: all format analyze test test-ci fix upgrade

all: format analyze test

format:
	go fmt ./... && goimports -w .

analyze:
	go vet ./... && staticcheck ./...

test:
	go test ./...

test-ci:
	go test -v ./...

fix:
	go fmt ./... && go fix ./...

upgrade:
	go mod tidy && go get -u ./... && go mod tidy
