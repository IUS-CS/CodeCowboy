package web

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (w *Web) setupDBUtilHandlers() chi.Router {
	r := chi.NewRouter()
	r.HandleFunc("/", w.handleDBRoot)
	return r
}

func (w *Web) handleDBRoot(wr http.ResponseWriter, r *http.Request) {
	w.Index("DB Utils", w.dbUtil()).Render(r.Context(), wr)
}
