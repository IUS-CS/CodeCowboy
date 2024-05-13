package web

import (
	"fmt"
	"net/http"

	"cso/codecowboy/classroom"

	"github.com/go-chi/chi/v5"
)

func (w *Web) setupCourseHandlers() chi.Router {
	router := chi.NewRouter()
	router.Get("/", w.handleCourseList)
	router.Get("/new", w.handleNewCourse)
	router.Get("/{course}", w.handleCourseDetails)
	router.Get("/{course}/assignments/{assignment}", w.handleAssignmentDetails)
	router.Post("/{course}/assignments", w.handleNewAssignment)
	return router
}

func (w *Web) handleCourseList(wr http.ResponseWriter, r *http.Request) {
	courses, err := classroom.All(w.db)
	if err != nil {
		w.renderErr(r.Context(), wr, err)
		return
	}
	w.Index("Courses", w.courseList(courses)).Render(r.Context(), wr)
}

func (w *Web) handleNewCourse(wr http.ResponseWriter, r *http.Request) {
	courses, err := classroom.All(w.db)
	if err != nil {
		w.renderErr(r.Context(), wr, err)
		return
	}
	w.Index("New Course", w.courseList(courses)).Render(r.Context(), wr)
}

func (w *Web) handleCourseDetails(wr http.ResponseWriter, r *http.Request) {
	courseName := chi.URLParam(r, "course")
	course, err := classroom.New(w.db, courseName)
	if err != nil {
		w.renderErr(r.Context(), wr, err)
		return
	}
	w.Index(courseName, w.courseDetails(course)).Render(r.Context(), wr)
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

func (w *Web) handleNewAssignment(wr http.ResponseWriter, r *http.Request) {
	cls, err := classroom.New(w.db, r.FormValue("course"))
	if err != nil {
		w.renderErr(r.Context(), wr, err)
		return
	}
	assign := classroom.AssignmentSpec{
		Name:      r.FormValue("name"),
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
