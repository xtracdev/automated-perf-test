package services

import (
	"encoding/xml"

	log "github.com/Sirupsen/logrus"
	"github.com/xtracdev/automated-perf-test/perfTestUtils"
	"os"
)

func writerXml(config perfTestUtils.Config, configPathDir string) bool {

	filename := configPathDir + config.APIName + ".xml"
	//validate that directory exists
	_, err := os.Stat(configPathDir)
	if os.IsNotExist(err) {
		log.Error("Can't find specified path")
		return false
	}

	configAsXml, err := xml.MarshalIndent(config, "  ", "    ")
	if err != nil {
		log.Error("Failed to marshal to XML. Error:", err)
		return false
	}

	file, err := os.Create(filename)
	if err != nil {
		log.Error("Failed to create output file. Error:", err)
		return false
	}
	if file != nil {
		defer file.Close()
		file.Write(configAsXml)
	}
	return true
}
