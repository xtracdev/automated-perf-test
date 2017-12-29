package testStrategies

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
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

// ServiceBasedTesting and SuiteBasedTesting are used as boolean to determine
// and control which strategy to use throughout the test run.
const (
	ServiceBasedTesting = "ServiceBased"
	SuiteBasedTesting   = "SuiteBased"
)

// Global Mutex
var mu sync.Mutex

// globalsMap contains parameter substitution across concurrent threads.
var globalsMap = make(map[string]map[string]interface{})

// Header defines request header key/value pair.
type Header struct {
	Value string `xml:",chardata"`
	Key   string `xml:"key,attr"`
}

// ResponseValue encapsulates the variable name (ExtractionKey) and the value
// from the call response. Used for variable substitution between calls.
type ResponseValue struct {
	Value         string `xml:",chardata"`
	ExtractionKey string `xml:"extractionKey,attr"`
}

// TestDefinition encapsulates the XML data.
type TestDefinition struct {
	XMLName             xml.Name             `xml:"testDefinition"`
	TestName            string               `xml:"testName"`
	OverrideHost        string               `xml:"overrideHost"`
	OverridePort        string               `xml:"overridePort"`
	HTTPMethod          string               `xml:"httpMethod"`
	BaseURI             string               `xml:"baseUri"`
	Multipart           bool                 `xml:"multipart"`
	Payload             string               `xml:"payload"`
	MultipartPayload    []multipartFormField `xml:"multipartPayload>multipartFormField"`
	ResponseStatusCode  int                  `xml:"responseStatusCode"`
	ResponseContentType string               `xml:"responseContentType"`
	Headers             []Header             `xml:"headers>header"`
	ResponseValues      []ResponseValue      `xml:"responseProperties>value"`
	PreThinkTime        int64
	PostThinkTime       int64
	ExecWeight          string
}

// TestSuite fields get populated from the TestSuiteDefinition after the XML
// unmarshal is complete. (See TestDefinition above).
type TestSuite struct {
	XMLName         xml.Name   `xml:"testSuite"`
	Name            string     `xml:"name"`
	TestStrategy    string     `xml:"testStrategy"`
	TestCases       []TestCase `xml:"testCases>testCase"`
	TestDefinitions []*TestDefinition
}

// TestCase is used to encapsulate and marshal a <testCase> tag from the
// <testSuite> XML file. This data will then be consolidated into
// the TestSuite/TestDefinition data structure for usage.
type TestCase struct {
	XMLName       xml.Name `xml:"testCase"`
	Name          string   `xml:",chardata"`
	PreThinkTime  int64    `xml:"preThinkTime,attr"`
	PostThinkTime int64    `xml:"postThinkTime,attr"`
	ExecWeight    string   `xml:"execWeight,attr"`
}

type multipartFormField struct {
	FieldName   string `xml:"fieldName"`
	FieldValue  string `xml:"fieldValue"`
	FileName    string `xml:"fileName"`
	FileContent []byte `xml:"fileContent"`
}

// BuildTestSuite sets TestStrategy and puts together the test suite
// accordingly. For ServiceBasedTesting, the test suite is built from test
// cases in the TestCaseDir, unordered. For SuiteBasedTesting, the test suite
// is built according to the testCases listed in the test suite definition.
func (ts *TestSuite) BuildTestSuite(configurationSettings *perfTestUtils.Config) {
	log.Info("Building Test Suite ....")
	// Default to ServiceBased testing:
	ts.TestStrategy = ServiceBasedTesting

	if configurationSettings.TestSuite == "" {
		ts.Name = "DefaultSuite"

		// If no test suite has been defined, treat and all test case files
		// found in the TestCaseDir as the suite.
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

			testDefinition, err := loadTestDefinition(bs)
			if err != nil {
				log.Error("Failed to load test definition. Error:", err)
				os.Exit(1)
			}
			ts.TestDefinitions = append(ts.TestDefinitions, testDefinition)
		}
	} else {
		// Flag as SuiteBased testing:
		ts.TestStrategy = SuiteBasedTesting

		// Get the testSuite filename ...
		bs, err := ioutil.ReadFile(configurationSettings.TestSuiteDir + "/" + configurationSettings.TestSuite)
		if err != nil {
			log.Error("Failed to read test suite definition file. Filename: ", configurationSettings.TestSuiteDir+"/"+configurationSettings.TestSuite, " ", err)
			os.Exit(1)
		}

		// ... and unmarshal the testSuite XML.
		err = ts.loadTestSuiteDefinition(bs)
		if err != nil {
			log.Errorf("Failed to load the test suite: %v", err)
			os.Exit(1)
		}

		// Populate ts.TestDefinitions array with test definitions.
		for _, testCase := range ts.TestCases {
			bs, err := ioutil.ReadFile(configurationSettings.TestCaseDir + "/" + testCase.Name)
			if err != nil {
				log.Error("Failed to read test file. Filename: ", testCase.Name, err)
				continue
			}

			testDefinition, err := loadTestDefinition(bs)
			if err != nil {
				log.Error("Failed to load test definition. Error:", err)
			}

			// Add the testCase attributes to the TestDefinition (thinktime, etc).
			// This effectively flattens the fields into TestDefinitions allowing
			// us to ignore TestCases.
			testDefinition.PreThinkTime = testCase.PreThinkTime
			testDefinition.PostThinkTime = testCase.PostThinkTime
			testDefinition.ExecWeight = testCase.ExecWeight

			// Append the testDefinition to the testSuite
			ts.TestDefinitions = append(ts.TestDefinitions, testDefinition)
		}
	}
}

func loadTestDefinition(bs []byte) (*TestDefinition, error) {
	td := &TestDefinition{}
	err := xml.Unmarshal(bs, &td)
	if err != nil {
		log.Errorf("Error occurred loading XML testCase definition file: %v\n", err)
		return nil, err
	}
	return td, nil
}

func (ts *TestSuite) loadTestSuiteDefinition(bs []byte) error {
	//ts := &TestSuite{}
	err := xml.Unmarshal(bs, ts)
	if err != nil {
		log.Errorf("Error occurred loading XML testSuite definition file: %v\n", err)
		return err
	}
	return nil
}

// BuildAndSendRequest builds a request from the test definition, performs
// variable substitutions, sends the request, and returns the response time
// of the call, or 0 if failure. Note: Response time does not include
// RequestDelay or ThinkTime.
func (testDefinition *TestDefinition) BuildAndSendRequest(
	delay int,
	targetHost string,
	targetPort string,
	uniqueTestRunID string,
) int64 {
	log.Debugf("BEGIN \"%s\" testDefinition\n-----\n%+v\n-----\nEND \"%s\" testDefinition\n",
		testDefinition.TestName,
		testDefinition,
		testDefinition.TestName,
	)

	randomDelay := rand.Intn(delay)
	time.Sleep(time.Duration(randomDelay) * time.Millisecond)

	//Execute the PreThinkTime, if any.
	if testDefinition.PreThinkTime > 0 {
		tt := float64(testDefinition.PreThinkTime) / 1000
		log.Infof("Think time: [%.2f] seconds.", tt)
	}
	time.Sleep(time.Duration(testDefinition.PreThinkTime) * time.Millisecond)

	var req *http.Request
	reqbody := "N/A" //for debug

	//Retrieve requestBaseURI and perform any necessary substitution
	requestBaseURI := substituteRequestValues(&testDefinition.BaseURI, uniqueTestRunID)

	if !testDefinition.Multipart {
		log.Debug("Building non-Multipart request.")
		if testDefinition.Payload != "" {
			//Retrieve Payload and perform any necessary substitution
			payload := testDefinition.Payload
			newPayload := substituteRequestValues(&payload, uniqueTestRunID)
			reqbody = newPayload
			req, _ = http.NewRequest(testDefinition.HTTPMethod, "http://"+targetHost+":"+targetPort+requestBaseURI, strings.NewReader(newPayload))
		} else {
			req, _ = http.NewRequest(testDefinition.HTTPMethod, "http://"+targetHost+":"+targetPort+requestBaseURI, nil)
		}
	} else {
		log.Debug("Building Multipart request.")
		if testDefinition.HTTPMethod != "POST" {
			log.Error("Multipart request must be 'POST' method.")
		} else {
			body := new(bytes.Buffer)
			writer := multipart.NewWriter(body)
			for _, field := range testDefinition.MultipartPayload {
				if field.FileName == "" {
					writer.WriteField(field.FieldName, substituteRequestValues(&field.FieldValue, uniqueTestRunID))
				} else {
					part, _ := writer.CreateFormFile(field.FieldName, field.FileName)
					io.Copy(part, bytes.NewReader(field.FileContent))
				}
			}
			writer.Close()
			req, _ = http.NewRequest(testDefinition.HTTPMethod, "http://"+targetHost+":"+targetPort+requestBaseURI, body)
			req.Header.Set("Content-Type", writer.FormDataContentType())

			// For debug output
			reqbody = body.String()
		}
	}

	// Add headers and perform and necessary substitution
	for _, header := range testDefinition.Headers {
		req.Header.Add(header.Key, substituteRequestValues(&header.Value, uniqueTestRunID))
	}

	log.Debugf(
		"BEGIN \"%s\" Request:\n-----\nHEADER:%+v\nURL:%s\nREQ_BODY:%s\n-----\nEND [%s] Request",
		testDefinition.TestName,
		req.Header,
		req.URL,
		reqbody,
		testDefinition.TestName,
	)

	var resp *http.Response
	var err error
	startTime := time.Now()
	if resp, err = (&http.Client{}).Do(req); err != nil {
		log.Errorf("Connection failed for request [Name:%s]: %+v", testDefinition.TestName, err)
		return 0
	}
	// Mark response time.
	timeTaken := time.Since(startTime)
	// Gather the response.
	body, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	log.Debugf(
		"BEGIN \"%s\" Response:\n-----\nSTATUSCODE:%d\nHEADER:%+v\nRESP_BODY:%s\n-----\nEND [%s] Response",
		testDefinition.TestName,
		resp.StatusCode,
		resp.Header,
		string(body),
		testDefinition.TestName,
	)

	//Validate service response
	responseCodeOk := perfTestUtils.ValidateResponseStatusCode(resp.StatusCode, testDefinition.ResponseStatusCode, testDefinition.TestName)
	responseTimeOK := perfTestUtils.ValidateServiceResponseTime(timeTaken.Nanoseconds(), testDefinition.TestName)

	if !responseCodeOk || !responseTimeOK {
		return 0
	}

	contentType := detectContentType(resp.Header, body, testDefinition.ResponseContentType)
	extractResponseValues(testDefinition.TestName, body, testDefinition.ResponseValues, uniqueTestRunID, contentType)

	//Execute the PostThinkTime, if any.
	if testDefinition.PostThinkTime > 0 {
		tt := float64(testDefinition.PostThinkTime) / 1000
		log.Infof("Think time: [%.2f] seconds.", tt)
	}
	time.Sleep(time.Duration(testDefinition.PostThinkTime) * time.Millisecond)

	return timeTaken.Nanoseconds()
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

func substituteRequestValues(requestBody *string, uniqueTestRunID string) string {
	// Make a value copy for substitution.
	requestPayloadCopy := *requestBody

	// Lock global data structure for the duration of this function and
	// those function that branch.
	mu.Lock()
	defer mu.Unlock()

	// Get the properties for this iteration.
	testRunGlobals := globalsMap[uniqueTestRunID]

	if testRunGlobals != nil {
		r := regexp.MustCompile("{{([^}]+)}}")
		res := r.FindAllString(*requestBody, -1)

		if len(res) > 0 {
			for _, propertyPlaceHolder := range res {
				//remove placeholder syntax
				cleanedPropertyName := strings.TrimPrefix(propertyPlaceHolder, "{{")
				cleanedPropertyName = strings.TrimSuffix(cleanedPropertyName, "}}")

				propertyPlaceHolderName := cleanedPropertyName
				propertyPlaceHolderIndex := 0

				// Extract the index number if array notation encountered.
				if strings.Contains(cleanedPropertyName, "[") && strings.Contains(cleanedPropertyName, "]") {
					propertyPlaceHolderName, propertyPlaceHolderIndex = getArrayNameAndIndex(testRunGlobals, cleanedPropertyName)
				}

				// Get the value based on the property name.
				placeHolderValue := testRunGlobals[propertyPlaceHolderName]
				valueAsString := convertStoredValueToRequestFormat(placeHolderValue, propertyPlaceHolderIndex)

				if valueAsString != "" {
					requestPayloadCopy = strings.Replace(requestPayloadCopy, propertyPlaceHolder, valueAsString, 1)
				}
			}

		}
	}
	return requestPayloadCopy
}

// getArrayNameAndIndex
func getArrayNameAndIndex(testRunGlobals map[string]interface{}, propPathPart string) (string, int) {

	propPathPart = strings.Replace(propPathPart, "[", "::", 1)
	propPathPart = strings.Replace(propPathPart, "]", "", 1)
	propertyNameParts := strings.Split(propPathPart, "::")

	// Get index or set to 0 if not a valid number.
	index, _ := strconv.Atoi(propertyNameParts[1])

	// Test to see if the property exists from a previous call.
	if testRunGlobals[propertyNameParts[0]] == nil {
		// Return index of 0 in cases where the property failed to populate on
		// a previous call: eg. connection error, timeout, http 502/503, etc.
		// Note: this may cause all subsequent calls within the
		// same uniqueTestRunID iteration to fail.
		return propertyNameParts[0], index
	}

	// A value of '?' rather than a number in the index notation indicates
	// the user would like a random record returned from the data set.
	if propertyNameParts[1] == "?" {
		// Seed the rand call.
		randIdx := rand.New(rand.NewSource(time.Now().UnixNano()))
		// Get the length of the property array to serve as boundary
		// for rand.
		arylen := len(testRunGlobals[propertyNameParts[0]].([]interface{}))
		// Check to ensure the array is not empty. We are not able to
		// continue in this case. The user must fix the data issue
		// before proceeding.
		if arylen == 0 {
			log.Errorf("FATAL: Unable to substitute property [%s]: Result array of size 0. Check data criteria for service call.", propertyNameParts[0])
			os.Exit(1)
		}
		// Set the index to a random value.
		index = randIdx.Intn(arylen)
	}

	return propertyNameParts[0], index
}

func convertStoredValueToRequestFormat(storedValue interface{}, requiredIndex int) string {
	requestFormattedValue := ""
	switch objectType := storedValue.(type) {

	case map[string]interface{}:
		jsonValue, _ := json.Marshal(objectType)
		requestFormattedValue = string(jsonValue)
	case []interface{}:
		value := objectType[requiredIndex]
		requestFormattedValue = convertStoredValueToRequestFormat(value, 0)
	case string:
		requestFormattedValue = string(objectType)
	case int64:
		requestFormattedValue = string(strconv.FormatInt( objectType, 10))
	case float64:
		requestFormattedValue = string(strconv.FormatFloat( objectType, 'f', 0, 64))
	default:
		requestFormattedValue = fmt.Sprintf("%v", objectType)
	}
	return requestFormattedValue
}

func extractResponseValues(testCaseName string, body []byte, responseValues []ResponseValue, uniqueTestRunID string, contentType string) {
	// Short-circuit if call returned empty response body.
	if string(body) == "" {
		return
	}

	if strings.Contains(contentType, "json") {
		extractJSONResponseValues(testCaseName, body, responseValues, uniqueTestRunID)
	} else if strings.Contains(contentType, "xml") {
		extractXMLResponseValues(testCaseName, body, responseValues, uniqueTestRunID)
	} else {
		log.Warn("Unsupported response content type of:", contentType)
	}
}

//----- extractJSONResponseValues --------------------------------------------
// Extract the response values from the JSON result based on the JMESPath
// query provided by the user.
func extractJSONResponseValues(testCaseName string, body []byte, responseValues []ResponseValue, uniqueTestRunID string) {
	// Get Global Properties for this test run.
	mu.Lock()
	testRunGlobals := globalsMap[uniqueTestRunID]
	mu.Unlock()

	if testRunGlobals == nil {
		testRunGlobals = make(map[string]interface{})
		mu.Lock()
		globalsMap[uniqueTestRunID] = testRunGlobals
		mu.Unlock()
	}

	for _, propPath := range responseValues {

		var data interface{}
		json.Unmarshal(body, &data)

		result, _ := jmespath.Search(propPath.Value, data)

		if testRunGlobals[testCaseName+"."+propPath.ExtractionKey] == nil {
			testRunGlobals[testCaseName+"."+propPath.ExtractionKey] = result
		}
	}
}

func extractXMLResponseValues(testCaseName string, body []byte, responseValues []ResponseValue, uniqueTestRunID string) {
	// Get Global Properties for this test run.
	mu.Lock()
	testRunGlobals := globalsMap[uniqueTestRunID]
	mu.Unlock()

	if testRunGlobals == nil {
		testRunGlobals = make(map[string]interface{})
		mu.Lock()
		globalsMap[uniqueTestRunID] = testRunGlobals
		mu.Unlock()
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
