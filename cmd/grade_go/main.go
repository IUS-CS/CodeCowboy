package main

import (
	"cso/codecowboy/store"
	"cso/codecowboy/students"
	"fmt"
	"github.com/charmbracelet/log"
)

const DBNAME = "codecowboy"

func main() {
	db, err := store.New(DBNAME)
	if err != nil {
		log.Fatal(err)
	}
	students := students.New(db)
	fmt.Println(students.String())
}
