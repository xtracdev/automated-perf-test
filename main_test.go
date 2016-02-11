package main

import (
	"flag"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/xtracdev/automated-perf-test/perfTestUtils"
	"io"
	"os"
	"strings"
	"sync"
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
    <testDefinitionsDir>./testDefinitions</testDefinitionsDir>

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

func TestInitConfigFileNotFound(t *testing.T) {
	willCallOsExit := false
	exit := func(i int) { willCallOsExit = true }

	args := []string{"-gbs", "-reBaseMemory", "-reBaseAll", "-configFilePath=test"}

	initConfig(args, osFileSystem, exit)
	assert.NotNil(t, configurationSettings)
	assert.True(t, willCallOsExit)
}

func TestInitConfigFileUnmarshalErr(t *testing.T) {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	willCallOsExit := false
	exit := func(i int) { willCallOsExit = true }

	args := []string{"-gbs", "-reBaseMemory", "-reBaseAll", "-configFilePath=FAIL"}

	initConfig(args, mockedFs, exit)
	assert.NotNil(t, configurationSettings)
	assert.True(t, willCallOsExit)
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
	assert.Equal(t, "./testDefinitions", configurationSettings.TestDefinitionsDir)
}

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

func TestAggregateResponseTimes(t *testing.T) {
	var wg sync.WaitGroup

	srtChan := make(chan perfTestUtils.RspTimes)
	testChan := make(chan perfTestUtils.RspTimes)
	times := &[]int64{1, 2, 3, 4}
	go func() {
		for i := 0; i < 5; i++ {
			wg.Add(1)
			go func() {
				srtChan <- []int64{10, 20}
			}()
			go aggregateResponseTimes(times, srtChan, &wg)
		}
		wg.Wait()
		testChan <- *times
	}()

	go func() {
		toTest := <-testChan
		assert.Equal(t, 14, len(toTest))
	}()
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
