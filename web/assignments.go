package web

import (
	"cso/codecowboy/classroom"
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
	"time"
)

func (w *Web) setupAssignmentHandlers() chi.Router {
	router := chi.NewRouter()

	router.Get("/newAssignment", w.handleNewAssignmentForm)
	router.Get("/{assignment}", w.handleAssignmentDetails)
	router.Delete("/{assignment}", w.handleRmAssignment)
	router.Post("/{assignment}/run", w.handleRunAssignment)
	router.Post("/", w.handleNewAssignment)

	return router
}

func (w *Web) handleAssignmentDetails(wr http.ResponseWriter, r *http.Request) {
	courseName := chi.URLParam(r, "course")
	assignmentName := chi.URLParam(r, "assignment")
	course, err := classroom.New(w.db, courseName)
	if err != nil {
		w.renderErr(r.Context(), wr, err)
	}
	for _, a := range course.Assignments {
		if a.Name == assignmentName {
			if a.Expr == "" {
				a.Expr = classroom.DEFAULT_EXPR
			}
			w.Index(assignmentName, w.assignmentDetails(a)).Render(r.Context(), wr)
			return
		}
	}
	w.renderErr(r.Context(), wr, fmt.Errorf("could not find assignment"))
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
		GitHubID:  r.FormValue("GitHubID"),
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
	cls, err := classroom.New(w.db, course)
	if err != nil {
		w.renderErr(r.Context(), wr, err)
		return
	}
	for _, a := range cls.Assignments {
		if a.Name == assignment {
			out, err := a.CloneAndRun(func() (string, error) {
				return "", nil //run(w.db, a)
			})
			if err != nil {
				w.renderErr(r.Context(), wr, err)
			}
			wr.Header().Set("Content-Disposition",
				fmt.Sprintf("attachment; filename=grade_%s_%s_%s.json",
					a.Course, a.Name, time.Now().Format(time.RFC3339)))
			wr.Header().Set("Content-Type", "text/csv")
			wr.Write([]byte(out))
		}
	}
	w.renderErr(r.Context(), wr, fmt.Errorf("could not find assignment"))
}

// func run(db *store.DB, a classroom.AssignmentSpec) (string, error) {
// 	grader := graders.GetGrader(a.Type, db)
// 	if grader == nil {
// 		return "", fmt.Errorf("unknown grader type: %s", a.Type)
// 	}
//
// 	out := bytes.NewBuffer([]byte{})
// 	err := grader.Grade(a, out)
// 	return out.String(), err
// }
