package services

import (
	"bytes"
	"encoding/json"
	"github.com/Sirupsen/logrus"
	"github.com/xeipuuv/gojsonschema"
	"github.com/xtracdev/automated-perf-test/testStrategies"
	"net/http"
	"os"
	"strings"
	"fmt"
)

func TestSuiteCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})

}
func getTestSuiteHeader(req *http.Request) string {
	testSuitePathDir := req.Header.Get("testSuitePathDir")

	if !strings.HasSuffix(testSuitePathDir, "/") {
		testSuitePathDir = testSuitePathDir + "/"
	}
	return testSuitePathDir
}

func postTestSuites(rw http.ResponseWriter, req *http.Request) {
	testSuitePathDir := getTestSuiteHeader(req)
	buf := new(bytes.Buffer)
	buf.ReadFrom(req.Body)

	testSuite := testStrategies.TestSuite{}
	err := json.Unmarshal(buf.Bytes(), &testSuite)
	if err != nil {
		logrus.Error("Failed to unmarshall json body", err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	if !ValidateFileNameAndHeader(rw,req, testSuitePathDir,testSuite.Name){
		return
	}

	if !ValidateJsonWithSchema(buf.Bytes(), "testSuite_schema.json", "TestSuite") {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	if !FilePathExist(testSuitePathDir) {
		logrus.Error("Directory path does not exist", err)
		rw.WriteHeader(http.StatusBadRequest)
		return

	}

	if FilePathExist(fmt.Sprintf("%s%s.xml",testSuitePathDir, testSuite.Name)) {
		logrus.Error("File already exists")
		rw.WriteHeader(http.StatusBadRequest)
		return

	}

	if !testSuiteWriterXml(testSuite, testSuitePathDir +testSuite.Name+".xml") {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusCreated)
}

func ValidateJsonWithSchema(testSuite []byte, schemaName, structType string) bool {
	goPath := os.Getenv("GOPATH")
	schemaLoader := gojsonschema.NewReferenceLoader("file:///" + goPath + "/src/github.com/xtracdev/automated-perf-test/ui-src/src/assets/"+ schemaName)
	documentLoader := gojsonschema.NewBytesLoader(testSuite)
	logrus.Info(schemaLoader)
	result, error := gojsonschema.Validate(schemaLoader, documentLoader)

	if error != nil {
		return false
	}

	if !result.Valid() {
		logrus.Error("**** The "+structType+" document is not valid. see errors :")
		for _, desc := range result.Errors() {
			logrus.Error("- ", desc)
			return false
		}
	}
	logrus.Info("**** The "+structType+" document is valid *****")
	return true
}
