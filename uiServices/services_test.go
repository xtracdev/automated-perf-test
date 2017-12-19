package services

import (
	"testing"
	"os"
	"net/http/httptest"
	"strings"
	"github.com/go-chi/chi"
	"io"
	"net/http"
	"github.com/Sirupsen/logrus"
	"github.com/xtracdev/automated-perf-test/perfTestUtils"
	"github.com/stretchr/testify/assert"
)


var (
	reader  io.Reader
)
func TestFilePathExist(t *testing.T) {
	path :=os.Getenv("GOPATH")+"/src/github.com/xtracdev/automated-perf-test/config/"
	actual := false
	actual = FilePathExist(path)
	expected := true
	assert.Equal (t, expected, actual)
}

func TestFilePathDoesNotExist(t *testing.T) {
	path := "C:/github.com/xxx"
	actual := FilePathExist(path)
	expected := false
	assert.Equal (t, expected, actual)
}

func TestValidJsonPost(t *testing.T) {
	r := chi.NewRouter()
	r.Mount("/", getIndexPage())
	
	configJson := `{
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

	reader= strings.NewReader(configJson)
	r.HandleFunc("/configs", configsHandler)

	filePath :=os.Getenv("GOPATH")+"/src/github.com/xtracdev/automated-perf-test/config/"

	request, err := http.NewRequest(http.MethodPost, "/configs", reader)
	request.Header.Set("configPathDir", filePath)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, request)

	if err != nil {
		t.Error(err)
	}

	if http.StatusCreated != 201 {
		t.Errorf("Error. Did not succesfully post")
	}else{
	logrus.Print("Successfully created XML File")
	}
}

func TestInvalidJsonPost(t *testing.T) {
	r := chi.NewRouter()
	r.Mount("/", getIndexPage())
	
	configJson := `{
        "apiName": 0,
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
	
	reader= strings.NewReader(configJson)
	r.HandleFunc("/configs", configsHandler)

	request, err := http.NewRequest(http.MethodPost, "/configs", reader)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, request)

	if err != nil {
		t.Error(err)
	}

	if http.StatusCreated != 201 {
		t.Errorf("Error. Did not succesfully post")
	}
}

func TestPostWithNoFilePath(t *testing.T) {
	r := chi.NewRouter()
	r.Mount("/", getIndexPage())

	configJson := `{
        "apiName": "Test",
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

	reader= strings.NewReader(configJson)

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

func TestInvalidFileName(t *testing.T) {
	r := chi.NewRouter()
	r.Mount("/", getIndexPage())

	configJson := `{
        "apiName": "//*()()()",
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

	reader= strings.NewReader(configJson)

	r.HandleFunc("/configs", configsHandler)
	request, err := http.NewRequest(http.MethodPost, "/configs", reader)

	filePath :=os.Getenv("GOPATH")+"/src/github.com/xtracdev/automated-perf-test/config/"
	request.Header.Set("configPathDir", filePath)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, request)

	if err != nil {
		t.Error(err)
	}

	if http.StatusCreated != 201 {
		t.Errorf("Error. Did not succesfully post")
	}

}