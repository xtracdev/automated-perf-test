package services

import (
	"testing"
	"github.com/go-chi/chi"
	"os"
	"net/http/httptest"
	"github.com/stretchr/testify/assert"
	"net/http"
)


func TestSuccessfulDelete(t *testing.T) {
	r := chi.NewRouter()
	r.Mount("/", GetIndexPage())

	filePath := os.Getenv("GOPATH") + "/src/github.com/xtracdev/automated-perf-test/uiServices/test/"
	request, err := http.NewRequest(http.MethodDelete, "/test-cases/Case1", nil)

	request.Header.Set("testCasePathDir", filePath)
	request.Header.Get("testCasePathDir")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, request)

	if err != nil {
		t.Error(err)
	}

	assert.Equal(t,http.StatusNoContent, w.Code, "Error. Did not successfully Delete")
}

func TestDeleteFileNotFound(t *testing.T) {
	r := chi.NewRouter()
	r.Mount("/", GetIndexPage())

	filePath := os.Getenv("GOPATH") + "/src/github.com/xtracdev/automated-perf-test/uiServices/test/"
	request, err := http.NewRequest(http.MethodDelete, "/test-suites/xxx", nil)

	request.Header.Set("testSuitePathDir", filePath)
	request.Header.Get("testSuitePathDir")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, request)

	if err != nil {
		t.Error(err)
	}

	assert.Equal(t,http.StatusNotFound, w.Code, "Should not have successfully deleted")
}

func TestDeleteWithNoHeader(t *testing.T) {
	r := chi.NewRouter()
	r.Mount("/", GetIndexPage())

	filePath := ""
	request, err := http.NewRequest(http.MethodDelete, "/test-suites/TestSuiteService", nil)

	request.Header.Set("testSuitePathDir", filePath)
	request.Header.Get("testSuitePathDir")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, request)

	if err != nil {
		t.Error(err)
	}

	assert.Equal(t,http.StatusBadRequest, w.Code, "Should not have successfully deleted")
}

