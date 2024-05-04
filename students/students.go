package students

import (
	"cso/codecowboy/store"
	"fmt"
)

type Student struct {
	GitHubID string
	LMSID    string
}

type Students struct {
	db      *store.DB
	Members []Student
}

func New(db *store.DB) *Students {
	return &Students{
		db:      db,
		Members: []Student{},
	}
}

func (s Student) String() string {
	return fmt.Sprintf("%s: %s", s.GitHubID, s.LMSID)
}

func (s Students) String() string {
	out := "Student list:\n\nGitHub\tLMS\n"
	for _, s := range s.Members {
		out += s.String()
	}
	if len(s.Members) == 0 {
		return "No students exist."
	}
	return out
}
