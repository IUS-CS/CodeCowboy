package main

import (
	"cso/codecowboy/canvasfmt"
	"cso/codecowboy/classroom"
	"cso/codecowboy/githubfmt"
	"cso/codecowboy/store"
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
	assignPath = flag.String("assignments", "", "Assignments JSON")
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

	cls, err := classroom.New(db, *course)
	if err != nil {
		log.Fatal(err)
	}
	roster := cls.Students

	if *ghPath != "" {
		roster, err = githubfmt.ParseFile(*ghPath, roster)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		log.Debug("Not importing GitHub")
	}

	if *canvasPath != "" {
		roster, err = canvasfmt.ParseFile(*canvasPath, roster)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		log.Debug("Not importing Canvas export")
	}

	if *assignPath != "" {
		assignments, err := classroom.ParseAssignmentsFile(*assignPath, cls.Name)
		if err != nil {
			log.Fatal(err)
		}
		cls.Assignments = assignments
	}

	cls.Students = roster
	err = cls.Save()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%d students saved\n", len(roster))
}
