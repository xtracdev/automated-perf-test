package main

import (
"net/http"
///"github.com/go-chi/chi"
//"fmt"
)

/*func main() {
	fmt.Print("testq")
	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("xxxx"))
	})
	http.ListenAndServe(":9191", r)


}*/




func main() {
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir("./ui")))
	http.ListenAndServe(":9191", mux)
}




