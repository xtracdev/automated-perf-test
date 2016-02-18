package perfTestUtils

import (
	"encoding/xml"
	"fmt"
	//log "github.com/Sirupsen/logrus"
	"net/http"
	"runtime"
	"strings"
	"time"
)

const (
	defaultAPIName                              = "Default_API_NAME"
	defaultTargetHost                           = "localhost"
	defaultTargetPort                           = "8080"
	defaultNumIterations                        = 1000
	defaultAllowablePeakMemoryVariance          = float64(15)
	defaultAllowableServiceResponseTimeVariance = float64(15)
	defaultTestDefinitionsDir                   = "./testDefinitions"
	defaultBaseStatsOutputDir                   = "./envStats"
	defaultReportOutputDir                      = "./"
	defaultConcurrentUsers                      = 1
)

type Config struct {
	APIName                              string  `xml:"apiName"`
	TargetHost                           string  `xml:"targetHost"`
	TargetPort                           string  `xml:"targetPort"`
	NumIterations                        int     `xml:"numIterations"`
	AllowablePeakMemoryVariance          float64 `xml:"allowablePeakMemoryVariance"`
	AllowableServiceResponseTimeVariance float64 `xml:"allowableServiceResponseTimeVariance"`
	TestDefinitionsDir                   string  `xml:"testDefinitionsDir"`
	BaseStatsOutputDir                   string  `xml:"baseStatsOutputDir"`
	ReportOutputDir                      string  `xml:"reportOutputDir"`
	ConcurrentUsers                      int     `xml:"concurrentUsers"`

	//These value can only be set by command line arguments as they control each training and test run.
	GBS          bool
	ReBaseMemory bool
	ReBaseAll    bool

	//This value is determined by the environment/machine on which the test is being run.
	ExecutionHost string

	//template file
	ReportTemplateFile string `xml:"reportTemplateFile,omitempty"`
	ConfigFileFormat   string
	TestFileFormat     string
}

func (c *Config) SetDefaults() {
	c.APIName = defaultAPIName
	c.TargetHost = defaultTargetHost
	c.TargetPort = defaultTargetPort
	c.NumIterations = defaultNumIterations
	c.AllowablePeakMemoryVariance = defaultAllowablePeakMemoryVariance
	c.AllowableServiceResponseTimeVariance = defaultAllowableServiceResponseTimeVariance
	c.TestDefinitionsDir = defaultTestDefinitionsDir
	c.BaseStatsOutputDir = defaultBaseStatsOutputDir
	c.ReportOutputDir = defaultReportOutputDir
	c.ConcurrentUsers = defaultConcurrentUsers
}

func (c Config) PrintAndValidateConfig(exit func(code int)) {
	isConfigValid := true
	configOutput := []byte("")
	configOutput = append(configOutput, []byte("\n============== Configuration Settings =========\n")...)
	configOutput = append(configOutput, []byte(fmt.Sprintf("%-45s %-90s %2s", "apiName", c.APIName, "\n"))...)
	configOutput = append(configOutput, []byte(fmt.Sprintf("%-45s %-90s %2s", "targetHost", c.TargetHost, "\n"))...)
	configOutput = append(configOutput, []byte(fmt.Sprintf("%-45s %-90s %2s", "targetPort", c.TargetPort, "\n"))...)
	configOutput = append(configOutput, []byte(fmt.Sprintf("%-45s %-90d %2s", "numIterations", c.NumIterations, "\n"))...)
	configOutput = append(configOutput, []byte(fmt.Sprintf("%-45s %-90d %2s", "concurrentUsers", c.ConcurrentUsers, "\n"))...)
	configOutput = append(configOutput, []byte(fmt.Sprintf("%-45s %-90.2f %2s", "allowablePeakMemoryVariance", c.AllowablePeakMemoryVariance, "\n"))...)
	configOutput = append(configOutput, []byte(fmt.Sprintf("%-45s %-90.2f %2s", "allowableServiceResponseTimeVariance", c.AllowableServiceResponseTimeVariance, "\n"))...)
	configOutput = append(configOutput, []byte(fmt.Sprintf("%-45s %-90s %2s", "testDefinitionsDir", c.TestDefinitionsDir, "\n"))...)
	configOutput = append(configOutput, []byte(fmt.Sprintf("%-45s %-90s %2s", "baseStatsOutputDir", c.BaseStatsOutputDir, "\n"))...)
	configOutput = append(configOutput, []byte(fmt.Sprintf("%-45s %-90s %2s", "reportOutputDir", c.ReportOutputDir, "\n"))...)
	configOutput = append(configOutput, []byte(fmt.Sprintf("%-45s %-90t %2s", "gbs", c.GBS, "\n"))...)
	configOutput = append(configOutput, []byte(fmt.Sprintf("%-45s %-90t %2s", "reBaseMemory", c.ReBaseMemory, "\n"))...)
	configOutput = append(configOutput, []byte(fmt.Sprintf("%-45s %-90t %2s", "reBaseAll", c.ReBaseAll, "\n"))...)
	configOutput = append(configOutput, []byte(fmt.Sprintf("%-45s %-90s %2s", "ExecutionHost", c.ExecutionHost, "\n"))...)
	configOutput = append(configOutput, []byte("\n=================================================\n")...)
	//log.Info(string(configOutput))
	fmt.Println(string(configOutput))

	if strings.TrimSpace(c.APIName) == "" {
		fmt.Println("CONFIG ERROR: apiName must be set in config file")
		isConfigValid = false
	}
	if strings.TrimSpace(c.TargetHost) == "" {
		fmt.Println("CONFIG ERROR: targetHost must be set in config file")
		isConfigValid = false
	}
	if strings.TrimSpace(c.TargetPort) == "" {
		fmt.Println("CONFIG ERROR: targetPort must be set in config file")
		isConfigValid = false
	}
	if c.NumIterations < 1 {
		fmt.Println("CONFIG ERROR: numIterations must be set in config file and must be greater than 1")
		isConfigValid = false
	}
	if c.ConcurrentUsers < 1 {
		fmt.Println("CONFIG ERROR: concurrentUsers must be set in config file and must be greater than 1")
		isConfigValid = false
	}
	if c.AllowablePeakMemoryVariance <= 0.0 {
		fmt.Println("CONFIG ERROR: allowablePeakMemoryVariance must be set in config file and must be greater than 0.0")
		isConfigValid = false
	}
	if c.AllowableServiceResponseTimeVariance <= 0.0 {
		fmt.Println("CONFIG ERROR: allowableServiceResponseTimeVariance must be set in config file and must be greater than 0.0")
		isConfigValid = false
	}
	if strings.TrimSpace(c.TestDefinitionsDir) == "" {
		fmt.Println("CONFIG ERROR: testDefinitionsDir must be set in config file")
		isConfigValid = false
	}
	if strings.TrimSpace(c.BaseStatsOutputDir) == "" {
		fmt.Println("CONFIG ERROR: baseStatsOutputDir must be set in config file")
		isConfigValid = false
	}
	if strings.TrimSpace(c.ReportOutputDir) == "" {
		fmt.Println("CONFIG ERROR: reportOutputDir must be set in config file")
		isConfigValid = false
	}
	if !isConfigValid {
		exit(1)
	}
}

type Header struct {
	Value string `xml:",chardata"`
	Key   string `xml:"key,attr"`
}

//This struct defines the base performance statistics
type TestDefinition struct {
	XMLName            xml.Name             `xml:"testDefinition"`
	TestName           string               `xml:"testName" toml:"testName"`
	HttpMethod         string               `xml:"httpMethod"`
	BaseUri            string               `xml:"baseUri"`
	Multipart          bool                 `xml:"multipart"`
	Payload            string               `xml:"payload"`
	MultipartPayload   []multipartFormField `xml:"multipartPayload>multipartFormField"`
	ResponseStatusCode int                  `xml:"responseStatusCode"`
	Headers            []Header             `xml:"headers>header"`
}

//TomlTestDefinition defines the test in TOML language
type TomlTestDefinition struct {
	TestName           string               `toml:"testName"`
	HttpMethod         string               `toml:"httpMethod"`
	BaseUri            string               `toml:"baseUri"`
	Multipart          bool                 `toml:"multipart"`
	Payload            string               `toml:"payload"`
	MultipartPayload   []multipartFormField `toml:"multipartFormField"`
	ResponseStatusCode int                  `toml:"responseStatusCode"`
	Headers            http.Header          `toml:"headers"`
}

type multipartFormField struct {
	FieldName   string `xml:"fieldName" toml:"fieldName"`
	FieldValue  string `xml:"fieldValue" toml:"fieldValue"`
	FileName    string `xml:"fileName" toml:"fileName"`
	FileContent []byte `xml:"fileContent" toml:"fileContent"`
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
	TestDate             time.Time
}

func (ps *PerfStats) GetTestTime() string {
	return ps.TestDate.Format(time.RFC850)
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
