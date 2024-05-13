all: generate server

generate:
	go generate

server: generate
	go build -o server cmd/web/main.go

@PHONY: server generate
