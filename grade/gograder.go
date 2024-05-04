package grade

import (
	"cso/codecowboy/canvasfmt"
	"cso/codecowboy/store"
	"cso/codecowboy/students"
	"encoding/json"
	"fmt"
	"github.com/charmbracelet/log"
	"os"
	"os/exec"
	"strings"
	"time"
)

type GoGrader struct {
	db *store.DB
}

func NewGoGrader(db *store.DB) Grader {
	return GoGrader{db}
}

func (g GoGrader) Grade(path, course, assignment, out string) error {
	studentList := students.New(g.db, course)

	getwd, err := os.Getwd()
	if err != nil {
		return err
	}

	err = os.Chdir(path)

	list, err := os.ReadDir(".")
	if err != nil {
		return err
	}

	grades := map[string]float64{}

	for _, d := range list {
		getwd, err := os.Getwd()
		if err != nil {
			return err
		}
		err = os.Chdir(d.Name())
		if err != nil {
			return err
		}

		cmd := exec.Command("go", "test", "-cover", "-json")
		var stdOut strings.Builder
		var stdErr strings.Builder
		cmd.Stdout = &stdOut
		cmd.Stderr = &stdErr
		err = cmd.Run()
		if err != nil && stdErr.Len() > 0 {
			log.Error("error executing", "stdout", stdOut.String(), "stderr", stdErr.String())
		}

		outputs := g.fromJSONLines(stdOut.String())
		passes, _ := g.getKind(outputs, KindPASS)
		fails, _ := g.getKind(outputs, KindFAIL)
		cover := g.getCoverage(outputs)

		grade := passes / (passes + fails) * 100
		gradeStr := fmt.Sprintf("%2.1f%%", grade)
		who := sisNameFromDirName(studentList, d.Name())

		grades[who] = grade

		log.Info("Finished grading", "user", who, "passes", passes, "fails", fails, "cover", cover, "grade", gradeStr)

		err = os.Chdir(getwd)
		if err != nil {
			return err
		}
	}

	err = os.Chdir(getwd)
	if err != nil {
		return err
	}

	w := os.Stdout
	if out != "stdout" {
		w, err = os.OpenFile(out, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		defer w.Close()
	}
	return canvasfmt.WriteCSV(w, assignment, studentList, grades)
}

type goTestOutput struct {
	Time    time.Time `json:"Time"`
	Action  string    `json:"Action"`
	Package string    `json:"Package"`
	Output  string    `json:"Output"`
}

func (g GoGrader) fromJSONLines(input string) []goTestOutput {
	var out []goTestOutput
	for _, line := range strings.Split(input, "\n") {
		if line == "" {
			continue
		}
		var o goTestOutput
		err := json.Unmarshal([]byte(line), &o)
		if err != nil {
			log.Debug(line)
		}
		if err != nil {
			log.Error(err)
		}
		out = append(out, o)
	}
	return out
}

func (g GoGrader) getCoverage(out []goTestOutput) string {
	for _, o := range out {
		if o.Action == "output" && strings.Contains(o.Output, "coverage:") {
			return o.Output
		}
	}
	return "coverage: none"
}

const (
	KindPASS = "PASS:"
	KindFAIL = "FAIL:"
)

func (g GoGrader) getKind(out []goTestOutput, kind string) (float64, []string) {
	n := 0.0
	output := []string{}
	for _, o := range out {
		if o.Action == "output" && strings.Contains(o.Output, kind) {
			n++
			output = append(output, o.Output)
		}
	}
	return n, output
}
