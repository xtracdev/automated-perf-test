package services

import (
	"github.com/go-chi/chi"

	"net/http"
	"net/http/httptest"
	"os"

	"testing"
	"github.com/stretchr/testify/assert"
	"strings"
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

const invalidTestSuite = `
{
  "name": 123,
  "testStrategy": "SuiteBased",
  "testCases": [
    {
      "name":"file1.xml",
      "preThinkTime": "number",
      "postThinkTime": 2000,
      "execWeight": "Infrequent"
    },
    {
      "name":"file2.xml",
      "preThinkTime": 1,
      "postThinkTime": 10,
      "execWeight": "Sparse"
    }
  ]
}
`
const TestSuiteMissingRequired = `
{
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



func TestValidTestSuitePut(t *testing.T) {
	r := chi.NewRouter()
	r.Mount("/", GetIndexPage())

	reader := strings.NewReader(validTestSuite)
	r.HandleFunc("/test-suites", putTestSuites)

	filePath := os.Getenv("GOPATH") + "/src/github.com/xtracdev/automated-perf-test/uiServices/test/"
	request, err := http.NewRequest(http.MethodPut, "/test-suites/TestSuiteService", reader)
	request.Header.Set("testSuitePathDir", filePath)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, request)

	if err != nil {
		t.Error(err)
	}

	assert.Equal(t,  http.StatusNoContent, w.Code, "Did Not successfully Update")
}

func TestTestSuiteMissingFieldPut(t *testing.T) {
	r := chi.NewRouter()
	r.Mount("/", GetIndexPage())

	reader := strings.NewReader(TestSuiteMissingRequired)
	r.HandleFunc("/test-suites", putTestSuites)

	filePath := os.Getenv("GOPATH") + "/src/github.com/xtracdev/automated-perf-test/uiServices/test/"
	request, err := http.NewRequest(http.MethodPut, "/test-suites/TestSuiteService", reader)
	request.Header.Set("testSuitePathDir", filePath)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, request)

	if err != nil {
		t.Error(err)
	}

	assert.Equal(t,http.StatusBadRequest, w.Code,  "Sucessfully updated. Field Should be missing so update shouldn't occur")
}


func TestInvalidTestSuitePut(t *testing.T) {
	r := chi.NewRouter()
	r.Mount("/", GetIndexPage())

	reader := strings.NewReader(invalidTestSuite)
	r.HandleFunc("/test-suites", putTestSuites)

	filePath := os.Getenv("GOPATH") + "/src/github.com/xtracdev/automated-perf-test/uiServices/test/"
	request, err := http.NewRequest(http.MethodPut, "/test-suites/TestSuiteService", reader)
	request.Header.Set("testSuitePathDir", filePath)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, request)

	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, http.StatusBadRequest, w.Code,  "Sucessfully updated. Field data type should have been incorrect so update should occur")
}


func TestInvalidUrlTestSuitePut(t *testing.T) {
	r := chi.NewRouter()
	r.Mount("/", GetIndexPage())

	reader := strings.NewReader(validTestSuite)
	r.HandleFunc("/test-suites", putTestSuites)

	filePath := os.Getenv("GOPATH") + "/src/github.com/xtracdev/automated-perf-test/uiServices/test/"
	request, err := http.NewRequest(http.MethodPut, "/test-suites/xxx", reader)
	request.Header.Set("testSuitePathDir", filePath)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, request)

	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, http.StatusNotFound, w.Code, "Sucessfully updated. Should have have worked using /test-suites/xxx")
}

func TestNoUrlTestSuitePut(t *testing.T) {
	r := chi.NewRouter()
	r.Mount("/", GetIndexPage())

	reader := strings.NewReader(validTestSuite)
	r.HandleFunc("/test-suites", putTestSuites)

	filePath := os.Getenv("GOPATH") + "/src/github.com/xtracdev/automated-perf-test/uiServices/test/"
	request, err := http.NewRequest(http.MethodPut, "", reader)
	request.Header.Set("testSuitePathDir", filePath)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, request)

	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, http.StatusNotFound, w.Code,  "Sucessfully updated. Should not have worked with no URL")
}

func TestPutWithNoPathSlash(t *testing.T) {
	r := chi.NewRouter()
	r.Mount("/", GetIndexPage())

	reader := strings.NewReader(validTestSuite)
	r.HandleFunc("/test-suites", putTestSuites)

	filePath := os.Getenv("GOPATH") + "/src/github.com/xtracdev/automated-perf-test/uiServices/test"
	request, err := http.NewRequest(http.MethodPut, "/test-suites/TestSuiteService", reader)
	request.Header.Set("testSuitePathDir", filePath)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, request)

	if err != nil {
		t.Error(err)
	}

	assert.Equal(t,http.StatusNoContent, w.Code,  "Did not update. Should have added '/' to path")
}

func TestNoPathTestSuitePut(t *testing.T) {
	r := chi.NewRouter()
	r.Mount("/", GetIndexPage())

	reader := strings.NewReader(validTestSuite)
	r.HandleFunc("/test-suites", putTestSuites)

	filePath := ""
	request, err := http.NewRequest(http.MethodPut, "/test-suites/TestSuiteService", reader)
	request.Header.Set("testSuitePathDir", filePath)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, request)

	if err != nil {
		t.Error(err)
	}

	assert.Equal(t,  http.StatusBadRequest, w.Code, "Successfully updated. Should not have worked due to no filepath")
}

func TestNoFileNameTestSuitePut(t *testing.T) {
	r := chi.NewRouter()
	r.Mount("/", GetIndexPage())

	reader := strings.NewReader(validTestSuite)
	r.HandleFunc("/test-suites", putTestSuites)

	filePath := os.Getenv("GOPATH") + "/src/github.com/xtracdev/automated-perf-test/uiServices/test"
	request, err := http.NewRequest(http.MethodPut, "/test-suites", reader)
	request.Header.Set("testSuitePathDir", filePath)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, request)

	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, http.StatusBadRequest, w.Code, "Successfully updated. Should not have worked due to no file name given")
}




func TestSuccessfulGetTestSuite(t *testing.T) {
	r := chi.NewRouter()
	r.Mount("/", GetIndexPage())

	filePath := os.Getenv("GOPATH") + "/src/github.com/xtracdev/automated-perf-test/uiServices/test/"
	request, err := http.NewRequest(http.MethodGet, "/test-suites/TestSuiteService", nil)

	request.Header.Set("testSuitePathDir", filePath)
	request.Header.Get("testSuitePathDir")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, request)

	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, w.Code, http.StatusOK, "Error. Did not successfully GET")
}

func TestGetTestSuiteNoPath(t *testing.T) {
	r := chi.NewRouter()
	r.Mount("/", GetIndexPage())

	filePath := ""
	request, err := http.NewRequest(http.MethodGet, "/test-suites/TestSuiteService.xml", nil)

	request.Header.Set("testSuitePathDir", filePath)
	request.Header.Get("testSuitePathDir")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, request)

	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, w.Code, http.StatusBadRequest, "Retrived file but should not have as there is no path")
}

func TestGetTestSuiteFileNotFound(t *testing.T) {
	r := chi.NewRouter()
	r.Mount("/", GetIndexPage())

	filePath := os.Getenv("GOPATH") + "/src/github.com/xtracdev/automated-perf-test/uiServices/test/"
	request, err := http.NewRequest(http.MethodGet, "/test-suites/xxx", nil)

	request.Header.Set("testSuitePathDir", filePath)
	request.Header.Get("testSuitePathDir")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, request)

	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, w.Code, http.StatusNotFound, "Retrived a file but should not have as there is no file")
}
