package main

import (
	"cso/codecowboy/canvasfmt"
	"cso/codecowboy/githubfmt"
	"cso/codecowboy/store"
	"cso/codecowboy/students"
	"flag"
	"fmt"
	"github.com/charmbracelet/log"
)

const DBNAME = "codecowboy"

var (
	course     = flag.String("course", "", "Course to grade")
	ghPath     = flag.String("ghpath", "", "Path to GitHub export")
	canvasPath = flag.String("canvaspath", "", "Path to Canvas export")
	debug      = flag.Bool("debug", false, "Enable debug mode")
)

func main() {
	flag.Parse()

	if *debug {
		log.SetLevel(log.DebugLevel)
	}

	db, err := store.New(DBNAME)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	roster := students.New(db, *course).Members

	if *ghPath != "" {
		roster, err = githubfmt.Parse(*ghPath, roster)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		log.Debug("Not importing GitHub")
	}

	if *canvasPath != "" {
		roster, err = canvasfmt.Parse(*canvasPath, roster)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		log.Debug("Not importing Canvas export")
	}

	list := students.NewFromList(db, *course, roster)
	err = list.Save()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%d students saved\n", len(roster))
}
