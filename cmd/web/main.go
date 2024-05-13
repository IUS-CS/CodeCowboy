package main

import (
	"cso/codecowboy/store"
	"cso/codecowboy/web"
	"flag"
	"github.com/charmbracelet/log"
)

const DBNAME = "codecowboy"

var addr = flag.String("addr", "127.0.0.1:8080", "http service address")

func main() {
	flag.Parse()

	db, err := store.New(DBNAME)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	log.Infof("Listening at http://%s", *addr)
	log.Fatalf("Error: %e", web.New(db, *addr).ListenAndServe())
}
