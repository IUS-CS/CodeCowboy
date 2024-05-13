package classroom

import (
	"encoding/json"
	"errors"
	"fmt"
)

type Student struct {
	Name           string
	ID             string
	SISLoginID     string
	Section        string
	GitHubUsername string
	GithubID       string
}

type Students []Student

func (s Student) String() string {
	return fmt.Sprintf("%s:\t%s", s.GitHubUsername, s.ID)
}

func (s Students) String() string {
	out := "Student list:\n\nGitHub:\tLMS\n"
	for _, s := range s {
		out += s.String() + "\n"
	}
	if len(s) == 0 {
		return "No students exist."
	}
	return out
}

func (s Students) ToJSON() (string, error) {
	out, err := json.MarshalIndent(s, "", "\t")
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func (s Student) Validate() error {
	errs := []error{}
	if s.Name == "" {
		errs = append(errs, errors.New("student name is required"))
	}
	if s.SISLoginID == "" {
		errs = append(errs, errors.New("student SISLoginID is required"))
	}
	if s.GitHubUsername == "" {
		errs = append(errs, errors.New("student GitHubUsername is required"))
	}
	return errors.Join(errs...)
}
