package net

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

type NetGrader struct {
	db *store.DB
}

func NewNetGrader(db *store.DB) NetGrader {
	return NetGrader{db}
}

func (n NetGrader) Grade(spec classroom.AssignmentSpec, out string) error {
	studentList, err := classroom.New(n.db, spec.Course)
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

		cmd := exec.Command("dotnet", "test", "--logger", "trx;logfilename=../../results.trx")
		var stdOut strings.Builder
		var stdErr strings.Builder
		cmd.Stdout = &stdOut
		cmd.Stderr = &stdErr
		err = cmd.Run()

		wd, _ := os.Getwd()
		score, err := readNetTestResults(path.Join(wd, "results.trx"))
		if err != nil {
			return err
		}

		who := canvasfmt.SISNameFromDirName(studentList.Students, d.Name())

		log.Debugf("grade for %s: %.2f", who, score*100)
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

func readNetTestResults(reportPath string) (float64, error) {
	contents, err := os.ReadFile(reportPath)
	if err != nil {
		return 0, err
	}
	suite := netTestRun{}
	err = xml.Unmarshal(contents, &suite)
	if err != nil {
		return 0, err
	}

	counters := suite.ResultSummary.Counters
	return float64(counters.Passed) / float64(counters.Total), nil
}

type netTestRun struct {
	XMLName xml.Name `xml:"TestRun"`
	Text    string   `xml:",chardata"`
	ID      string   `xml:"id,attr"`
	Name    string   `xml:"name,attr"`
	Xmlns   string   `xml:"xmlns,attr"`
	Times   struct {
		Text     string `xml:",chardata"`
		Creation string `xml:"creation,attr"`
		Queuing  string `xml:"queuing,attr"`
		Start    string `xml:"start,attr"`
		Finish   string `xml:"finish,attr"`
	} `xml:"Times"`
	TestSettings struct {
		Text       string `xml:",chardata"`
		Name       string `xml:"name,attr"`
		ID         string `xml:"id,attr"`
		Deployment struct {
			Text              string `xml:",chardata"`
			RunDeploymentRoot string `xml:"runDeploymentRoot,attr"`
		} `xml:"Deployment"`
	} `xml:"TestSettings"`
	Results struct {
		Text           string `xml:",chardata"`
		UnitTestResult []struct {
			Text                     string `xml:",chardata"`
			ExecutionId              string `xml:"executionId,attr"`
			TestId                   string `xml:"testId,attr"`
			TestName                 string `xml:"testName,attr"`
			ComputerName             string `xml:"computerName,attr"`
			Duration                 string `xml:"duration,attr"`
			StartTime                string `xml:"startTime,attr"`
			EndTime                  string `xml:"endTime,attr"`
			TestType                 string `xml:"testType,attr"`
			Outcome                  string `xml:"outcome,attr"`
			TestListId               string `xml:"testListId,attr"`
			RelativeResultsDirectory string `xml:"relativeResultsDirectory,attr"`
			Output                   struct {
				Text      string `xml:",chardata"`
				ErrorInfo struct {
					Text       string `xml:",chardata"`
					Message    string `xml:"Message"`
					StackTrace string `xml:"StackTrace"`
				} `xml:"ErrorInfo"`
			} `xml:"Output"`
		} `xml:"UnitTestResult"`
	} `xml:"Results"`
	TestDefinitions struct {
		Text     string `xml:",chardata"`
		UnitTest []struct {
			Text      string `xml:",chardata"`
			Name      string `xml:"name,attr"`
			Storage   string `xml:"storage,attr"`
			ID        string `xml:"id,attr"`
			Execution struct {
				Text string `xml:",chardata"`
				ID   string `xml:"id,attr"`
			} `xml:"Execution"`
			TestMethod struct {
				Text            string `xml:",chardata"`
				CodeBase        string `xml:"codeBase,attr"`
				AdapterTypeName string `xml:"adapterTypeName,attr"`
				ClassName       string `xml:"className,attr"`
				Name            string `xml:"name,attr"`
			} `xml:"TestMethod"`
		} `xml:"UnitTest"`
	} `xml:"TestDefinitions"`
	TestEntries struct {
		Text      string `xml:",chardata"`
		TestEntry []struct {
			Text        string `xml:",chardata"`
			TestId      string `xml:"testId,attr"`
			ExecutionId string `xml:"executionId,attr"`
			TestListId  string `xml:"testListId,attr"`
		} `xml:"TestEntry"`
	} `xml:"TestEntries"`
	TestLists struct {
		Text     string `xml:",chardata"`
		TestList []struct {
			Text string `xml:",chardata"`
			Name string `xml:"name,attr"`
			ID   string `xml:"id,attr"`
		} `xml:"TestList"`
	} `xml:"TestLists"`
	ResultSummary struct {
		Text     string `xml:",chardata"`
		Outcome  string `xml:"outcome,attr"`
		Counters struct {
			Text                int `xml:",chardata"`
			Total               int `xml:"total,attr"`
			Executed            int `xml:"executed,attr"`
			Passed              int `xml:"passed,attr"`
			Failed              int `xml:"failed,attr"`
			Error               int `xml:"error,attr"`
			Timeout             int `xml:"timeout,attr"`
			Aborted             int `xml:"aborted,attr"`
			Inconclusive        int `xml:"inconclusive,attr"`
			PassedButRunAborted int `xml:"passedButRunAborted,attr"`
			NotRunnable         int `xml:"notRunnable,attr"`
			NotExecuted         int `xml:"notExecuted,attr"`
			Disconnected        int `xml:"disconnected,attr"`
			Warning             int `xml:"warning,attr"`
			Completed           int `xml:"completed,attr"`
			InProgress          int `xml:"inProgress,attr"`
			Pending             int `xml:"pending,attr"`
		} `xml:"Counters"`
	} `xml:"ResultSummary"`
}
