package services

import (
	"net/http/httptest"
	"strings"
	"github.com/go-chi/chi"
	"testing"
	"net/http"
	"github.com/xtracdev/automated-perf-test/perfTestUtils"
)


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

	w := httptest.NewRecorder()
	r.ServeHTTP(w, request)

	if err != nil {
		t.Error(err)
	}

	if http.StatusCreated != 201 {
		t.Errorf("Error. Did not succesfully post")
	}

}

func TestInvalidURL(t *testing.T) {
	pt := perfTestUtils.Config{}
	writerXml(pt, "/path/xxx")
}