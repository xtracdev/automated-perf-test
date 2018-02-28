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