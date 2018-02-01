package services

import (
       "fmt"

       "github.com/go-chi/chi"
       "github.com/stretchr/testify/assert"
       "github.com/xtracdev/automated-perf-test/perfTestUtils"
       "net/http"
       "net/http/httptest"
       "os"
       "strings"
       "testing"
)

const validJson = `{
        "apiName": "ServiceTestConfig",
       "targetHost": "localhost",
       "targetPort": "9191",
       "numIterations": 1000,
       "allowablePeakMemoryVariance": 30,
       "allowableServiceResponseTimeVariance": 30,
       "testCaseDir": "./definitions/testCases",
       "testSuiteDir": "./definitions/testSuites",
        "baseStatsOutputDir": "./envStats",
       "reportOutputDir": "./report",
       "concurrentUsers": 50,
       "testSuite": "Default-1",
       "memoryEndpoint": "/alt/debug/vars",
       "requestDelay": 5000,
       "TPSFreq": 30,
       "rampUsers": 5,
       "rampDelay": 15
       }`

const invalidJsonMissingFields = `{
        "apiName": "ServiceTestConfi",
       "targetHost": "localhost",
       "targetPort": "9191",
       "numIterations": 1000,
       "allowablePeakMemoryVariance": -1,
       "allowableServiceResponseTimeVariance": 30,
       "testCaseDir": "./definitions/testCases",
       "testSuiteDir": "./definitions/testSuites",
        "baseStatsOutputDir": "./envStats",
       "reportOutputDir": "./report",
       "concurrentUsers": 50,
       "testSuite": "",
       "memoryEndpoint": "/alt/debug/vars",
       "requestDelay": 5000,
       "TPSFreq": 30,
       "rampUsers": 5,
       "rampDelay": 0
       }`

const invalidJson = `{
        "apiName"://*()()(),
       "targetHost": 0,
       "targetPort": 0,
       "numIterations": "x",
       "allowablePeakMemoryVariance": 30,
       "allowableServiceResponseTimeVariance": 30,
       "testCaseDir": "./definitions/testCases",
       "testSuiteDir": "./definitions/testSuites",
        "baseStatsOutputDir": "./envStats",
       "reportOutputDir": "./report",
       "concurrentUsers": 50,
       "testSuite": "suiteFileName.xml",
       "memoryEndpoint": "/alt/debug/vars",
       "requestDelay": 5000,
       "TPSFreq": 30,
       "rampUsers": 5,
       "rampDelay": 15
       }`

func TestFilePathExist(t *testing.T) {
       path := os.Getenv("GOPATH") + "/src/github.com/xtracdev/automated-perf-test/uiServices/test/"
       actual := false
       fmt.Println(path)
       actual = FilePathExist(path)
       expected := true
       assert.Equal(t, expected, actual)
}

func TestFilePathDoesNotExist(t *testing.T) {
       path := "((((((("
       actual := FilePathExist(path)
       expected := false
       assert.Equal(t, expected, actual)
}

func TestInvalidJsonPostMissingRequiredField(t *testing.T) {
       r := chi.NewRouter()
       r.Mount("/", GetIndexPage())

       reader := strings.NewReader(invalidJsonMissingFields)
       r.HandleFunc("/configs", postConfigs)

       filePath := os.Getenv("GOPATH") + "/src/github.com/xtracdev/automated-perf-test/uiServices/test/"
       request, err := http.NewRequest(http.MethodPost, "/configs", reader)
       request.Header.Set("configPathDir", filePath)

       w := httptest.NewRecorder()
       r.ServeHTTP(w, request)

       if err != nil {
              t.Error(err)
       }

       if w.Code != http.StatusBadRequest {
              t.Error("TestValidJsonPost. Expected:", http.StatusBadRequest, " Got:", w.Code, "  Error. Did not succesfully post")
       }

}

func TestValidJsonPost(t *testing.T) {
       r := chi.NewRouter()
       r.Mount("/", GetIndexPage())

       reader := strings.NewReader(validJson)
       r.HandleFunc("/configs", postConfigs)

       filePath := os.Getenv("GOPATH") + "/src/github.com/xtracdev/automated-perf-test/uiServices/test"
       request, err := http.NewRequest(http.MethodPost, "/configs", reader)
       request.Header.Set("configPathDir", filePath)

       w := httptest.NewRecorder()
       r.ServeHTTP(w, request)

       if err != nil {
              t.Error(err)
       }

       if w.Code != http.StatusCreated {
              t.Error("TestValidJsonPost. Expected:", http.StatusCreated, " Got:", w.Code, "  Error. Did not succesfully post")
       }
}

func TestPostWithInvalidHeader(t *testing.T) {
       r := chi.NewRouter()
       r.Mount("/", GetIndexPage())

       reader := strings.NewReader(validJson)
       r.HandleFunc("/configs", postConfigs)

       filePath := "xxxxxx"
       request, err := http.NewRequest(http.MethodPost, "/configs", reader)
       request.Header.Set("configPathDir", filePath)

       w := httptest.NewRecorder()
       r.ServeHTTP(w, request)

       if err != nil {
              t.Error(err)
       }

       if w.Code != http.StatusBadRequest {
              t.Error("TestValidJsonPost. Expected:", http.StatusBadRequest, " Got:", w.Code, "  Error. Did not succesfully post")
       }
}

func TestInvalidJsonPost(t *testing.T) {
       r := chi.NewRouter()
       r.Mount("/", GetIndexPage())

       reader := strings.NewReader(invalidJson)
       r.HandleFunc("/configs", postConfigs)

       request, err := http.NewRequest(http.MethodPost, "/configs", reader)

       w := httptest.NewRecorder()
       r.ServeHTTP(w, request)

       if err != nil {
              t.Error(err)
       }

       if w.Code != http.StatusBadRequest {
              t.Error("TestInvalidJsonPost.  Expected:", http.StatusBadRequest, " Got:", w.Code, "Error. Did not succesfully post ")
       }
}

func TestWhenConfigPathDirEmpty(t *testing.T) {
       r := chi.NewRouter()
       r.Mount("/", GetIndexPage())

       reader := strings.NewReader(validJson)
       r.HandleFunc("/configs", postConfigs)

       request, err := http.NewRequest(http.MethodPost, "/configs", reader)

       w := httptest.NewRecorder()
       r.ServeHTTP(w, request)

       if err != nil {
              t.Error(err)
       }

       if w.Code != http.StatusBadRequest {
              t.Error("TestWhenConfigPathDirEmpty.  Expected:", http.StatusBadRequest, " Got:", w.Code, "Error. ConfigPathDir is Empty ")
       }
}

func TestInvalidURL(t *testing.T) {
       pt := perfTestUtils.Config{}
       writerXml(pt, "/path/xxx")
}

func TestSuccessfulGet(t *testing.T) {
       r := chi.NewRouter()
       r.Mount("/", GetIndexPage())

       // create file to GET
       reader := strings.NewReader(validJson)
       r.HandleFunc("/configs", postConfigs)

       filePath := os.Getenv("GOPATH") + "/src/github.com/xtracdev/automated-perf-test/uiServices/test/"
       request, err := http.NewRequest(http.MethodPost, "/configs", reader)
       request.Header.Set("configPathDir", filePath)

       w := httptest.NewRecorder()
       r.ServeHTTP(w, request)

       //prepare GET request
       filePath = os.Getenv("GOPATH") + "/src/github.com/xtracdev/automated-perf-test/uiServices/test/"
       request, err = http.NewRequest(http.MethodGet, "/configs/ServiceTestConfig", nil)

       request.Header.Set("configPathDir", filePath)
       request.Header.Get("configPathDir")

       w = httptest.NewRecorder()
       r.ServeHTTP(w, request)

       if err != nil {
              t.Error(err)
       }

       if w.Code != http.StatusOK {
              t.Error("TestSuccessfulGET. Expected:", http.StatusOK, " Got:", w.Code, "  Error. Did not succesfully get")
       }
}

func TestSuccessfulGetPathWihoutSlash(t *testing.T) {
       r := chi.NewRouter()
       r.Mount("/", GetIndexPage())

       r.HandleFunc("/configs", getConfigs)
       //no slash at end of filepath header
       filePath := os.Getenv("GOPATH") + "/src/github.com/xtracdev/automated-perf-test/uiServices/test"
       request, err := http.NewRequest(http.MethodGet, "/configs/ServiceTestConfig", nil)

       request.Header.Set("configPathDir", filePath)
       request.Header.Get("configPathDir")

       w := httptest.NewRecorder()
       r.ServeHTTP(w, request)

       if err != nil {
              t.Error(err)
       }

       if w.Code != http.StatusOK {
              t.Error("Test Get Path ends with backslash. Expected:", http.StatusOK, " Got:", w.Code, "  Error. Did not succesfully get")
       }
}

func TestGetNoHeaderPath(t *testing.T) {
       r := chi.NewRouter()
       r.Mount("/", GetIndexPage())

       r.HandleFunc("/configs", getConfigs)

       filePath := ""
       request, err := http.NewRequest(http.MethodGet, "/configs/serviceTestConfig", nil)

       request.Header.Set("configPathDir", filePath)
       request.Header.Get("configPathDir")

       w := httptest.NewRecorder()
       r.ServeHTTP(w, request)

       if err != nil {
              t.Error(err)
       }

       if w.Code != http.StatusBadRequest {
              t.Error("Test No-Header-Get. Expected:", http.StatusBadRequest, " Got:", w.Code)
       }
}

func TestGetFileNotFound(t *testing.T) {
       r := chi.NewRouter()
       r.Mount("/", GetIndexPage())

       r.HandleFunc("/configs", getConfigs)

       filePath := os.Getenv("GOPATH") + "/src/github.com/xtracdev/automated-perf-test/uiServices/test/"
       request, err := http.NewRequest(http.MethodGet, "/configs/xxx.java", nil)

       request.Header.Set("configPathDir", filePath)
       request.Header.Get("configPathDir")

       w := httptest.NewRecorder()
       r.ServeHTTP(w, request)

       if err != nil {
              t.Error(err)
       }

       if w.Code != http.StatusNotFound {
              t.Error("Test File Not Found. Expected:", http.StatusNotFound, " Got:", w.Code)
       }
}

func TestValidJsonPut(t *testing.T) {
       r := chi.NewRouter()
       r.Mount("/", GetIndexPage())

       reader := strings.NewReader(validJson)
       r.HandleFunc("/configs", putConfigs)

      filePath := os.Getenv("GOPATH") + "/src/github.com/xtracdev/automated-perf-test/uiServices/test/"
       request, err := http.NewRequest(http.MethodPut, "/configs/ServiceTestConfig", reader)
       request.Header.Set("configPathDir", filePath)

       w := httptest.NewRecorder()
       r.ServeHTTP(w, request)

       if err != nil {
              t.Error(err)
       }

       assert.Equal(t, w.Code, http.StatusNoContent, "Did Not successfully Update")
}

func TestMissingFieldPut(t *testing.T) {
       r := chi.NewRouter()
       r.Mount("/", GetIndexPage())

       reader := strings.NewReader(invalidJsonMissingFields)
       r.HandleFunc("/configs", putConfigs)

       filePath := os.Getenv("GOPATH") + "/src/github.com/xtracdev/automated-perf-test/uiServices/test/"
       request, err := http.NewRequest(http.MethodPut, "/configs/ServiceTestConfig", reader)
       request.Header.Set("configPathDir", filePath)

       w := httptest.NewRecorder()
       r.ServeHTTP(w, request)

       if err != nil {
              t.Error(err)
       }

       assert.Equal(t, w.Code, http.StatusBadRequest, "Sucessfully updated. Field Should be missing so update shouldn't occur")
}

func TestInvalidJsonPut(t *testing.T) {
       r := chi.NewRouter()
       r.Mount("/", GetIndexPage())

       reader := strings.NewReader(invalidJson)
       r.HandleFunc("/configs", putConfigs)

       filePath := os.Getenv("GOPATH") + "/src/github.com/xtracdev/automated-perf-test/uiServices/test/"
       request, err := http.NewRequest(http.MethodPut, "/configs/ServiceTestConfig", reader)
       request.Header.Set("configPathDir", filePath)

       w := httptest.NewRecorder()
       r.ServeHTTP(w, request)

       if err != nil {
              t.Error(err)
       }

       assert.Equal(t, w.Code, http.StatusBadRequest, "Sucessfully updated. Field data type should have been incorrect so update should occur")
}

func TestInvalidUrlPut(t *testing.T) {
       r := chi.NewRouter()
       r.Mount("/", GetIndexPage())

       reader := strings.NewReader(validJson)
       r.HandleFunc("/configs", putConfigs)

       filePath := os.Getenv("GOPATH") + "/src/github.com/xtracdev/automated-perf-test/uiServices/test/"
       request, err := http.NewRequest(http.MethodPut, "/configs/xxx", reader)
       request.Header.Set("configPathDir", filePath)

       w := httptest.NewRecorder()
       r.ServeHTTP(w, request)

       if err != nil {
              t.Error(err)
       }

       assert.Equal(t, w.Code, http.StatusConflict, "Sucessfully updated. Should have have worked using /configs/xxx")
}

func TestNoUrlPut(t *testing.T) {
       r := chi.NewRouter()
       r.Mount("/", GetIndexPage())

       reader := strings.NewReader(validJson)
       r.HandleFunc("/configs", putConfigs)

       filePath := os.Getenv("GOPATH") + "/src/github.com/xtracdev/automated-perf-test/uiServices/test/"
       request, err := http.NewRequest(http.MethodPut, "", reader)
       request.Header.Set("configPathDir", filePath)

       w := httptest.NewRecorder()
       r.ServeHTTP(w, request)

       if err != nil {
              t.Error(err)
       }

       assert.Equal(t, w.Code, http.StatusNotFound, "Sucessfully updated. Should not have worked with no URL")
}

func TestSuccessfulPutWithNoPathSlash(t *testing.T) {
       r := chi.NewRouter()
       r.Mount("/", GetIndexPage())

       reader := strings.NewReader(validJson)
       r.HandleFunc("/configs", putConfigs)

       filePath := os.Getenv("GOPATH") + "/src/github.com/xtracdev/automated-perf-test/uiServices/test"
       request, err := http.NewRequest(http.MethodPut, "/configs/ServiceTestConfig", reader)
       request.Header.Set("configPathDir", filePath)

       w := httptest.NewRecorder()
       r.ServeHTTP(w, request)

       if err != nil {
              t.Error(err)
       }

       assert.Equal(t, w.Code, http.StatusNoContent, "Did not update. Should have added '/' to path")
}
func TestNoPathPut(t *testing.T) {
       r := chi.NewRouter()
       r.Mount("/", GetIndexPage())

       reader := strings.NewReader(validJson)
       r.HandleFunc("/configs", putConfigs)

       filePath := ""
       request, err := http.NewRequest(http.MethodPut, "/configs/ServiceTestConfig", reader)
       request.Header.Set("configPathDir", filePath)

       w := httptest.NewRecorder()
       r.ServeHTTP(w, request)

       if err != nil {
              t.Error(err)
       }

       assert.Equal(t, w.Code, http.StatusBadRequest, "Successfully updated. Should not have worked due to no filepath")
}

func TestNoFileNamePut(t *testing.T) {
       r := chi.NewRouter()
       r.Mount("/", GetIndexPage())

       reader := strings.NewReader(validJson)
       r.HandleFunc("/configs", putConfigs)

       filePath := os.Getenv("GOPATH") + "/src/github.com/xtracdev/automated-perf-test/uiServices/test"
       request, err := http.NewRequest(http.MethodPut, "/configs", reader)
       request.Header.Set("configPathDir", filePath)

       w := httptest.NewRecorder()
       r.ServeHTTP(w, request)

       if err != nil {
              t.Error(err)
       }

       assert.Equal(t, w.Code, http.StatusNotFound, "Successfully updated. Should not have worked due to no file name given")
}


From: Quinn, Frank 
Sent: 01 February 2018 14:12
To: Rooney, Angela
Subject: arger

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

func getConfigHeader(req *http.Request) string {
       configPathDir := req.Header.Get("configPathDir")

       if !strings.HasSuffix(configPathDir, "/") {
              configPathDir = configPathDir + "/"
       }
       return configPathDir
}

func validateFileNameAndHeader(rw http.ResponseWriter, req *http.Request, header, name string) bool {

       if len(name) <= 1 {
              logrus.Error("File Not Found")
              rw.WriteHeader(http.StatusNotFound)
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
       if !writerXml(config, configPathDir+config.APIName) {

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

       configPathDir := getConfigHeader(req)
       configName := chi.URLParam(req, "configName")

       if !validateFileNameAndHeader(rw, req, configPathDir, configName) {
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

       if !validateFileNameAndHeader(rw, req, path, configName) {
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
              rw.WriteHeader(http.StatusConflict)
              return

       }

       if !writerXml(config, path+configName) {
              rw.WriteHeader(http.StatusInternalServerError)
              return
       }

       rw.WriteHeader(http.StatusNoContent)
}
