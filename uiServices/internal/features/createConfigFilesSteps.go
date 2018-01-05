package uiServices

import (
. "github.com/gucumber/gucumber"
"github.com/xtracdev/automated-perf-test/uiServices"
	"net/http"
	"strings"
	"os"

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


var httpClient *http.Client

func init() {

	Given(`^the automated performance ui server is available`, func() {
		services.StartUiMode()
	})

	When(`^the user makes a request for (.+?) http://localhost:9191/configs with payload$`, func() {
		//get payload (JSON)
		reader := strings.NewReader(validJson)
		//path to save file
		filePath := os.Getenv("GOPATH") + "/src/github.com/xtracdev/automated-perf-test/config/"
		//set URI and header
		request, _ := http.NewRequest(http.MethodPost, "/configs", reader)
		request.Header.Set("configPathDir", filePath)
		//run request
		httpClient.Do(request)
	})

	Then(`^the POST configuration service returns 201 HTTP status`, func() {
		//check the response
	})
}
