package web

import (
	"cso/codecowboy/classroom"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func (w *Web) setupStudentHandlers() chi.Router {
	router := chi.NewRouter()

	router.Delete("/{sisID}", w.handleRmStudent)
	router.Post("/", w.handleNewStudent)
	router.Get("/", w.handleStudentForm)

	return router
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
