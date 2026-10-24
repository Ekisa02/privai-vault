BINARY=privai-server

.PHONY: build run test docker

build:
	go build -o $(BINARY) ./cmd/server

run: build
	./$(BINARY)

test:
	go test ./...

docker:
	docker build -t privai-vault:latest .
