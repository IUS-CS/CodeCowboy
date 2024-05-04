package main

import (
	"cso/codecowboy/githubfmt"
	"cso/codecowboy/store"
	"cso/codecowboy/students"
	"flag"
	"github.com/charmbracelet/log"
)

const DBNAME = "codecowboy"

var (
	course = flag.String("course", "", "Course to grade")
	path   = flag.String("path", "classroom_roster.csv", "Path to GitHub export")
)

func main() {
	flag.Parse()
	db, err := store.New(DBNAME)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	roster, err := githubfmt.Parse(*path)
	if err != nil {
		log.Fatal(err)
	}
	list := students.NewFromList(db, *course, roster)
	err = list.Save()
	if err != nil {
		log.Fatal(err)
	}
}
