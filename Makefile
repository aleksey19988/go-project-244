init:
	go mod tidy
build: # Сборка бинарника
	go build -o bin/gendiff ./cmd/gendiff;
lint:
	golangci-lint run
test:
	go test ./...