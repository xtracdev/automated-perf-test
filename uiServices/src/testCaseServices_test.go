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

func TestSuccessfulGetTestCase(t *testing.T) {
	r := chi.NewRouter()
	r.Mount("/", GetIndexPage())

	filePath := os.Getenv("GOPATH") + "/src/github.com/xtracdev/automated-perf-test/uiServices/test/"
	request, err := http.NewRequest(http.MethodGet, "/test-cases/TestCaseService", nil)

	request.Header.Set("testCasePathDir", filePath)
	request.Header.Get("testCasePathDir")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, request)

	if err != nil {
		t.Error(err)
	}

	assert.Equal(t,http.StatusOK, w.Code, "Error. Did not successfully GET")
}

func TestGetTestCaseNoHeader(t *testing.T) {
	r := chi.NewRouter()
	r.Mount("/", GetIndexPage())

	filePath := ""
	request, err := http.NewRequest(http.MethodGet, "/test-cases/TestCaseService", nil)

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

	filePath := os.Getenv("GOPATH") + "/src/github.com/xtracdev/automated-perf-test/uiServices/test/cases"
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

func TestDeleteAllCasesSuccess(t *testing.T) {
	r := chi.NewRouter()
	r.Mount("/", GetIndexPage())

	DirectoryPath := os.Getenv("GOPATH") + "/src/github.com/xtracdev/automated-perf-test/uiServices/test/cases"
	err := ioutil.WriteFile(fmt.Sprintf("%s%s.xml", DirectoryPath, "test"), nil, 0666)
	if err != nil {
		logrus.Errorf("Error trying to create a file: %s", err)
	}

	request, err := http.NewRequest(http.MethodDelete, "/test-cases/all", nil)
	if err != nil {
		logrus.Warnf("Error creating the request %s", err)
	}
	request.Header.Set("testCasePathDir", DirectoryPath)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, request)

	assert.Equal(t, http.StatusOK, w.Code, "All files have been DELETED")

}

func TestDeleteAllCasesNoHeader(t *testing.T) {
	r := chi.NewRouter()
	r.Mount("/", GetIndexPage())

	DirectoryPath := ""
	request, err := http.NewRequest(http.MethodDelete, "/test-cases/all", nil)
	if err != nil {
		logrus.Warnf("Error creating the request %s", err)
	}

	request.Header.Set("testCasePathDir", DirectoryPath)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, request)

	assert.Equal(t, http.StatusBadRequest, w.Code, "Did not DELETE the files")
}

func TestDeleteAllCasesEmptyDirectory(t *testing.T) {
	r := chi.NewRouter()
	r.Mount("/", GetIndexPage())

	DirectoryPath := os.Getenv("GOPATH") + "/src/github.com/xtracdev/automated-perf-test/uiServices/test/cases"
	request, err := http.NewRequest(http.MethodDelete, "/test-cases/", nil)
	if err != nil {
		logrus.Warnf("Error creating the request %s", err)
	}

	request.Header.Set("testCasePathDir", DirectoryPath)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, request)

	assert.Equal(t, http.StatusNotFound, w.Code, "Empty Directory")
}
