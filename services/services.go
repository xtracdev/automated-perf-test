package services
import (
"encoding/json"
"net/http"

	"github.com/xtracdev/automated-perf-test/perfTestUtils"
	"github.com/Sirupsen/logrus"
)


func jsonHandler(rw http.ResponseWriter, req *http.Request) {
	//decoder := json.NewDecoder(req.Body)
	config, err:= json.Unmarshal(req.Body, &perfTestUtils.Config{})
	//err2 := decoder.Decode(&t)
	if err != nil {
		logrus.Error("Failed to unmarshall json body", err)
		respons
	}

	writerXml(config)
	defer req.Body.Close()

}
