package services

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/Sirupsen/logrus"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
)

const validTestSuite = `
{
  "name": "TestSuiteService",
  "testStrategy": "SuiteBased",
  "description": "Services for XYZ",
  "testCases": [
    {
      "name":"file1",
      "preThinkTime": 1000,
      "postThinkTime": 2000,
      "execWeight": "Infrequent"
    },
    {
      "name":"file2"
    }
  ]
}
`

const invalidTestSuite = `
{
  "name": 123,
  "testStrategy": "SuiteBased",
  "testCases": [
    {
      "name":"file1.xml",
      "preThinkTime": "number",
      "postThinkTime": 2000,
      "execWeight": "Infrequent"
    },
    {
      "name":"file2.xml",
      "preThinkTime": 1,
      "postThinkTime": 10,
      "execWeight": "Sparse"
    }
  ]
}
`
const TestSuiteMissingRequired = `
{
  "name":"",
  "testCases": [
    {
      "name":"file1.xml",
      "preThinkTime": 1000,
      "postThinkTime": 2000,
      "execWeight": "Infrequent"
    }
  ]
}
`
const TestSuiteNoName = `
{
  "name": "",
  "testStrategy": "SuiteBased",
   "description": "Services for XYZ",
  "testCases": [
    {
      "name":"file1.xml",
      "preThinkTime": "xxxx"
      "postThinkTime": 2000,
      "execWeight": 123
    }
  ]
}
`

func TestValidTestSuitePost(t *testing.T) {
	r := chi.NewRouter()
	r.Mount("/", GetIndexPage())

	reader := strings.NewReader(validTestSuite)
	r.HandleFunc("/test-suites", postTestSuites)

	os.Remove(os.Getenv("GOPATH") + "/src/github.com/xtracdev/automated-perf-test/uiServices/test/TestSuiteService.xml")

	filePath := os.Getenv("GOPATH") + "/src/github.com/xtracdev/automated-perf-test/uiServices/test"
	request, err := http.NewRequest(http.MethodPost, "/test-suites", reader)
	request.Header.Set("testSuitePathDir", filePath)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, request)

	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, http.StatusCreated, w.Code, "Error: Did Not Successfully Post")
}

func TestFileExistsPost(t *testing.T) {
	r := chi.NewRouter()
	r.Mount("/", GetIndexPage())

	reader := strings.NewReader(validTestSuite)
	r.HandleFunc("/test-suites", postTestSuites)

	filePath := os.Getenv("GOPATH") + "/src/github.com/xtracdev/automated-perf-test/uiServices/test"
	request, err := http.NewRequest(http.MethodPost, "/test-suites", reader)
	request.Header.Set("testSuitePathDir", filePath)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, request)

	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, http.StatusBadRequest, w.Code, "Created, but should not have. File Should Alreday Exist")
}

func TestMissingRequiredTestSuitePost(t *testing.T) {
	r := chi.NewRouter()
	r.Mount("/", GetIndexPage())

	reader := strings.NewReader(TestSuiteMissingRequired)
	r.HandleFunc("/test-suites", postTestSuites)

	filePath := os.Getenv("GOPATH") + "/src/github.com/xtracdev/automated-perf-test/uiServices/test"
	request, err := http.NewRequest(http.MethodPost, "/test-suites", reader)
	request.Header.Set("testSuitePathDir", filePath)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, request)

	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, http.StatusBadRequest, w.Code, "Created TestSuite but should not have due to missing fields")

}

func TestIncorrectDataTypePost(t *testing.T) {
	r := chi.NewRouter()
	r.Mount("/", GetIndexPage())

	reader := strings.NewReader(TestSuiteNoName)
	r.HandleFunc("/test-suites", postTestSuites)

	filePath := os.Getenv("GOPATH") + "/src/github.com/xtracdev/automated-perf-test/uiServices/test"
	request, err := http.NewRequest(http.MethodPost, "/test-suites", reader)
	request.Header.Set("testSuitePathDir", filePath)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, request)

	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, http.StatusBadRequest, w.Code, "Created TestSuite but should not have due to incorrect data types")

}

func TestValidTestSuitePostNoConfigPathDir(t *testing.T) {
	r := chi.NewRouter()
	r.Mount("/", GetIndexPage())

	reader := strings.NewReader(validTestSuite)
	r.HandleFunc("/test-suites", postTestSuites)

	filePath := ""
	request, err := http.NewRequest(http.MethodPost, "/test-suites", reader)
	request.Header.Set("testSuitePathDir", filePath)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, request)

	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, http.StatusBadRequest, w.Code, "successfully created but should not have due to no config path specified")

}

func TestValidTestSuitePostConfigPathDirNotExist(t *testing.T) {
	r := chi.NewRouter()
	r.Mount("/", GetIndexPage())

	reader := strings.NewReader(validTestSuite)
	r.HandleFunc("/test-suites", postTestSuites)

	filePath := "C:/a/b/c/d/////"
	request, err := http.NewRequest(http.MethodPost, "/test-suites", reader)
	request.Header.Set("testSuitePathDir", filePath)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, request)

	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, http.StatusBadRequest, w.Code, "successfully created but should not have due to no config path specified")

}

func TestValidTestSuitePut(t *testing.T) {
	r := chi.NewRouter()
	r.Mount("/", GetIndexPage())

	reader := strings.NewReader(validTestSuite)
	r.HandleFunc("/test-suites", putTestSuites)

	filePath := os.Getenv("GOPATH") + "/src/github.com/xtracdev/automated-perf-test/uiServices/test/"
	request, err := http.NewRequest(http.MethodPut, "/test-suites/TestSuiteService", reader)
	request.Header.Set("testSuitePathDir", filePath)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, request)

	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, http.StatusNoContent, w.Code, "Did Not successfully Update")
}

func TestTestSuiteMissingFieldPut(t *testing.T) {
	r := chi.NewRouter()
	r.Mount("/", GetIndexPage())

	reader := strings.NewReader(TestSuiteMissingRequired)
	r.HandleFunc("/test-suites", putTestSuites)

	filePath := os.Getenv("GOPATH") + "/src/github.com/xtracdev/automated-perf-test/uiServices/test/"
	request, err := http.NewRequest(http.MethodPut, "/test-suites/TestSuiteService", reader)
	request.Header.Set("testSuitePathDir", filePath)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, request)

	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, http.StatusBadRequest, w.Code, "Sucessfully updated. Field Should be missing so update shouldn't occur")
}

func TestInvalidTestSuitePut(t *testing.T) {
	r := chi.NewRouter()
	r.Mount("/", GetIndexPage())

	reader := strings.NewReader(invalidTestSuite)
	r.HandleFunc("/test-suites", putTestSuites)

	filePath := os.Getenv("GOPATH") + "/src/github.com/xtracdev/automated-perf-test/uiServices/test/"
	request, err := http.NewRequest(http.MethodPut, "/test-suites/TestSuiteService", reader)
	request.Header.Set("testSuitePathDir", filePath)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, request)

	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, http.StatusBadRequest, w.Code, "Sucessfully updated. Field data type should have been incorrect so update should occur")
}

func TestInvalidUrlTestSuitePut(t *testing.T) {
	r := chi.NewRouter()
	r.Mount("/", GetIndexPage())

	reader := strings.NewReader(validTestSuite)
	r.HandleFunc("/test-suites", putTestSuites)

	filePath := os.Getenv("GOPATH") + "/src/github.com/xtracdev/automated-perf-test/uiServices/test/"
	request, err := http.NewRequest(http.MethodPut, "/test-suites/xxx", reader)
	request.Header.Set("testSuitePathDir", filePath)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, request)

	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, http.StatusNotFound, w.Code, "Sucessfully updated. Should have have worked using /test-suites/xxx")
}

func TestNoUrlTestSuitePut(t *testing.T) {
	r := chi.NewRouter()
	r.Mount("/", GetIndexPage())

	reader := strings.NewReader(validTestSuite)
	r.HandleFunc("/test-suites", putTestSuites)

	filePath := os.Getenv("GOPATH") + "/src/github.com/xtracdev/automated-perf-test/uiServices/test/"
	request, err := http.NewRequest(http.MethodPut, "", reader)
	request.Header.Set("testSuitePathDir", filePath)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, request)

	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, http.StatusNotFound, w.Code, "Sucessfully updated. Should not have worked with no URL")
}

func TestPutWithNoPathSlash(t *testing.T) {
	r := chi.NewRouter()
	r.Mount("/", GetIndexPage())

	reader := strings.NewReader(validTestSuite)
	r.HandleFunc("/test-suites", putTestSuites)

	filePath := os.Getenv("GOPATH") + "/src/github.com/xtracdev/automated-perf-test/uiServices/test"
	request, err := http.NewRequest(http.MethodPut, "/test-suites/TestSuiteService", reader)
	request.Header.Set("testSuitePathDir", filePath)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, request)

	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, http.StatusNoContent, w.Code, "Did not update. Should have added '/' to path")
}

func TestNoPathTestSuitePut(t *testing.T) {
	r := chi.NewRouter()
	r.Mount("/", GetIndexPage())

	reader := strings.NewReader(validTestSuite)
	r.HandleFunc("/test-suites", putTestSuites)

	filePath := ""
	request, err := http.NewRequest(http.MethodPut, "/test-suites/TestSuiteService", reader)
	request.Header.Set("testSuitePathDir", filePath)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, request)

	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, http.StatusBadRequest, w.Code, "Successfully updated. Should not have worked due to no filepath")
}

func TestNoFileNameTestSuitePut(t *testing.T) {
	r := chi.NewRouter()
	r.Mount("/", GetIndexPage())

	reader := strings.NewReader(validTestSuite)
	r.HandleFunc("/test-suites", putTestSuites)

	filePath := os.Getenv("GOPATH") + "/src/github.com/xtracdev/automated-perf-test/uiServices/test"
	request, err := http.NewRequest(http.MethodPut, "/test-suites", reader)
	request.Header.Set("testSuitePathDir", filePath)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, request)

	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, http.StatusBadRequest, w.Code, "Successfully updated. Should not have worked due to no file name given")
}

func TestSuccessfulGetTestSuite(t *testing.T) {
	r := chi.NewRouter()
	r.Mount("/", GetIndexPage())

	filePath := os.Getenv("GOPATH") + "/src/github.com/xtracdev/automated-perf-test/uiServices/test/"
	request, err := http.NewRequest(http.MethodGet, "/test-suites/TestSuiteService", nil)

	request.Header.Set("testSuitePathDir", filePath)
	request.Header.Get("testSuitePathDir")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, request)

	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, http.StatusOK, w.Code, "Error. Did not successfully GET")
}

func TestGetTestSuiteNoPath(t *testing.T) {
	r := chi.NewRouter()
	r.Mount("/", GetIndexPage())

	filePath := ""
	request, err := http.NewRequest(http.MethodGet, "/test-suites/TestSuiteService.xml", nil)

	request.Header.Set("testSuitePathDir", filePath)
	request.Header.Get("testSuitePathDir")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, request)

	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, http.StatusBadRequest, w.Code, "Retrived file but should not have as there is no path")
}

func TestGetTestSuiteFileNotFound(t *testing.T) {
	r := chi.NewRouter()
	r.Mount("/", GetIndexPage())

	filePath := os.Getenv("GOPATH") + "/src/github.com/xtracdev/automated-perf-test/uiServices/test/"
	request, err := http.NewRequest(http.MethodGet, "/test-suites/xxx", nil)

	request.Header.Set("testSuitePathDir", filePath)
	request.Header.Get("testSuitePathDir")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, request)

	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, http.StatusNotFound, w.Code, "Retrived a file but should not have as there is no file")
}

func TestSuccessfulGetAllSuites(t *testing.T) {
	r := chi.NewRouter()
	r.Mount("/", GetIndexPage())

	DirectoryPath := os.Getenv("GOPATH") + "/src/github.com/xtracdev/automated-perf-test/uiServices/test/"
	request, err := http.NewRequest(http.MethodGet, "/test-suites", nil)

	request.Header.Set("testSuitePathDir", DirectoryPath)
	request.Header.Get("testSuitePathDir")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, request)

	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, http.StatusOK, w.Code, "Did not get all test suites")
}

func TestGetAllSuitesNoHeader(t *testing.T) {
	r := chi.NewRouter()
	r.Mount("/", GetIndexPage())

	DircetoryPath := ""
	request, err := http.NewRequest(http.MethodGet, "/test-suites", nil)

	request.Header.Set("testSuitePathDir", DircetoryPath)
	request.Header.Get("testSuitePathDir")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, request)

	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, http.StatusBadRequest, w.Code, "Did not get all test suites")
}

func TestDeleteAllSuites(t *testing.T) {
	r := chi.NewRouter()
	r.Mount("/", GetIndexPage())

	t.Run("400 - Bad request no header", func(t *testing.T) {
		directoryPath := ""
		request, err := http.NewRequest(http.MethodDelete, "/test-suites/testParam", nil)
		if err != nil {
			logrus.Warnf("Error creating the request %s", err)
		}

		request.Header.Set("testSuitePathDir", directoryPath)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, request)

		assert.Equal(t, http.StatusBadRequest, w.Code, "Did not DELETE the file")
	})

	t.Run("failure - 500", func(t *testing.T) {
		filePath := os.Getenv("GOPATH") + "C:/Usersa622123/go/src/github.com/xtracdev/automated-perf-test/config"
		request, err := http.NewRequest(http.MethodDelete, "/test-suites/test", nil)

		request.Header.Set("testSuitePathDir", filePath)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, request)

		if err != nil {
			t.Error(err)
		}

		assert.Equal(t, http.StatusInternalServerError, w.Code, "Internal error: failed to DELETE the file")
	})

	t.Run("failure - 404", func(t *testing.T) {
		filePath := os.Getenv("GOPATH") + "/src/github.com/xtracdev/automated-perf-test/config"
		request, err := http.NewRequest(http.MethodDelete, "/test-suites/abc", nil)

		request.Header.Set("testSuitePathDir", filePath)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, request)

		if err != nil {
			t.Error(err)
		}

		assert.Equal(t, http.StatusNotFound, w.Code, "File was found and NOT DELETED")
	})

	t.Run("success - 204", func(t *testing.T) {
		filePath := os.Getenv("GOPATH") + "/src/github.com/xtracdev/automated-perf-test/config"
		err := ioutil.WriteFile(fmt.Sprintf("%s/%s.xml", filePath, "test"), nil, 0666)
		if err != nil {
			logrus.Errorf("Error trying to create the file %s", err)
		}

		request, err := http.NewRequest(http.MethodDelete, "/test-suites/test", nil)
		if err != nil {
			logrus.Errorf("error trying to delete the file %s", err)
		}

		request.Header.Set("testSuitePathDir", filePath)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, request)

		assert.Equal(t, http.StatusNoContent, w.Code, "File was NOT DELETED")
	})
}
