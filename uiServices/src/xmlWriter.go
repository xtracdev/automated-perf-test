package services

import (
	"encoding/xml"

	"fmt"
	"os"

	log "github.com/Sirupsen/logrus"
)

func writeXML(definition interface{}, path string) bool {
	filename := fmt.Sprintf("%s", path)

	payloadXML, err := xml.MarshalIndent(definition, "  ", "    ")
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
		file.Write(payloadXML)
	}
	return true
}
