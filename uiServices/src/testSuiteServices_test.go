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
  "description": "Services for XYZ",
  "testCases": [
    {
      "name":"file1",
      "preThinkTime": 1000,
      "postThinkTime": 2000,
      "execWeight": "Infrequent",
       "description": "Desc1"
    },
    {
      "name":"file2",
       "description": "Desc2"
    }
  ]
}
`
const TestSuiteMissingRequired = `
{
  "name":"",
  "testCases": [
    {
      "name":"file1.xml",
      "preThinkTime": 1000,
      "postThinkTime": 2000,
      "execWeight": "Infrequent"
    }
  ]
}
`
const TestSuiteNoName = `
{
  "name": "",
  "testStrategy": "SuiteBased",
   "description": "Services for XYZ",
  "testCases": [
    {
      "name":"file1.xml",
      "preThinkTime": "xxxx"
      "postThinkTime": 2000,
      "execWeight": 123,
       "description": "Desc"
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
	request.Header.Set("testSuitePathDir", filePath)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, request)

	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, http.StatusCreated, w.Code, "Error: Did Not Successfully Post")
}

func TestFileExistsPost(t *testing.T) {
	r := chi.NewRouter()
	r.Mount("/", GetIndexPage())

	reader := strings.NewReader(validTestSuite)
	r.HandleFunc("/test-suites", postTestSuites)

	filePath := os.Getenv("GOPATH") + "/src/github.com/xtracdev/automated-perf-test/uiServices/test"
	request, err := http.NewRequest(http.MethodPost, "/test-suites", reader)
	request.Header.Set("testSuitePathDir", filePath)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, request)

	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, http.StatusBadRequest, w.Code, "Created, but should not have. File Should Alreday Exist")
}

func TestMissingRequiredTestSuitePost(t *testing.T) {
	r := chi.NewRouter()
	r.Mount("/", GetIndexPage())

	reader := strings.NewReader(TestSuiteMissingRequired)
	r.HandleFunc("/test-suites", postTestSuites)

	filePath := os.Getenv("GOPATH") + "/src/github.com/xtracdev/automated-perf-test/uiServices/test"
	request, err := http.NewRequest(http.MethodPost, "/test-suites", reader)
	request.Header.Set("testSuitePathDir", filePath)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, request)

	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, http.StatusBadRequest, w.Code, "Created TestSuite but should not have due to missing fields")

}

func TestIncorrectDataTypePost(t *testing.T) {
	r := chi.NewRouter()
	r.Mount("/", GetIndexPage())

	reader := strings.NewReader(TestSuiteNoName)
	r.HandleFunc("/test-suites", postTestSuites)

	filePath := os.Getenv("GOPATH") + "/src/github.com/xtracdev/automated-perf-test/uiServices/test"
	request, err := http.NewRequest(http.MethodPost, "/test-suites", reader)
	request.Header.Set("testSuitePathDir", filePath)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, request)

	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, http.StatusBadRequest, w.Code, "Created TestSuite but should not have due to incorrect data types")

}

func TestValidTestSuitePostNoConfigPathDir(t *testing.T) {
	r := chi.NewRouter()
	r.Mount("/", GetIndexPage())

	reader := strings.NewReader(validTestSuite)
	r.HandleFunc("/test-suites", postTestSuites)

	filePath := ""
	request, err := http.NewRequest(http.MethodPost, "/test-suites", reader)
	request.Header.Set("testSuitePathDir", filePath)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, request)

	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, http.StatusBadRequest, w.Code, "successfully created but should not have due to no config path specified")

}

func TestValidTestSuitePostConfigPathDirNotExist(t *testing.T) {
	r := chi.NewRouter()
	r.Mount("/", GetIndexPage())

	reader := strings.NewReader(validTestSuite)
	r.HandleFunc("/test-suites", postTestSuites)

	filePath := "C:/a/b/c/d/////"
	request, err := http.NewRequest(http.MethodPost, "/test-suites", reader)
	request.Header.Set("testSuitePathDir", filePath)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, request)

	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, http.StatusBadRequest, w.Code, "successfully created but should not have due to no config path specified")

}
