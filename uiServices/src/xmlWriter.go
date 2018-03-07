package services

import (
	"encoding/xml"

	"fmt"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/xtracdev/automated-perf-test/perfTestUtils"
	"github.com/xtracdev/automated-perf-test/testStrategies"
)

func configWriterXML(config perfTestUtils.Config, configPathDir string) bool {
	filename := fmt.Sprintf("%s", configPathDir)

	configAsXML, err := xml.MarshalIndent(config, "  ", "    ")
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
		file.Write(configAsXML)
	}
	return true
}

func testSuiteWriterXML(testSuite testStrategies.TestSuite, configPathDir string) bool {
	filename := fmt.Sprintf("%s", configPathDir)

	testSuiteAsXML, err := xml.MarshalIndent(testSuite, "  ", "    ")
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
		file.Write(testSuiteAsXML)
	}
	return true
}

func testCaseWriterXml(testSuite testStrategies.TestDefinition, path string) bool {
	filename := fmt.Sprintf("%s", path)

	testCaseAsXml, err := xml.MarshalIndent(testSuite, "  ", "    ")
	if err != nil {
		log.Error("Failed to marshal to XML. Error:", err)
		return false
	}

	file, err := os.Create(filename)
	if err != nil {
		log.Error("Failed to create output file. Error:", err)
		return false
	}
	defer file.Close()
	file.Write(testCaseAsXml)
	return true
}
