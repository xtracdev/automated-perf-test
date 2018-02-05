package services

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/go-chi/chi"
	"github.com/xeipuuv/gojsonschema"
	"github.com/xtracdev/automated-perf-test/perfTestUtils"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func GetConfigHeader(req *http.Request) string {
	configPathDir := req.Header.Get("configPathDir")

	if !strings.HasSuffix(configPathDir, "/") {
		configPathDir = configPathDir + "/"
	}
	return configPathDir
}

func ValidateFileNameAndHeader(rw http.ResponseWriter, req *http.Request, header, name string) bool {

	if len(name) < 1 {
		logrus.Error("File Name is Empty")
		rw.WriteHeader(http.StatusBadRequest)
		return false
	}

	if len(header) <= 1 {
		logrus.Error("No Header Path Found")
		rw.WriteHeader(http.StatusBadRequest)
		return false
	}
	return true
}

func ConfigCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}

func postConfigs(rw http.ResponseWriter, req *http.Request) {
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

	if FilePathExist(configPathDir + config.APIName + ".xml") {
		logrus.Error("File already exists", err)
		rw.WriteHeader(http.StatusBadRequest)
		return

	}
	//Create file once checks are complete
	if !configWriterXml(config, configPathDir+config.APIName) {

		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusCreated)
}

func FilePathExist(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func validateJsonWithSchema(config []byte) bool {
	goPath := os.Getenv("GOPATH")
	schemaLoader := gojsonschema.NewReferenceLoader("file:///" + goPath + "/src/github.com/xtracdev/automated-perf-test/ui-src/src/assets/schema.json")
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
		logrus.Error("**** The document is not valid. see errors :")
		for _, desc := range result.Errors() {
			logrus.Error("- ", desc)
			return false
		}
	}
	return true
}

func getConfigs(rw http.ResponseWriter, req *http.Request) {

	configPathDir := GetConfigHeader(req)
	configName := chi.URLParam(req, "configName")

	if !ValidateFileNameAndHeader(rw, req, configPathDir, configName) {
		return
	}

	file, err := os.Open(fmt.Sprintf("%s%s.xml", configPathDir, configName))
	if err != nil {
		logrus.Error("Configuration Name Not Found: " + configPathDir + configName)
		rw.WriteHeader(http.StatusNotFound)
		return
	}

	defer file.Close()

	var config perfTestUtils.Config

	byteValue, err := ioutil.ReadAll(file)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		logrus.Error("Cannot Read File")
		return
	}

	err = xml.Unmarshal(byteValue, &config)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		logrus.Error("Cannot Unmarshall")
		return
	}

	configJson, err := json.MarshalIndent(config, "", "")
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		logrus.Error("Cannot Marshall")
		return
	}

	rw.WriteHeader(http.StatusOK)
	rw.Write(configJson)
	logrus.Println(string(configJson))

}

func putConfigs(rw http.ResponseWriter, req *http.Request) {
	path := GetConfigHeader(req)
	configName := chi.URLParam(req, "configName")

	if !ValidateFileNameAndHeader(rw, req, path, configName) {
		return
	}

	configPathDir := fmt.Sprintf("%s%s.xml", path, configName)
	buf := new(bytes.Buffer)
	buf.ReadFrom(req.Body)

	if !validateJsonWithSchema(buf.Bytes()) {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	config := perfTestUtils.Config{}
	err := json.Unmarshal(buf.Bytes(), &config)
	if err != nil {
		logrus.Error("Cannot Unmarshall Json")
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	if !FilePathExist(configPathDir) {
		logrus.Error("File path does not exist", err)
		rw.WriteHeader(http.StatusNotFound)
		return

	}

	if !configWriterXml(config, path+configName) {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusNoContent)
}
