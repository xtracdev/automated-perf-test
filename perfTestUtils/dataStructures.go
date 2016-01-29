package perfTestUtils

import (
	"encoding/xml"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"os"
	"runtime"
	"strings"
)

type Config struct {
	APIName                              string  `xml:"apiName"`
	TargetHost                           string  `xml:"targetHost"`
	TargetPort                           string  `xml:"targetPort"`
	NumIterations                        int     `xml:"numIterations"`
	AllowablePeakMemoryVariance          float64 `xml:"allowablePeakMemoryVariance"`
	AllowableServiceResponseTimeVariance float64 `xml:"allowableServiceResponseTimeVariance"`
	TestDefinationsDir                   string  `xml:"testDefinationsDir"`
	BaseStatsOutputDir                   string  `xml:"baseStatsOutputDir"`
	ReportOutputDir                      string  `xml:"reportOutputDir"`

	//These value can only be set by command line arguments as they control each training and test run.
	GBS             bool
	ResetPeakMemory bool

	//This value is determined by the environment/machine on which the test is being run.
	ExecutionHost string
}

func (c Config) PrintAndValidateConfig() {
	isConfigValid := true
	configOutput := []byte("")
	configOutput = append(configOutput, []byte("\n============== Configuration Settings =========\n")...)
	configOutput = append(configOutput, []byte(fmt.Sprintf("%-45s %-90s %2s", "apiName", c.APIName, "\n"))...)
	configOutput = append(configOutput, []byte(fmt.Sprintf("%-45s %-90s %2s", "targetHost", c.TargetHost, "\n"))...)
	configOutput = append(configOutput, []byte(fmt.Sprintf("%-45s %-90s %2s", "targetPort", c.TargetPort, "\n"))...)
	configOutput = append(configOutput, []byte(fmt.Sprintf("%-45s %-90d %2s", "numIterations", c.NumIterations, "\n"))...)
	configOutput = append(configOutput, []byte(fmt.Sprintf("%-45s %-90.2f %2s", "allowablePeakMemoryVariance", c.AllowablePeakMemoryVariance, "\n"))...)
	configOutput = append(configOutput, []byte(fmt.Sprintf("%-45s %-90.2f %2s", "allowableServiceResponseTimeVariance", c.AllowableServiceResponseTimeVariance, "\n"))...)
	configOutput = append(configOutput, []byte(fmt.Sprintf("%-45s %-90s %2s", "testDefinationsDir", c.TestDefinationsDir, "\n"))...)
	configOutput = append(configOutput, []byte(fmt.Sprintf("%-45s %-90s %2s", "baseStatsOutputDir", c.BaseStatsOutputDir, "\n"))...)
	configOutput = append(configOutput, []byte(fmt.Sprintf("%-45s %-90s %2s", "reportOutputDir", c.ReportOutputDir, "\n"))...)
	configOutput = append(configOutput, []byte(fmt.Sprintf("%-45s %-90t %2s", "gbs", c.GBS, "\n"))...)
	configOutput = append(configOutput, []byte(fmt.Sprintf("%-45s %-90t %2s", "resetPeakMemory", c.ResetPeakMemory, "\n"))...)
	configOutput = append(configOutput, []byte(fmt.Sprintf("%-45s %-90s %2s", "ExecutionHost", c.ExecutionHost, "\n"))...)
	configOutput = append(configOutput, []byte("\n=================================================\n")...)
	log.Info(string(configOutput))

	if strings.TrimSpace(c.APIName) == "" {
		log.Error("CONFIG ERROR: apiName must be set in config file")
		isConfigValid = false
	}
	if strings.TrimSpace(c.TargetHost) == "" {
		log.Error("CONFIG ERROR: targetHost must be set in config file")
		isConfigValid = false
	}
	if strings.TrimSpace(c.TargetPort) == "" {
		log.Error("CONFIG ERROR: targetPort must be set in config file")
		isConfigValid = false
	}
	if c.NumIterations < 1 {
		log.Error("CONFIG ERROR: numIterations must be set in config file and must be greater than 1")
		isConfigValid = false
	}
	if c.AllowablePeakMemoryVariance <= 0.0 {
		log.Error("CONFIG ERROR: allowablePeakMemoryVariance must be set in config file and must be greater than 0.0")
		isConfigValid = false
	}
	if c.AllowableServiceResponseTimeVariance <= 0.0 {
		log.Error("CONFIG ERROR: allowableServiceResponseTimeVariance must be set in config file and must be greater than 0.0")
		isConfigValid = false
	}
	if strings.TrimSpace(c.TestDefinationsDir) == "" {
		log.Error("CONFIG ERROR: testDefinationsDir must be set in config file")
		isConfigValid = false
	}
	if strings.TrimSpace(c.BaseStatsOutputDir) == "" {
		log.Error("CONFIG ERROR: baseStatsOutputDir must be set in config file")
		isConfigValid = false
	}
	if strings.TrimSpace(c.ReportOutputDir) == "" {
		log.Error("CONFIG ERROR: reportOutputDir must be set in config file")
		isConfigValid = false
	}
	if !isConfigValid {
		os.Exit(1)
	}
}

type Header struct {
	Value string `xml:",chardata"`
	Key   string `xml:"key,attr"`
}

//This struct defines the base performance statistics
type TestDefination struct {
	XMLName            xml.Name             `xml:"testDefination"`
	TestName           string               `xml:"testName"`
	HttpMethod         string               `xml:"httpMethod"`
	BaseUri            string               `xml:"baseUri"`
	Multipart          bool                 `xml:"multipart"`
	Payload            string               `xml:"payload"`
	MultipartPayload   []multipartFormField `xml:"multipartPayload>multipartFormField"`
	Scenario           string               `xml:"scenario"`
	ResponseStatusCode int                  `xml:"responseStatusCode"`
	XtracToken         string               `xml:"xtracToken"`
	Headers            []Header             `xml:"headers>header"`
}

type multipartFormField struct {
	FieldName   string `xml:"fieldName"`
	FieldValue  string `xml:"fieldValue"`
	FileName    string `xml:"fileName"`
	FileContent []byte `xml:"fileContent"`
}

//This struct defines the base performance statistics
type BasePerfStats struct {
	GenerationDate           string           `json:"GenerationDate"`
	ModifiedDate             string           `json:"ModifiedDate"`
	BasePeakMemory           uint64           `json:"BasePeakMemory"`
	BaseServiceResponseTimes map[string]int64 `json:"BaseServiceResponseTimes"`
	MemoryAudit              []uint64         `json:"MemoryAudit"`
}

//This struct defines the performance statistics for this test run
type PerfStats struct {
	PeakMemory           uint64
	ServiceResponseTimes map[string]int64
	MemoryAudit          []uint64
	TestPartitions       []TestPartition
}

type TestPartition struct {
	Count    int
	TestName string
}

type Entry struct {
	Cmdline  []string         `json:"cmdline"`
	Memstats runtime.MemStats `json:"memstats"`
}

type RspTimes []int64
