package services

import (
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
	"github.com/xtracdev/automated-perf-test/testStrategies"
	"github.com/go-chi/chi"
)

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
		if filepath.Ext(testCasePathDir +file.Name()) == ".xml" {

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
	if err != nil{
		logrus.Error("Could not enocde Test Cases")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
}

func getTestCase(rw http.ResponseWriter, req *http.Request) {

	testCasePathDir := getTestCaseHeader(req)
	testCaseName := chi.URLParam(req, "testCaseName")

	ValidateFileNameAndHeader(rw, req, testCasePathDir, testCaseName)

	file, err := os.Open(fmt.Sprintf("%s%s.xml", testCasePathDir, testCaseName))
	if err != nil {
		logrus.Error("Test Case Name Not Found: " + testCasePathDir + testCaseName)
		rw.WriteHeader(http.StatusNotFound)
		return
	}
	defer file.Close()

	var testCase testStrategies.TestDefinition

	byteValue, err := ioutil.ReadAll(file)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		logrus.Error("Cannot Read File")
		return
	}

	err = xml.Unmarshal(byteValue, &testCase)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		logrus.Error("Cannot Unmarshall from XML")
		return
	}

	testSuiteJSON, err := json.MarshalIndent(testCase, "", "")
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		logrus.Error("Cannot marshall to JSON")
		return
	}

	rw.WriteHeader(http.StatusOK)
	rw.Write(testSuiteJSON)

}

func deleteTestCase(rw http.ResponseWriter, req *http.Request){
	testCasePathDir := getTestCaseHeader(req)
	testCaseName := chi.URLParam(req, "testCaseName")
	ValidateFileNameAndHeader(rw, req, testCasePathDir, testCaseName)

	filepath := fmt.Sprintf("%s%s.xml", testCasePathDir, testCaseName)

	if _, err := os.Stat(filepath); err != nil {
		if os.IsNotExist(err) {
			logrus.Println("File Not Found")
			rw.WriteHeader(http.StatusNotFound)
			return
		}
	}

	err := os.Remove(filepath)
	if  err != nil{
		logrus.Println("File was not deleted")
		rw.WriteHeader(http.StatusInternalServerError)
		return

	}

	rw.WriteHeader(http.StatusNoContent)

}