package main

import (
"encoding/xml"
"fmt"
"os"
"io"


)

func writerXml(t Config_Struct) {

	ApiName := t.ApiName
	TargetHost := t.TargetHost
	TargetPort := t.TargetPort
	NumIterations := t.NumIterations
	AllowablePeakMemoryVariance := t.AllowablePeakMemoryVariance
	AllowableServiceResponseTimeVariance := t.AllowableServiceResponseTimeVariance
	TestCaseDir := t.TestCaseDir
	TestSuiteDir := t.TestSuiteDir
	BaseStatsOutputDir := t.BaseStatsOutputDir
	ReportOutputDir := t.ReportOutputDir
	ConcurrentUsers := t.ConcurrentUsers
	TestSuite := t.TestSuite
	MemoryEndpoint := t.MemoryEndpoint
	RequestDelay := t.RequestDelay
	TPSFreq := t.TPSFreq
	RampUsers := t.RampUsers
	RampDelay := t.RampDelay
	ReportTemplateFile := t.ReportTemplateFile

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

v.Configs = append(v.Configs, Config{ApiName:ApiName , TargetHost:TargetHost, TargetPort:TargetPort,
	NumIterations:NumIterations, AllowablePeakMemoryVariance:AllowablePeakMemoryVariance, AllowableServiceResponseTimeVariance:AllowableServiceResponseTimeVariance,
	TestCaseDir:TestCaseDir, TestSuiteDir:TestSuiteDir,
	BaseStatsOutputDir:BaseStatsOutputDir, ReportOutputDir:ReportOutputDir,ConcurrentUsers:ConcurrentUsers,
	TestSuite:TestSuite,MemoryEndpoint:MemoryEndpoint, RequestDelay:RequestDelay, TPSFreq:TPSFreq,
	RampUsers:RampUsers, RampDelay:RampDelay, ReportTemplateFile:ReportTemplateFile,
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