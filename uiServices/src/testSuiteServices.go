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

	"github.com/Sirupsen/logrus"
	"github.com/go-chi/chi"
	"github.com/xtracdev/automated-perf-test/testStrategies"
)

const (
	schemaFile = "testSuite_schema.json"
	structType = "testSuite"
)

type Suite struct {
	File         string                    `json:"file"`
	Name         string                    `json:"name"`
	Description  string                    `json:"description"`
	TestStrategy string                    `json:"testStrategy"`
	TestCases    []testStrategies.TestCase `json:"testCases"`
}

func TestSuiteCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}

func postTestSuites(rw http.ResponseWriter, req *http.Request) {
	testSuitePathDir := getPathHeader(req)
	buf := new(bytes.Buffer)
	buf.ReadFrom(req.Body)

	testSuite := testStrategies.TestSuite{}
	err := json.Unmarshal(buf.Bytes(), &testSuite)
	if err != nil {
		logrus.Error("Failed to unmarshall json body", err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := IsHeaderValid(testSuitePathDir); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := IsNameValid(testSuite.Name); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	if !validateJSONWithSchema(buf.Bytes(), schemaFile, structType) {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	if !FilePathExist(testSuitePathDir) {
		logrus.Error("Directory path does not exist", err)
		rw.WriteHeader(http.StatusBadRequest)
		return

	}

	filePath := fmt.Sprintf("%s%s.xml", testSuitePathDir, testSuite.Name)

	if FilePathExist(filePath) {
		logrus.Error("File already exists")
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	if !testSuiteWriterXml(testSuite, filePath) {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusCreated)
}

func putTestSuites(rw http.ResponseWriter, req *http.Request) {
	path := getPathHeader(req)
	testSuiteName := chi.URLParam(req, "testSuiteName")

	if err := IsHeaderValid(path); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := IsNameValid(testSuiteName); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
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

	if !validateJSONWithSchema(buf.Bytes(), schemaFile, structType) {
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

func deleteTestSuite(rw http.ResponseWriter, req *http.Request) {
	testSuitePathDir := getPathHeader(req)
	testSuiteName := chi.URLParam(req, "testSuiteName")

	if err := IsHeaderValid(testSuitePathDir); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := IsNameValid(testSuiteName); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	filepath := fmt.Sprintf("%s%s.xml", testSuitePathDir, testSuiteName)

	if _, err := os.Stat(filepath); err != nil {
		if os.IsNotExist(err) {
			logrus.Error("Test Suite File Not Found: ", err)
			rw.WriteHeader(http.StatusNotFound)
			return
		}
	}

	filePath := fmt.Sprintf("%s%s.xml", testSuitePathDir, testSuiteName)
	err := os.Remove(filePath)
	if err != nil {
		logrus.Errorf("Error deleting the file from filesystem: %s", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusNoContent)
}

func getTestSuite(rw http.ResponseWriter, req *http.Request) {
	testSuitePathDir := getPathHeader(req)
	testSuiteName := chi.URLParam(req, "testSuiteName")

	if err := IsHeaderValid(testSuitePathDir); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := IsNameValid(testSuiteName); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	file, err := os.Open(fmt.Sprintf("%s%s.xml", testSuitePathDir, testSuiteName))
	if err != nil {
		logrus.Error("Test Suite Name Not Found: " + testSuitePathDir + testSuiteName)
		rw.WriteHeader(http.StatusNotFound)
		return
	}
	defer file.Close()

	var testSuite testStrategies.TestSuite

	byteValue, err := ioutil.ReadAll(file)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		logrus.Error("Cannot Read File")
		return
	}

	err = xml.Unmarshal(byteValue, &testSuite)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		logrus.Error("Cannot Unmarshall from XML", err)
		return
	}

	testSuiteJSON, err := json.MarshalIndent(testSuite, "", "")
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		logrus.Error("Cannot marshall to JSON", err)
		return
	}

	rw.WriteHeader(http.StatusOK)
	rw.Write(testSuiteJSON)

}

func getAllTestSuites(rw http.ResponseWriter, req *http.Request) {
	testSuitePathDir := getPathHeader(req)
	if len(testSuitePathDir) <= 1 {
		logrus.Error("No file directory entered")
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	files, err := ioutil.ReadDir(testSuitePathDir)
	if err != nil {
		log.Fatal(err)
	}

	testSuites := make([]Suite, 0)

	for _, file := range files {
		if filepath.Ext(testSuitePathDir+file.Name()) == ".xml" {

			testSuite := new(testStrategies.TestSuite)

			filename := file.Name()

			file, err := os.Open(fmt.Sprintf("%s%s", testSuitePathDir, filename))
			if err != nil {
				logrus.Error("Cannot open file: ", filename)
			}

			byteValue, err := ioutil.ReadAll(file)
			if err != nil {
				logrus.Error("Cannot Read File: ", filename)
			}

			err = xml.Unmarshal(byteValue, testSuite)
			if err != nil {
				logrus.Error("Cannot Unmarshall: ", filename)

			}

			//if a Test Suite Name can't be assigned, it isn't a Test Suite object
			if testSuite.Name != "" {
				testSuites = append(testSuites, Suite{
					Name:         testSuite.Name,
					Description:  testSuite.Description,
					File:         filename,
					TestStrategy: testSuite.TestStrategy,
					TestCases:    testSuite.TestCases,
				})
			}
		}
	}

	json.NewEncoder(rw).Encode(testSuites)

	rw.WriteHeader(http.StatusOK)
}
