package perfTestUtils

import (
	"github.com/stretchr/testify/assert"
	"math/rand"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestGenerateTemplate(t *testing.T) {

	ps := &PerfStats{
		TestDate:   time.Now(),
		PeakMemory: 10e6,
	}

	bs := &BasePerfStats{
		BasePeakMemory: 10e6 - 10e3,
	}

	c := &Config{
		AllowablePeakMemoryVariance:          15,
		AllowableServiceResponseTimeVariance: 15,
		APIName:    "TEST",
		TargetHost: "localhost",
		TargetPort: "8080",
	}

	bma := make([]uint64, 0)
	for i := 0; i < 100; i++ {
		bma = append(bma, bs.BasePeakMemory+uint64(10e5*rand.Int()))
	}
	bs.MemoryAudit = bma
	bsrt := make(map[string]int64)
	bsrt["service 1"] = 3e6
	bsrt["service 2"] = 2e6
	bsrt["service 3"] = 4e6
	bs.BaseServiceResponseTimes = bsrt

	pma := make([]uint64, 0)
	for i := 0; i < 100; i++ {
		pma = append(pma, ps.PeakMemory+uint64(5e5*rand.Int()))
	}
	ps.MemoryAudit = pma

	ps.TestPartitions = []TestPartition{TestPartition{Count: 0, TestName: "StartUp"}, TestPartition{Count: 30, TestName: "service 1"}, TestPartition{Count: 60, TestName: "service 2"}, TestPartition{Count: 90, TestName: "service 3"}}

	psrt := make(map[string]int64)
	psrt["service 1"] = 3e5
	psrt["service 2"] = 2e5
	ps.ServiceResponseTimes = psrt

	gopath := os.Getenv("GOPATH")
	t.Logf("$GOPATH = %v\n", gopath)

	tf := gopath + `/src/github.com/xtracdev/automated-perf-test/report/test.tmpl`
	t.Logf("template = %v\n", filepath.Base(tf))

	err := generateTemplate(bs, ps, c, os.Stdout, gopath+`/src/github.com/xtracdev/automated-perf-test/report/test.tmpl`)
	if err != nil {
		t.Errorf("Expected to be nil: %v", err)
	}

}

func TestGenerateTemplateBuiltin(t *testing.T) {

	ps := &PerfStats{
		TestDate:   time.Now(),
		PeakMemory: 10e6,
	}

	bs := &BasePerfStats{
		BasePeakMemory: 10e6 - 10e3,
	}

	c := &Config{
		AllowablePeakMemoryVariance:          15,
		AllowableServiceResponseTimeVariance: 15,
		APIName:    "TEST",
		TargetHost: "localhost",
		TargetPort: "8080",
	}

	bma := make([]uint64, 0)
	for i := 0; i < 100; i++ {
		bma = append(bma, bs.BasePeakMemory+uint64(10e5*rand.Int()))
	}
	bs.MemoryAudit = bma
	bsrt := make(map[string]int64)
	bsrt["service 1"] = 3e6
	bsrt["service 2"] = 2e6
	bsrt["service 3"] = 4e6
	bs.BaseServiceResponseTimes = bsrt

	pma := make([]uint64, 0)
	for i := 0; i < 100; i++ {
		pma = append(pma, ps.PeakMemory+uint64(5e5*rand.Int()))
	}
	ps.MemoryAudit = pma

	ps.TestPartitions = []TestPartition{TestPartition{Count: 0, TestName: "StartUp"}, TestPartition{Count: 30, TestName: "service 1"}, TestPartition{Count: 60, TestName: "service 2"}, TestPartition{Count: 90, TestName: "service 3"}}

	psrt := make(map[string]int64)
	psrt["service 1"] = 3e5
	psrt["service 2"] = 2e5
	ps.ServiceResponseTimes = psrt

	gopath := os.Getenv("GOPATH")
	t.Logf("$GOPATH = %v\n", gopath)

	err := generateTemplate(bs, ps, c, os.Stdout, "")
	if err != nil {
		t.Errorf("Expected to be nil: %v", err)
	}
}

func TestGenerateTemplateNoFile(t *testing.T) {
	err := generateTemplate(nil, nil, nil, os.Stdout, "XXX")
	assert.NotNil(t, err)
	assert.Equal(t, "Error loading template files: open XXX: no such file or directory", err.Error())
}

func TestIsServiceTimePassNoRespTimes(t *testing.T) {
	pm := &perfStatsModel{}
	pm.Config = &Config{AllowableServiceResponseTimeVariance: float64(15)}
	pm.PerfStats = &PerfStats{
		ServiceResponseTimes: map[string]int64{},
	}
	assert.False(t, pm.IsServiceTimePass("test"))
}

func TestIsServiceTimePassNotPassing(t *testing.T) {
	pm := &perfStatsModel{}
	pm.Config = &Config{AllowableServiceResponseTimeVariance: float64(15)}
	pm.PerfStats = &PerfStats{
		ServiceResponseTimes: map[string]int64{"service 1": 100},
	}
	pm.BasePerfStats = &BasePerfStats{
		BaseServiceResponseTimes: map[string]int64{"service 1": 10},
	}
	assert.False(t, pm.IsServiceTimePass("service 1"))
}

func TestIsServiceTimePassOk(t *testing.T) {
	pm := &perfStatsModel{}
	pm.Config = &Config{AllowableServiceResponseTimeVariance: float64(15)}
	pm.PerfStats = &PerfStats{
		ServiceResponseTimes: map[string]int64{"service 1": 100},
	}
	pm.BasePerfStats = &BasePerfStats{
		BaseServiceResponseTimes: map[string]int64{"service 1": 101},
	}
	assert.True(t, pm.IsServiceTimePass("service 1"))
}

func TestIsTimePassOk(t *testing.T) {
	pm := &perfStatsModel{}
	pm.Config = &Config{AllowableServiceResponseTimeVariance: float64(15)}
	pm.PerfStats = &PerfStats{
		ServiceResponseTimes: map[string]int64{"service 1": 100},
	}
	pm.BasePerfStats = &BasePerfStats{
		BaseServiceResponseTimes: map[string]int64{"service 1": 101},
	}
	assert.True(t, pm.IsTimePass())
}

func TestIsTimePassFail(t *testing.T) {
	pm := &perfStatsModel{}
	pm.Config = &Config{AllowableServiceResponseTimeVariance: float64(15)}
	pm.PerfStats = &PerfStats{
		ServiceResponseTimes: map[string]int64{"service 1": 100, "service 2": 100},
	}
	pm.BasePerfStats = &BasePerfStats{
		BaseServiceResponseTimes: map[string]int64{"service 1": 101, "service 2": 80},
	}
	assert.False(t, pm.IsTimePass())
}

func TestGenerateTemplateReport(t *testing.T) {
	c := &Config{AllowableServiceResponseTimeVariance: float64(15)}
	ps := &PerfStats{
		ServiceResponseTimes: map[string]int64{"service 1": 100, "service 2": 100},
	}
	bs := &BasePerfStats{
		BaseServiceResponseTimes: map[string]int64{"service 1": 101, "service 2": 80},
	}
	GenerateTemplateReport(bs, ps, c, mockedFs, "TestSuiteName")
}

func TestGenerateTemplateReportErrorCreate(t *testing.T) {
	c := &Config{AllowableServiceResponseTimeVariance: float64(15), ReportOutputDir: "FAIL"}
	ps := &PerfStats{
		ServiceResponseTimes: map[string]int64{"service 1": 100, "service 2": 100},
	}
	bs := &BasePerfStats{
		BaseServiceResponseTimes: map[string]int64{"service 1": 101, "service 2": 80},
	}
	GenerateTemplateReport(bs, ps, c, mockedFs, "TestSuiteName")
}
