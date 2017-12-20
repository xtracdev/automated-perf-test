package services

import (
	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
	"github.com/xtracdev/automated-perf-test/perfTestUtils"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

var (
	reader io.Reader
)
const validJson = `{
        "apiName": "TestConfig",
       "targetHost": "localhost",
       "targetPort": "9191",
       "numIterations": 1000,
       "allowablePeakMemoryVariance": 30,
       "allowableServiceResponseTimeVariance": 30,
       "testCaseDir": "./definitions/testCases",
       "testSuiteDir": "./definitions/testSuites",
        "baseStatsOutputDir": "./envStats",
       "reportOutputDir": "./report",
       "concurrentUsers": 50,
       "testSuite": "suiteFileName.xml",
       "memoryEndpoint": "/alt/debug/vars",
       "requestDelay": 5000,
       "TPSFreq": 30,
       "rampUsers": 5,
       "rampDelay": 15
       }`

const invalidJson = `{
        "apiName"://*()()(),
       "targetHost": 0,
       "targetPort": 0,
       "numIterations": "x",
       "allowablePeakMemoryVariance": 30,
       "allowableServiceResponseTimeVariance": 30,
       "testCaseDir": "./definitions/testCases",
       "testSuiteDir": "./definitions/testSuites",
        "baseStatsOutputDir": "./envStats",
       "reportOutputDir": "./report",
       "concurrentUsers": 50,
       "testSuite": "suiteFileName.xml",
       "memoryEndpoint": "/alt/debug/vars",
       "requestDelay": 5000,
       "TPSFreq": 30,
       "rampUsers": 5,
       "rampDelay": 15
       }`

func TestFilePathExist(t *testing.T) {
	path := os.Getenv("GOPATH") + "/src/github.com/xtracdev/automated-perf-test/config/"
	actual := false
	actual = FilePathExist(path)
	expected := true
	assert.Equal(t, expected, actual)
}

func TestFilePathDoesNotExist(t *testing.T) {
	path := "C:/github.com/xxx"
	actual := FilePathExist(path)
	expected := false
	assert.Equal(t, expected, actual)
}

func TestValidJsonPost(t *testing.T) {
	r := chi.NewRouter()
	r.Mount("/", getIndexPage())

	reader = strings.NewReader(validJson)
	r.HandleFunc("/configs", configsHandler)

	filePath := os.Getenv("GOPATH") + "/src/github.com/xtracdev/automated-perf-test/config/"

	request, err := http.NewRequest(http.MethodPost, "/configs", reader)
	request.Header.Set("configPathDir", filePath)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, request)

	if err != nil {
		t.Error(err)
	}

	if w.Code != http.StatusCreated  {
		t.Errorf("TestValidJsonPost. Expected:",http.StatusCreated, " Got:",w.Code,"  Error. Did not succesfully post",)
	}
}

func TestInvalidJsonPost(t *testing.T) {
	r := chi.NewRouter()
	r.Mount("/", getIndexPage())

	reader = strings.NewReader(invalidJson)
	r.HandleFunc("/configs", configsHandler)

	request, err := http.NewRequest(http.MethodPost, "/configs", reader)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, request)

	if err != nil {
		t.Error(err)
	}

	if w.Code != http.StatusBadRequest {
		t.Errorf("TestInvalidJsonPost.  Expected:",http.StatusBadRequest," Got:",w.Code,  "Created XML")
	}
}

func TestPostWithNoFilePath(t *testing.T) {
	r := chi.NewRouter()
	r.Mount("/", getIndexPage())

	reader = strings.NewReader(validJson)
	r.HandleFunc("/configs", configsHandler)

	request, err := http.NewRequest(http.MethodPost, "/configs", reader)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, request)

	if err != nil {
		t.Error(err)
	}
}

func TestInvalidURL(t *testing.T) {
	pt := perfTestUtils.Config{}
	writerXml(pt, "/path/xxx")
}