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
				logrus.Error("Cannot open file: ", filename)
			}

			byteValue, err := ioutil.ReadAll(file)
			if err != nil {
				logrus.Error("Cannot Read File: ", filename)
			}

			err = xml.Unmarshal(byteValue, testCase)
			if err != nil {
				logrus.Error("Cannot Unmarshall: ", filename)

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
