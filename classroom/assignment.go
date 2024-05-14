package classroom

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/expr-lang/expr"
	"io"
	"os"
	"os/exec"
	"strings"
)

const DEFAULT_EXPR = `passed / (passed+failed)`

var Languages = []string{"go", "java", "net"}

type Assignments []AssignmentSpec

func ParseAssignmentsFile(path, courseName string) (Assignments, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	return ParseAssignments(f, courseName)
}

func ParseAssignments(r io.Reader, courseName string) (Assignments, error) {
	assignments := Assignments{}
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, &assignments)
	if err != nil {
		return nil, err
	}
	for i := range assignments {
		assignments[i].Course = courseName
	}
	return assignments, nil
}

func (a AssignmentSpec) Score(passed, failed, cover float64) (float64, error) {
	env := map[string]any{
		"passed": passed,
		"failed": failed,
		"cover":  cover,
	}
	if a.Expr == "" {
		a.Expr = DEFAULT_EXPR
	}
	pgm, err := expr.Compile(a.Expr, expr.Env(env))
	if err != nil {
		return 0.0, err
	}
	result, err := expr.Run(pgm, env)
	if err != nil {
		return 0.0, err
	}
	return result.(float64), nil
}

type AssignmentSpec struct {
	Name      string
	GitHubID  string
	Type      string
	Path      string
	Course    string
	ExtrasSrc string
	ExtrasDst string
	Expr      string
}

func in(input string, list []string) bool {
	for _, e := range list {
		if input == e {
			return true
		}
	}
	return false
}

func (a AssignmentSpec) Validate() error {
	errs := make([]error, 0)
	if !in(a.Type, Languages) {
		errs = append(errs, fmt.Errorf("language type %s not found", a.Type))
	}
	if a.GitHubID == "" {
		errs = append(errs, fmt.Errorf("missing GitHub Assignment ID"))
	}
	if a.Name == "" {
		errs = append(errs, fmt.Errorf("missing assignment name"))
	}
	return errors.Join(errs...)
}

var cloner = `gh classroom clone student-repos -d "%s" -a $(gh classroom assignments -c $(gh classroom list | tail -n +4 | grep "%s" | cut -w -f1)|tail -n +4 | grep "%s" | cut -w -f1)`

func stripDanger(input string) string {
	strips := []string{";", "&", "!"}
	for _, s := range strips {
		input = strings.ReplaceAll(input, s, "")
	}
	return input
}

func (a AssignmentSpec) CloneAndRun(runner func() (string, error)) (string, error) {
	if err := a.Validate(); err != nil {
		return "", err
	}
	path, err := os.MkdirTemp("", "*-repos")
	if err != nil {
		return "", err
	}
	defer os.RemoveAll(path)
	log.Debugf("Created tmp dir: %s", path)
	clone := fmt.Sprintf(cloner, path, stripDanger(a.Course), stripDanger(a.Name))
	log.Debugf("Running command: %s", clone)
	cmd := exec.Command("sh -c", clone)
	if err = cmd.Run(); err != nil {
		return "", err
	}
	output, err := runner()
	if err != nil {
		return "", err
	}
	return output, nil
}
