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
    <apiName>Xtrac API</apiName>

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
	n, err = m.r.Read(p)
	return
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
