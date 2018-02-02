package services

import (
	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

const validTestSuite = `
{
  "name": "TestSuiteService",
  "testStrategy": "SuiteBased",
  "testCases": [
    {
      "name":"file1.xml",
      "preThinkTime": 1000,
      "postThinkTime": 2000,
      "execWeight": "infrequent"
    },
    {
      "name":"file2.xml",
      "preThinkTime": 1,
      "postThinkTime": 10,
      "execWeight": "sparce"
    }
  ]
}
`
const InvalidTestSuiteMissingRequired = `
{
  "testStrategy": "",
  "testCases": [
    {
      "name":"file1.xml",
      "preThinkTime": 1000,
      "postThinkTime": 2000,
      "execWeight": "infrequent"
    },
    {
      "name":"file2.xml",
      "preThinkTime": 1,
      "postThinkTime": 10,
      "execWeight": "sparce"
    }
  ]
}
`
const InvalidTestSuiteIncorrectData = `
{
  "name": 9,
  "testStrategy": 9,
  "testCases": [
    {
      "name":"file1.xml",
      "preThinkTime": "NUMBER",
      "postThinkTime": "NUMBER",
      "execWeight": "infrequent"
    }
  ]
}
`

func TestValidTestSuitePost(t *testing.T) {
	r := chi.NewRouter()
	r.Mount("/", GetIndexPage())

	reader := strings.NewReader(validTestSuite)
	r.HandleFunc("/test-suites", postTestSuites)

	os.Remove(os.Getenv("GOPATH") + "/src/github.com/xtracdev/automated-perf-test/uiServices/test/TestSuiteService.xml")

	filePath := os.Getenv("GOPATH") + "/src/github.com/xtracdev/automated-perf-test/uiServices/test"
	request, err := http.NewRequest(http.MethodPost, "/test-suites", reader)
	request.Header.Set("configPathDir", filePath)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, request)

	if err != nil {
		t.Error(err)

		assert.Equal(t, w.Code, http.StatusCreated, "Error: Did Not Successfully Post")
	}
}

func TestFileExistsPost(t *testing.T) {
	r := chi.NewRouter()
	r.Mount("/", GetIndexPage())

	reader := strings.NewReader(validTestSuite)
	r.HandleFunc("/test-suites", postTestSuites)

	filePath := os.Getenv("GOPATH") + "/src/github.com/xtracdev/automated-perf-test/uiServices/test"
	request, err := http.NewRequest(http.MethodPost, "/test-suites", reader)
	request.Header.Set("configPathDir", filePath)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, request)

	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, w.Code, http.StatusBadRequest, "Created, but should not have. File Should Alreday Exist")
}

func TestMissingRequiredTestSuitePost(t *testing.T) {
	r := chi.NewRouter()
	r.Mount("/", GetIndexPage())

	reader := strings.NewReader(InvalidTestSuiteMissingRequired)
	r.HandleFunc("/test-suites", postTestSuites)

	filePath := os.Getenv("GOPATH") + "/src/github.com/xtracdev/automated-perf-test/uiServices/test"
	request, err := http.NewRequest(http.MethodPost, "/test-suites", reader)
	request.Header.Set("configPathDir", filePath)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, request)

	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, w.Code, http.StatusBadRequest, "Created TestSuite but should not have due to missing fields")

}

func TestIncorrectDataTestSuitePost(t *testing.T) {
	r := chi.NewRouter()
	r.Mount("/", GetIndexPage())

	reader := strings.NewReader(InvalidTestSuiteIncorrectData)
	r.HandleFunc("/test-suites", postTestSuites)

	filePath := os.Getenv("GOPATH") + "/src/github.com/xtracdev/automated-perf-test/uiServices/test"
	request, err := http.NewRequest(http.MethodPost, "/test-suites", reader)
	request.Header.Set("configPathDir", filePath)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, request)

	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, w.Code, http.StatusBadRequest, "Created TestSuite but should not have due to incorrect data types")

}

func TestValidTestSuitePostNoConfigPathDir(t *testing.T) {
	r := chi.NewRouter()
	r.Mount("/", GetIndexPage())

	reader := strings.NewReader(validTestSuite)
	r.HandleFunc("/test-suites", postTestSuites)

	filePath := ""
	request, err := http.NewRequest(http.MethodPost, "/test-suites", reader)
	request.Header.Set("configPathDir", filePath)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, request)

	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, w.Code, http.StatusBadRequest, "successfully created but should not have due to no config path specified")

}
