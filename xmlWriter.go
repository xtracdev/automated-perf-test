package main

import (
"encoding/xml"
"fmt"
"os"
"io"
)

func main() {

type Config struct {
XMLName xml.Name `xml:"config"`
ApiName string `xml:"apiName"`
TargetHost string `xml:"targetHost"`
TargetPort int `xml:"targetPort"`
NumIterations int `xml:"numIterations"`
AllowablePeakMemoryVariance float64 `xml:"allowablePeakMemoryVariance"`
AllowableServiceResponseTimeVariance float64 `xml:"allowableServiceResponseTimeVariance"`
TestCaseDir string `xml:"testCaseDir"`
TestSuiteDir string `xml:"testSuiteDir"`
BaseStatsOutputDir string `xml:"baseStatsOutputDir"`
ReportOutputDir string `xml:"reportOutputDir"`
ConcurrentUsers int `xml:"concurrentUsers"`
TestSuite string `xml:"testSuite"`
MemoryEndpoint string `xml:"memoryEndpoint"`
RequestDelay int `xml:"requestDelay"`
TPSFreq int `xml:"TPSFreq"`
RampUsers int `xml:"rampUsers"`
RampDelay int `xml:"rampDelay"`
ReportTemplateFile string `xml:"reportTemplateFile"`
}

type ConfigFiles struct {
Configs []Config
}


v := &ConfigFiles{}

v.Configs = append(v.Configs, Config{ApiName: "Xtrac API", TargetHost:"localhost", TargetPort:9191,
	NumIterations:1000, AllowablePeakMemoryVariance:15, AllowableServiceResponseTimeVariance:15,
	TestCaseDir:"./definitions/testCases", TestSuiteDir:"./definitions/testSuites",
	BaseStatsOutputDir:"./envStats", ReportOutputDir:"./report",ConcurrentUsers:50,
	TestSuite:"suiteFileName.xml",MemoryEndpoint:"/alt/debug/vars", RequestDelay:5000, TPSFreq:30,
	RampUsers:5, RampDelay:15, ReportTemplateFile:"???",
})

filename := "ConfigFile.xml"
file, _ := os.Create(filename)

xmlWriter := io.Writer(file)

enc := xml.NewEncoder(xmlWriter)
enc.Indent("  ", "    ")
if err := enc.Encode(v); err != nil {
fmt.Printf("error: %v\n", err)
}

}