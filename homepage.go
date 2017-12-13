package main

import (
"net/http"
"github.com/go-chi/chi"
)

const contentTypeHeader = `Content-Type`
const htmlType = `text/html`

func startUiMode() *chi.Mux{
	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./ui/index.html")
		w.Header().Set(contentTypeHeader, htmlType)

	})
	r.HandleFunc("/post", jsonHandler)
	http.ListenAndServe(":9191", r)
	return r
}