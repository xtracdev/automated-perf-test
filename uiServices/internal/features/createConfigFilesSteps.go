package features

import (
	. "github.com/gucumber/gucumber"

	"net/http"
	"strings"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/xtracdev/automated-perf-test/uiServices"

)

var httpClient *http.Client

func init() {

	var code int
	Given(`^the automated performance ui server is available`, func() {
		httpClient = &http.Client{}

		services.StartUiMode()
	})

	When(`^the user makes a request for (.+?) http://localhost:9191/configs with payload$`, func(method, data string) {

		//get payload (JSON)
		reader := strings.NewReader(data)
		//path to save file
		filePath := os.Getenv("GOPATH") + "/src/github.com/xtracdev/automated-perf-test/config/"
		//set URI and header
		request, err := http.NewRequest(method, "http://localhost:9191/configs", reader)
		logrus.Error(err) //log and (TODO) handle errors
		request.Header.Set("configPathDir", filePath)
		//run request
		resp, err := httpClient.Do(request)
		logrus.Error(err) //log and (TODO) handle errors
		code = resp.StatusCode
	})

	Then(`^the POST configuration service returns (.+?) HTTP status`, func(code int) {
		//check the response
		if code != http.StatusCreated {
			T.Errorf("incorrect error code")
		}
	})
}
