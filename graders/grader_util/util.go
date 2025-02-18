package util

import (
	"cso/codecowboy/canvasfmt"
	"cso/codecowboy/classroom"
	"cso/codecowboy/graders/types"
	"cso/codecowboy/store"
	"github.com/charmbracelet/log"
	cp "github.com/otiai10/copy"
	"io"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"
)

func CopyExtras(from string, to string) error {
	d, err := os.ReadDir(from)
	if err != nil {
		return err
	}
	for _, entry := range d {
		err = cp.Copy(path.Join(from, entry.Name()), path.Join(to, entry.Name()))
		if err != nil {
			return err
		}
	}
	return nil
}

type TestResultFunc func(stdOut string, dueDate time.Duration) (types.GraderReturn, error)

func Grade(db *store.DB, command []string, spec classroom.AssignmentSpec, dueDate time.Time, testResult func(stdOut string, timeLate time.Duration) (types.GraderReturn, error), out io.Writer) error {
	studentList, err := classroom.New(db, spec.Course)
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

		if spec.ExtrasSrc != "" {
			err := CopyExtras(spec.ExtrasSrc, path.Join(getwd, d.Name(), spec.ExtrasDst))
			if err != nil {
				return err
			}
		}

		cmd := exec.Command(command[0], command[1:]...)
		var stdOut strings.Builder
		var stdErr strings.Builder
		cmd.Stdout = &stdOut
		cmd.Stderr = &stdErr
		err = cmd.Run()
		if err != nil && stdErr.Len() > 0 {
			log.Error("error executing", "stdout", stdOut.String(), "stderr", stdErr.String())
		}

		//timeLate, err := spec.CheckSubmissionDate(getwd, dueDate)
		//if err != nil {
		//	return err
		//}

		timeLate := time.Duration(0)

		result, err := testResult(stdOut.String(), timeLate)
		if err != nil {
			return err
		}

		score, err := spec.Score(result)
		if err != nil {
			return err
		}

		who := canvasfmt.SISNameFromDirName(studentList.Students, d.Name())

		grades[who] = score

		log.Info("Finished grading", "user", who, "result", result)

		err = os.Chdir(getwd)
		if err != nil {
			return err
		}
	}

	err = os.Chdir(getwd)
	if err != nil {
		return err
	}

	return canvasfmt.WriteCSV(out, spec.Name, studentList.Students, grades)
}
