package graders

import (
	"cso/codecowboy/classroom"
	"cso/codecowboy/graders/golang"
	"cso/codecowboy/graders/java"
	"cso/codecowboy/graders/net"
	"cso/codecowboy/store"
)

type Grader interface {
	Grade(spec classroom.AssignmentSpec, out string) error
}

func GetGrader(language string, db *store.DB) Grader {
	switch language {
	case "go":
		return golang.NewGoGrader(db)
	case "java":
		return java.NewJavaGrader(db)
	case "net":
		return net.NewNetGrader(db)
	}
	return nil
}

type TestResult func(stdOut string) (float64, float64, float64, error)
