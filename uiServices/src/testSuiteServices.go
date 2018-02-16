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
	"path/filepath"
	"log"
)

var schemaFile string = "testSuite_schema.json"
var structType string = "TestSuite"

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

	if !ValidateJsonWithSchema(buf.Bytes(), schemaFile, structType) {
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
		logrus.Errorf("%sdocument is not valid. see errors :", structType)
		for _, desc := range result.Errors() {
			logrus.Error("- ", desc)
			return false
		}
	}

	logrus.Infof("%s document is valid", structType)
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

	if !FilePathExist(testSuitePathDir) {
		logrus.Error("File path does not exist")
		rw.WriteHeader(http.StatusNotFound)
		return

	}

	if !ValidateJsonWithSchema(buf.Bytes(),schemaFile, structType) {
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

	file, err := os.Open(fmt.Sprintf("%s%s.xml", testSuitePathDir, testSuiteName))
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

func getAllTestSuites(rw http.ResponseWriter, req *http.Request){
	var testSuite testStrategies.TestSuite
	var filename string

	type Suite struct {
		File string `json:"file"`
		Name string `json:"name"`
		Description string `json:"description"`
	}

	var suite Suite

	type Suites[]Suite

	suites := Suites{}

	testSuitePathDir := getTestSuiteHeader(req)
	if len(testSuitePathDir) <= 1 {
		logrus.Error("No file directory entered")
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	files, err := ioutil.ReadDir(testSuitePathDir)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if filepath.Ext(testSuitePathDir + file.Name()) == ".xml" {
			filename = file.Name()

			file, err := os.Open(fmt.Sprintf("%s%s", testSuitePathDir, filename))
			if err != nil {
				logrus.Error("Cannot open file: ", filename)
			}

			byteValue, err := ioutil.ReadAll(file)
			if err != nil{
				logrus.Error("Cannot Read File: ", filename)
			}

			err = xml.Unmarshal(byteValue, &testSuite)
			if err != nil {
				logrus.Error("Cannot Unmarshall: ", filename)
			}

			suite.File = filename
			suite.Name = testSuite.Name
			suite.Description = testSuite.Description

			//if a Test Suite Name can't be assigned, it isn't a Test Suite object
			if suite.Name != "" {
				suites = append(suites, suite)
			}

			//ensure values are reset every iteration
			filename = ""
			testSuite.Name = ""
			testSuite.Description = ""

		}
	}

	json.NewEncoder(rw).Encode(suites)
	logrus.Println(suites)

	rw.WriteHeader(http.StatusOK)
}