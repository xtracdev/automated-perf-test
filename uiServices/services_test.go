package services

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"os"
	"net/http/httptest"
	"strings"
	"github.com/go-chi/chi"
	"io"
	"net/http"
)



func TestFilePathExist(t *testing.T) {
	path := os.Getenv("GOPATH")
	actual := false
	actual = FilePathExist(path)
	expected := true
	os.Getenv("GOPATH")
	assert.Equal (t, expected, actual)
}

func TestFilePathDoesNotExist(t *testing.T) {
	path := "C:/Users/a615194/go/src/github.com/xtrac"
	actual := FilePathExist(path)
	expected := false
	assert.Equal (t, expected, actual)
}


var (
	reader  io.Reader
)


func TestValidJsonPost(t *testing.T) {
	//start server
	r := chi.NewRouter()
	r.Mount("/", getIndexPage())

	//create JSON
	configJson := `{
        "apiName": "TestConfigFile",
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
	//reader
	reader= strings.NewReader(configJson)
	//prepare POST method
	r.HandleFunc("/configs", configsHandler)
	//perfrom POTS method using URI and read JSON
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

func TestInvalidJsonPost(t *testing.T) {
	//start server
	r := chi.NewRouter()
	r.Mount("/", getIndexPage())

	//create JSON
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
	//reader
	reader= strings.NewReader(configJson)
	//prepare POST method
	r.HandleFunc("/configs", configsHandler)
	//perfrom POTS method using URI and read JSON
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


//
//func TestFileSave(t *testing.T) {
//
//	path :=  "C:/Users/a615194/go/src/github.com/xtracdev/automated-perf-test/config/ConfigTest.xml"
//	actual := FilePathExist(path)
//	expected := false
//	assert.Equal (t, expected, actual)
//}
