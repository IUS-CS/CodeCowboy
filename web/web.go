package web

import (
	"context"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/go-chi/chi/v5/middleware"

	"cso/codecowboy/store"

	"github.com/go-chi/chi/v5"

	"sync"
)

type Web struct {
	listenAddr string
	db         *store.DB

	runLog map[string]map[string]string
	mu     sync.Mutex // Adding Mutex to protect runLog
}

func New(db *store.DB, listenAddr string) *Web {
	return &Web{db: db, listenAddr: listenAddr, runLog: map[string]map[string]string{}}
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

	router.Use(middleware.DefaultLogger)

	router.Get("/", func(wr http.ResponseWriter, r *http.Request) {
		http.Redirect(wr, r, "/courses", http.StatusFound)
	})

	router.Mount("/import", w.setupImportHandlers())
	router.Mount("/courses", w.setupCourseHandlers())
	router.Mount("/db", w.setupDBUtilHandlers())

	workDir, _ := os.Getwd()
	filesDir := http.Dir(filepath.Join(workDir, "static"))
	fileServer(router, "/static", filesDir)

	return http.ListenAndServe(w.listenAddr, router)
}

func fileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}

func (w *Web) renderErr(ctx context.Context, wr http.ResponseWriter, err error) {
	log.Error("Controller error", "err", err)
	w.Error(err.Error()).Render(ctx, wr)
}
