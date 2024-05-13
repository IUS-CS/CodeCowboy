package web

import (
	"context"
	"github.com/charmbracelet/log"
	"net/http"

	"cso/codecowboy/store"

	"github.com/go-chi/chi/v5"
)

type Web struct {
	listenAddr string
	db         *store.DB
}

func New(db *store.DB, listenAddr string) *Web {
	return &Web{db: db, listenAddr: listenAddr}
}

func (w *Web) SiteName() string {
	return "CodeCowboy 🤠"
}

type NavItem struct {
	Name string
	URL  string
}

func (w *Web) Navs() []NavItem {
	return []NavItem{
		{"Courses", "/courses"},
		{"Import", "/import"},
		{"DB Utils", "/db"},
	}
}

func (w *Web) ListenAndServe() error {
	router := chi.NewRouter()
	router.Get("/", func(wr http.ResponseWriter, r *http.Request) {
		http.Redirect(wr, r, "/courses", http.StatusFound)
	})

	router.Mount("/import", w.setupImportHandlers())
	router.Mount("/courses", w.setupCourseHandlers())
	router.Mount("/db", w.setupDBUtilHandlers())

	return http.ListenAndServe(w.listenAddr, router)
}

func (w *Web) renderErr(ctx context.Context, wr http.ResponseWriter, err error) {
	log.Error("Controller error", "err", err)
	w.Error(err.Error()).Render(ctx, wr)
}
