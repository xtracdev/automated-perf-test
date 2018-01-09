package perfTestUtils

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
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
	defaultTestCaseDir                          = "./definitions/testCases"
	defaultTestSuiteDir                         = "./definitions/testSuites"
	defaultBaseStatsOutputDir                   = "./envStats"
	defaultReportOutputDir                      = "./report"
	defaultConcurrentUsers                      = 1
	defaultTestSuite                            = ""
	defaultMemoryEndpoint                       = "/debug/vars"
	defaultRequestDelay                         = 1
	defaultTPSFreq                              = 30
	defaultRampUsers                            = 0
	defaultRampDelay                            = 10
)

// Config struct contains all values set by the config.xml file. Most, if not
// all, can be overridden from command line.
type Config struct {
	APIName                              string  `xml:"apiName" json:"apiName"`
	TargetHost                           string  `xml:"targetHost" json:"targetHost"`
	TargetPort                           string  `xml:"targetPort" json:"targetPort"`
	NumIterations                        int     `xml:"numIterations" json:"numIterations"`
	AllowablePeakMemoryVariance          float64 `xml:"allowablePeakMemoryVariance" json:"allowablePeakMemoryVariance"`
	AllowableServiceResponseTimeVariance float64 `xml:"allowableServiceResponseTimeVariance" json:"allowableServiceResponseTimeVariance"`
	TestCaseDir                          string  `xml:"testCaseDir" json:"testCaseDir"`
	TestSuiteDir                         string  `xml:"testSuiteDir" json:"testSuiteDir"`
	BaseStatsOutputDir                   string  `xml:"baseStatsOutputDir" json:"baseStatsOutputDir"`
	ReportOutputDir                      string  `xml:"reportOutputDir" json:"reportOutputDir"`
	ConcurrentUsers                      int     `xml:"concurrentUsers" json:"concurrentUsers"`
	TestSuite                            string  `xml:"testSuite" json:"testSuite"`
	MemoryEndpoint                       string  `xml:"memoryEndpoint" json:"memoryEndpoint"`
	RequestDelay                         int     `xml:"requestDelay" json:"requestDelay"`
	TPSFreq                              int     `xml:"TPSFreq" json:"TPSFreq"`
	RampUsers                            int     `xml:"rampUsers" json:"rampUsers"`
	RampDelay                            int     `xml:"rampDelay" json:"rampDelay"`

	//These value can only be set by command line arguments as they control each training and test run.
	GBS          bool `xml:"gbs,omitempty"`
	ReBaseMemory bool `xml:"reBaseMemory,omitempty"`
	ReBaseAll    bool `xml:"reBaseAll,omitempty"`

	//This value is determined by the environment/machine on which the test is being run.
	ExecutionHost string `xml:"executionHost,omitempty"`

	//Template file
	ReportTemplateFile string `xml:"reportTemplateFile,omitempty"`
}

// SetDefaults initializes the Config struct.
func (c *Config) SetDefaults() {
	c.APIName = defaultAPIName
	c.TargetHost = defaultTargetHost
	c.TargetPort = defaultTargetPort
	c.NumIterations = defaultNumIterations
	c.AllowablePeakMemoryVariance = defaultAllowablePeakMemoryVariance
	c.AllowableServiceResponseTimeVariance = defaultAllowableServiceResponseTimeVariance
	c.TestCaseDir = defaultTestCaseDir
	c.TestSuiteDir = defaultTestSuiteDir
	c.BaseStatsOutputDir = defaultBaseStatsOutputDir
	c.ReportOutputDir = defaultReportOutputDir
	c.ConcurrentUsers = defaultConcurrentUsers
	c.TestSuite = defaultTestSuite
	c.MemoryEndpoint = defaultMemoryEndpoint
	c.RequestDelay = defaultRequestDelay
	c.TPSFreq = defaultTPSFreq
	c.RampUsers = defaultRampUsers
	c.RampDelay = defaultRampDelay

	c.GBS = false
	c.ReBaseMemory = false
	c.ReBaseAll = false
}

// PrintAndValidateConfig sets any out of bounds value of the Config struct to the default.
func (c *Config) PrintAndValidateConfig() {

	if strings.TrimSpace(c.APIName) == "" {
		c.APIName = defaultAPIName
	}
	if strings.TrimSpace(c.TargetHost) == "" {
		c.TargetHost = defaultTargetHost
	}
	if strings.TrimSpace(c.TargetPort) == "" {
		c.TargetPort = defaultTargetPort
	}
	if c.NumIterations < 1 {
		c.NumIterations = defaultNumIterations
	}
	if c.ConcurrentUsers < 1 {
		c.ConcurrentUsers = defaultConcurrentUsers
	}
	if c.AllowablePeakMemoryVariance < 0.0 {
		c.AllowablePeakMemoryVariance = defaultAllowablePeakMemoryVariance
	}
	if c.AllowableServiceResponseTimeVariance < 0.0 {
		c.AllowableServiceResponseTimeVariance = defaultAllowableServiceResponseTimeVariance
	}
	if strings.TrimSpace(c.TestCaseDir) == "" {
		c.TestCaseDir = defaultTestCaseDir
	}
	if strings.TrimSpace(c.BaseStatsOutputDir) == "" {
		c.BaseStatsOutputDir = defaultBaseStatsOutputDir
	}
	if strings.TrimSpace(c.ReportOutputDir) == "" {
		c.ReportOutputDir = defaultReportOutputDir
	}
	if strings.TrimSpace(c.MemoryEndpoint) == "" {
		c.MemoryEndpoint = defaultMemoryEndpoint
	}
	if c.RequestDelay < 1 {
		c.RequestDelay = defaultRequestDelay
	}
	if c.TPSFreq < 1 {
		c.TPSFreq = defaultTPSFreq
	}
	if c.RampUsers < 0 {
		c.RampUsers = defaultRampUsers
	}
	if c.RampDelay < 1 {
		c.RampDelay = defaultRampDelay
	}

	configOutput := []byte("")
	configOutput = append(configOutput, []byte("\n============== Configuration Settings =========\n")...)
	configOutput = append(configOutput, []byte(fmt.Sprintf("%-45s %-90s %2s", "apiName", c.APIName, "\n"))...)
	configOutput = append(configOutput, []byte(fmt.Sprintf("%-45s %-90s %2s", "targetHost", c.TargetHost, "\n"))...)
	configOutput = append(configOutput, []byte(fmt.Sprintf("%-45s %-90s %2s", "targetPort", c.TargetPort, "\n"))...)
	configOutput = append(configOutput, []byte(fmt.Sprintf("%-45s %-90d %2s", "numIterations", c.NumIterations, "\n"))...)
	configOutput = append(configOutput, []byte(fmt.Sprintf("%-45s %-90d %2s", "concurrentUsers", c.ConcurrentUsers, "\n"))...)
	configOutput = append(configOutput, []byte(fmt.Sprintf("%-45s %-90.2f %2s", "allowablePeakMemoryVariance", c.AllowablePeakMemoryVariance, "\n"))...)
	configOutput = append(configOutput, []byte(fmt.Sprintf("%-45s %-90.2f %2s", "allowableServiceResponseTimeVariance", c.AllowableServiceResponseTimeVariance, "\n"))...)
	configOutput = append(configOutput, []byte(fmt.Sprintf("%-45s %-90s %2s", "testCaseDir", c.TestCaseDir, "\n"))...)
	configOutput = append(configOutput, []byte(fmt.Sprintf("%-45s %-90s %2s", "testSuiteDir", c.TestSuiteDir, "\n"))...)
	configOutput = append(configOutput, []byte(fmt.Sprintf("%-45s %-90s %2s", "testSuite", c.TestSuite, "\n"))...)
	configOutput = append(configOutput, []byte(fmt.Sprintf("%-45s %-90s %2s", "memoryEndpoint", c.MemoryEndpoint, "\n"))...)
	configOutput = append(configOutput, []byte(fmt.Sprintf("%-45s %-90s %2s", "baseStatsOutputDir", c.BaseStatsOutputDir, "\n"))...)
	configOutput = append(configOutput, []byte(fmt.Sprintf("%-45s %-90s %2s", "reportOutputDir", c.ReportOutputDir, "\n"))...)
	configOutput = append(configOutput, []byte(fmt.Sprintf("%-45s %-90t %2s", "gbs", c.GBS, "\n"))...)
	configOutput = append(configOutput, []byte(fmt.Sprintf("%-45s %-90t %2s", "reBaseMemory", c.ReBaseMemory, "\n"))...)
	configOutput = append(configOutput, []byte(fmt.Sprintf("%-45s %-90t %2s", "reBaseAll", c.ReBaseAll, "\n"))...)
	configOutput = append(configOutput, []byte(fmt.Sprintf("%-45s %-90s %2s", "executionHost", c.ExecutionHost, "\n"))...)
	configOutput = append(configOutput, []byte(fmt.Sprintf("%-45s %-90d %2s", "requestDelay", c.RequestDelay, "\n"))...)
	configOutput = append(configOutput, []byte(fmt.Sprintf("%-45s %-90d %2s", "tpsFreq", c.TPSFreq, "\n"))...)
	configOutput = append(configOutput, []byte(fmt.Sprintf("%-45s %-90d %2s", "rampUsers", c.RampUsers, "\n"))...)
	configOutput = append(configOutput, []byte(fmt.Sprintf("%-45s %-90d %2s", "rampDelay", c.RampDelay, "\n"))...)
	configOutput = append(configOutput, []byte("\n=================================================\n")...)
	log.Info(string(configOutput))
}

// BasePerfStats struct defines the base performance statistics
type BasePerfStats struct {
	GenerationDate           string           `json:"GenerationDate"`
	ModifiedDate             string           `json:"ModifiedDate"`
	BasePeakMemory           uint64           `json:"BasePeakMemory"`
	BaseServiceResponseTimes map[string]int64 `json:"BaseServiceResponseTimes"`
	MemoryAudit              []uint64         `json:"MemoryAudit"`
}

// PerfStats struct defines the performance statistics for this test run
type PerfStats struct {
	PeakMemory           uint64
	ServiceResponseTimes map[string]int64
	ServiceTransCount    map[string]*uint64
	ServiceErrorCount    map[string]*uint64
	ServiceTPS           map[string]float64
	OverAllTransCount    uint64
	OverAllErrorCount    uint64
	OverAllTPS           float64
	MemoryAudit          []uint64
	TestPartitions       []TestPartition
	TestTimeStart        time.Time
	TestTimeEnd          time.Time
}

// GetTestTimeStart returns the start time of the test in RFC850 format.
func (ps *PerfStats) GetTestTimeStart() string {
	return ps.TestTimeStart.Format(time.RFC850)
}

// GetTestTimeEnd returns the end time of the test in RFC850 format.
func (ps *PerfStats) GetTestTimeEnd() string {
	return ps.TestTimeEnd.Format(time.RFC850)
}

// TestPartition struct combines the test name with a count for use on the report.
type TestPartition struct {
	Count    int
	TestName string
}

// Entry struct combines a serivce call with memory statistics for use on the report.
type Entry struct {
	Cmdline  []string         `json:"cmdline"`
	Memstats runtime.MemStats `json:"memstats"`
}

// RspTimes contains a sequential array of response time integers.
type RspTimes []int64
