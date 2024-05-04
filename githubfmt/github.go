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

func Parse(path string, current []students.Student) ([]students.Student, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	reader := csv.NewReader(f)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}
	canvasMap := make(map[string]students.Student)
	for _, s := range current {
		if s.SISLoginID != "" {
			canvasMap[s.SISLoginID] = s
		}
	}
	out := []students.Student{}
	for _, r := range records {
		student := students.Student{
			SISLoginID:     r[Identifier],
			GitHubUsername: r[GithubUsername],
			GithubID:       r[GithubID],
			Name:           r[Name],
		}
		student = Update(canvasMap, student)
		out = append(out, student)
	}
	return out, nil
}

func Update(canvasMap map[string]students.Student, student students.Student) students.Student {
	if cs, ok := canvasMap[student.SISLoginID]; ok {
		student.Name = cs.Name
		student.ID = cs.ID
		student.SISLoginID = cs.SISLoginID
		student.Section = cs.Section
	}
	return student
}
