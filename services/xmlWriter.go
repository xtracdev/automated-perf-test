package services

import (
"encoding/xml"
"fmt"
"os"
"io"

	"github.com/xtracdev/automated-perf-test/perfTestUtils"
)

func writerXml(t perfTestUtils.Config) {




type ConfigFiles struct {
Configs []perfTestUtils.Config
}



filename := "ConfigFile.xml"
file, _ := os.Create(filename)

xmlWriter := io.Writer(file)

enc := xml.NewEncoder(xmlWriter)
enc.Indent("  ", "    ")
if err := enc.Encode(t); err != nil {
fmt.Printf("error: %v\n", err)
}

}