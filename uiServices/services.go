package services

import (
	"bytes"
	"encoding/json"
	"github.com/Sirupsen/logrus"
	"github.com/xeipuuv/gojsonschema"
	"github.com/xtracdev/automated-perf-test/perfTestUtils"
	"net/http"
	"os"
	"strings"
)

func ConfigsHandler(rw http.ResponseWriter, req *http.Request) {
	configPathDir := req.Header.Get("configPathDir")
	buf := new(bytes.Buffer)
	buf.ReadFrom(req.Body)

	if !validateJsonWithSchema(buf.Bytes()) {
		rw.WriteHeader(http.StatusBadRequest)
		return

	}

	config := perfTestUtils.Config{}
	err := json.Unmarshal(buf.Bytes(), &config)

	if err != nil {
		logrus.Error("Failed to unmarshall json body", err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	//error check to ensure file path ends with "\"
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
	if !writerXml(config, configPathDir) {

		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusCreated)

}

// exists returns whether the given file or directory exists or not
func FilePathExist(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func validateJsonWithSchema(config []byte) bool {
	goPath := os.Getenv("GOPATH")
	schemaLoader := gojsonschema.NewReferenceLoader("file:///" + goPath + "/src/github.com/xtracdev/automated-perf-test/schema.json")
	documentLoader := gojsonschema.NewBytesLoader(config)
	//logrus.Info(schemaLoader)
	result, error := gojsonschema.Validate(schemaLoader, documentLoader)

	if error != nil {
		return false
	}
	if result.Valid() {
		logrus.Info("**** The document is valid *****")

		return true
	}
	if !result.Valid() {
		logrus.Error("**** The document is not valid. see errors :")
		for _, desc := range result.Errors() {
			logrus.Error("- ", desc)
			return false
		}
	}
	return true

}
