package java

import (
	"cso/codecowboy/classroom"
	util "cso/codecowboy/graders/grader_util"
	"cso/codecowboy/graders/types"
	"cso/codecowboy/store"
	"encoding/xml"
	"github.com/charmbracelet/log"
	"io"
	"os"
	"path"
	"strings"
	"time"
)

type JavaGrader struct {
	db *store.DB
}

func (j JavaGrader) Grade(spec classroom.AssignmentSpec, dueDate time.Time, out io.Writer) error {
	return util.Grade(j.db, []string{"./gradlew", "test"}, spec, dueDate, func(stdOut string, timeLate time.Duration) (types.GraderReturn, error) {
		wd, _ := os.Getwd()
		reportPath := path.Join(wd, "build", "test-results", "test")
		log.Debug("reading test output", "reportPath", reportPath)
		return readJavaTestResults(reportPath, timeLate)
	}, out)
}

func NewJavaGrader(db *store.DB) JavaGrader {
	return JavaGrader{db}
}

type javaTestsuite struct {
	XMLName    xml.Name `xml:"testsuite"`
	Text       string   `xml:",chardata"`
	Name       string   `xml:"name,attr"`
	Tests      int      `xml:"tests,attr"`
	Skipped    int      `xml:"skipped,attr"`
	Failures   int      `xml:"failures,attr"`
	Errors     int      `xml:"errors,attr"`
	Timestamp  string   `xml:"timestamp,attr"`
	Hostname   string   `xml:"hostname,attr"`
	Time       string   `xml:"time,attr"`
	Properties string   `xml:"properties"`
	Testcase   struct {
		Text      string `xml:",chardata"`
		Name      string `xml:"name,attr"`
		Classname string `xml:"classname,attr"`
		Time      string `xml:"time,attr"`
		Failure   struct {
			Text    string `xml:",chardata"`
			Message string `xml:"message,attr"`
			Type    string `xml:"type,attr"`
		} `xml:"failure"`
	} `xml:"testcase"`
	SystemOut string `xml:"system-out"`
	SystemErr string `xml:"system-err"`
}

func readJavaTestResults(reportPath string, timeLate time.Duration) (types.GraderReturn, error) {
	tests, failures := 0.0, 0.0
	files, err := os.ReadDir(reportPath)
	if err != nil {
		return types.GraderReturn{}, err
	}
	for _, f := range files {
		if !strings.HasSuffix(f.Name(), "xml") {
			continue
		}
		contents, err := os.ReadFile(path.Join(reportPath, f.Name()))
		if err != nil {
			return types.GraderReturn{}, err
		}
		suite := javaTestsuite{}
		err = xml.Unmarshal(contents, &suite)
		if err != nil {
			return types.GraderReturn{}, err
		}
		tests += float64(suite.Tests)
		failures += float64(suite.Failures)
		log.Debug("found test file", "tests", tests, "failures", failures)
	}
	return types.GraderReturn{
		Passed:   tests - failures,
		Failed:   failures,
		Coverage: 0,
		TimeLate: timeLate,
	}, nil
}
