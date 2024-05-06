package graders

import (
	"cso/codecowboy/graders/golang"
	"cso/codecowboy/graders/java"
	"cso/codecowboy/graders/net"
	"cso/codecowboy/store"
)

type Grader interface {
	Grade(path, course, assignment, out string) error
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
