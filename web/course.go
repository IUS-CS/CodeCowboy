package web

import (
	"net/http"

	"cso/codecowboy/classroom"

	"github.com/go-chi/chi/v5"
)

func (w *Web) setupCourseHandlers() chi.Router {
	router := chi.NewRouter()
	router.Get("/", w.handleCourseList)
	router.Get("/new", w.handleNewCourse)
	router.Get("/{course}", w.handleCourseDetails)
	router.Delete("/{course}", w.handleRmCourse)

	router.Mount("/{course}/students", w.setupStudentHandlers())
	router.Mount("/{course}/assignments", w.setupAssignmentHandlers())

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

func (w *Web) handleRmCourse(wr http.ResponseWriter, r *http.Request) {
	courseName := chi.URLParam(r, "course")
	err := w.db.Delete(courseName)
	if err != nil {
		w.renderErr(r.Context(), wr, err)
	}
	wr.Header().Set("HX-Redirect", "/courses/")
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
