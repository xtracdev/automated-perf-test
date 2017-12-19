package services

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/xeipuuv/gojsonschema"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
)

const contentTypeHeader = `Content-Type`
const htmlType = `text/html`

func StartUiMode() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	fmt.Println(gojsonschema.KEY_ALL_OF)
	r.Mount("/", getIndexPage())
	log.Print("http:\\localhost:9191")
	http.ListenAndServe(":9191", r)
}

func getIndexPage() *chi.Mux {
	router := chi.NewRouter()

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		absPath, err := filepath.Abs("./ui/index.html")

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

	return router
}
