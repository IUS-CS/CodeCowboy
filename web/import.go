package web

import (
	"errors"
	"net/http"

	"cso/codecowboy/canvasfmt"
	"cso/codecowboy/classroom"
	"cso/codecowboy/githubfmt"

	"github.com/go-chi/chi/v5"
)

func (w *Web) setupImportHandlers() chi.Router {
	router := chi.NewRouter()

	router.Get("/", w.handleImportRoot)
	router.Post("/", w.handleImportData)

	return router
}

func (w *Web) handleImportRoot(wr http.ResponseWriter, r *http.Request) {
	currentUser := w.getCurrentUser(r)
	w.Index("Import", currentUser, w.importForm()).Render(r.Context(), wr)
}

func (w *Web) handleImportData(wr http.ResponseWriter, r *http.Request) {
	course := r.FormValue("course")
	if course == "" {
		w.renderErr(r.Context(), wr, errors.New("course is required"))
	}
	cls, err := classroom.New(w.db, course)
	if err != nil {
		w.renderErr(r.Context(), wr, err)
		return
	}
	roster := cls.Students
	file, _, err := r.FormFile("assignments")
	if !errors.Is(err, http.ErrMissingFile) && err != nil {
		w.renderErr(r.Context(), wr, err)
		return
	} else if !errors.Is(err, http.ErrMissingFile) {
		assignments, err := classroom.ParseAssignments(file, cls.Name)
		if err != nil {
			w.renderErr(r.Context(), wr, err)
			return
		}
		cls.Assignments = assignments
	}
	file, _, err = r.FormFile("github")
	if !errors.Is(err, http.ErrMissingFile) && err != nil {
		w.renderErr(r.Context(), wr, err)
		return
	} else if !errors.Is(err, http.ErrMissingFile) {
		roster, err = githubfmt.Parse(file, roster)
		if err != nil {
			w.renderErr(r.Context(), wr, err)
			return
		}
	}
	file, _, err = r.FormFile("canvas")
	if !errors.Is(err, http.ErrMissingFile) && err != nil {
		w.renderErr(r.Context(), wr, err)
		return
	} else if !errors.Is(err, http.ErrMissingFile) {
		roster, err = canvasfmt.Parse(file, roster)
		if err != nil {
			w.renderErr(r.Context(), wr, err)
			return
		}
	}
	cls.Students = roster

	currentUser := w.getCurrentUser(r)
	cls.Instructors = []string{currentUser}

	err = cls.Save()
	if err != nil {
		w.renderErr(r.Context(), wr, err)
		return
	}
	wr.Header().Set("HX-Redirect", "/courses/"+course)
	w.courseDetails(cls).Render(r.Context(), wr)
}
