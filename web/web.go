package web

import (
	"cso/codecowboy/classroom"
	"cso/codecowboy/store"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type Web struct {
	listenAddr string
	db         *store.DB
}

func New(db *store.DB, listenAddr string) *Web {
	return &Web{db: db, listenAddr: listenAddr}
}

func (w *Web) SiteName() string {
	return "CodeCowboy"
}

type NavItem struct {
	Name string
	URL  string
}

func (w *Web) Navs() []NavItem {
	return nil
}

func (w *Web) ListenAndServe() error {
	router := chi.NewRouter()
	router.HandleFunc("/", func(wr http.ResponseWriter, r *http.Request) {
		w.Index("Index", nil).Render(r.Context(), wr)
	})
	router.HandleFunc("/courses", func(wr http.ResponseWriter, r *http.Request) {
		courses, err := classroom.All(w.db)
		if err != nil {
			w.Index("Error", w.Error(err.Error())).Render(r.Context(), wr)
		}
		w.Index("Courses", w.courseList(courses)).Render(r.Context(), wr)
	})
	return http.ListenAndServe(w.listenAddr, router)
}
