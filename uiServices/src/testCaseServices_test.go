package services

import (
	"testing"
	"github.com/go-chi/chi"
	"strings"
	"os"
	"net/http/httptest"
	"github.com/stretchr/testify/assert"
	"net/http"
)

const validTestCase = `
{
   "testname":"TestCaseService",
   "description":"desc",
   "overrideHost":"host",
   "overridePort":"9191",
   "HttpMethod":"GET",
   "BaseURI": "path/to/URI",
   "multipart":false,
   "payload": "payload",
   "responseStatusCode":200,
   "responseContentType": "JSON" ,
   "preThinkTime": 1000,
   "postThinkTime":2000,
   "execWeight": "Sparse",
   "Headers":[{
   	 "Key": "Authorization",
     "Value" :"Header-Value"
   }],
  "ResponseValues":[{
     "Value":"Res-Value",
     "ExtractionKey": "Res-Key"
  }],
  "MultipartPayload":[{
     "fieldName": "F-Name",
   	 "FieldValue":"PayloadName",
     "FileName": "file-name"
  }]
}
`

const TestCaseNoName = `
{
   "testname":"",
   "description":"desc",
   "overrideHost":"host",
   "overridePort":"9191",
   "HttpMethod":"GET",
   "BaseURI": "path/to/URI",
   "multipart":false,
   "payload": "payload",
   "responseStatusCode":200,
   "responseContentType": "JSON" ,
   "preThinkTime": 1000,
   "postThinkTime":2000,
   "execWeight": "Sparse",
   "Headers":[{
   	 "Key": "Authorization",
     "Value" :"Header-Value"
   }],
  "ResponseValues":[{
     "Value":"Res-Value",
     "ExtractionKey": "Res-Key"
  }],
  "MultipartPayload":[{
     "fieldName": "F-Name",
   	 "FieldValue":"PayloadName",
     "FileName": "file-name"
  }]
}
`

const TestCaseMissingRequired = `
{
   "testname":"",
   "description":"",
   "overrideHost":"",
   "overridePort":"",
   "BaseURI": "path/to/URI",
   "multipart":false,
   "payload": "payload",
   "responseStatusCode":200,
   "responseContentType": "JSON" ,
   "preThinkTime": 1000,
   "postThinkTime":2000,
   "execWeight": "Sparse"
}
`


func TestValidTestCasePost(t *testing.T) {
	r := chi.NewRouter()
	r.Mount("/", GetIndexPage())

	reader := strings.NewReader(validTestCase)
	r.HandleFunc("/test-cases", postTestCase)

	os.Remove(os.Getenv("GOPATH") + "/src/github.com/xtracdev/automated-perf-test/uiServices/test/TestCaseService.xml")

	filePath := os.Getenv("GOPATH") + "/src/github.com/xtracdev/automated-perf-test/uiServices/test"
	request, err := http.NewRequest(http.MethodPost, "/test-cases", reader)
	request.Header.Set("testCasePathDir", filePath)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, request)

	assert.NoError(t, err)

	assert.Equal(t, http.StatusCreated, w.Code, "Error: Did Not Successfully Post")
}

func TestCasePostWithExistingFileName(t *testing.T) {
	r := chi.NewRouter()
	r.Mount("/", GetIndexPage())

	reader := strings.NewReader(validTestCase)
	r.HandleFunc("/test-cases", postTestCase)

	filePath := os.Getenv("GOPATH") + "/src/github.com/xtracdev/automated-perf-test/uiServices/test"
	request, err := http.NewRequest(http.MethodPost, "/test-cases", reader)
	request.Header.Set("testCasePathDir", filePath)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, request)

	assert.NoError(t, err)

	assert.Equal(t, http.StatusBadRequest, w.Code, "Should not have Successfully posted")
}

func TestCasePostNoHeader(t *testing.T) {
	r := chi.NewRouter()
	r.Mount("/", GetIndexPage())

	reader := strings.NewReader(validTestCase)
	r.HandleFunc("/test-cases", postTestCase)

	filePath := ""
	request, err := http.NewRequest(http.MethodPost, "/test-cases", reader)
	request.Header.Set("testCasePathDir", filePath)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, request)

	assert.NoError(t, err)

	assert.Equal(t, http.StatusBadRequest, w.Code, "Should not have Successfully posted")
}

func TestPostTestCaseMissingRequiredValues(t *testing.T) {
	r := chi.NewRouter()
	r.Mount("/", GetIndexPage())

	reader := strings.NewReader(TestCaseMissingRequired)
	r.HandleFunc("/test-cases", postTestCase)

	filePath := os.Getenv("GOPATH") + "/src/github.com/xtracdev/automated-perf-test/uiServices/test"
	request, err := http.NewRequest(http.MethodPost, "/test-cases", reader)
	request.Header.Set("testCasePathDir", filePath)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, request)

	assert.NoError(t, err)

	assert.Equal(t, http.StatusBadRequest, w.Code, "Should not have Successfully posted")
}

func TestValidTestCasePut(t *testing.T) {
	r := chi.NewRouter()
	r.Mount("/", GetIndexPage())

	reader := strings.NewReader(validTestCase)
	r.HandleFunc("/test-cases", putTestCase)

	filePath := os.Getenv("GOPATH") + "/src/github.com/xtracdev/automated-perf-test/uiServices/test/"
	request, err := http.NewRequest(http.MethodPut, "/test-cases/TestCaseService", reader)
	request.Header.Set("testCasePathDir", filePath)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, request)

	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, http.StatusNoContent, w.Code, "Did Not successfully Update")
}

func TestTestCasePutMissingRequired(t *testing.T) {
	r := chi.NewRouter()
	r.Mount("/", GetIndexPage())

	reader := strings.NewReader(TestCaseMissingRequired)
	r.HandleFunc("/test-cases", putTestCase)

	filePath := os.Getenv("GOPATH") + "/src/github.com/xtracdev/automated-perf-test/uiServices/test/"
	request, err := http.NewRequest(http.MethodPut, "/test-cases/TestCaseService", reader)
	request.Header.Set("testCasePathDir", filePath)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, request)

	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, http.StatusBadRequest, w.Code, "Should not have successfully updated")
}


func TestInvalidUrlTestCasePut(t *testing.T) {
	r := chi.NewRouter()
	r.Mount("/", GetIndexPage())

	reader := strings.NewReader(validTestCase)
	r.HandleFunc("/test-suites", putTestCase)

	filePath := os.Getenv("GOPATH") + "/src/github.com/xtracdev/automated-perf-test/uiServices/test/"
	request, err := http.NewRequest(http.MethodPut, "/test-cases/xxxxxxxxxxxzzx", reader)
	request.Header.Set("testCasePathDir", filePath)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, request)

	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, http.StatusNotFound, w.Code, "Sucessfully updated. Should not have updated")
}

func TestNoUrlTestCasePut(t *testing.T) {
	r := chi.NewRouter()
	r.Mount("/", GetIndexPage())

	reader := strings.NewReader(validTestCase)
	r.HandleFunc("/test-cases", putTestCase)

	filePath := os.Getenv("GOPATH") + "/src/github.com/xtracdev/automated-perf-test/uiServices/test/"
	request, err := http.NewRequest(http.MethodPut, "", reader)
	request.Header.Set("testCasePathDir", filePath)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, request)

	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, http.StatusNotFound, w.Code, "Sucessfully updated. Should not have updated")
}

func TestCasePutWithNoPathSlash(t *testing.T) {
	r := chi.NewRouter()
	r.Mount("/", GetIndexPage())

	reader := strings.NewReader(validTestCase)
	r.HandleFunc("/test-cases", putTestCase)

	filePath := os.Getenv("GOPATH") + "/src/github.com/xtracdev/automated-perf-test/uiServices/test"
	request, err := http.NewRequest(http.MethodPut, "/test-cases/TestCaseService", reader)
	request.Header.Set("testCasePathDir", filePath)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, request)

	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, http.StatusNoContent, w.Code, "Did not update")
}

func TestNoPathTestCasePut(t *testing.T) {
	r := chi.NewRouter()
	r.Mount("/", GetIndexPage())

	reader := strings.NewReader(validTestCase)
	r.HandleFunc("/test-cases", putTestCase)

	filePath := ""
	request, err := http.NewRequest(http.MethodPut, "/test-cases/TestCaseService", reader)
	request.Header.Set("testCasePathDir", filePath)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, request)

	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, http.StatusBadRequest, w.Code, "Successfully updated. Should not have worked due to no filepath")
}

func TestNoNameTestCasePut(t *testing.T) {
	r := chi.NewRouter()
	r.Mount("/", GetIndexPage())

	reader := strings.NewReader(TestCaseNoName)
	r.HandleFunc("/test-cases", putTestCase)

	filePath := os.Getenv("GOPATH") + "/src/github.com/xtracdev/automated-perf-test/uiServices/test"
	request, err := http.NewRequest(http.MethodPut, "/test-cases/TestCaseService", reader)
	request.Header.Set("testCasePathDir", filePath)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, request)

	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, http.StatusBadRequest, w.Code, "Successfully updated. Should not have worked due to no filepath")
}



func TestSuccessfulGetAllCases(t *testing.T) {
	r := chi.NewRouter()
	r.Mount("/", GetIndexPage())

	DirectoryPath := os.Getenv("GOPATH") + "/src/github.com/xtracdev/automated-perf-test/uiServices/test/"
	request, err := http.NewRequest(http.MethodGet, "/test-cases", nil)

	request.Header.Set("testCasePathDir", DirectoryPath)
	request.Header.Get("testCasePathDir")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, request)

	if assert.NoError(t,err) {
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

	if assert.NoError(t,err) {
		assert.Equal(t, http.StatusBadRequest, w.Code, "Did not get all test cases")
	}
}


