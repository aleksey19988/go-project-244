init:
	go mod tidy
build: # Сборка бинарника
	go build -o bin/gendiff ./cmd/gendiff;