package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/xeipuuv/gojsonschema"
	"github.com/xtracdev/automated-perf-test/perfTestUtils"
	"net/http"
	"os"
	"strings"
)

func configsHandler(rw http.ResponseWriter, req *http.Request) {
	configPathDir := req.Header.Get("configPathDir")

	config := perfTestUtils.Config{}
	buf := new(bytes.Buffer)
	buf.ReadFrom(req.Body)

	defer req.Body.Close()
	err := json.Unmarshal(buf.Bytes(), &config)

	if len(configPathDir) < 1 {
		logrus.Error("File path is length too short", err)
		rw.WriteHeader(http.StatusBadRequest)
		return

	}
	//error check to ensure file path ends with "\"
	if !strings.HasSuffix(configPathDir, "/") {
		configPathDir = configPathDir + "/"
	}

	if err != nil {
		logrus.Error("Failed to unmarshall json body", err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	if !FilePathExist(configPathDir) {
		logrus.Error("File path does not exist", err)
		rw.WriteHeader(http.StatusBadRequest)
		return

	}
	if validateJsonWithSchema(buf.Bytes()) {
		isSuccessful := writerXml(config, configPathDir)
		if !isSuccessful {
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
		rw.WriteHeader(http.StatusCreated)
		return

	}

	return

}

// exists returns whether the given file or directory exists or not
func FilePathExist(path string) bool {
	FileExist := true
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		FileExist = false

	}
	return FileExist
}

func validateJsonWithSchema(config []byte) bool {
	goPath := os.Getenv("GOPATH")
	fmt.Println(goPath)
	schemaLoader := gojsonschema.NewReferenceLoader("file:///" + goPath + "/src/github.com/xtracdev/automated-perf-test/schema.json")
	documentLoader := gojsonschema.NewBytesLoader(config)
	logrus.Info(schemaLoader)
	result, error := gojsonschema.Validate(schemaLoader, documentLoader)
	if error != nil {
		return false
	}
	if result.Valid() {
		logrus.Info("**** The document is valid *****")

		return true
	}
	if !result.Valid() {
		fmt.Println("**** The document is not valid. see errors :\n ****")
		for _, desc := range result.Errors() {
			logrus.Info("- %s\n", desc)
			return false
		}
	}
	return true

}
