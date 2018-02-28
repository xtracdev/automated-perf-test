package services

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
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

func getAllTestCases(rw http.ResponseWriter, req *http.Request) {

	testCasePathDir := getTestCaseHeader(req)
	if len(testCasePathDir) <= 1 {
		logrus.Error("No file directory entered")
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	files, err := ioutil.ReadDir(testCasePathDir)
	if err != nil {
		log.Fatal(err)
	}

	testCases := make([]Case, 0)

	for _, file := range files {
		if filepath.Ext(testCasePathDir+file.Name()) == ".xml" {

			testCase := new(testStrategies.TestDefinition)

			filename := file.Name()

			file, err := os.Open(fmt.Sprintf("%s%s", testCasePathDir, filename))
			if err != nil {
				continue
			}

			byteValue, err := ioutil.ReadAll(file)
			if err != nil {
				continue
			}

			err = xml.Unmarshal(byteValue, testCase)
			if err != nil {
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
