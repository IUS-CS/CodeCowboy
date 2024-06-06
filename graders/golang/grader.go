package golang

import (
	"cso/codecowboy/classroom"
	util "cso/codecowboy/graders/grader_util"
	"cso/codecowboy/store"
	"encoding/json"
	"github.com/charmbracelet/log"
	"io"
	"strconv"
	"strings"
	"time"
)

type GoGrader struct {
	db *store.DB
}

func NewGoGrader(db *store.DB) GoGrader {
	return GoGrader{db}
}

func (g GoGrader) Grade(spec classroom.AssignmentSpec, dueDate time.Time, out io.Writer) error {
	return util.Grade(g.db, []string{"go", "test", "-cover", "-json"}, spec, dueDate, g.readGoResults, out)
}

func (g GoGrader) readGoResults(testOutput string, dueDate time.Time) (float64, float64, float64, time.Duration, error) {
	outputs := g.fromJSONLines(testOutput)
	passes, _ := g.getKind(outputs, KindPASS)
	fails, _ := g.getKind(outputs, KindFAIL)
	cover := g.getCoverage(outputs)
	// todo need to calculate late time here?
	return passes, fails, cover, time.Duration(0), nil
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

func (g GoGrader) getCoverage(out []goTestOutput) float64 {
	for _, o := range out {
		if o.Action == "output" && strings.Contains(o.Output, "coverage:") {
			out, _ := strconv.ParseFloat(o.Output, 64)
			return out
		}
	}
	return 0.0
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
