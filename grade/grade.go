package grade

import (
	"cso/codecowboy/store"
	"cso/codecowboy/students"
	"strings"
)

type Grader interface {
	Grade(path, course, assignment, out string) error
}

func sisNameFromDirName(students *students.Students, dirName string) string {
	fields := strings.Split(dirName, "-")
	sName := fields[len(fields)-1]
	for _, s := range students.Members {
		if s.GitHubUsername == sName {
			return s.SISLoginID
		}
	}
	return "Unknown"
}

type GraderFunc func(*store.DB) Grader

var Graders = map[string]GraderFunc{
	"go": NewGoGrader,
}
