package students

import (
	"cso/codecowboy/store"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/dgraph-io/badger/v3"
)

type Student struct {
	Name           string
	ID             string
	SISLoginID     string
	Section        string
	GitHubUsername string
	GithubID       string
}

type Students struct {
	db      *store.DB
	Course  string
	Members []Student
}

func New(db *store.DB, course string) *Students {
	s := &Students{
		db:      db,
		Members: []Student{},
		Course:  course,
	}
	err := s.Populate()
	if err != nil {
		log.Fatal(err)
	}
	return s
}

func NewFromList(db *store.DB, course string, list []Student) *Students {
	s := &Students{
		db:      db,
		Course:  course,
		Members: list,
	}
	return s
}

func (s *Students) Populate() error {
	err := s.db.Unmarshal(s.Course, &s.Members)
	if !errors.Is(err, badger.ErrKeyNotFound) && err != nil {
		return err
	}
	return nil
}

func (s *Students) Save() error {
	return s.db.Set(s.Course, &s.Members)
}

func (s Student) String() string {
	return fmt.Sprintf("%s:\t%s", s.GitHubUsername, s.ID)
}

func (s *Students) String() string {
	out := "Student list:\n\nGitHub:\tLMS\n"
	for _, s := range s.Members {
		out += s.String() + "\n"
	}
	if len(s.Members) == 0 {
		return "No students exist."
	}
	return out
}

func (s *Students) ToJSON() (string, error) {
	out, err := json.MarshalIndent(s, "", "\t")
	if err != nil {
		return "", err
	}
	return string(out), nil
}
