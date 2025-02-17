package githubfmt

import (
	"cso/codecowboy/classroom"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	csvData := `12345,johndoe,67890,John Doe
67890,janedoe,54321,Jane Doe`
	r := strings.NewReader(csvData)

	current := classroom.Students{
		{SISLoginID: "12345", Name: "Johnathan Doe", ID: "1", Section: "A"},
		{SISLoginID: "67890", Name: "Jane Doe", ID: "2", Section: "B"},
	}

	parsed, err := Parse(r, current)
	assert.NoError(t, err)

	expected := classroom.Students{
		{SISLoginID: "12345", GitHubUsername: "johndoe", GithubID: "67890", Name: "Johnathan Doe", ID: "1", Section: "A"},
		{SISLoginID: "67890", GitHubUsername: "janedoe", GithubID: "54321", Name: "Jane Doe", ID: "2", Section: "B"},
	}

	assert.Equal(t, expected, parsed)
}

func TestParse_EmptyCSV(t *testing.T) {
	r := strings.NewReader("")
	parsed, err := Parse(r, nil)
	assert.NoError(t, err)
	assert.Empty(t, parsed)
}

func TestParse_NoCurrentStudents(t *testing.T) {
	csvData := `11111,alice,22222,Alice Doe`
	r := strings.NewReader(csvData)
	parsed, err := Parse(r, nil)
	assert.NoError(t, err)

	expected := classroom.Students{
		{SISLoginID: "11111", GitHubUsername: "alice", GithubID: "22222", Name: "Alice Doe"},
	}

	assert.Equal(t, expected, parsed)
}
