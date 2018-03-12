package services

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/go-chi/chi"
	"github.com/xtracdev/automated-perf-test/perfTestUtils"
)

var schemaFilename string = "schema.json"
var structName string = "Config"

func getConfigHeader(req *http.Request) string {
	configPathDir := req.Header.Get("configPathDir")

	if !strings.HasSuffix(configPathDir, "/") {
		configPathDir = configPathDir + "/"
	}
	return configPathDir
}

func IsHeaderValid(header string) error {
	if len(header) <= 1 {
		logrus.Error("No Header Path Found")
		return fmt.Errorf("No Header Path Found")
	}
	return nil
}

func IsNameValid(name string) error {
	if len(name) < 1 {
		logrus.Error("File Name is Empty")
		return fmt.Errorf("File Name is Empty")
	}
	return nil
}

func IsPathDirValid(name string, rw http.ResponseWriter) bool {
	if len(name) <= 1 {
		logrus.Error("No file directory entered")
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
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	configPathDir := req.Header.Get("configPathDir")
	buf := new(bytes.Buffer)
	buf.ReadFrom(req.Body)

	if !ValidateJsonWithSchema(buf.Bytes(), "schema.json", "Configurations") {
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

	if !configWriterXML(config, configPathDir+config.APIName+".xml") {

		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusCreated)
}

func FilePathExist(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func getConfigs(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	configPathDir := getConfigHeader(req)
	configName := chi.URLParam(req, "configName")

	if err := IsHeaderValid(configPathDir); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := IsNameValid(configName); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
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
	path := getConfigHeader(req)
	configName := chi.URLParam(req, "configName")

	if err := IsHeaderValid(path); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := IsNameValid(configName); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	configPathDir := fmt.Sprintf("%s%s.xml", path, configName)
	buf := new(bytes.Buffer)
	buf.ReadFrom(req.Body)

	if !ValidateJsonWithSchema(buf.Bytes(), schemaFilename, structName) {
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

	if !configWriterXml(config, configPathDir) {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusNoContent)
}

func putConfigFileName(rw http.ResponseWriter, req *http.Request) {
	path := getConfigHeader(req)
	configFileName := chi.URLParam(req, "configFileName")
	newConfigFileName := chi.URLParam(req, "newConfigFileName")

	if err := IsHeaderValid(path); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := IsNameValid(configFileName); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := IsNameValid(newConfigFileName); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	configFilePath := fmt.Sprintf("%s%s.xml", path, configFileName)
	newConfigFilePath := fmt.Sprintf("%s%s.xml", path, newConfigFileName)

	err := os.Rename(configFilePath, newConfigFilePath)
	if err != nil {
		logrus.Error("File was not updated", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	rw.WriteHeader(http.StatusNoContent)
}
