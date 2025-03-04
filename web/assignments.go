package web

import (
	"bytes"
	"cso/codecowboy/classroom"
	"cso/codecowboy/graders/golang"
	"cso/codecowboy/graders/java"
	"cso/codecowboy/graders/net"
	"cso/codecowboy/store"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/charmbracelet/log"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

const STATUS_RUNNING = "running"

func (w *Web) setupAssignmentHandlers() chi.Router {
	router := chi.NewRouter()

	router.Get("/newAssignment", w.handleNewAssignmentForm)
	router.Get("/{assignment}", w.handleAssignmentDetails)
	router.Delete("/{assignment}", w.handleRmAssignment)
	router.Post("/runAll", w.handleRunAllAssignments)
	router.Post("/{assignment}/run", w.handleRunAssignment)
	router.Get("/{assignment}/download/{id}", w.handleDownloadResult)
	router.Get("/{assignment}/view/{id}", w.handleViewResult)
	router.Get("/{assignment}/status", w.handleExecutionList)
	router.Post("/", w.handleNewAssignment)

	return router
}

func (w *Web) handleAssignmentDetails(wr http.ResponseWriter, r *http.Request) {
	courseName := chi.URLParam(r, "course")
	assignmentName := chi.URLParam(r, "assignment")
	if courseName == "" || assignmentName == "" {
		w.renderErr(r.Context(), wr, fmt.Errorf("handleAssignmentDetails missing information, course: %s, assignment: %s", courseName, assignmentName))
		return
	}
	course, err := classroom.New(w.db, courseName)
	if err != nil {
		w.renderErr(r.Context(), wr, err)
		return
	}
	for _, a := range course.Assignments {
		log.Debugf("handleAssignmentDetails checking %s against %s", a.Name, assignmentName)
		if a.Name == assignmentName {
			if a.Expr == "" {
				a.Expr = classroom.DEFAULT_EXPR
			}
			w.Index(assignmentName, w.assignmentDetails(a, w.runLog[courseName+assignmentName])).Render(r.Context(), wr)
			return
		}
	}
	log.Debugf("handleAssignmentDetails checking %s against %s, course: %+v", assignmentName, courseName, course)
	w.renderErr(r.Context(), wr, fmt.Errorf("handleAssignmentDetails could not find assignment %s for %s", assignmentName, courseName))
}

func (w *Web) handleRmAssignment(wr http.ResponseWriter, r *http.Request) {
	assignmentName := chi.URLParam(r, "assignment")
	courseName := chi.URLParam(r, "course")
	cls, err := classroom.New(w.db, courseName)
	if err != nil {
		w.renderErr(r.Context(), wr, err)
		return
	}
	assignments := classroom.Assignments{}
	for _, a := range cls.Assignments {
		if a.Name != assignmentName {
			assignments = append(assignments, a)
		}
	}
	cls.Assignments = assignments
	err = cls.Save()
	if err != nil {
		w.renderErr(r.Context(), wr, err)
		return
	}
	wr.Header().Set("HX-Redirect", "/courses/"+cls.Name)
}

func (w *Web) handleNewAssignmentForm(wr http.ResponseWriter, r *http.Request) {
	courseName := chi.URLParam(r, "course")
	w.newAssignmentForm(courseName).Render(r.Context(), wr)
}

func (w *Web) handleNewAssignment(wr http.ResponseWriter, r *http.Request) {
	cls, err := classroom.New(w.db, chi.URLParam(r, "course"))
	if err != nil {
		w.renderErr(r.Context(), wr, err)
		return
	}
	assign := classroom.AssignmentSpec{
		Name:      r.FormValue("name"),
		Type:      r.FormValue("type"),
		Path:      r.FormValue("path"),
		Course:    cls.Name,
		ExtrasSrc: r.FormValue("extrasSrc"),
		ExtrasDst: r.FormValue("extrasDst"),
		Expr:      r.FormValue("expr"),
	}

	cls.Assignments = append(cls.Assignments, assign)

	err = cls.Save()
	if err != nil {
		w.renderErr(r.Context(), wr, err)
		return
	}

	wr.Header().Set("HX-Redirect", "/courses/"+cls.Name)
	w.courseDetails(cls).Render(r.Context(), wr)
}

func (w *Web) handleRunAssignment(wr http.ResponseWriter, r *http.Request) {
	course := chi.URLParam(r, "course")
	assignment := chi.URLParam(r, "assignment")
	rawDueDate := r.FormValue("duedate")
	dueDate, err := time.Parse(time.DateTime, rawDueDate)
	if err != nil {
		dueDate = time.Now()
	}
	cls, err := classroom.New(w.db, course)
	if err != nil {
		w.renderErr(r.Context(), wr, err)
		return
	}
	// TODO: UUID is a weird thing to use for this.. executions start jumping around and getting unordered.
	id := uuid.New().String()
	for _, a := range cls.Assignments {
		if a.Name == assignment {
			go func() {
				if _, ok := w.runLog[course+assignment]; !ok {
					w.runLog[course+assignment] = map[string]string{}
				}
				w.runLog[course+assignment][id] = STATUS_RUNNING

				wd, _ := os.Getwd()
				tmpDir, assnPath, err := a.Clone()
				defer a.Cleanup(wd, tmpDir)

				a.Path = assnPath

				out, err := run(w.db, a, dueDate)

				if err != nil {
					w.runLog[course+assignment][id] = err.Error()
				} else {
					w.runLog[course+assignment][id] = out
				}
			}()
			wr.Header().Set("HX-Location", "/courses/"+cls.Name+"/assignments/"+assignment)
			return
		}
	}
	w.renderErr(r.Context(), wr, fmt.Errorf("handleRunAssignment could not find assignment"))
}

// Runs all the assignments for a course sequentially
func (w *Web) handleRunAllAssignments(wr http.ResponseWriter, r *http.Request) {
	course := chi.URLParam(r, "course")
	cls, err := classroom.New(w.db, course)
	if err != nil {
		w.renderErr(r.Context(), wr, err)
		return
	}

	done := make(chan bool)
	idToAssignment := make(map[string]string) // Store assignment names per ID

	go func() {
		for _, a := range cls.Assignments {
			id := uuid.New().String()
			idToAssignment[id] = a.Name // Store assignment name for this run

			rawDueDate := r.FormValue("duedate")
			dueDate, err := time.Parse(time.DateTime, rawDueDate)
			if err != nil {
				dueDate = time.Now()
			}

			if _, ok := w.runLog[course+a.Name]; !ok {
				w.runLog[course+a.Name] = map[string]string{}
			}
			w.runLog[course+a.Name][id] = STATUS_RUNNING

			wd, _ := os.Getwd()
			tmpDir, assnPath, err := a.Clone()
			defer a.Cleanup(wd, tmpDir)

			a.Path = assnPath
			out, err := run(w.db, a, dueDate)

			if err != nil {
				w.runLog[course+a.Name][id] = err.Error()
			} else {
				w.runLog[course+a.Name][id] = out
			}
		}
		done <- true
	}()
	<-done

	wr.Header().Set("Content-Type", "text/html")
	wr.Write([]byte(`<span id="run-all-status">Complete</span>`)) // TODO: link to results file
}

func (w *Web) handleExecutionList(wr http.ResponseWriter, r *http.Request) {
	course := chi.URLParam(r, "course")
	assignment := chi.URLParam(r, "assignment")
	w.listExecutions(course, assignment, w.runLog[course+assignment]).Render(r.Context(), wr)
}

func (w *Web) handleViewResult(wr http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	course := chi.URLParam(r, "course")
	assignment := chi.URLParam(r, "assignment")
	if w.runLog[course+assignment] != nil && w.runLog[course+assignment][id] == "" {
		w.renderErr(r.Context(), wr, fmt.Errorf("could not find command execution"))
		return
	}
	if w.runLog[course+assignment] != nil && w.runLog[course+assignment][id] == STATUS_RUNNING {
		w.renderErr(r.Context(), wr, fmt.Errorf("command still running"))
		return
	}
	w.viewResult(w.runLog[course+assignment][id]).Render(r.Context(), wr)
}

func (w *Web) handleDownloadResult(wr http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	course := chi.URLParam(r, "course")
	assignment := chi.URLParam(r, "assignment")
	if w.runLog[course+assignment] != nil && w.runLog[course+assignment][id] == "" {
		w.renderErr(r.Context(), wr, fmt.Errorf("could not find command execution"))
		return
	}
	if w.runLog[course+assignment] != nil && w.runLog[course+assignment][id] == STATUS_RUNNING {
		w.renderErr(r.Context(), wr, fmt.Errorf("command still running"))
		return
	}
	wr.Header().Set("Content-Disposition",
		fmt.Sprintf("attachment; filename=grade_%s_%s_%s.csv",
			course, assignment, time.Now().Format(time.RFC3339)))
	wr.Header().Set("Content-Type", "text/csv")
	wr.Write([]byte(w.runLog[course+assignment][id]))
}

func run(db *store.DB, a classroom.AssignmentSpec, dueDate time.Time) (string, error) {

	var grader Grader

	switch a.Type {
	case "go":
		grader = golang.NewGoGrader(db)
	case "java":
		grader = java.NewJavaGrader(db)
	case "net":
		grader = net.NewNetGrader(db)
	}
	if grader == nil {
		return "", fmt.Errorf("unknown grader type: %s", a.Type)
	}

	out := bytes.NewBuffer([]byte{})
	err := grader.Grade(a, dueDate, out)
	return out.String(), err
}

type Grader interface {
	Grade(spec classroom.AssignmentSpec, timeLate time.Time, out io.Writer) error
}