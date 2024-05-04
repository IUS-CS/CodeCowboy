package main

import (
	"cso/codecowboy/store"
	"cso/codecowboy/students"
	"flag"
	"fmt"
	"github.com/charmbracelet/log"
)

const DBNAME = "codecowboy"

var (
	course     = flag.String("course", "", "Course to grade")
	assignment = flag.String("assignment", "", "Assignment to grade")
)

func main() {
	flag.Parse()
	db, err := store.New(DBNAME)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	students := students.New(db, *course)
	fmt.Println(students.String())
}
