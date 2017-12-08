package main

import (
"net/http"
)

func startUiMode() {
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir("./ui")))
	http.ListenAndServe(":9191", mux)
}




