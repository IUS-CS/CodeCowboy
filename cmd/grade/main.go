package main

import (
	"cso/codecowboy/classroom"
	"cso/codecowboy/graders"
	"cso/codecowboy/store"
	"flag"
	"github.com/charmbracelet/log"
	"os"
)

const DBNAME = "codecowboy"

var (
	course     = flag.String("course", "", "Course to grade")
	assignment = flag.String("assignment", "", "Assignment to grade")
	dir        = flag.String("dir", "", "Directory to grade")
	debug      = flag.Bool("debug", false, "Enable debug mode")
	out        = flag.String("output", "stdout", "Output file")
	graderType = flag.String("type", "go", "Language: go, java, net")
)

func main() {
	flag.Parse()
	if *debug {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.ErrorLevel)
	}

	if *dir == "" || *course == "" {
		flag.Usage()
		os.Exit(1)
	}

	db, err := store.New(DBNAME)
	checkErr(err)

	grader := graders.GetGrader(*graderType, db)
	if grader == nil {
		log.Error("Unknown grader type: ", *graderType)
	}

	w := os.Stdout
	if *out != "stdout" && *out != "" {
		w, err = os.Open(*out)
		checkErr(err)
	}

	checkErr(grader.Grade(classroom.AssignmentSpec{
		Name:   *assignment,
		Path:   *dir,
		Course: *course,
	}, w))
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
