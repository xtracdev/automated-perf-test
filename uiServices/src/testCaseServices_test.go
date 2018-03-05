package services

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/go-chi/chi"
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

	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, w.Code, "Did not get all test cases")
	}
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

	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusBadRequest, w.Code, "Did not get all test cases")
	}
}

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

	assert.Equal(t, http.StatusBadRequest, w.Code, "Should not return data")
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

	assert.Equal(t, http.StatusNotFound, w.Code, "Should not return data")
}

func testDeleteAllCases(t *testing.T) {
	r := chi.NewRouter()
	r.Mount("/", GetIndexPage())

	DirectoryPath := os.Getenv("GOPATH") + "/src/github.com/xtracdev/automated-perf-test/uiServices/test/"
	err := ioutil.WriteFile(fmt.Sprintf("%s%s.xml", DirectoryPath, "test"), nil, 0666)
	if err != nil {
		logrus.Errorf("Error trying to create a file", err)
	}

	request, err := http.NewRequest(http.MethodDelete, "/test-cases/all", nil)
	if err != nil {
		logrus.Errorf("Error trying to delete files", err)
	}

	request.Header.Set("testCasePathDir", DirectoryPath)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, request)
	assert.Equal(t, http.StatusNoContent, w.Code, "all files deleted : 204")

}
