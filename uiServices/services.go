package services

import (
	"bytes"
	"encoding/json"
	"github.com/Sirupsen/logrus"
	"github.com/xtracdev/automated-perf-test/perfTestUtils"
	"net/http"
	"os"
	"strings"
	"fmt"
	"github.com/xeipuuv/gojsonschema"
)

func configsHandler(rw http.ResponseWriter, req *http.Request) {

		configPathDir := req.Header.Get("configPathDir")
		//error check to ensure file path ends with "\"
		if !strings.HasSuffix(configPathDir, "/") {
			configPathDir = configPathDir + "/"
		}


		config := perfTestUtils.Config{}
		buf := new(bytes.Buffer)
		buf.ReadFrom(req.Body)
		defer req.Body.Close()

		err := json.Unmarshal(buf.Bytes(), &config)

	if err != nil {
		logrus.Error("Failed to unmarshall json body", err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	if FilePathExist(configPathDir) == false {
		logrus.Error("File path does not exist", err)
		rw.WriteHeader(http.StatusNotFound)
		return

	}

	if validdateJsonWithSchema(buf.Bytes()) == true {
		isSuccessful := writerXml(config, configPathDir)
		if !isSuccessful {
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
		rw.WriteHeader(http.StatusCreated)
		return
	}

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

func validdateJsonWithSchema(config []byte) bool {
	goPath := os.Getenv("GOPATH")
	schemaLoader := gojsonschema.NewReferenceLoader("file:///" + goPath + "/src/github.com/xtracdev/automated-perf-test/schema.json")
	documentLoader := gojsonschema.NewBytesLoader(config)
	logrus.Print(schemaLoader)
	result, error := gojsonschema.Validate(schemaLoader, documentLoader)
	logrus.Print(error)
	if error != nil {
		return false
	}
	if result.Valid() {
		fmt.Printf("**** The document is valid *****")

		return true
	} else {
		logrus.Println("**** The document is not valid. see errors :\n ****")
		for _, desc := range result.Errors() {
			fmt.Printf("- %s\n", desc)
			return false
		}
	}
	return true

}