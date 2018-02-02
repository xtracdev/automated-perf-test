package services

import (
	"encoding/xml"

	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/xtracdev/automated-perf-test/perfTestUtils"
	"github.com/xtracdev/automated-perf-test/testStrategies"
	"os"
)

func configWriterXml(config perfTestUtils.Config, configPathDir string) bool {
	filename := fmt.Sprintf("%s.xml", configPathDir)

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

func testSuiteWriterXml(testSuite testStrategies.TestSuite, configPathDir string) bool {
	filename := fmt.Sprintf("%s.xml", configPathDir)

	testSuiteAsXml, err := xml.MarshalIndent(testSuite, "  ", "    ")
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
		file.Write(testSuiteAsXml)
	}
	return true
}
