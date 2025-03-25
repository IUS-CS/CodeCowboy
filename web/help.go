package web

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (w *Web) setupHelpHandlers() chi.Router {
	r := chi.NewRouter()
	r.Get("/", w.handleHelpRoot)
	return r
}

func (w *Web) handleHelpRoot(wr http.ResponseWriter, r *http.Request) {
	w.Index("Help", w.help()).Render(r.Context(), wr)
}

