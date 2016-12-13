package testStrategies

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/BurntSushi/toml"
	log "github.com/Sirupsen/logrus"
	"github.com/jmespath/go-jmespath"
	"github.com/xtracdev/automated-perf-test/perfTestUtils"
	"io"
	"io/ioutil"
	"math/rand"
	"mime/multipart"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	SERVICE_BASED_TESTING = "ServiceBased"
	SUITE_BASED_TESTING   = "SuiteBased"
)

type GlobalsMaps struct {
	sync.RWMutex
	m map[string]map[string]interface{}
}

var GlobalsLockCounter = GlobalsMaps{m: make(map[string]map[string]interface{})}

type Header struct {
	Value string `xml:",chardata"`
	Key   string `xml:"key,attr"`
}
type ResponseValue struct {
	Value         string `xml:",chardata" toml:"value"`
	ExtractionKey string `xml:"extractionKey,attr" toml:"extractionKey"`
}

//This struct defines the base performance statistics
type XmlTestDefinition struct {
	XMLName             xml.Name             `xml:"testDefinition"`
	TestName            string               `xml:"testName" toml:"testName"`
	OverrideHost        string               `xml:"overrideHost" toml:"overrideHost"`
	OverridePort        string               `xml:"overridePort" toml:"overridePort"`
	HttpMethod          string               `xml:"httpMethod"`
	BaseUri             string               `xml:"baseUri"`
	Multipart           bool                 `xml:"multipart"`
	Payload             string               `xml:"payload"`
	MultipartPayload    []multipartFormField `xml:"multipartPayload>multipartFormField"`
	ResponseStatusCode  int                  `xml:"responseStatusCode"`
	ResponseContentType string               `xml:"responseContentType"`
	Headers             []Header             `xml:"headers>header"`
	ResponseValues      []ResponseValue      `xml:"responseProperties>value"`
}

//TomlTestDefinition defines the test in TOML language
type TestDefinition struct {
	TestName            string               `toml:"testName"`
	OverrideHost        string               `toml:"overrideHost"`
	OverridePort        string               `toml:"overridePort"`
	HttpMethod          string               `toml:"httpMethod"`
	BaseUri             string               `toml:"baseUri"`
	Multipart           bool                 `toml:"multipart"`
	Payload             string               `toml:"payload"`
	MultipartPayload    []multipartFormField `toml:"multipartFormField"`
	ResponseStatusCode  int                  `toml:"responseStatusCode"`
	ResponseContentType string               `xml:"responseContentType"`
	Headers             http.Header          `toml:"headers"`
	ResponseValues      []ResponseValue      `toml:"responseProperties"`
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
	// Default to ServiceBased testing:
	ts.TestStrategy = SERVICE_BASED_TESTING

	if configurationSettings.TestSuite == "" {
		ts.Name = "DefaultSuite"

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
		// Flag as SuiteBased testing:
		ts.TestStrategy = SUITE_BASED_TESTING

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

func (testDefinition *TestDefinition) BuildAndSendRequest(delay int, targetHost string, targetPort string, uniqueTestRunId string, globalsMap GlobalsMaps) int64 {

	randomDelay := rand.Intn(delay)
	time.Sleep(time.Duration(randomDelay) * time.Millisecond)

	var req *http.Request

	//Retrieve requestBaseURI and perform any necessary substitution
	requestBaseURI := substituteRequestValues(&testDefinition.BaseUri, uniqueTestRunId, globalsMap)

	if !testDefinition.Multipart {
		if testDefinition.Payload != "" {
			//Retrieve Payload and perform any necessary substitution
			paylaod := testDefinition.Payload
			newPayload := substituteRequestValues(&paylaod, uniqueTestRunId, globalsMap)

			req, _ = http.NewRequest(testDefinition.HttpMethod, "http://"+targetHost+":"+targetPort+requestBaseURI, strings.NewReader(newPayload))
		} else {
			req, _ = http.NewRequest(testDefinition.HttpMethod, "http://"+targetHost+":"+targetPort+requestBaseURI, nil)
		}
	} else {
		if testDefinition.HttpMethod != "POST" {
			log.Error("Multipart request has to be 'POST' method.")
		} else {
			body := new(bytes.Buffer)
			writer := multipart.NewWriter(body)
			for _, field := range testDefinition.MultipartPayload {
				if field.FileName == "" {
					writer.WriteField(field.FieldName, substituteRequestValues(&field.FieldValue, uniqueTestRunId, globalsMap))
				} else {
					part, _ := writer.CreateFormFile(field.FieldName, field.FileName)
					io.Copy(part, bytes.NewReader(field.FileContent))
				}
			}
			writer.Close()
			req, _ = http.NewRequest(testDefinition.HttpMethod, "http://"+targetHost+":"+targetPort+requestBaseURI, body)
			req.Header.Set("Content-Type", writer.FormDataContentType())
		}
	}

	//add headers and perform and necessary substitution
	for k, v := range testDefinition.Headers {
		for _, hv := range v {
			req.Header.Add(k, substituteRequestValues(&hv, uniqueTestRunId, globalsMap))
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
		responseCodeOk := perfTestUtils.ValidateResponseStatusCode(resp.StatusCode, testDefinition.ResponseStatusCode, testDefinition.TestName)
		responseTimeOK := perfTestUtils.ValidateServiceResponseTime(timeTaken.Nanoseconds(), testDefinition.TestName)

		if responseCodeOk && responseTimeOK {
			contentType := detectContentType(resp.Header, body, testDefinition.ResponseContentType)
			extractResponseValues(testDefinition.TestName, body, testDefinition.ResponseValues, uniqueTestRunId, globalsMap, contentType)
			return timeTaken.Nanoseconds()
		} else {
			return 0
		}
	}
}
func detectContentType(respHeaders http.Header, respBody []byte, respContentType string) string {
	if respHeaders.Get("Content-Type") != "" {
		return respHeaders.Get("Content-Type")
	} else if respContentType != "" {
		return respContentType
	} else {
		return http.DetectContentType(respBody)
	}
}

func determineHostandPortforRequest(testDefinition *TestDefinition, configurationSettings *perfTestUtils.Config) (string, string) {

	var targetHost = configurationSettings.TargetHost
	var targetPort = configurationSettings.TargetPort

	if testDefinition.OverrideHost != "" {
		targetHost = testDefinition.OverrideHost
	}
	if testDefinition.OverridePort != "" {
		targetPort = testDefinition.OverridePort
	}
	return targetHost, targetPort
}

func substituteRequestValues(requestBody *string, uniqueTestRunId string, globalsMap GlobalsMaps) string {
	requestPayloadCopy := *requestBody

	//Get Global Properties for this test run
	globalsMap.RLock()
	testRunGlobals := globalsMap.m[uniqueTestRunId]
	globalsMap.RUnlock()

	if testRunGlobals != nil {
		r := regexp.MustCompile("{{(.[^ ]+)?}}")
		res := r.FindAllString(*requestBody, -1)

		if len(res) > 0 {
			for _, propertyPlaceHolder := range res {
				//remove placeholder syntax
				cleanedPropertyName := strings.TrimPrefix(propertyPlaceHolder, "{{")
				cleanedPropertyName = strings.TrimSuffix(cleanedPropertyName, "}}")

				propertyPlaceHolderName := cleanedPropertyName
				propertyPlaceHolderIndex := 0

				if strings.Contains(cleanedPropertyName, "[") && strings.Contains(cleanedPropertyName, "]") {
					propertyPlaceHolderName, propertyPlaceHolderIndex = getArrayNameAndIndex(cleanedPropertyName)
				}

				placeHolderValue := testRunGlobals[propertyPlaceHolderName]

				valueAsString := convertStoredValuetoRequestFormat(placeHolderValue, propertyPlaceHolderIndex)

				if valueAsString != "" {
					requestPayloadCopy = strings.Replace(requestPayloadCopy, propertyPlaceHolder, valueAsString, 1)
				}
			}

		}
	}
	return requestPayloadCopy
}
func convertStoredValuetoRequestFormat(storedValue interface{}, requiredIndex int) string {
	requestFormattedValue := ""
	switch objectType := storedValue.(type) {

	case map[string]interface{}:
		jsonValue, _ := json.Marshal(objectType)
		requestFormattedValue = string(jsonValue)
	case []interface{}:
		value := objectType[requiredIndex]
		requestFormattedValue = convertStoredValuetoRequestFormat(value, 0)
	case string:
		requestFormattedValue = string(objectType)
	default:
		requestFormattedValue = fmt.Sprintf("%v", objectType)
	}
	return requestFormattedValue
}

func extractResponseValues(testCaseName string, body []byte, responseValues []ResponseValue, uniqueTestRunId string, globalsMap GlobalsMaps, contentType string) {
	if strings.Contains(contentType, "json") {
		extractJSONResponseValues(testCaseName, body, responseValues, uniqueTestRunId, globalsMap)
	} else if strings.Contains(contentType, "xml") {
		extractXMLResponseValues(testCaseName, body, responseValues, uniqueTestRunId, globalsMap)
	} else {
		log.Warn("Unsupported resposne content type of:", contentType)
	}
}

func extractJSONResponseValues(testCaseName string, body []byte, responseValues []ResponseValue, uniqueTestRunId string, globalsMap GlobalsMaps) {
	//Get Global Properties for this test run
	globalsMap.RLock()
	testRunGlobals := globalsMap.m[uniqueTestRunId]
	globalsMap.RUnlock()

	if testRunGlobals == nil {
		testRunGlobals = make(map[string]interface{})
		globalsMap.Lock()
		globalsMap.m[uniqueTestRunId] = testRunGlobals
		globalsMap.Unlock()
	}

	for _, propPath := range responseValues {

		var data interface{}
		json.Unmarshal(body, &data)

		result, _ := jmespath.Search(propPath.Value, data) //Todo handle error

		if testRunGlobals[testCaseName+"."+propPath.ExtractionKey] == nil {
			testRunGlobals[testCaseName+"."+propPath.ExtractionKey] = result
		}
	}
}

func getArrayNameAndIndex(propPathPart string) (string, int) {

	propPathPart = strings.Replace(propPathPart, "[", "::", 1)
	propPathPart = strings.Replace(propPathPart, "]", "", 1)
	propertyNameParts := strings.Split(propPathPart, "::")

	index, _ := strconv.Atoi(propertyNameParts[1]) //todo, handle error

	return propertyNameParts[0], index
}

func extractXMLResponseValues(testCaseName string, body []byte, responseValues []ResponseValue, uniqueTestRunId string, globalsMap GlobalsMaps) {
	//Get Global Properties for this test run
	globalsMap.RLock()
	testRunGlobals := globalsMap.m[uniqueTestRunId]
	globalsMap.RUnlock()

	if testRunGlobals == nil {
		testRunGlobals = make(map[string]interface{})
		globalsMap.Lock()
		globalsMap.m[uniqueTestRunId] = testRunGlobals
		globalsMap.Unlock()
	}

	for _, responseValue := range responseValues {
		extractionKey := responseValue.ExtractionKey
		if extractionKey == "" {
			extractionKey = responseValue.Value
		}

		if testRunGlobals[testCaseName+"."+responseValue.Value] == "" {
			r := regexp.MustCompile("<(.+)?:" + responseValue.Value + ">(.+)?</(.+)?:" + responseValue.Value + ">")
			res := r.FindStringSubmatch(string(body))
			testRunGlobals[testCaseName+"."+responseValue.ExtractionKey] = res[2]
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
			TestName:            td.TestName,
			OverrideHost:        td.OverrideHost,
			OverridePort:        td.OverridePort,
			HttpMethod:          td.HttpMethod,
			BaseUri:             td.BaseUri,
			Multipart:           td.Multipart,
			Payload:             td.Payload,
			MultipartPayload:    td.MultipartPayload,
			ResponseStatusCode:  td.ResponseStatusCode,
			ResponseContentType: td.ResponseContentType,
			Headers:             tomlHeaders(td.Headers),
			ResponseValues:      td.ResponseValues,
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
