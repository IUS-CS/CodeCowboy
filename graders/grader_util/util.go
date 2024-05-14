package util

import (
	"cso/codecowboy/canvasfmt"
	"cso/codecowboy/classroom"
	"cso/codecowboy/graders"
	"cso/codecowboy/store"
	"fmt"
	"github.com/charmbracelet/log"
	cp "github.com/otiai10/copy"
	"io"
	"os"
	"os/exec"
	"path"
	"strings"
)

func CopyExtras(from string, to string) error {
	d, err := os.ReadDir(from)
	if err != nil {
		return err
	}
	for _, entry := range d {
		err = cp.Copy(path.Join(from, entry.Name()), to)
		if err != nil {
			return err
		}
	}
	return nil
}

func Grade(db *store.DB, command []string, spec classroom.AssignmentSpec, testResult graders.TestResult, out io.Writer) error {
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
			err := CopyExtras(spec.ExtrasSrc, path.Join(getwd, spec.ExtrasDst))
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

		passes, fails, cover, err := testResult(stdOut.String())
		if err != nil {
			return err
		}

		score, err := spec.Score(passes, fails, cover)
		if err != nil {
			return err
		}

		gradeStr := fmt.Sprintf("%2.1f%%", score)
		who := canvasfmt.SISNameFromDirName(studentList.Students, d.Name())

		grades[who] = score

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

	return canvasfmt.WriteCSV(out, spec.Name, studentList.Students, grades)
}
