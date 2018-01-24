package services

import (
    "bytes"
    "encoding/json"
    "net/http"
    "os"
    "strings"

    "github.com/Sirupsen/logrus"
    "github.com/xeipuuv/gojsonschema"
    "github.com/xtracdev/automated-perf-test/perfTestUtils"
    "fmt"
    "encoding/xml"
    "io/ioutil"
    "github.com/go-chi/chi"
)

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

func getConfigs(rw http.ResponseWriter, req *http.Request){

    configPathDir := req.Header.Get("configPathDir")
    configName := chi.URLParam(req, "configName")

    if !strings.HasSuffix(configPathDir, "/"){
        configPathDir = configPathDir + "/"
    }

    file, err := os.Open(configPathDir + configName)

    if len(configPathDir) <= 1{
        rw.WriteHeader(http.StatusBadRequest)
        logrus.Println("No Header Path Found")
        return
    }

    if err != nil {
        logrus.Println("Configuration Name Not Found: "+configPathDir + configName)
        rw.WriteHeader(http.StatusNotFound)
        return
    }

    defer file.Close()

    var config perfTestUtils.Config

    byteValue, err := ioutil.ReadAll(file)
    xml.Unmarshal(byteValue, &config)
    if err != nil{
        logrus.Println("Cannot Unmarshall")
        return
    }

    configJson, err := json.MarshalIndent(config,"","")
    if err != nil {
        logrus.Println("Cannot Marshall")
        return
    }

    rw.WriteHeader(http.StatusOK)
    fmt.Println(string(configJson))

}