package java

import (
	"cso/codecowboy/classroom"
	util "cso/codecowboy/graders/grader_util"
	"cso/codecowboy/store"
	"encoding/xml"
	"github.com/charmbracelet/log"
	"os"
	"path"
	"strings"
)

type JavaGrader struct {
	db *store.DB
}

func (j JavaGrader) Grade(spec classroom.AssignmentSpec, out string) error {
	return util.Grade(j.db, []string{"./gradlew", "test"}, spec, func(stdOut string) (float64, float64, float64, error) {
		wd, _ := os.Getwd()
		reportPath := path.Join(wd, "build", "test-results", "test")
		log.Debug("reading test output", "reportPath", reportPath)
		return readJavaTestResults(reportPath)
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

func readJavaTestResults(reportPath string) (float64, float64, float64, error) {
	tests, failures := 0.0, 0.0
	files, err := os.ReadDir(reportPath)
	if err != nil {
		return 0, 0, 0, err
	}
	for _, f := range files {
		if !strings.HasSuffix(f.Name(), "xml") {
			continue
		}
		contents, err := os.ReadFile(path.Join(reportPath, f.Name()))
		if err != nil {
			return 0, 0, 0, err
		}
		suite := javaTestsuite{}
		err = xml.Unmarshal(contents, &suite)
		if err != nil {
			return 0, 0, 0, err
		}
		tests += float64(suite.Tests)
		failures += float64(suite.Failures)
		log.Debug("found test file", "tests", tests, "failures", failures)
	}
	return tests - failures, failures, 0.0, nil
}
