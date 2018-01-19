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
	"github.com/xtracdev/automated-perf-test/perfTestUtils"
	"encoding/json"
	"reflect"
)

type apiFeature struct {
	resp   *http.Response
	client *http.Client
	requestbody string
	header string
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

func (a *apiFeature) theHeaderConfigsDirPathIs(path string) error{
	a.header = path
	expectedConfigsPathDir :=os.Getenv("GOPATH") + "/src/github.com/xtracdev/automated-perf-test/uiServices/test/GodogConfig.xml"
	actualConfigsPathDir := os.Getenv("GOPATH") + "/src/github.com/xtracdev/automated-perf-test" + a.header

	if expectedConfigsPathDir != actualConfigsPathDir{
		return fmt.Errorf("Error: expected response code to be: "+expectedConfigsPathDir+" but actual is: "+actualConfigsPathDir)
	}
	return nil
}

func (a *apiFeature) theConfigFileWasCreatedAtLocationDefinedByConfigsPathDir() error {
	configsPathDir := os.Getenv("GOPATH") + "/src/github.com/xtracdev/automated-perf-test"+a.header

	fileExists := services.FilePathExist(configsPathDir)

	if (!fileExists){
		return fmt.Errorf("File Does Not Exist")
	}

	logrus.Println("File Exists")
	if !IsValidXml(a.requestbody, a.header){
		return fmt.Errorf("File is not a valid XML file")
	}
	logrus.Println("File is a valid XML file")
	return nil
}

func (a *apiFeature) iSendRequestToWithABody(method, endpoint string, body *gherkin.DocString) error {

	response, err := makeRequest(a.client, method, endpoint, body.Content)
	a.requestbody = body.Content
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
	s.Step(`^the header configsDirPath is "([^"]*)"$`, api.theHeaderConfigsDirPathIs)
	s.Step(`^the response should match json:$`, api.theResponseShouldMatchJSON)
	s.Step(`^the response body should be empty$`, api.theResponseBodyShouldBeEmpty)
	s.Step(`^the config file was created at location defined by configsPathDir$`, api.theConfigFileWasCreatedAtLocationDefinedByConfigsPathDir)
	s.Step(`^the automated performance ui server is available$`, theAutomatedPerformanceUiServerIsAvailable)
	s.Step(`^I send "([^"]*)" request to "([^"]*)" with a body:$`, api.iSendRequestToWithABody)
}

func theAutomatedPerformanceUiServerIsAvailable() error {
	go http.ListenAndServe(":9191", services.GetRouter())
	return nil
}

func  IsValidXml(config string, header string) bool{
	file, err := os.Open(os.Getenv("GOPATH") + "/src/github.com/xtracdev/automated-perf-test"+header)
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
	if err != nil{
		logrus.Println("Cannot Unmarshall into XML")
		return false
	}

	bytes, err:= []byte(config),err
	json.Unmarshal(bytes, &expectedConfig)
	if err != nil{
		logrus.Println("Cannot Unmarshall into JSON")
		return false
	}

	if !reflect.DeepEqual(&expectedConfig,&actualConfig) {
		logrus.Println("Error: Values Not Equal")
		logrus.Println("Expected :", expectedConfig," ,but actual was :", actualConfig)
		return false
	}
	return true
}
