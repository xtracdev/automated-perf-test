package services

import (
	"github.com/Sirupsen/logrus"
	"github.com/xeipuuv/gojsonschema"
	"os"
	"strings"
	"encoding/json"
	"github.com/xtracdev/automated-perf-test/testStrategies"
	"bytes"
	"net/http"
)

func TestSuiteCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})

}

func postTestSuites(rw http.ResponseWriter, req *http.Request) {
	configPathDir := req.Header.Get("configPathDir")
	buf := new(bytes.Buffer)
	buf.ReadFrom(req.Body)

	testSuite := testStrategies.TestSuite{}
	err := json.Unmarshal(buf.Bytes(), &testSuite)

	if !validateTestSuiteJsonWithSchema(buf.Bytes()) {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	if err != nil {
		logrus.Error("Failed to unmarshall json body", err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	if !strings.HasSuffix(configPathDir, "/") {
		configPathDir = configPathDir + "/"
	}

	if len(configPathDir) <= 1 {
		logrus.Error("File path is length too short", err)
		rw.WriteHeader(http.StatusBadRequest)
		return

	}

	if !FilePathExist(configPathDir) {
		logrus.Error("File path does not exist", err)
		rw.WriteHeader(http.StatusBadRequest)
		return

	}
	//Create file once checks are complete
	if !testSuiteWriterXml(testSuite, configPathDir + testSuite.Name) {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusCreated)
}

func validateTestSuiteJsonWithSchema(testSuite []byte) bool {
	goPath := os.Getenv("GOPATH")
	schemaLoader := gojsonschema.NewReferenceLoader("file:///" + goPath + "/src/github.com/xtracdev/automated-perf-test/ui-src/src/assets/testSuite_schema.json")
	documentLoader := gojsonschema.NewBytesLoader(testSuite)
	logrus.Info(schemaLoader)
	result, error := gojsonschema.Validate(schemaLoader, documentLoader)

	if error != nil {
		return false
	}
	if result.Valid() {
		logrus.Info("**** The TestSuite document is valid *****")

		return true
	}
	if !result.Valid() {
		logrus.Error("**** The TestSuite document is not valid. see errors :")
		for _, desc := range result.Errors() {
			logrus.Error("- ", desc)
			return false
		}
	}
	return true
}