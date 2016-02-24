package testStrategies

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"github.com/BurntSushi/toml"
	log "github.com/Sirupsen/logrus"
	"github.com/xtracdev/automated-perf-test/perfTestUtils"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

const (
	SERVICE_BASED_TESTING = "ServiceBased"
	SUITE_BASED_TESTING   = "SuiteBased"
)

var globals map[string]map[string]string

func init() {
	//Initilize globals map
	globals = make(map[string]map[string]string)

}

type Header struct {
	Value string `xml:",chardata"`
	Key   string `xml:"key,attr"`
}

//This struct defines the base performance statistics
type XmlTestDefinition struct {
	XMLName            xml.Name             `xml:"testDefinition"`
	TestName           string               `xml:"testName" toml:"testName"`
	HttpMethod         string               `xml:"httpMethod"`
	BaseUri            string               `xml:"baseUri"`
	Multipart          bool                 `xml:"multipart"`
	Payload            string               `xml:"payload"`
	MultipartPayload   []multipartFormField `xml:"multipartPayload>multipartFormField"`
	ResponseStatusCode int                  `xml:"responseStatusCode"`
	Headers            []Header             `xml:"headers>header"`
	ResponseProperties []string             `xml:"responseProperties>value"`
}

//TomlTestDefinition defines the test in TOML language
type TestDefinition struct {
	TestName           string               `toml:"testName"`
	HttpMethod         string               `toml:"httpMethod"`
	BaseUri            string               `toml:"baseUri"`
	Multipart          bool                 `toml:"multipart"`
	Payload            string               `toml:"payload"`
	MultipartPayload   []multipartFormField `toml:"multipartFormField"`
	ResponseStatusCode int                  `toml:"responseStatusCode"`
	Headers            http.Header          `toml:"headers"`
	ResponseProperties []string             `toml:"responseProperties"`
}

//This struct defines a load test scenario
type TestSuiteDefinition struct {
	XMLName      xml.Name `xml:"testSuite"`
	Name         string   `xml:"name" toml:"name"`
	TestStrategy string   `xml:"testStrategy" toml:"testStrategy"`
	TestCases    []string `xml:"testCases>testCase" toml:"testCases"`
}

//This struct defines a load test scenario //fixme xml flags are needed?
type TestSuite struct {
	XMLName      xml.Name          `xml:"testSuite"`
	Name         string            `xml:"name"`
	TestStrategy string            `xml:"testStrategy"`
	TestCases    []*TestDefinition `xml:"testCases>testCase"`
}

type multipartFormField struct {
	FieldName   string `xml:"fieldName" toml:"fieldName"`
	FieldValue  string `xml:"fieldValue" toml:"fieldValue"`
	FileName    string `xml:"fileName" toml:"fileName"`
	FileContent []byte `xml:"fileContent" toml:"fileContent"`
}

func (ts *TestSuite) BuildTestSuite(configurationSettings *perfTestUtils.Config) {
	log.Info("Building Test Suite ....")

	if configurationSettings.TestSuite == "" {
		ts.Name = "DefaultSuite"
		ts.TestStrategy = SERVICE_BASED_TESTING

		//If no test suite has been defined, treat and all test case files as the suite
		d, err := os.Open(configurationSettings.TestCaseDir)
		if err != nil {
			log.Error("Failed to open test definitions directory. Error:", err)
			os.Exit(1)
		}
		defer d.Close()

		fi, err := d.Readdir(-1)
		if err != nil {
			log.Error("Failed to read files in test definitions directory. Error:", err)
			os.Exit(1)
		}
		if len(fi) == 0 {
			log.Error("No test case files found in specified directory ", configurationSettings.TestCaseDir)
			os.Exit(1)
		}

		for _, fi := range fi {
			bs, err := ioutil.ReadFile(configurationSettings.TestCaseDir + "/" + fi.Name())
			if err != nil {
				log.Error("Failed to read test file. Filename: ", fi.Name(), err)
				continue
			}

			testDefinition, err := loadTestDefinition(bs, configurationSettings)
			if err != nil {
				log.Error("Failed to load test definition. Error:", err)
				os.Exit(1)
			}
			ts.TestCases = append(ts.TestCases, testDefinition)
		}
	} else {
		//If a test suite has been defined, load in all tests associated with the test suite.
		bs, err := ioutil.ReadFile(configurationSettings.TestSuiteDir + "/" + configurationSettings.TestSuite)
		if err != nil {
			log.Error("Failed to read test suite definition file. Filename: ", configurationSettings.TestSuiteDir+"/"+configurationSettings.TestSuite, " ", err)
			os.Exit(1)
		}

		testSuiteDefinition, err := loadTestSuiteDefinition(bs, configurationSettings)
		if err != nil {
			log.Errorf("Failed to load the test suite: %v", err)
			os.Exit(1)
		}

		ts.Name = testSuiteDefinition.Name
		ts.TestStrategy = testSuiteDefinition.TestStrategy
		for _, fi := range testSuiteDefinition.TestCases {
			bs, err := ioutil.ReadFile(configurationSettings.TestCaseDir + "/" + fi)
			if err != nil {
				log.Error("Failed to read test file. Filename: ", fi, err)
				continue
			}

			testDefinition, err := loadTestDefinition(bs, configurationSettings)
			if err != nil {
				log.Error("Failed to load test definition. Error:", err)
			}
			ts.TestCases = append(ts.TestCases, testDefinition)
		}

	}
}

func (testDefinition *TestDefinition) BuildAndSendRequest(targetHost string, targetPort string, uniqueTestRunId string) int64 {

	var req *http.Request

	if !testDefinition.Multipart {
		if testDefinition.Payload != "" {
			paylaod := testDefinition.Payload
			newPayload := substituteRequestValues(&paylaod, uniqueTestRunId)
			req, _ = http.NewRequest(testDefinition.HttpMethod, "http://"+targetHost+":"+targetPort+testDefinition.BaseUri, strings.NewReader(newPayload))
		} else {
			req, _ = http.NewRequest(testDefinition.HttpMethod, "http://"+targetHost+":"+targetPort+testDefinition.BaseUri, nil)
		}
	} else {
		if testDefinition.HttpMethod != "POST" {
			log.Error("Multipart request has to be 'POST' method.")
		} else {
			body := new(bytes.Buffer)
			writer := multipart.NewWriter(body)
			for _, field := range testDefinition.MultipartPayload {
				if field.FileName == "" {
					writer.WriteField(field.FieldName, field.FieldValue)
				} else {
					part, _ := writer.CreateFormFile(field.FieldName, field.FileName)
					io.Copy(part, bytes.NewReader(field.FileContent))
				}
			}
			writer.Close()
			req, _ = http.NewRequest(testDefinition.HttpMethod, "http://"+targetHost+":"+targetPort+testDefinition.BaseUri, body)
			req.Header.Set("Content-Type", writer.FormDataContentType())
		}
	}

	//add headers
	/*	for _, v := range testDefinition.Headers {
		req.Header.Add(v.Key, v.Value)
	}*/
	for k, v := range testDefinition.Headers {
		for _, hv := range v {
			req.Header.Add(k, hv)
		}
	}
	startTime := time.Now()
	if resp, err := (&http.Client{}).Do(req); err != nil {
		log.Error("Error by firing request: ", req, "Error:", err)
		return 0
	} else {

		timeTaken := time.Since(startTime)

		body, _ := ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()

		//Validate service response
		contentLengthOk := perfTestUtils.ValidateResponseBody(body, testDefinition.TestName)
		responseCodeOk := perfTestUtils.ValidateResponseStatusCode(resp.StatusCode, testDefinition.ResponseStatusCode, testDefinition.TestName)
		responseTimeOK := perfTestUtils.ValidateServiceResponseTime(timeTaken.Nanoseconds(), testDefinition.TestName)

		if contentLengthOk && responseCodeOk && responseTimeOK {
			extracResponseValues(testDefinition.TestName, body, testDefinition.ResponseProperties, uniqueTestRunId)
			return timeTaken.Nanoseconds()
		} else {
			return 0
		}
	}
}

func substituteRequestValues(requestBody *string, uniqueTestRunId string) string {
	requestPayloadCopy := *requestBody

	//Get Global Properties for this test run
	testRunGlobals := globals[uniqueTestRunId]
	if testRunGlobals != nil {
		r := regexp.MustCompile("{{(.+)?}}")
		res := r.FindAllString(*requestBody, -1)

		if len(res) > 0 {
			for _, property := range res {
				//remove placeholder syntax
				cleanedPropertyName := strings.TrimPrefix(property, "{{")
				cleanedPropertyName = strings.TrimSuffix(cleanedPropertyName, "}}")

				//lookup value in the test run map
				value := testRunGlobals[cleanedPropertyName]
				if value != "" {
					requestPayloadCopy = strings.Replace(requestPayloadCopy, property, value, 1)
				}
			}

		}
	}
	return requestPayloadCopy
}

func extracResponseValues(testCaseName string, body []byte, resposneProperties []string, uniqueTestRunId string) {
	//Get Global Properties for this test run
	testRunGlobals := globals[uniqueTestRunId]
	if testRunGlobals == nil {
		testRunGlobals = make(map[string]string)
		globals[uniqueTestRunId] = testRunGlobals
	}

	for _, name := range resposneProperties {
		if testRunGlobals[testCaseName+"."+name] == "" {
			r := regexp.MustCompile("<(.+)?:" + name + ">(.+)?</(.+)?:" + name + ">")
			res := r.FindStringSubmatch(string(body))
			testRunGlobals[testCaseName+"."+name] = res[2]
		}
	}
}

func loadTestDefinition(bs []byte, configurationSettings *perfTestUtils.Config) (*TestDefinition, error) {
	testDefinition := &TestDefinition{}
	switch configurationSettings.TestFileFormat {
	case "toml":
		err := toml.Unmarshal(bs, testDefinition)
		if err != nil {
			fmt.Printf("Error occurred loading test definition file: %v\n", err)
			return nil, err
		}
	default:
		td := &XmlTestDefinition{}
		err := xml.Unmarshal(bs, &td)
		if err != nil {
			fmt.Printf("Error occurred loading test definition file: %v\n", err)
			return nil, err
		}
		testDefinition = &TestDefinition{
			TestName:           td.TestName,
			HttpMethod:         td.HttpMethod,
			BaseUri:            td.BaseUri,
			Multipart:          td.Multipart,
			Payload:            td.Payload,
			MultipartPayload:   td.MultipartPayload,
			ResponseStatusCode: td.ResponseStatusCode,
			Headers:            tomlHeaders(td.Headers),
			ResponseProperties: td.ResponseProperties,
		}
	}
	return testDefinition, nil
}

func tomlHeaders(headers []Header) http.Header {
	h := make(http.Header)
	for _, v := range headers {
		h.Add(v.Key, v.Value)
	}
	return h
}

func loadTestSuiteDefinition(bs []byte, configurationSettings *perfTestUtils.Config) (*TestSuiteDefinition, error) {
	ts := &TestSuiteDefinition{}
	switch configurationSettings.TestFileFormat {
	case "toml":
		err := toml.Unmarshal(bs, ts)
		if err != nil {
			fmt.Printf("Error occurred loading test definition file: %v\n", err)
			return nil, err
		}
	default:
		err := xml.Unmarshal(bs, ts)
		if err != nil {
			fmt.Printf("Error occurred loading test definition file: %v\n", err)
			return nil, err
		}
	}
	return ts, nil
}
