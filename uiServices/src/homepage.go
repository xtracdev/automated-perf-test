package services

import (
	"github.com/Sirupsen/logrus"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"io/ioutil"
	"net/http"
	"path/filepath"

	"os"
)

const contentTypeHeader = `Content-Type`
const htmlType = `text/html`

func StartUiMode() {
	http.ListenAndServe(":9191", GetRouter())
}

func GetRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Mount("/", GetIndexPage())

	logrus.Print("http://localhost:9191")

	return r
}

func GetIndexPage() *chi.Mux {
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
		return
	})

	router.Mount("/configs", routeConfigs())

	router.Get("/*", func(w http.ResponseWriter, r *http.Request) {

		resourceUrl := r.URL.String()

		path := os.Getenv("GOPATH") + "/src/github.com/xtracdev/automated-perf-test/ui/" + resourceUrl

		absPath, err := filepath.Abs(path)

		if err != nil {
			logrus.Error("Unable to get absolute path", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if _, err := os.Stat(absPath); os.IsNotExist(err) {
			logrus.Error("Resource cannot be found", err)
			w.WriteHeader(http.StatusNotFound)
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
		return
	})

	return router
}

func routeConfigs() http.Handler {
	router := chi.NewRouter()
	router.Use(ConfigCtx)
	router.Post("/", postConfigs)
	router.Get("/", getConfigs)
	router.Put("/", putConfigs)

	return router
}
