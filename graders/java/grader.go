package java

import (
	"cso/codecowboy/canvasfmt"
	"cso/codecowboy/classroom"
	"cso/codecowboy/store"
	"encoding/xml"
	"github.com/charmbracelet/log"
	"os"
	"os/exec"
	"path"
	"strings"
)

type JavaGrader struct {
	db *store.DB
}

func (j JavaGrader) Grade(spec classroom.AssignmentSpec, out string) error {
	studentList, err := classroom.New(j.db, spec.Course)
	if err != nil {
		return err
	}

	getwd, err := os.Getwd()
	if err != nil {
		return err
	}

	err = os.Chdir(spec.Path)

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

		cmd := exec.Command("./gradlew", "test")
		var stdOut strings.Builder
		var stdErr strings.Builder
		cmd.Stdout = &stdOut
		cmd.Stderr = &stdErr
		err = cmd.Run()

		wd, _ := os.Getwd()
		reportPath := path.Join(wd, "build", "test-results", "test")
		log.Debug("reading test output", "reportPath", reportPath)
		passes, fails, cover, err := readJavaTestResults(reportPath)
		if err != nil {
			return err
		}

		score, err := spec.Score(passes, fails, cover)
		if err != nil {
			return err
		}

		who := canvasfmt.SISNameFromDirName(studentList.Students, d.Name())

		log.Debugf("grade for %s: %.2f", who, score)
		grades[who] = score * 100

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
	return canvasfmt.WriteCSV(w, spec.Name, studentList.Students, grades)
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
