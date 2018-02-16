package services

import (
	"github.com/go-chi/chi"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"github.com/stretchr/testify/assert"
)


func TestSuccessfulGetAllCases(t *testing.T) {
	r := chi.NewRouter()
	r.Mount("/", GetIndexPage())

	DirectoryPath := os.Getenv("GOPATH") + "/src/github.com/xtracdev/automated-perf-test/uiServices/test/"
	request, err := http.NewRequest(http.MethodGet, "/test-cases", nil)

	request.Header.Set("testCasePathDir", DirectoryPath)
	request.Header.Get("testCasePathDir")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, request)

	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, http.StatusOK, w.Code, "Did not get all test casess")
}

func TestGetAllCasesNoHeader(t *testing.T) {
	r := chi.NewRouter()
	r.Mount("/", GetIndexPage())

	DircetoryPath := ""
	request, err := http.NewRequest(http.MethodGet, "/test-cases", nil)

	request.Header.Set("testCasePathDir", DircetoryPath)
	request.Header.Get("testCasePathDir")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, request)

	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, http.StatusBadRequest, w.Code, "Should not have retrieved all test cases")
}


//func TestSuccessfulGetTestCase(t *testing.T) {
//	r := chi.NewRouter()
//	r.Mount("/", GetIndexPage())
//
//	filePath := os.Getenv("GOPATH") + "/src/github.com/xtracdev/automated-perf-test/uiServices/test/"
//	request, err := http.NewRequest(http.MethodGet, "/test-cases/Case1", nil)
//
//	request.Header.Set("testCasePathDir", filePath)
//	request.Header.Get("testCasePathDir")
//
//	w := httptest.NewRecorder()
//	r.ServeHTTP(w, request)
//
//	if err != nil {
//		t.Error(err)
//	}
//
//	assert.Equal(t,http.StatusOK, w.Code, "Error. Did not successfully GET")
//}

func TestGetTestCaseNoHeader(t *testing.T) {
	r := chi.NewRouter()
	r.Mount("/", GetIndexPage())

	filePath := ""
	request, err := http.NewRequest(http.MethodGet, "/test-cases/Case1", nil)

	request.Header.Set("testCasePathDir", filePath)
	request.Header.Get("testCasePathDir")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, request)

	if err != nil {
		t.Error(err)
	}

	assert.Equal(t,http.StatusBadRequest, w.Code, "Should not return data")
}

func TestGetTestCaseFileNotFound(t *testing.T) {
	r := chi.NewRouter()
	r.Mount("/", GetIndexPage())

	filePath := os.Getenv("GOPATH") + "/src/github.com/xtracdev/automated-perf-test/uiServices/test/"
	request, err := http.NewRequest(http.MethodGet, "/test-cases/xxx", nil)

	request.Header.Set("testCasePathDir", filePath)
	request.Header.Get("testCasePathDir")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, request)

	if err != nil {
		t.Error(err)
	}

	assert.Equal(t,http.StatusNotFound, w.Code, "Should not return data")
}