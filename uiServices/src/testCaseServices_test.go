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

	request.Header.Set("testSuitePathDir", DirectoryPath)
	request.Header.Get("testSuitePathDir")

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

	request.Header.Set("testSuitePathDir", DircetoryPath)
	request.Header.Get("testSuitePathDir")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, request)

	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, http.StatusBadRequest, w.Code, "Should not have retrieved all test cases")
}
