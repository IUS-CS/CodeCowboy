package classroom

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseAssignments_ValidJSON(t *testing.T) {
	jsonData := `[
		{"Name": "HW1", "Type": "Homework", "Path": "./hw1"},
		{"Name": "HW2", "Type": "Homework", "Path": "./hw2"}
	]`
	r := strings.NewReader(jsonData)
	courseName := "CS101"
	assignments, err := ParseAssignments(r, courseName)
	assert.Nil(t, err)
	assert.Len(t, assignments, 2)
	for _, a := range assignments {
		assert.Equal(t, courseName, a.Course)
	}
}

func TestParseAssignments_InvalidJSON(t *testing.T) {
	invalidJSON := `[{"Name": "HW1", "Type": "Homework", "Path": "./hw1"`
	r := strings.NewReader(invalidJSON)
	_, err := ParseAssignments(r, "CS101")
	assert.NotNil(t, err)
}

func TestParseAssignments_EmptyInput(t *testing.T) {
	r := strings.NewReader("")
	_, err := ParseAssignments(r, "CS101")
	assert.NotNil(t, err)
}
