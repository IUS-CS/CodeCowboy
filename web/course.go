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
	router.Delete("/{course}/students/{sisID}", w.handleRmStudent)
	router.Post("/{course}/students", w.handleNewStudent)
	router.Get("/{course}/students", w.handleStudentForm)
	router.Get("/{course}/assignments/{assignment}", w.handleAssignmentDetails)
	router.Delete("/{course}/assignments/{assignment}", w.handleRmAssignment)
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

func (w *Web) handleRmStudent(wr http.ResponseWriter, r *http.Request) {
	identifier := chi.URLParam(r, "sisID")
	course := chi.URLParam(r, "course")
	cls, err := classroom.New(w.db, course)
	if err != nil {
		w.renderErr(r.Context(), wr, err)
		return
	}
	roster := classroom.Students{}
	for _, s := range cls.Students {
		if s.SISLoginID != identifier {
			roster = append(roster, s)
		}
	}
	cls.Students = roster
	err = cls.Save()
	if err != nil {
		w.renderErr(r.Context(), wr, err)
		return
	}
	w.studentTable(cls).Render(r.Context(), wr)
}

func (w *Web) handleStudentForm(wr http.ResponseWriter, r *http.Request) {
	cls, err := classroom.New(w.db, chi.URLParam(r, "course"))
	if err != nil {
		w.renderErr(r.Context(), wr, err)
		return
	}
	w.newStudentForm(cls).Render(r.Context(), wr)
}

func (w *Web) handleNewStudent(wr http.ResponseWriter, r *http.Request) {
	// Note: there is no current expectation that the form includes all of these elements
	student := classroom.Student{
		Name:           r.FormValue("name"),
		ID:             r.FormValue("id"),
		SISLoginID:     r.FormValue("sisloginid"),
		Section:        r.FormValue("section"),
		GitHubUsername: r.FormValue("githubusername"),
		GithubID:       r.FormValue("githubid"),
	}
	cls, err := classroom.New(w.db, chi.URLParam(r, "course"))
	if err != nil {
		w.renderErr(r.Context(), wr, err)
		return
	}
	cls.Students = append(cls.Students, student)
	err = cls.Save()
	if err != nil {
		w.renderErr(r.Context(), wr, err)
		return
	}
	w.studentTable(cls).Render(r.Context(), wr)
}
