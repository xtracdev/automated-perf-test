package services

import (
	"fmt"
	"net/http"
	"strings"
	"github.com/Sirupsen/logrus"
	"github.com/go-chi/chi"
	"os"
)

func TestCaseCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})

}
func getTestCaseHeader(req *http.Request) string {
	testCasePathDir := req.Header.Get("testCasePathDir")

	if !strings.HasSuffix(testCasePathDir, "/") {
		testCasePathDir = testCasePathDir + "/"
	}
	return testCasePathDir
}

func deleteTestCase(rw http.ResponseWriter, req *http.Request){
	testCasePathDir := getTestCaseHeader(req)
	testCaseName := chi.URLParam(req, "testCaseName")
	ValidateFileNameAndHeader(rw, req, testCasePathDir, testCaseName)

	filepath := fmt.Sprintf("%s%s.xml", testCasePathDir, testCaseName)

	if _, err := os.Stat(filepath); err != nil {
		if os.IsNotExist(err) {
			logrus.Println("File Not Found")
			rw.WriteHeader(http.StatusNotFound)
			return
		}
	}

	os.Remove(filepath)

	rw.WriteHeader(http.StatusNoContent)

}