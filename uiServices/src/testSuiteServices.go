package services

import (
	"bytes"
	"encoding/json"
	"github.com/Sirupsen/logrus"
	"github.com/xeipuuv/gojsonschema"
	"github.com/xtracdev/automated-perf-test/testStrategies"
	"net/http"
	"os"
	"github.com/go-chi/chi"
	"fmt"
	"strings"
	"encoding/xml"
	"io/ioutil"
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

	if FilePathExist(testSuitePathDir + testSuite.Name + ".xml") {
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


func putTestSuites(rw http.ResponseWriter, req *http.Request) {
	path := getTestSuiteHeader(req)
	testSuiteName := chi.URLParam(req, "testSuiteName")

	if !ValidateFileNameAndHeader(rw, req, path, testSuiteName) {
		return
	}

	testSuitePathDir := fmt.Sprintf("%s%s.xml", path, testSuiteName)
	buf := new(bytes.Buffer)
	buf.ReadFrom(req.Body)

	if !ValidateJsonWithSchema(buf.Bytes(),"testSuite_schema.json", "TestSuite") {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	testSuite := testStrategies.TestSuite{}
	err := json.Unmarshal(buf.Bytes(), &testSuite)
	if err != nil {
		logrus.Error("Cannot Unmarshall Json")
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	if !FilePathExist(testSuitePathDir) {
		logrus.Error("File path does not exist", err)
		rw.WriteHeader(http.StatusNotFound)
		return

	}

	if !testSuiteWriterXml(testSuite, testSuitePathDir) {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusNoContent)
}

func getTestSuite(rw http.ResponseWriter, req *http.Request){

	testSuitePathDir := getTestSuiteHeader(req)
	testSuiteName := chi.URLParam(req, "testSuiteName")

	ValidateFileNameAndHeader(rw,req,testSuitePathDir, testSuiteName)

	file, err := os.Open(fmt.Sprintf("%s%s", testSuitePathDir, testSuiteName))
	if err != nil {
		logrus.Error("Test Suite Name Not Found: "+testSuitePathDir + testSuiteName)
		rw.WriteHeader(http.StatusNotFound)
		return
	}
	defer file.Close()


	var testSuite testStrategies.TestSuite

	byteValue, err := ioutil.ReadAll(file)
	if err != nil{
		rw.WriteHeader(http.StatusInternalServerError)
		logrus.Error("Cannot Read File")
		return
	}

	err = xml.Unmarshal(byteValue, &testSuite)
	if err != nil{
		rw.WriteHeader(http.StatusInternalServerError)
		logrus.Error("Cannot Unmarshall")
		return
	}

	testSuiteJSON, err := json.MarshalIndent(testSuite,"","")
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		logrus.Error("Cannot Marshall")
		return
	}

	rw.WriteHeader(http.StatusOK)
	rw.Write(testSuiteJSON)
	logrus.Println(string(testSuiteJSON))

}
