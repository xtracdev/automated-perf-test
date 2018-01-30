package services

import (
	"encoding/xml"

	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/xtracdev/automated-perf-test/perfTestUtils"
	"os"
	"github.com/xtracdev/automated-perf-test/testStrategies"
)
//configurations writer
func writerXml(config perfTestUtils.Config, configPathDir string) bool {
	filename := fmt.Sprintf("%s%s.xml", configPathDir, config.APIName)

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

//test suite writer
func testSuiteWriterXml(testSuite testStrategies.TestSuite, configPathDir string) bool {
	filename := fmt.Sprintf("%s%s.xml", configPathDir, testSuite.Name)

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