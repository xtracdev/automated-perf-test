package main

import (
"encoding/json"
"net/http"

)

type Config_Struct struct {
	ApiName string
	TargetHost string
	TargetPort int
	NumIterations int
	AllowablePeakMemoryVariance float64
	AllowableServiceResponseTimeVariance float64
	TestCaseDir string
	TestSuiteDir string
	BaseStatsOutputDir string
	ReportOutputDir string
	ConcurrentUsers int
	TestSuite string
	MemoryEndpoint string
	RequestDelay int
	TPSFreq int
	RampUsers int
	RampDelay int
	ReportTemplateFile string
}

func test(rw http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var t Config_Struct
	err := decoder.Decode(&t)
	if err != nil {
		panic(err)
	}
	writerXml(t)
	defer req.Body.Close()

}
