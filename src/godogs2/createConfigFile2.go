package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"github.com/go-chi/chi"
	"github.com/xtracdev/automated-perf-test/uiServices"
	"strings"
	"os"
	"net/http/httptest"
)
// valid json file
const validJson = `{
        "apiName": "GodogConfig",
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


func main() {
	// start server
	services.StartUiMode()
}

func createFile(w http.ResponseWriter, r *http.Request) {
	// ensure only POST methods are accepted
	if r.Method != "POST" {
		// handle any other method calls
		HandleInvalidHttpRequest(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// valid config file object/structure and initial values
	data := struct {
		TPSFreq                              int     `xml:"TPSFreq" json:"TPSFreq"`
		AllowablePeakMemoryVariance          float64 `xml:"allowablePeakMemoryVariance" json:"allowablePeakMemoryVariance"`
		AllowableServiceResponseTimeVariance float64 `xml:"allowableServiceResponseTimeVariance" json:"allowableServiceResponseTimeVariance"`
		APIName                              string  `xml:"apiName" json:"apiName"`
		BaseStatsOutputDir                   string  `xml:"baseStatsOutputDir" json:"baseStatsOutputDir"`
		ConcurrentUsers                      int     `xml:"concurrentUsers" json:"concurrentUsers"`
		MemoryEndpoint                       string  `xml:"memoryEndpoint" json:"memoryEndpoint"`
		NumIterations                        int     `xml:"numIterations" json:"numIterations"`
		RampDelay                            int     `xml:"rampDelay" json:"rampDelay"`
		RampUsers                            int     `xml:"rampUsers" json:"rampUsers"`
		ReportOutputDir                      string  `xml:"reportOutputDir" json:"reportOutputDir"`
		RequestDelay                         int     `xml:"requestDelay" json:"requestDelay"`
		TargetHost                           string  `xml:"targetHost" json:"targetHost"`
		TargetPort                           string  `xml:"targetPort" json:"targetPort"`
		TestCaseDir                          string  `xml:"testCaseDir" json:"testCaseDir"`
		TestSuite                            string  `xml:"testSuite" json:"testSuite"`
		TestSuiteDir                         string  `xml:"testSuiteDir" json:"testSuiteDir"`
	}{30, 30,30,"GodogConfig","./envStats",50, "/alt/debug/vars",1000,15, 5,"./report",5000,"localhost","9191","./definitions/testCases","suiteFileName.xml", "./definitions/testSuites"}

	// method to create file if data is valid
	createNewFile(w, data)
}


// fail writes a json response with error msg and status header
func HandleInvalidHttpRequest(w http.ResponseWriter, msg string, status int) {

	// json object to be output if http method not POST
	data := struct {
		Error string `json:"error"`
	}{Error: msg}
	// marshall json
	resp, _ := json.Marshal(data)
	w.WriteHeader(status)
	// print and compare json file to json in feature file
	fmt.Fprintf(w, string(resp))
}

// writes data to response with status
func createNewFile(w http.ResponseWriter, data interface{}) {
	// router
	r:= chi.NewRouter()
	r.Mount("/", services.GetIndexPage())

	// read json
	reader := strings.NewReader(validJson)
	r.HandleFunc("/configs", services.ConfigsHandler)

	// create filepath
	filePath := os.Getenv("GOPATH") + "/src/github.com/xtracdev/automated-perf-test/config/"

	// append endpoint for POST method
	request, err := http.NewRequest(http.MethodPost, "/configs", reader)

	// set header
	request.Header.Set("configPathDir", filePath)

	// marshall data
	res, err := json.Marshal(data)
	//error handling
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		HandleInvalidHttpRequest(w, "Internal Server Error", 500)
		return
	}

	// print marshalled json to compare it to feature file
	fmt.Fprintf(w, string(res))

	//carry out command
	w = httptest.NewRecorder()
	r.ServeHTTP(w, request)
	// error handling
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		HandleInvalidHttpRequest(w, "Internal Server Error", 500)
		return
	}
}