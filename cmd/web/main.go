package main

import (
	"cso/codecowboy/store"
	"cso/codecowboy/web"
	"flag"

	"github.com/charmbracelet/log"
)

var (
	addr   = flag.String("addr", "127.0.0.1:8080", "http service address")
	debug  = flag.Bool("debug", false, "Enable debug mode")
	dbPath = flag.String("path", "codecowboy.db", "db path")
)

func main() {
	flag.Parse()

	if *debug {
		log.SetLevel(log.DebugLevel)
	}

	db, err := store.New(*dbPath)
	if err != nil {
		log.Fatal(err)
	}

	log.Infof("Listening at http://%s", *addr)
	log.Fatalf("Error: %e", web.New(db, *addr).ListenAndServe())
}
