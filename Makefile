build:
	go build -o bin/web-go main.go

run: build
	./bin/web-go
