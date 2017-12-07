package main

import (
	"net/http"
	"github.com/go-chi/chi"
	"fmt"
)

func main() {
	fmt.Print("http:\\localhost:9191")
	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("HWLLO WORLD"))
	})
	http.ListenAndServe(":9191", r)


}