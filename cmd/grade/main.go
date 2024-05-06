package main

import (
	"cso/codecowboy/grade"
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
	defer db.Close()

	graderFunc, ok := grade.Graders[*graderType]
	if !ok {
		log.Error("Unknown grader type: ", *graderType)
	}

	checkErr(graderFunc(db).Grade(*dir, *course, *assignment, *out))
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
