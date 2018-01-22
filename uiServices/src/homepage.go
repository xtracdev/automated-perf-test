package services

import (
    "io/ioutil"
    "net/http"
    "path/filepath"
    "github.com/Sirupsen/logrus"
    "github.com/go-chi/chi"
    "github.com/go-chi/chi/middleware"

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

    logrus.Print("http:\\localhost:9191")

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
    })

    router.Mount("/configs", routeConfigs())

    return router

    router.Get("/*", func(w http.ResponseWriter, r *http.Request) {
        requestUri := r.RequestURI
        goPath := os.Getenv("GOPATH")
        path := goPath + "/src/github.com/xtracdev/automated-perf-test/ui/index.html"

        if _, err := os.Stat(goPath + "/src/github.com/xtracdev/automated-perf-test/ui/" + requestUri); err == nil {
            path = goPath + "/src/github.com/xtracdev/automated-perf-test/ui/" + requestUri
        }

        absPath, err := filepath.Abs(path)

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
    router.Mount("/configs", routeConfigs())

    return router
}

func routeConfigs() http.Handler {
    router := chi.NewRouter()
    router.Use(ConfigCtx)
    router.Post("/", postConfigs)

    return router
}

