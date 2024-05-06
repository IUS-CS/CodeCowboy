package canvasfmt

import (
	"cso/codecowboy/classroom"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
)

const (
	Student = iota
	ID
	SISLoginID
	Section
)

func Parse(path string, current []classroom.Student) ([]classroom.Student, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	reader := csv.NewReader(f)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}
	ghMap := make(map[string]classroom.Student)
	for _, s := range current {
		if s.SISLoginID != "" {
			ghMap[s.SISLoginID] = s
		}
	}
	out := make([]classroom.Student, 0)
	for i, r := range records {
		if i < 3 {
			continue
		}
		student := classroom.Student{
			Name:       r[Student],
			ID:         r[ID],
			SISLoginID: r[SISLoginID],
			Section:    r[Section],
		}
		student = Update(ghMap, student)
		out = append(out, student)
	}
	return out, nil
}

func Update(ghMap map[string]classroom.Student, student classroom.Student) classroom.Student {
	if ghs, ok := ghMap[student.SISLoginID]; ok {
		student.GithubID = ghs.GithubID
		student.GitHubUsername = ghs.GitHubUsername
	}
	return student
}

func WriteCSV(out io.Writer, assignment string, students classroom.Students, studentGrades map[string]float64) error {
	w := csv.NewWriter(out)
	err := w.Write([]string{"Student", "ID", "SIS Login ID", "Section", assignment})
	if err != nil {
		return err
	}
	for _, s := range students {
		grade := fmt.Sprintf("%2.1f", studentGrades[s.SISLoginID])
		err = w.Write([]string{s.Name, s.ID, s.SISLoginID, s.Section, grade})
		if err != nil {
			return err
		}
	}
	w.Flush()
	return nil
}

func SISNameFromDirName(students classroom.Students, dirName string) string {
	fields := strings.Split(dirName, "-")
	sName := fields[len(fields)-1]
	for _, s := range students {
		if s.GitHubUsername == sName {
			return s.SISLoginID
		}
	}
	return "Unknown"
}
