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

	"github.com/Sirupsen/logrus"
	"github.com/go-chi/chi"
	"github.com/xtracdev/automated-perf-test/testStrategies"
)

const (
	testCaseSchema = "testCase_Schema.json"
	testCaseStruct = "testDefinition"
)

type Case struct {
	HttpMethod    string `json:"httpMethod"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	PreThinkTime  int64  `json:"preThinkTime"`
	PostThinkTime int64  `json:"postThinkTime"`
	ExecWeight    string `json:"execWeight"`
}

func TestCaseCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}

func postTestCase(rw http.ResponseWriter, req *http.Request) {
	testCasePathDir := getPathHeader(req)
	buf := new(bytes.Buffer)
	buf.ReadFrom(req.Body)

	testCase := testStrategies.TestDefinition{}
	err := json.Unmarshal(buf.Bytes(), &testCase)
	if err != nil {
		logrus.Error("Failed to unmarshall json body", err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := IsHeaderValid(testCasePathDir); err != nil {
		logrus.Error("No header path found", err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := IsNameValid(testCase.TestName); err != nil {
		logrus.Error("File name is empty", err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	if !FilePathExist(testCasePathDir) {
		logrus.Error("No directory path entered")
		rw.WriteHeader(http.StatusBadRequest)
		return

	}

	if FilePathExist(fmt.Sprintf("%s%s.xml", testCasePathDir, testCase.TestName)) {
		logrus.Error("File already exists")
		rw.WriteHeader(http.StatusBadRequest)
		return

	}

	if !validateJSONWithSchema(buf.Bytes(), testCaseSchema, testCaseStruct) {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	if !testCaseWriterXml(testCase, testCasePathDir+testCase.TestName+".xml") {
		logrus.Error("Error writing to the file")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusCreated)
}

func putTestCase(rw http.ResponseWriter, req *http.Request) {
	path := getPathHeader(req)
	testCaseName := chi.URLParam(req, "testCaseName")

	if err := IsHeaderValid(path); err != nil {
		logrus.Error("No header path found", err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := IsNameValid(testCaseName); err != nil {
		logrus.Error("File name is empty", err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	testCasePathDir := fmt.Sprintf("%s%s.xml", path, testCaseName)

	logrus.Error("testCasePathDir", testCasePathDir)

	if !FilePathExist(testCasePathDir) {
		logrus.Error("File path does not exist")
		rw.WriteHeader(http.StatusNotFound)
		return

	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(req.Body)

	if !validateJSONWithSchema(buf.Bytes(), testCaseSchema, testCaseStruct) {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	testCase := testStrategies.TestDefinition{}
	err := json.Unmarshal(buf.Bytes(), &testCase)

	if err != nil {
		logrus.Error("Error unmarshalling Json", err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	if len(testCase.TestName) < 1 {
		logrus.Error("No TestName Entered")
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	if !testCaseWriterXml(testCase, testCasePathDir) {
		logrus.Error("Error writing to the file")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusNoContent)
}

func getAllTestCases(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	testCasePathDir := getPathHeader(req)

	if err := IsPathDirValid(testCasePathDir); err != nil {
		logrus.Error("Path Directory is not valid", err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	files, err := ioutil.ReadDir(testCasePathDir)
	if err != nil {
		logrus.Error("Error reading the directory ", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	testCases := make([]Case, 0)

	for _, file := range files {
		if filepath.Ext(testCasePathDir+file.Name()) == ".xml" {

			testCase := new(testStrategies.TestDefinition)

			filename := file.Name()

			file, err := os.Open(fmt.Sprintf("%s%s", testCasePathDir, filename))
			if err != nil {
				logrus.Error("Cannot open file: ", filename, err)
				logrus.Error("Error opening the file: " + filename)
				continue
			}

			byteValue, err := ioutil.ReadAll(file)
			if err != nil {
				logrus.Error("Cannot Read File: ", filename, err)
				logrus.Error("Error reading the file: " + filename)
				continue
			}

			err = xml.Unmarshal(byteValue, testCase)
			if err != nil {
				logrus.Error("Cannot Unmarshall: ", filename, err)

				logrus.Error("Error unmarshalling the file: " + filename)
				continue
			}

			if testCase.TestName != "" {
				testCases = append(testCases, Case{
					Name:          testCase.TestName,
					Description:   testCase.Description,
					HttpMethod:    testCase.HTTPMethod,
					PostThinkTime: testCase.PostThinkTime,
					PreThinkTime:  testCase.PreThinkTime,
					ExecWeight:    testCase.ExecWeight,
				})
			}
		}
	}

	err = json.NewEncoder(rw).Encode(testCases)
	if err != nil {
		logrus.Error("Error encoding Test Cases", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
}

func getTestCase(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	testCasePathDir := getPathHeader(req)
	testCaseName := chi.URLParam(req, "testCaseName")

	if err := IsHeaderValid(testCasePathDir); err != nil {
		logrus.Error("No header path found", err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := IsNameValid(testCaseName); err != nil {
		logrus.Error("File name is empty", err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	filePath := fmt.Sprintf("%s%s.xml", testCasePathDir, testCaseName)
	logrus.Info("FilePath: ", filePath)

	if _, err := os.Stat(fmt.Sprintf("%s%s.xml", testCasePathDir, testCaseName)); err != nil {
		if os.IsNotExist(err) {
			logrus.Error("Test Case File Not Found: " + testCaseName)
			rw.WriteHeader(http.StatusNotFound)
			return
		}
	}

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
		logrus.Error("Error reading from file", err)
		rw.WriteHeader(http.StatusInternalServerError)
		logrus.Error("Cannot Read File")
		return
	}

	err = xml.Unmarshal(byteValue, &testCase)
	if err != nil {
		logrus.Error("Error unmarshalling from XML", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	testSuiteJSON, err := json.MarshalIndent(testCase, "", "")
	if err != nil {
		logrus.Error("Error marshalling to JSON", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
	rw.Write(testSuiteJSON)

}

func deleteTestCase(rw http.ResponseWriter, req *http.Request) {
	testCasePathDir := getPathHeader(req)
	testCaseName := chi.URLParam(req, "testCaseName")

	if err := IsHeaderValid(testCasePathDir); err != nil {
		logrus.Error("No header path found", err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := IsNameValid(testCaseName); err != nil {
		logrus.Error("File name is empty", err)
		rw.WriteHeader(http.StatusBadRequest)
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
		logrus.Println("Error deleting the file", err)
		rw.WriteHeader(http.StatusNoContent)
		return

	}

	rw.WriteHeader(http.StatusNoContent)

}
func deleteAllTestCases(rw http.ResponseWriter, req *http.Request) {
	testCasePathDir := getPathHeader(req)

	if err := IsHeaderValid(testCasePathDir); err != nil {
		logrus.Error("No heared path found", err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	files, err := ioutil.ReadDir(testCasePathDir)
	if err != nil {
		logrus.Error("Error reading the directory ", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	for _, file := range files {
		filepath := fmt.Sprintf("%s%s", testCasePathDir, file.Name())

		err := os.Remove(filepath)
		if err != nil {
			logrus.Error("Error removing the files from directory", err)
			rw.WriteHeader(http.StatusNoContent)
			return
		}
	}
	rw.WriteHeader(http.StatusOK)
}
