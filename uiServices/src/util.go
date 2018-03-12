package services

import (
	"net/http"
	"os"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/xeipuuv/gojsonschema"
)

func getPathHeader(req *http.Request) string {
	pathHeader := req.Header.Get("path")

	if !strings.HasSuffix(pathHeader, "/") {
		pathHeader = pathHeader + "/"
	}
	return pathHeader
}

func validateJSONWithSchema(testSuite []byte, schemaName, structType string) bool {
	goPath := os.Getenv("GOPATH")
	schemaLoader := gojsonschema.NewReferenceLoader("file:///" + goPath + "/src/github.com/xtracdev/automated-perf-test/ui-src/src/assets/" + schemaName)
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
func FilePathExist(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
