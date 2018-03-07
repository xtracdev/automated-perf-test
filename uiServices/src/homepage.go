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

func StartUIMode() {
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
	router.Mount("/test-suites", routeTestSuites())
	router.Mount("/test-cases", routeTestCases())

	router.Get("/*", func(w http.ResponseWriter, r *http.Request) {

		resourceURL := r.URL.String()

		path := os.Getenv("GOPATH") + "/src/github.com/xtracdev/automated-perf-test/ui/" + resourceURL

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
	router.Get("/{configName}", getConfigs)
	router.Put("/{configName}", putConfigs)
	router.Put("/filename/{configFileName}/{newConfigFileName}", putConfigFileName)

	return router
}

func routeTestSuites() http.Handler {
	router := chi.NewRouter()
	router.Use(TestSuiteCtx)
	router.Post("/", postTestSuites)
	router.Put("/{testSuiteName}", putTestSuites)
	router.Get("/{testSuiteName}", getTestSuite)
	router.Get("/", getAllTestSuites)
	router.Delete("/{testSuiteName}", deleteTestSuite)

	return router
}

func routeTestCases() http.Handler {
	router := chi.NewRouter()
	router.Use(TestCaseCtx)
	router.Post("/", postTestCase)
	router.Put("/{testCaseName}", putTestCase)
	router.Get("/", getAllTestCases)
	router.Get("/{testCaseName}", getTestCase)
	router.Delete("/{testCaseName}", deleteTestCase)
	router.Delete("/all", deleteAllTestCases)

	return router
}
