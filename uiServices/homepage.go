package services

import (
	"github.com/Sirupsen/logrus"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"

	"os"
)

const contentTypeHeader = `Content-Type`
const htmlType = `text/html`

func StartUiMode() {

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Mount("/", getIndexPage())

	log.Print("http:\\localhost:9191")

	http.ListenAndServe(":9191", r)
}

func getIndexPage() *chi.Mux {
	router := chi.NewRouter()

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		goPath := os.Getenv("GOPATH")
		absPath, err := filepath.Abs(goPath + "/src/github.com/xtracdev/automated-perf-test/ui/index.html")

		if err != nil {
			logrus.Error("Unable to find homepage", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		htmlBytes, err := ioutil.ReadFile(absPath)

		if err != nil {
			logrus.Error("Unable to read file: ", absPath, err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set(contentTypeHeader, htmlType)
		w.Write([]byte(htmlBytes))
	})
	router.HandleFunc("/configs", configsHandler)

	return router
}
