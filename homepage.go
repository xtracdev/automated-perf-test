package main

import (
"net/http"
	"github.com/go-chi/chi"
)

func startUiMode()  {

	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./ui/index.html")
	})
	r.HandleFunc("/post", jsonHandler)
	http.ListenAndServe(":9191", r)
}


