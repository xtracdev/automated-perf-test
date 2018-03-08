package services

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/go-chi/chi"
	"github.com/xtracdev/automated-perf-test/testStrategies"
)

const testCaseSchema string = "testCase_schema.json"
const structTypeName string = "TestCase "

type Case struct {
	HttpMethod  string `json:"httpMethod"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func TestCaseCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})

}
func getTestCaseHeader(req *http.Request) string {
	testCasePathDir := req.Header.Get("testCasePathDir")

	if !strings.HasSuffix(testCasePathDir, "/") {
		testCasePathDir = testCasePathDir + "/"
	}
	return testCasePathDir
}

func postTestCase(rw http.ResponseWriter, req *http.Request) {
	testCasePathDir := getTestCaseHeader(req)
	buf := new(bytes.Buffer)
	buf.ReadFrom(req.Body)

	testCase := testStrategies.TestDefinition{}
	err := json.Unmarshal(buf.Bytes(), &testCase)
	if err != nil {
		logrus.Error("Failed to unmarshall json body")
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	if !ValidateFileNameAndHeader(rw, req, testCasePathDir, testCase.TestName) {
		return
	}

	if !FilePathExist(testCasePathDir) {
		logrus.Error("Directory path does not exist")
		rw.WriteHeader(http.StatusBadRequest)
		return

	}

	if FilePathExist(fmt.Sprintf("%s%s.xml", testCasePathDir, testCase.TestName)) {
		logrus.Error("File already exists")
		rw.WriteHeader(http.StatusBadRequest)
		return

	}

	if !testCaseWriterXml(testCase, testCasePathDir+testCase.TestName+".xml") {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusCreated)
}

func putTestCase(rw http.ResponseWriter, req *http.Request) {
	path := getTestCaseHeader(req)
	testCaseName := chi.URLParam(req, "testCaseName")

	if !ValidateFileNameAndHeader(rw, req, path, testCaseName) {
		return
	}

	testCasePathDir := fmt.Sprintf("%s%s.xml", path, testCaseName)
	buf := new(bytes.Buffer)
	buf.ReadFrom(req.Body)

	if !FilePathExist(testCasePathDir) {
		logrus.Error("File path does not exist")
		rw.WriteHeader(http.StatusNotFound)
		return

	}

	if !ValidateJSONWithSchema(buf.Bytes(), testCaseSchema, structTypeName) {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	testCase := testStrategies.TestDefinition{}
	err := json.Unmarshal(buf.Bytes(), &testCase)

	if err != nil {
		logrus.Error("Cannot Unmarshall Json")
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	if len(testCase.TestName) < 1 {
		logrus.Error("No TestName Entered")
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	if !testCaseWriterXml(testCase, testCasePathDir) {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusNoContent)
}

func getAllTestCases(rw http.ResponseWriter, req *http.Request) {

	testCasePathDir := getTestCaseHeader(req)

	if !IsPathDirValid(testCasePathDir, rw) {
		return
	}

	// if len(testCasePathDir) <= 1 {
	// 	logrus.Error("No file directory entered")
	// 	rw.WriteHeader(http.StatusBadRequest)
	// 	return
	// }

	files, err := ioutil.ReadDir(testCasePathDir)
	if err != nil {
		logrus.Error(err)
	}

	testCases := make([]Case, 0)

	for _, file := range files {
		if filepath.Ext(testCasePathDir+file.Name()) == ".xml" {

			testCase := new(testStrategies.TestDefinition)

			filename := file.Name()

			file, err := os.Open(fmt.Sprintf("%s%s", testCasePathDir, filename))
			if err != nil {
				logrus.Error("Cannot Open File: " + filename)
				continue
			}

			byteValue, err := ioutil.ReadAll(file)
			if err != nil {
				logrus.Error("Cannot Read File: " + filename)
				continue
			}

			err = xml.Unmarshal(byteValue, testCase)
			if err != nil {
				logrus.Error("Cannot Unmarshall File: " + filename)
				continue
			}

			//if a Test Case Name can't be assigned, it isn't a Test Case object
			if testCase.TestName != "" {
				testCases = append(testCases, Case{
					Name:        testCase.TestName,
					Description: testCase.Description,
					HttpMethod:  testCase.HTTPMethod,
				})
			}
		}
	}

	err = json.NewEncoder(rw).Encode(testCases)
	if err != nil {
		logrus.Error("Could not enocde Test Cases")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
}

func getTestCase(rw http.ResponseWriter, req *http.Request) {

	testCasePathDir := getTestCaseHeader(req)
	testCaseName := chi.URLParam(req, "testCaseName")

	if !ValidateFileNameAndHeader(rw, req, testCasePathDir, testCaseName) {
		return
	}

	if _, err := os.Stat(fmt.Sprintf("%s%s.xml", testCasePathDir, testCaseName)); err != nil {
		if os.IsNotExist(err) {
			logrus.Error("Test Case File Not Found: " + testCaseName)
			rw.WriteHeader(http.StatusNotFound)
			return
		}
	}

	file, err := os.Open(fmt.Sprintf("%s%s.xml", testCasePathDir, testCaseName))
	if err != nil {
		logrus.Error("Cannot open: " + testCaseName)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer file.Close()

	var testCase testStrategies.TestDefinition

	byteValue, err := ioutil.ReadAll(file)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		logrus.Error("Cannot Read File", err)
		return
	}

	err = xml.Unmarshal(byteValue, &testCase)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		logrus.Error("Cannot Unmarshall from XML", err)
		return
	}

	testSuiteJSON, err := json.MarshalIndent(testCase, "", "")
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		logrus.Error("Cannot marshall to JSON", err)
		return
	}

	rw.WriteHeader(http.StatusOK)
	rw.Write(testSuiteJSON)

}

func deleteTestCase(rw http.ResponseWriter, req *http.Request) {
	testCasePathDir := getTestCaseHeader(req)
	testCaseName := chi.URLParam(req, "testCaseName")

	if !ValidateFileNameAndHeader(rw, req, testCasePathDir, testCaseName) {
		return
	}

	filepath := fmt.Sprintf("%s%s.xml", testCasePathDir, testCaseName)

	if _, err := os.Stat(filepath); err != nil {
		if os.IsNotExist(err) {
			logrus.Println("File Not Found", err)
			rw.WriteHeader(http.StatusNotFound)
			return
		}
	}

	err := os.Remove(filepath)
	if err != nil {
		logrus.Println("File was not deleted", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return

	}

	rw.WriteHeader(http.StatusNoContent)

}
