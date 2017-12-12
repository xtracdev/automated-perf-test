package main

import (
"net/http"
)

func startUiMode() (int) {
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir("./ui/index.html")))
	mux.HandleFunc("/test", test)

	http.ListenAndServe(":9191", mux)

	return http.StatusOK
}





