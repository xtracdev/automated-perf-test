package main

import (
	"errors"
	"fmt"
	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	"io/ioutil"
	"net/http"
	"io"
	"strings"
	"os"
	"github.com/Sirupsen/logrus"
	"github.com/xtracdev/automated-perf-test/uiServices/src"
	"encoding/xml"
)

type Config struct {
	APIName                              string  `xml:"apiName"`
	TargetHost                           string  `xml:"targetHost"`
	TargetPort                           string  `xml:"targetPort"`
	NumIterations                        int     `xml:"numIterations"`
	AllowablePeakMemoryVariance          float64 `xml:"allowablePeakMemoryVariance"`
	AllowableServiceResponseTimeVariance float64 `xml:"allowableServiceResponseTimeVariance"`
	TestCaseDir                          string  `xml:"testCaseDir"`
	TestSuiteDir                         string  `xml:"testSuiteDir"`
	BaseStatsOutputDir                   string  `xml:"baseStatsOutputDir"`
	ReportOutputDir                      string  `xml:"reportOutputDir"`
	ConcurrentUsers                      int     `xml:"concurrentUsers"`
	TestSuite                            string  `xml:"testSuite"`
	MemoryEndpoint                       string  `xml:"memoryEndpoint"`
	RequestDelay                         int     `xml:"requestDelay"`
	TPSFreq                              int     `xml:"TPSFreq"`
	RampUsers                            int     `xml:"rampUsers"`
	RampDelay                            int     `xml:"rampDelay"`
}

const validJsonConfig = `{
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

type apiFeature struct {
	resp   *http.Response
	client *http.Client
}

func (a *apiFeature) resetResponse() {
	a.client = &http.Client{}
	a.resp = nil
}

func (a *apiFeature) iSendrequestTo(method, endpoint string) (err error) {
	response, err := makeRequest(a.client, method, endpoint, "")
	if err != nil {
		return err
	}
	a.resp = response
	return nil
}

func (a *apiFeature) theResponseCodeShouldBe(code int) error {

	if code != a.resp.StatusCode {
		return fmt.Errorf("expected response code to be: %d, but actual is: %d", code, a.resp.StatusCode)
	}
	return nil
}

//not currently needed for POST request. Will be needed for GET request in future sprints
func (a *apiFeature) theResponseShouldMatchJSON(body *gherkin.DocString) (err error) {
	return nil
}

func (a *apiFeature) theResponseBodyShouldBeEmpty() error {
	defer a.resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(a.resp.Body)

	if err != nil {
		logrus.Error(err)
		return err
	}

	if len(bodyBytes) > 0 {
		return errors.New("Body should be empty")
	}
	return nil
}

func theHeaderConfigsDirPathIs(path string) error{
	expectedConfigsPathDir :=os.Getenv("GOPATH") + "/src/github.com/xtracdev/automated-perf-test/uiServices/test/GodogConfig.xml"
	actualConfigsPathDir := os.Getenv("GOPATH") + "/src/github.com/xtracdev/automated-perf-test" + path

	if expectedConfigsPathDir != actualConfigsPathDir{
		return fmt.Errorf("Error: expected response code to be: "+expectedConfigsPathDir+" but actual is: "+actualConfigsPathDir)
	}
	return nil
}

func (a *apiFeature) theConfigFileWasCreatedAtLocationDefinedByConfigsPathDir() error {
	configsPathDir := os.Getenv("GOPATH") + "/src/github.com/xtracdev/automated-perf-test/uiServices/test/GodogConfig.xml"

	fileExists := services.FilePathExist(configsPathDir)

	if (!fileExists){
		return fmt.Errorf("File Does Not Exist")
	}

	logrus.Println("File Exists")
	checkXmlValidity()
	return nil
}

func (a *apiFeature) iSendRequestToWithABody(method, endpoint string) error {

	response, err := makeRequest(a.client, method, endpoint, validJsonConfig)
	if err != nil {
		return err
	}
	a.resp = response
	return nil
}

func makeRequest(client *http.Client, method, endpoint, body string) (*http.Response, error) {

	var reqBody io.Reader
	if body != "" {
		reqBody = strings.NewReader(body)
	}

	req, err := http.NewRequest(method, "http://localhost:9191"+endpoint, reqBody)
	req.Header.Set("configPathDir", fmt.Sprintf("%s/src/github.com/xtracdev/automated-perf-test/uiServices/test/",os.Getenv("GOPATH")))

	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func FeatureContext(s *godog.Suite) {
	api := &apiFeature{}

	s.BeforeScenario(func(interface{}) {

		api.resetResponse()

	})

	s.Step(`^I send "(GET|POST|PUT|DELETE)" request to "([^"]*)"$`, api.iSendrequestTo)
	s.Step(`^the response code should be (\d+)$`, api.theResponseCodeShouldBe)
	s.Step(`^the header configsDirPath is "([^"]*)"$`, theHeaderConfigsDirPathIs)
	s.Step(`^the response should match json:$`, api.theResponseShouldMatchJSON)
	s.Step(`^the response body should be empty$`, api.theResponseBodyShouldBeEmpty)
	s.Step(`^the config file was created at location defined by configsPathDir$`, api.theConfigFileWasCreatedAtLocationDefinedByConfigsPathDir)
	s.Step(`^the automated performance ui server is available$`, theAutomatedPerformanceUiServerIsAvailable)
	s.Step(`^I send "([^"]*)" request to "([^"]*)" with a body$`, api.iSendRequestToWithABody)
}

func theAutomatedPerformanceUiServerIsAvailable() error {
	go http.ListenAndServe(":9191", services.GetRouter())
	return nil
}

func checkXmlValidity() error{
	//Open xmlFile
	xmlFile, err := os.Open(os.Getenv("GOPATH") + "/src/github.com/xtracdev/automated-perf-test/uiServices/test/GodogConfig.xml")

	// handle error if can't open file
	if err != nil {
		fmt.Println(err)
		return err
	}

	logrus.Println("Successfully Opened XML file")

	byteValue, _ := ioutil.ReadAll(xmlFile)
	var config Config

	// unmarshal our byteArray into 'config' struct
	xml.Unmarshal(byteValue, &config)

		logrus.Println("File is valid XML file.")
		fmt.Println("<ApiName>", config.APIName, "</ApiName>")
		fmt.Println("<TargetHost>", config.TargetHost, "</TargetHost>")
		fmt.Println("<TargetPort>", config.TargetPort, "</TargetPort>")
		fmt.Println("<NumIterations>", config.NumIterations, "</NumIterations>")
		fmt.Println("<AllowablePeakMemoryVariance>", config.AllowablePeakMemoryVariance, "</AllowablePeakMemoryVariance>")
		fmt.Println("<AllowableServiceResTimeVariance>", config.AllowableServiceResponseTimeVariance, "</AllowableServiceResTimeVariance>")
		fmt.Println("<TestCaseDir>", config.TestCaseDir, "</TestCaseDir>")
		fmt.Println("<TestSuiteDir>", config.TestSuiteDir, "</TestSuiteDir>")
		fmt.Println("<BaseStatsOutputDir>", config.BaseStatsOutputDir, "</BaseStatsOutputDir>")
		fmt.Println("<ReportOutputDir>", config.ReportOutputDir, "</ReportOutputDir>")
		fmt.Println("<ConcurrentUsers>", config.ConcurrentUsers, "</ConcurrentUsers>")
		fmt.Println("<TestSuite>", config.TestSuite, "</TestSuite>")
		fmt.Println("<MemoryEndpoint>", config.MemoryEndpoint, "</MemoryEndpoint>")
		fmt.Println("<RequestDelay>", config.RequestDelay, "</RequestDelay>")
		fmt.Println("<TPSFreq>", config.TPSFreq, "</TPSFreq>")
		fmt.Println("<RampUsers>", config.RampUsers, "</RampUsers>")
		fmt.Println("<RampDelay>", config.RampDelay, "</RampDelay>")

	xmlFile.Close()

	return nil
}
