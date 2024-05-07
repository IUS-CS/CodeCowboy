package main

import (
	"cso/codecowboy/web"
	"flag"
	"github.com/charmbracelet/log"
)

var addr = flag.String("addr", "127.0.0.1:8080", "http service address")

func main() {
	flag.Parse()
	log.Infof("Listening at http://%s", *addr)
	log.Fatalf("Error: %e", web.New(*addr).ListenAndServe())
}
