build:
	go build -o bin/gobit main.go

install:
	go install -v ./...

help:
	go run main.go help