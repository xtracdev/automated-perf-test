package services

import (
	"encoding/xml"

	"fmt"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/xtracdev/automated-perf-test/perfTestUtils"
	"github.com/xtracdev/automated-perf-test/testStrategies"
)

func writeXml(definition interface{}, path string) bool {
	filename := fmt.Sprintf("%s", path)

	configAsXml, err := xml.MarshalIndent(definition, "  ", "    ")
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

func configWriterXml(config perfTestUtils.Config, path string) bool {
	filename := fmt.Sprintf("%s", path)

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

func testSuiteWriterXml(testSuite testStrategies.TestSuite, path string) bool {
	filename := fmt.Sprintf("%s", path)

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

func testCaseWriterXml(testCase testStrategies.TestDefinition, path string) bool {
	filename := fmt.Sprintf("%s", path)

	testCaseAsXml, err := xml.MarshalIndent(testCase, "  ", "    ")
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
		file.Write(testCaseAsXml)
	}
	return true
}
