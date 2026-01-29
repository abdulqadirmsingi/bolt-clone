.PHONY: build run proto tidy test

build:
	go build -o bin/server ./cmd/server

run: build
	./bin/server

proto:
	@which protoc >/dev/null || (echo "install protoc and protoc-gen-go, protoc-gen-go-grpc"; exit 1)
	protoc --go_out=internal/proto/gen --go_opt=paths=source_relative \
		--go-grpc_out=internal/proto/gen --go-grpc_opt=paths=source_relative \
		-I internal/proto internal/proto/*.proto

tidy:
	go mod tidy

test:
	go test ./...
