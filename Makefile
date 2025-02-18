all: generate server

generate:
	go generate

server: generate
	go build -o server cmd/web/main.go

run: all
	./server

@PHONY: server generate run
