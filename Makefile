.PHONY: test mocks build

test:
	go test -failfast -v -coverpkg=./... -coverprofile=coverage.out ./...

cover:
	go tool cover -func=coverage.out

mocks:
	go generate ./...

build:
	go build -o bin/slice cmd/main.go
