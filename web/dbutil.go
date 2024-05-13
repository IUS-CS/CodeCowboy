package web

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

func (w *Web) setupDBUtilHandlers() chi.Router {
	r := chi.NewRouter()
	r.Get("/", w.handleDBRoot)
	r.Get("/export", w.handleDBExport)
	r.Post("/import", w.handleDBImport)
	return r
}

func (w *Web) handleDBImport(wr http.ResponseWriter, r *http.Request) {
	file, _, err := r.FormFile("db")
	if err != nil {
		w.renderErr(r.Context(), wr, err)
		return
	}
	fileContents, err := io.ReadAll(file)
	if err != nil {
		w.renderErr(r.Context(), wr, err)
		return
	}
	err = w.db.Import(fileContents)
	if err != nil {
		w.renderErr(r.Context(), wr, err)
		return
	}
	wr.Header().Set("HX-Redirect", "/")
	w.dbUtil().Render(r.Context(), wr)
}

func (w *Web) handleDBExport(wr http.ResponseWriter, r *http.Request) {
	data, err := w.db.Export()
	if err != nil {
		w.renderErr(r.Context(), wr, err)
	}
	wr.Header().Set("Content-Disposition",
		fmt.Sprintf("attachment; filename=export_%s.json", time.Now().Format(time.RFC3339)))
	wr.Header().Set("Content-Type", "application/json")
	wr.Write(data)
}

func (w *Web) handleDBRoot(wr http.ResponseWriter, r *http.Request) {
	w.Index("DB Utils", w.dbUtil()).Render(r.Context(), wr)
}
