package services

import (
	"encoding/json"
	"net/http"
	"bytes"
	"github.com/Sirupsen/logrus"
	"github.com/xtracdev/automated-perf-test/perfTestUtils"
	"os"
)


func configsHandler(rw http.ResponseWriter, req *http.Request) {

	configPathDir := req.Header.Get("configPathDir")

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

	if FilePathExist(configPathDir)== false{
		logrus.Error("File path does not exist", err)
		rw.WriteHeader(http.StatusNotFound)
		return

	}

	IsSuccessful := writerXml(config, configPathDir)

	if IsSuccessful {
		rw.WriteHeader(http.StatusOK)
		return
	} else {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

}

// exists returns whether the given file or directory exists or not
func FilePathExist(path string) (bool) {
	FileExist := true
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		FileExist = false
	}
return FileExist
}

