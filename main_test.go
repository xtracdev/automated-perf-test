package main

import (
	"flag"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/xtracdev/automated-perf-test/perfTestUtils"
	"io"
	"os"
	"strings"
	"testing"
)

const (
	configFile = `<?xml version="1.0" encoding="UTF-8"?>
<config>
    <apiName>Test API</apiName>

    <!--Target API under test-->
    <targetHost>localhost</targetHost>
    <targetPort>8282</targetPort>

    <!--Allowed vairance as a percentage over base values-->
    <allowablePeakMemoryVariance>15</allowablePeakMemoryVariance>
    <allowableServiceResponseTimeVariance>15</allowableServiceResponseTimeVariance>

    <!--THe number of time each test case is executed-->
    <numIterations>1000</numIterations>

    <!--Location of directory where test cases reside-->
    <testCaseDir>./definitions/testCases</testCaseDir>

    <!--Output locations for generated files-->
    <baseStatsOutputDir>./envStats</baseStatsOutputDir>
    <reportOutputDir>./report</reportOutputDir>

    <!-- Concurrent users per service test -->
    <concurrentUsers>1</concurrentUsers>

    <!-- optional template file location
    <reportTemplateFile/>-->
</config>`
)

var mockedFs perfTestUtils.FileSystem = mockFs{}

type mockFs struct{}

func (mockFs) Open(name string) (perfTestUtils.File, error) {
	if strings.Contains(name, "FAIL") {
		return &mockedFile{Content: []byte("requested mock FAIL!")}, nil
	}
	return &mockedFile{Content: []byte(configFile)}, nil
}
func (mockFs) Create(name string) (perfTestUtils.File, error) {
	if strings.Contains(name, "FAIL") {
		return nil, fmt.Errorf("requested mock FAIL!")
	}
	return &mockedFile{}, nil
}

type mockedFile struct {
	Content []byte
	r       *strings.Reader
}

func (*mockedFile) Readdir(n int) (fi []os.FileInfo, err error) {
	if n == -1 {
		return make([]os.FileInfo, 10), nil
	} else {
		return nil, fmt.Errorf("Mock dir error!")
	}
}
func (*mockedFile) Close() error { return nil }

func (*mockedFile) Write(p []byte) (n int, err error) { return io.WriteString(os.Stdout, string(p)) }
func (m *mockedFile) Read(p []byte) (n int, err error) {
	if m.r == nil {
		m.r = strings.NewReader(string(m.Content))
	}
	return m.r.Read(p)
}

func assertDefaultConfig(t *testing.T) {
	assert.NotNil(t, configurationSettings)
	assert.Equal(t, "Default_API_NAME", configurationSettings.APIName)
	assert.Equal(t, "localhost", configurationSettings.TargetHost)
	assert.Equal(t, "8080", configurationSettings.TargetPort)
	assert.Equal(t, 1000, configurationSettings.NumIterations)
	assert.Equal(t, float64(15), configurationSettings.AllowablePeakMemoryVariance)
	assert.Equal(t, float64(15), configurationSettings.AllowableServiceResponseTimeVariance)
	assert.Equal(t, "./definitions/testCases", configurationSettings.TestCaseDir)
	assert.Equal(t, "./definitions/testSuites", configurationSettings.TestSuiteDir)
	assert.Equal(t, "./envStats", configurationSettings.BaseStatsOutputDir)
	assert.Equal(t, "./report", configurationSettings.ReportOutputDir)
	assert.Equal(t, 1, configurationSettings.ConcurrentUsers)
	assert.Equal(t, "", configurationSettings.TestSuite)

	assert.True(t, configurationSettings.ReBaseMemory)
	assert.True(t, configurationSettings.ReBaseAll)
	assert.True(t, configurationSettings.GBS)
}

func TestOverrideConfig(t *testing.T) {
	configurationSettings = new(perfTestUtils.Config)
	configurationSettings.SetDefaults()

	configOverrides = new(perfTestUtils.Config)

	configOverrides.APIName = "1"
	configOverrides.TargetHost = "2"
	configOverrides.TargetPort = "3"
	configOverrides.NumIterations = 4
	configOverrides.AllowablePeakMemoryVariance = 5.0
	configOverrides.AllowableServiceResponseTimeVariance = 6.0
	configOverrides.TestCaseDir = "7"
	configOverrides.TestSuiteDir = "8"
	configOverrides.BaseStatsOutputDir = "9"
	configOverrides.ReportOutputDir = "10"
	configOverrides.ConcurrentUsers = 11
	configOverrides.TestSuite = "12"
	configOverrides.MemoryEndpoint = "13"
	configOverrides.RequestDelay = 14
	configOverrides.TPSFreq = 15
	configOverrides.RampUsers = 16
	configOverrides.RampDelay = 17

	overrideConfigOpts()

	assert.Equal(t, "1", configurationSettings.APIName)
	assert.Equal(t, "2", configurationSettings.TargetHost)
	assert.Equal(t, "3", configurationSettings.TargetPort)
	assert.Equal(t, 4, configurationSettings.NumIterations)
	assert.Equal(t, 5.0, configurationSettings.AllowablePeakMemoryVariance)
	assert.Equal(t, 6.0, configurationSettings.AllowableServiceResponseTimeVariance)
	assert.Equal(t, "7", configurationSettings.TestCaseDir)
	assert.Equal(t, "8", configurationSettings.TestSuiteDir)
	assert.Equal(t, "9", configurationSettings.BaseStatsOutputDir)
	assert.Equal(t, "10", configurationSettings.ReportOutputDir)
	assert.Equal(t, 11, configurationSettings.ConcurrentUsers)
	assert.Equal(t, "12", configurationSettings.TestSuite)
	assert.Equal(t, "13", configurationSettings.MemoryEndpoint)
	assert.Equal(t, 14, configurationSettings.RequestDelay)
	assert.Equal(t, 15, configurationSettings.TPSFreq)
	assert.Equal(t, 16, configurationSettings.RampUsers)
	assert.Equal(t, 17, configurationSettings.RampDelay)
}

func TestInitConfigFileNotFound(t *testing.T) {
	willCallOsExit := false
	exit := func(i int) { willCallOsExit = false }

	args := []string{"-gbs", "-reBaseMemory", "-reBaseAll", "-configFilePath=test"}

	initConfig(args, osFileSystem, exit)
	assertDefaultConfig(t)
	assert.False(t, willCallOsExit)
}

func TestInitConfigFileUnmarshalErr(t *testing.T) {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	willCallOsExit := false
	exit := func(i int) { willCallOsExit = false }

	args := []string{"-gbs", "-reBaseMemory", "-reBaseAll", "-configFilePath=FAIL"}

	initConfig(args, mockedFs, exit)
	assertDefaultConfig(t)
	assert.False(t, willCallOsExit)
}

func TestInitConfigFile(t *testing.T) {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	willCallOsExit := false
	exit := func(i int) { willCallOsExit = true }

	args := []string{"-gbs", "-reBaseMemory=true", "-configFilePath=test"}

	initConfig(args, mockedFs, exit)
	assert.NotNil(t, configurationSettings)
	assert.False(t, willCallOsExit)
	assert.Equal(t, "localhost", configurationSettings.TargetHost)
	assert.Equal(t, "8282", configurationSettings.TargetPort)
	assert.True(t, configurationSettings.ReBaseMemory)
	assert.False(t, configurationSettings.ReBaseAll)
	assert.True(t, configurationSettings.GBS)
	assert.Equal(t, "./definitions/testCases", configurationSettings.TestCaseDir)
}

func TestRunAssertions(t *testing.T) {
	bs := &perfTestUtils.BasePerfStats{
		BasePeakMemory:           100,
		BaseServiceResponseTimes: map[string]int64{"s1": 10, "s2": 20, "s3": 30},
	}

	ps := &perfTestUtils.PerfStats{
		PeakMemory:           150,
		ServiceResponseTimes: map[string]int64{"s1": 100, "s2": 20},
	}
	configurationSettings = new(perfTestUtils.Config)
	configurationSettings.SetDefaults()
	toTest := runAssertions(bs, ps)
	t.Logf("%v\n", toTest)
	assert.Equal(t, 3, len(toTest))
}
