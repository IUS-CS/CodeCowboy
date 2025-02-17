package classroom

import (
	"bytes"
	"cso/codecowboy/graders/types"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/log"
	"github.com/expr-lang/expr"
)

const DEFAULT_EXPR = `passed / (passed+failed) * (late == 0 ? 1 : 0)`

var Languages = []string{"go", "java", "net"}

type Assignments []AssignmentSpec

type AssignmentSpec struct {
	Name      string
	Type      string
	Path      string
	Course    string
	ExtrasSrc string
	ExtrasDst string
	Expr      string
}

func ParseAssignmentsFile(path, courseName string) (Assignments, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("ParseAssignmentsOpen: %w", err)
	}
	return ParseAssignments(f, courseName)
}

func ParseAssignments(r io.Reader, courseName string) (Assignments, error) {
	assignments := Assignments{}
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("ParseAssignments ReadAll: %w", err)
	}
	err = json.Unmarshal(data, &assignments)
	if err != nil {
		return nil, fmt.Errorf("ParseAssignments Unmarshal: %w", err)
	}
	for i := range assignments {
		assignments[i].Course = courseName
	}
	return assignments, nil
}

func (a AssignmentSpec) Score(testResult types.GraderReturn) (float64, error) {
	env := map[string]any{
		"passed": testResult.Passed,
		"failed": testResult.Failed,
		"cover":  testResult.Coverage,
		"late":   int(testResult.TimeLate),
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
	log.Debug("score finished", "env", env, "result", result)
	return result.(float64), nil
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
	if a.Name == "" {
		errs = append(errs, fmt.Errorf("missing assignment name"))
	}
	return errors.Join(errs...)
}

// TODO: Relying on CLI output is a really bad idea. Switch to API for reliable querying...
var cloner = `gh classroom clone student-repos -d "%s" -a $(gh classroom assignments -c $(gh classroom list | grep "%s" | cut -w -f1)|tail -n +4 | grep "%s" | cut -w -f1)`
var assignmentName = `gh classroom assignments -c $(gh classroom list | tail -n +4 | grep "%s" | cut -w -f1)|tail -n +4 | grep "%s" | cut -w -f2`

func stripDanger(input string) string {
	strips := []string{";", "&", "!"}
	for _, s := range strips {
		input = strings.ReplaceAll(input, s, "")
	}
	return input
}

var emptyTime = time.Duration(0)

func (a AssignmentSpec) CheckSubmissionDate(path string, dueDate time.Time) (time.Duration, error) {
	wd, err := os.Getwd()
	if err != nil {
		return emptyTime, err
	}
	defer os.Chdir(wd)
	err = os.Chdir(path)
	if err != nil {
		return emptyTime, err
	}
	cmd := exec.Command("git", "log", "-1", "--format=\"%at\"")
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	err = cmd.Run()
	if err != nil {
		return emptyTime, err
	}
	if stderr.Len() == 0 {
		return emptyTime, fmt.Errorf("error getting commit time: %s", stderr.String())
	}
	tstamp, err := strconv.ParseInt(stdout.String(), 10, 64)
	if err != nil {
		return emptyTime, err
	}
	commitTime := time.UnixMicro(tstamp)
	return commitTime.Sub(dueDate), nil
}

func errReturn(err error) (string, string, error) {
	return "", "", err
}

func (a AssignmentSpec) Clone() (string, string, error) { // TODO: This may be causing errors, look into this
	if err := a.Validate(); err != nil {
		return errReturn(err)
	}
	tmpPath, err := os.MkdirTemp("", "*-repos")
	if err != nil {
		return errReturn(err)
	}
	err = os.Chdir(tmpPath)
	if err != nil {
		return errReturn(err)
	}

	log.Debugf("Created tmp dir: %s", tmpPath)
	fmtCmd := fmt.Sprintf(cloner, tmpPath, stripDanger(a.Course), stripDanger(a.Name))
	log.Debugf("Running command: %s", fmtCmd)
	cmd := exec.Command("/bin/sh", "-c", fmtCmd)

	// TODO: if this errors out, the user is never notified.
	// may error if gh classroom extension isn't installed
	if err := cmd.Run(); err != nil {
		return errReturn(err)
	}
	fmtCmd = fmt.Sprintf(assignmentName, stripDanger(a.Course), stripDanger(a.Name))
	cmd = exec.Command("/bin/sh", "-c", fmtCmd)
	log.Debugf("Running command: %s", cmd.String())
	stdOut := strings.Builder{}
	cmd.Stdout = &stdOut
	if err := cmd.Run(); err != nil {
		return errReturn(err)
	}
	log.Debugf("result: %s", stdOut.String())
	ghAssignmentName := strings.ToLower(strings.TrimSpace(stdOut.String()))
	dir, err := os.ReadDir(".")
	if err != nil {
		return errReturn(err)
	}
	var assnPath string
	for _, d := range dir {
		if strings.HasPrefix(d.Name(), ghAssignmentName) && strings.HasSuffix(d.Name(), "-submissions") {
			log.Debugf("Inspecting %s similar to %s", d.Name(), ghAssignmentName)
			assnPath = filepath.Join(tmpPath, d.Name())
			break
		}
	}
	if assnPath == "" {
		return errReturn(fmt.Errorf("assignment path not found: %s/%s-submissions", tmpPath, ghAssignmentName))
	}
	return tmpPath, assnPath, nil
}

func (a AssignmentSpec) Cleanup(wd, tmpPath string) error {
	chBack := func() {
		if err := os.Chdir(wd); err != nil {
			panic(err)
		}
	}
	defer chBack()
	return os.RemoveAll(tmpPath)
}
