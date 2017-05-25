package perfTestUtils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSetDefaults(t *testing.T) {
	c := &Config{}
	c.SetDefaults()
	assert.Equal(t, defaultAPIName, c.APIName)
	assert.Equal(t, defaultTargetHost, c.TargetHost)
	assert.Equal(t, defaultTargetPort, c.TargetPort)
	assert.Equal(t, defaultNumIterations, c.NumIterations)
	assert.Equal(t, defaultAllowablePeakMemoryVariance, c.AllowablePeakMemoryVariance)
	assert.Equal(t, defaultAllowableServiceResponseTimeVariance, c.AllowableServiceResponseTimeVariance)
	assert.Equal(t, defaultTestCaseDir, c.TestCaseDir)
	assert.Equal(t, defaultBaseStatsOutputDir, c.BaseStatsOutputDir)
	assert.Equal(t, defaultReportOutputDir, c.ReportOutputDir)
	assert.Equal(t, defaultConcurrentUsers, c.ConcurrentUsers)
	assert.Equal(t, defaultTestSuite, c.TestSuite)
	assert.Equal(t, defaultMemoryEndpoint, c.MemoryEndpoint)
	assert.Equal(t, defaultRequestDelay, c.RequestDelay)
	assert.Equal(t, defaultTPSFreq, c.TPSFreq)
	assert.Equal(t, defaultRampUsers, c.RampUsers)
	assert.Equal(t, defaultRampDelay, c.RampDelay)
	assert.Equal(t, false, c.GBS)
	assert.Equal(t, false, c.ReBaseMemory)
	assert.Equal(t, false, c.ReBaseAll)
}

func TestPrintAndValidate(t *testing.T) {
	c := &Config{}

	// Set all to be out of range, which PrintAndValidateConfig() should fix.
	c.APIName = ""
	c.TargetHost = ""
	c.TargetPort = ""
	c.NumIterations = 0
	c.ConcurrentUsers = 0
	c.AllowablePeakMemoryVariance = -1.1
	c.AllowableServiceResponseTimeVariance = -1.1
	c.TestCaseDir = ""
	c.BaseStatsOutputDir = ""
	c.ReportOutputDir = ""
	c.MemoryEndpoint = ""
	c.RequestDelay = 0
	c.TPSFreq = 0
	c.RampUsers = -3
	c.RampDelay = 0

	c.PrintAndValidateConfig()

	// Assert that PrintAndValidateConfig() fixed the out of range values.
	assert.Equal(t, defaultAPIName, c.APIName)
	assert.Equal(t, defaultTargetHost, c.TargetHost)
	assert.Equal(t, defaultTargetPort, c.TargetPort)
	assert.Equal(t, defaultNumIterations, c.NumIterations)
	assert.Equal(t, defaultConcurrentUsers, c.ConcurrentUsers)
	assert.Equal(t, defaultAllowablePeakMemoryVariance, c.AllowablePeakMemoryVariance)
	assert.Equal(t, defaultAllowableServiceResponseTimeVariance, c.AllowableServiceResponseTimeVariance)
	assert.Equal(t, defaultTestCaseDir, c.TestCaseDir)
	assert.Equal(t, defaultBaseStatsOutputDir, c.BaseStatsOutputDir)
	assert.Equal(t, defaultReportOutputDir, c.ReportOutputDir)
	assert.Equal(t, defaultMemoryEndpoint, c.MemoryEndpoint)
	assert.Equal(t, defaultRequestDelay, c.RequestDelay)
	assert.Equal(t, defaultTPSFreq, c.TPSFreq)
	assert.Equal(t, defaultRampUsers, c.RampUsers)
	assert.Equal(t, defaultRampDelay, c.RampDelay)
}
