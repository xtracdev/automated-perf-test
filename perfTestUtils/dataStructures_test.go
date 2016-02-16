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
	assert.Equal(t, defaultTestDefinitionsDir, c.TestDefinitionsDir)
	assert.Equal(t, defaultBaseStatsOutputDir, c.BaseStatsOutputDir)
	assert.Equal(t, defaultReportOutputDir, c.ReportOutputDir)
	assert.Equal(t, defaultConcurrentUsers, c.ConcurrentUsers)
	assert.Equal(t, false, c.GBS)
	assert.Equal(t, false, c.ReBaseMemory)
	assert.Equal(t, false, c.ReBaseAll)
}

func TestPrintAndValidateConfig(t *testing.T) {
	willCallOsExit := false
	exit := func(i int) { willCallOsExit = true }
	c := &Config{}
	c.SetDefaults()
	c.PrintAndValidateConfig(exit)
	assert.False(t, willCallOsExit)
}

func TestPrintAndValidateConfigErr(t *testing.T) {
	willCallOsExit := false
	exit := func(i int) { willCallOsExit = true }
	c := &Config{}
	c.PrintAndValidateConfig(exit)
	assert.True(t, willCallOsExit)
}
