package main

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	"github.com/Sirupsen/logrus"
	"github.com/xtracdev/automated-perf-test/perfTestUtils"
	"github.com/xtracdev/automated-perf-test/uiServices/src"
)

type apiFeature struct {
	res         *httptest.ResponseRecorder
	resp        *http.Response
	client      *http.Client
	requestbody string
	headerPath  string
	filename    string
	headerName  string
}

func (a *apiFeature) resetResponse() {
	a.client = &http.Client{}
	a.resp = nil
}

func (a *apiFeature) iSendrequestTo(method, endpoint string) (err error) {
	response, err := makePostRequest(a.client, method, endpoint, "", "", "")
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

func (a *apiFeature) theResponseBodyShouldMatchJSON(body *gherkin.DocString) (err error) {
	var expectedConfig perfTestUtils.Config
	var actualConfig perfTestUtils.Config

	expectedJson := `"""
       {
              "apiName": "GodogConfig",
       "targetHost": "localhost",
       "targetPort":"9191",
       "memoryEndpoint": "/alt/debug/vars",
       "numIterations": 1000,
       "allowablePeakMemoryVariance": 30,
       "allowableServiceResponseTimeVariance": 30,
       "testCaseDir": "./definitions/testCases",
       "testSuiteDir": "./definitions/testSuites",
       "baseStatsOutputDir": "./envStats",
       "reportOutputDir": "./report",
       "concurrentUsers": 50,
       "testSuite": "Default-3",
       "requestDelay": 5000,
       "TPSFreq": 30,
       "rampUsers": 5,
       "rampDelay": 15
       }
       """`

	json.Unmarshal([]byte(body.Content), &actualConfig)
	json.Unmarshal([]byte(expectedJson), &expectedConfig)

	if !reflect.DeepEqual(&expectedConfig, &actualConfig) {
		return fmt.Errorf("Expected :%v ,but actually was :%v", expectedConfig, actualConfig)
	}

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

func (a *apiFeature) theHeaderIs(name string, path string) error {
	a.headerPath = path
	a.headerName = name

	if path == "" {
		fmt.Println("Error: No Header Defined")
		return nil
	}
	fmt.Println(path)
	return nil
}

func (a *apiFeature) theConfigFileWasCreatedAtLocationDefinedByConfigsPathDir() error {
	configsPathDir := os.Getenv("GOPATH") + "/src/github.com/xtracdev/automated-perf-test" + a.headerPath

	fileExists := services.FilePathExist(configsPathDir)

	if !fileExists {
		return fmt.Errorf("File Does Not Exist")
	}

	logrus.Println("File Exists")
	if !IsValidXml(a.requestbody, a.headerPath) {
		return fmt.Errorf("File is not a valid XML file")
	}
	logrus.Println("File is a valid XML file")
	return nil
}

func (a *apiFeature) iSendRequestToWithABody(method, endpoint string, body *gherkin.DocString) error {
	response, err := makePostRequest(a.client, method, endpoint, body.Content, a.headerPath, a.headerName)
	a.requestbody = body.Content
	if err != nil {
		return err
	}
	a.resp = response
	return nil
}

func makePostRequest(client *http.Client, method, endpoint, body string, headerPath string, headerName string) (*http.Response, error) {

	var reqBody io.Reader
	if body != "" {
		reqBody = strings.NewReader(body)
	}

	req, err := http.NewRequest(method, "http://localhost:9191"+endpoint, reqBody)

	if headerPath == "" {
		req.Header.Set(headerName, "")
	} else {
		req.Header.Set(headerName, fmt.Sprintf("%s/src/github.com/xtracdev/automated-perf-test/uiServices/test/", os.Getenv("GOPATH")))
	}
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
	s.Step(`^the header "([^"]*)" is "([^"]*)"$`, api.theHeaderIs)
	s.Step(`^the response body should match json:$`, api.theResponseBodyShouldMatchJSON)
	s.Step(`^the response body should be empty$`, api.theResponseBodyShouldBeEmpty)
	s.Step(`^the config file was created at location defined by configsPathDir$`, api.theConfigFileWasCreatedAtLocationDefinedByConfigsPathDir)
	s.Step(`^the automated performance ui server is available$`, theAutomatedPerformanceUiServerIsAvailable)
	s.Step(`^I send "([^"]*)" request to "([^"]*)" with a body:$`, api.iSendRequestToWithABody)
	s.Step(`^the file name is "([^"]*)"$`, api.theFileNameis)
	s.Step(`^I send a "([^"]*)" request to "([^"]*)"$`, api.iSendARequestTo)
	s.Step(`^the file "([^"]*)" exists at "([^"]*)"$`, theFileExistsAt)
	s.Step(`^I send "([^"]*)" request to "([^"]*)" with body:$`, api.iSendRequestToWithBody)
	s.Step(`^the updated file should match json:$`, api.theUpdatedFileShouldMatchJSON)
	s.Step(`^there is no existing test file "([^"]*)"$`, api.thereIsNoExistingTestFile)
}

func (a *apiFeature) thereIsNoExistingTestFile(file string) error {
	os.Remove(os.Getenv("GOPATH") + "/src/github.com/xtracdev/automated-perf-test/uiServices/test/" + file)
	return nil
}

func theAutomatedPerformanceUiServerIsAvailable() error {
	go http.ListenAndServe(":9191", services.GetRouter())
	return nil
}

func IsValidXml(config string, header string) bool {
	file, err := os.Open(os.Getenv("GOPATH") + "/src/github.com/xtracdev/automated-perf-test" + header)
	if err != nil {
		fmt.Println(err)
		return false
	}

	defer file.Close()
	logrus.Println("Successfully Opened XML file")

	var actualConfig perfTestUtils.Config
	var expectedConfig perfTestUtils.Config

	byteValue, err := ioutil.ReadAll(file)
	xml.Unmarshal(byteValue, &actualConfig)
	if err != nil {
		logrus.Println("Cannot Unmarshall into XML")
		return false
	}

	bytes, err := []byte(config), err
	json.Unmarshal(bytes, &expectedConfig)
	if err != nil {
		logrus.Println("Cannot Unmarshall into JSON")
		return false
	}

	if !reflect.DeepEqual(&expectedConfig, &actualConfig) {
		logrus.Println("Error: Values Not Equal")
		logrus.Println("Expected :", expectedConfig, " ,but actual was :", actualConfig)
		return false
	}
	return true
}

func (a *apiFeature) theFileNameis(filename string) error {
	a.filename = filename
	return nil
}

func (a *apiFeature) iSendARequestTo(method, endpoint string) error {
	response, err := makeGetRequest(a.client, method, endpoint, a.filename, a.headerPath, a.headerName)
	if err != nil {
		return err
	}

	a.resp = response
	return nil
}

func makeGetRequest(client *http.Client, method, endpoint string, filename string, headerPath, headerName string) (*http.Response, error) {
	req, err := http.NewRequest(method, "http://localhost:9191"+endpoint, nil)

	if headerPath == "" {
		req.Header.Set(headerName, "")
	} else {
		req.Header.Set(headerName, fmt.Sprintf("%s/src/github.com/xtracdev/automated-perf-test/uiServices/test/", os.Getenv("GOPATH")))
	}
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func theFileExistsAt(filename, path string) error {
	_, err := os.Stat(os.Getenv("GOPATH") + "/src/github.com/xtracdev/automated-perf-test" + path + filename)
	if err != nil {
		fmt.Println("Error. File Not Found at location : " + path + filename)
		return nil
	}
	return nil
}

func (a *apiFeature) iSendRequestToWithBody(method, endpoint string, body *gherkin.DocString) error {
	response, err := makePutRequest(a.client, method, endpoint, body.Content, a.headerPath, a.headerName)
	a.requestbody = body.Content
	if err != nil {
		return err
	}
	a.resp = response
	return nil
}

func makePutRequest(client *http.Client, method, endpoint, body string, headerPath, headerName string) (*http.Response, error) {

	var reqBody io.Reader
	if body != "" {
		reqBody = strings.NewReader(body)
	}

	req, err := http.NewRequest(method, "http://localhost:9191"+endpoint, reqBody)

	if headerPath == "" {
		req.Header.Set(headerName, "")
	} else {
		req.Header.Set(headerName, fmt.Sprintf("%s/src/github.com/xtracdev/automated-perf-test/uiServices/test/", os.Getenv("GOPATH")))
	}
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (a *apiFeature) theUpdatedFileShouldMatchJSON(body *gherkin.DocString) (err error) {
	var expectedConfig perfTestUtils.Config
	var actualConfig perfTestUtils.Config

	expectedJson := `"""
       {
       "apiName": "GodogConfig",
       "targetHost": "localhost2",
       "targetPort":"1001",
       "memoryEndpoint": "/alt/debug/vars",
       "numIterations": 4000,
       "allowablePeakMemoryVariance": 50,
       "allowableServiceResponseTimeVariance": 50,
       "testCaseDir": "./definitions/testCases",
       "testSuiteDir": "./definitions/testSuites",
       "baseStatsOutputDir": "./envStats",
       "reportOutputDir": "./report",
       "concurrentUsers": 50,
       "testSuite": "Default-3",
       "requestDelay": 1000,
       "TPSFreq": 10,
       "rampUsers": 10,
       "rampDelay": 10
       }
       """`

	json.Unmarshal([]byte(body.Content), &actualConfig)
	json.Unmarshal([]byte(expectedJson), &expectedConfig)

	if !reflect.DeepEqual(&expectedConfig, &actualConfig) {
		return fmt.Errorf("Expected : %v ,but actually was :%v", expectedConfig, actualConfig)

	}

	return nil
}
