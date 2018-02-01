package services

import (
	"fmt"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
	"github.com/xtracdev/automated-perf-test/perfTestUtils"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

const validJson = `{
        "apiName": "ServiceTestConfig",
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
       "testSuite": "Default-1",
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
	path := os.Getenv("GOPATH") + "/src/github.com/xtracdev/automated-perf-test/uiServices/test/"
	actual := false
	fmt.Println(path)
	actual = FilePathExist(path)
	expected := true
	assert.Equal(t, expected, actual)
}

func TestFilePathDoesNotExist(t *testing.T) {
	path := "((((((("
	actual := FilePathExist(path)
	expected := false
	assert.Equal(t, expected, actual)
}

func TestValidJsonPost(t *testing.T) {
	r := chi.NewRouter()
	r.Mount("/", GetIndexPage())

	reader := strings.NewReader(validJson)
	r.HandleFunc("/configs", postConfigs)

	filePath := os.Getenv("GOPATH") + "/src/github.com/xtracdev/automated-perf-test/uiServices/test/"
	request, err := http.NewRequest(http.MethodPost, "/configs", reader)
	request.Header.Set("configPathDir", filePath)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, request)

	if err != nil {
		t.Error(err)
	}

	if w.Code != http.StatusCreated {
		t.Errorf("TestValidJsonPost. Expected:", http.StatusCreated, " Got:", w.Code, "  Error. Did not succesfully post")
	}
}
func TestFilePathEndsWIthSlash(t *testing.T) {
	r := chi.NewRouter()
	r.Mount("/", GetIndexPage())

	reader := strings.NewReader(validJson)
	r.HandleFunc("/configs", postConfigs)

	filePath := os.Getenv("GOPATH") + "/src/github.com/xtracdev/automated-perf-test/uiServices/test"

	request, err := http.NewRequest(http.MethodPost, "/configs", reader)
	request.Header.Set("configPathDir", filePath)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, request)

	if err != nil {
		t.Error(err)
	}

	if w.Code != http.StatusCreated {
		t.Errorf("TestFilePathEndsWith'/'.  Expected:", http.StatusCreated, " Got:", w.Code, "  Error. Did not succesfully post")
	}
}

func TestInvalidJsonPost(t *testing.T) {
	r := chi.NewRouter()
	r.Mount("/", GetIndexPage())

	reader := strings.NewReader(invalidJson)
	r.HandleFunc("/configs", postConfigs)

	request, err := http.NewRequest(http.MethodPost, "/configs", reader)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, request)

	if err != nil {
		t.Error(err)
	}

	if w.Code != http.StatusBadRequest {
		t.Errorf("TestInvalidJsonPost.  Expected:", http.StatusBadRequest, " Got:", w.Code, "Error. Did not succesfully post ")
	}
}

func TestWhenConfigPathDirEmpty(t *testing.T) {
	r := chi.NewRouter()
	r.Mount("/", GetIndexPage())

	reader := strings.NewReader(validJson)
	r.HandleFunc("/configs", postConfigs)

	request, err := http.NewRequest(http.MethodPost, "/configs", reader)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, request)

	if err != nil {
		t.Error(err)
	}

	if w.Code != http.StatusBadRequest {
		t.Errorf("TestWhenConfigPathDirEmpty.  Expected:", http.StatusBadRequest, " Got:", w.Code, "Error. ConfigPathDir is Empty ")
	}
}

func TestInvalidURL(t *testing.T) {
	pt := perfTestUtils.Config{}
	writerXml(pt, "/path/xxx")
}
