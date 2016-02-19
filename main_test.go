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

	tomlConfigFile = `apiName = "Test API"

#Target API under test
targetHost = "localhost"
targetPort = "8282"

#Allowed vairance as a percentage over base values
allowablePeakMemoryVariance = 15.0
allowableServiceResponseTimeVariance = 15.0

#The number of time each test case is executed
numIterations = 1000

#Location of directory where test cases reside
testDefinitionsDir = "./testDefinitions"

#Output locations for generated files
baseStatsOutputDir = "./envStats"
reportOutputDir = "./report"

#Concurrent users per service test
concurrentUsers = 1

#optional template file location
#reportTemplateFile = .`
)

var mockedFs perfTestUtils.FileSystem = mockFs{}

type mockFs struct{}

func (mockFs) Open(name string) (perfTestUtils.File, error) {
	if strings.Contains(name, "FAIL") {
		return &mockedFile{Content: []byte("requested mock FAIL!")}, nil
	}
	switch name {
	case "toml":
		return &mockedFile{Content: []byte(tomlConfigFile)}, nil
	default:
		return &mockedFile{Content: []byte(configFile)}, nil
	}

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
	assert.Equal(t, "./", configurationSettings.ReportOutputDir)
	assert.Equal(t, 1, configurationSettings.ConcurrentUsers)
	assert.Equal(t, "", configurationSettings.TestSuite)
	assert.Equal(t, "xml", configurationSettings.TestFileFormat)

	assert.True(t, configurationSettings.ReBaseMemory)
	assert.True(t, configurationSettings.ReBaseAll)
	assert.True(t, configurationSettings.GBS)
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

/*func TestInitConfigFileToml(t *testing.T) {
	configurationSettings = &perfTestUtils.Config{}
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	willCallOsExit := false
	exit := func(i int) { willCallOsExit = true }

	args := []string{"-gbs", "-reBaseMemory=true", "-configFilePath=toml", "-configFileFormat=toml"}

	initConfig(args, mockedFs, exit)
	assert.NotNil(t, configurationSettings)
	assert.False(t, willCallOsExit)
	assert.Equal(t, "localhost", configurationSettings.TargetHost)
	assert.Equal(t, "8282", configurationSettings.TargetPort)
	assert.True(t, configurationSettings.ReBaseMemory)
	assert.False(t, configurationSettings.ReBaseAll)
	assert.True(t, configurationSettings.GBS)
	assert.Equal(t, "./testDefinitions", configurationSettings.TestCaseDir)
}*/

func TestValidateBasePerfStat(t *testing.T) {
	bs := &perfTestUtils.BasePerfStats{}
	assert.False(t, validateBasePerfStat(bs))

	bs.BaseServiceResponseTimes = map[string]int64{"service 1": 123, "service 2": -1}
	assert.False(t, validateBasePerfStat(bs))

	bs.BaseServiceResponseTimes = map[string]int64{"service 1": 123, "service 2": 321}
	bs.BasePeakMemory = 12
	bs.GenerationDate = "aaa"
	bs.ModifiedDate = "bbb"
	bs.MemoryAudit = []uint64{1, 2, 3}
	assert.True(t, validateBasePerfStat(bs))
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
