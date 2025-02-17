all: generate server

generate:
	go generate

server: generate
	go build -o server cmd/web/main.go

test: 
	go test -v ./... | grep -v '\[no test files\]' 

@PHONY: server generate
