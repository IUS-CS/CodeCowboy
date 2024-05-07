package web

import (
	"net/http"
)

type Web struct {
	listenAddr string
}

func New(listenAddr string) *Web {
	return &Web{listenAddr: listenAddr}
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
	http.HandleFunc("/", func(wr http.ResponseWriter, r *http.Request) {
		Index(w, "Index", nil).Render(r.Context(), wr)
	})
	return http.ListenAndServe(w.listenAddr, nil)
}
