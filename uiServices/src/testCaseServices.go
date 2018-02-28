package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"github.com/Sirupsen/logrus"
	"github.com/xtracdev/automated-perf-test/testStrategies"
	"bytes"
	"github.com/go-chi/chi"
)

const testCaseSchema string = "testCase_schema.json"
const structTypeName string = "TestCase "

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

	if !testCaseWriterXml(testCase, testCasePathDir + testCase.TestName+".xml") {
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

	if len(testCase.TestName) < 1{
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