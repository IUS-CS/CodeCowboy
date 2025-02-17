package canvasfmt

import (
	"cso/codecowboy/classroom"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse_ValidCSV(t *testing.T) {
	csvData := `
Student,ID,SIS Login ID,Section,Assignment 1 (15735942)
,,,,Manual Posting
    Points Possible,,,,100.0
John Doe,123,jdoe,Section A,
Jane Doe,456,jane.doe,Section B,\n`
	reader := strings.NewReader(csvData)
	currentStudents := []classroom.Student{
		{SISLoginID: "jdoe", GithubID: "gh123", GitHubUsername: "jdoeGit"},
	}

	students, err := Parse(reader, currentStudents)
	assert.NoError(t, err)
	_ = students
	assert.Len(t, students, 2)

	assert.Equal(t, "John Doe", students[0].Name)
	assert.Equal(t, "123", students[0].ID)
	assert.Equal(t, "jdoe", students[0].SISLoginID)
	assert.Equal(t, "Section A", students[0].Section)
	assert.Equal(t, "gh123", students[0].GithubID) // Ensuring Update applied
	assert.Equal(t, "jdoeGit", students[0].GitHubUsername)
}

func TestParse_InvalidCSV(t *testing.T) {
	csvData := `Invalid,\n,CSV\n,,,"Data\n` // This will cause a parsing error
	reader := strings.NewReader(csvData)

	students, err := Parse(reader, nil)
	assert.Nil(t, students)
	assert.Error(t, err)
}

func TestUpdate_ExistingStudent(t *testing.T) {
	ghMap := map[string]classroom.Student{
		"jdoe": {GithubID: "gh123", GitHubUsername: "jdoeGit"},
	}
	student := classroom.Student{Name: "John Doe", SISLoginID: "jdoe"}
	updatedStudent := Update(ghMap, student)

	assert.Equal(t, "gh123", updatedStudent.GithubID)
	assert.Equal(t, "jdoeGit", updatedStudent.GitHubUsername)
}

func TestUpdate_NewStudent(t *testing.T) {
	ghMap := map[string]classroom.Student{}
	student := classroom.Student{Name: "New Student", SISLoginID: "newid"}
	updatedStudent := Update(ghMap, student)

	assert.Empty(t, updatedStudent.GithubID)
	assert.Empty(t, updatedStudent.GitHubUsername)
}
