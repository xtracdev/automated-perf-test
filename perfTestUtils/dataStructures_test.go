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
	assert.Equal(t, false, c.GBS)
	assert.Equal(t, false, c.ReBaseMemory)
	assert.Equal(t, false, c.ReBaseAll)
}
