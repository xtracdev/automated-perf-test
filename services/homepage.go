package services

import (
"net/http"
	"github.com/go-chi/chi"
)

const contentTypeHeader = `Content-Type`
const htmlType = `text/html`

func StartUiMode() *chi.Mux{
	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./ui/index.html")
		w.Header().Set(contentTypeHeader, htmlType)

	})

	http.ListenAndServe(":9191", r)
	return r
}


