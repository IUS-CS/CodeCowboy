package githubfmt

import (
	"cso/codecowboy/students"
	"encoding/csv"
	"os"
)

const (
	Identifier = iota
	GithubUsername
	GithubID
	Name
)

func Parse(path string) ([]students.Student, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	reader := csv.NewReader(f)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}
	students := make([]students.Student, len(records))
	for i, r := range records {
		students[i].Identifier = r[Identifier]
		students[i].GitHubUsername = r[GithubUsername]
		students[i].GithubID = r[GithubID]
		students[i].Name = r[Name]
	}
	return students, nil
}
