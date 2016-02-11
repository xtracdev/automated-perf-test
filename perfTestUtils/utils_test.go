package perfTestUtils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"
)

var mockedFs FileSystem = mockFs{}

type mockFs struct{}

func (mockFs) Open(name string) (File, error) { return &mockedFile{}, nil }
func (mockFs) Create(name string) (File, error) {
	if strings.Contains(name, "FAIL") {
		return nil, fmt.Errorf("requested mock FAIL!")
	}
	return &mockedFile{}, nil
}

type mockedFile struct {
	Content []byte
	r       *strings.Reader
}

func (mockedFile) Readdir(n int) (fi []os.FileInfo, err error) {
	if n == -1 {
		return make([]os.FileInfo, 10), nil
	} else {
		return nil, fmt.Errorf("Mock dir error!")
	}
}
func (mockedFile) Close() error { return nil }

func (mockedFile) Write(p []byte) (n int, err error) { return io.WriteString(os.Stdout, string(p)) }
func (m *mockedFile) Read(p []byte) (n int, err error) {
	if m.r == nil {
		m.r = strings.NewReader(string(m.Content))
	}
	return m.r.Read(p)
}

func TestReadBasePerfFile(t *testing.T) {
	bs := &BasePerfStats{
		BasePeakMemory: 10e6 - 10e3,
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

	b, err := json.Marshal(bs)
	if err != nil {
		t.Errorf("expected to be nil: %v\n", err)
	}
	t.Logf("%s\n", b)
	toTest, err := ReadBasePerfFile(bytes.NewReader(b))
	assert.Nil(t, err)
	assert.NotNil(t, toTest)
	assert.IsType(t, new(BasePerfStats), toTest)
	assert.Equal(t, 100, len(toTest.MemoryAudit))
}

func TestReadBasePerfFileErrUnmarshal(t *testing.T) {
	toTest, err := ReadBasePerfFile(bytes.NewReader([]byte("test")))
	assert.NotNil(t, err)
	assert.NotNil(t, toTest)
	assert.IsType(t, new(BasePerfStats), toTest)
	assert.Equal(t, 0, len(toTest.MemoryAudit))
	assert.Equal(t, `invalid character 'e' in literal true (expecting 'r')`, err.Error())
}

func TestValidateTestDefinitionAmount(t *testing.T) {
	c := &Config{
		TestDefinitionsDir: "mockDir",
	}
	valid := ValidateTestDefinitionAmount(10, c, mockedFs)
	assert.True(t, valid)
}

func TestValidateTestDefinitionAmountFail(t *testing.T) {
	c := &Config{
		TestDefinitionsDir: "mockDir",
	}
	valid := ValidateTestDefinitionAmount(100, c, mockedFs)
	assert.False(t, valid)
}

func TestCalcPeakMemoryVariancePercentage(t *testing.T) {
	vp := CalcPeakMemoryVariancePercentage(100, 110)
	assert.Equal(t, float64(10), vp)

	vp = CalcPeakMemoryVariancePercentage(100, 90)
	assert.Equal(t, float64(-10), vp)
}

func BenchmarkCalcPeakMemoryVariancePercentage(t *testing.B) {
	for i := 0; i < t.N; i++ {
		CalcPeakMemoryVariancePercentage(100, 90)
	}
}

func TestCalcAverageResponseTime(t *testing.T) {
	times := make([]int64, 0)
	for i := int64(200); i >= 0; i-- {
		times = append(times, i*1000)
	}
	avg := CalcAverageResponseTime(times, 100)
	assert.Equal(t, int64(201157), avg)
}

func TestCalcAverageResponseVariancePercentage(t *testing.T) {
	vp := CalcAverageResponseVariancePercentage(110, 100)
	assert.Equal(t, float64(10), vp)

	vp = CalcAverageResponseVariancePercentage(90, 100)
	assert.Equal(t, float64(-10), vp)
}

func BenchmarkCalcAverageResponseVariancePercentage(t *testing.B) {
	for i := 0; i < t.N; i++ {
		CalcAverageResponseVariancePercentage(100, 90)
	}
}

func TestPopulateBasePerfStats(t *testing.T) {
	ps := &PerfStats{
		TestDate:   time.Now(),
		PeakMemory: 10e6,
	}

	bs := &BasePerfStats{
		BaseServiceResponseTimes: make(map[string]int64),
	}

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

	populateBasePerfStats(ps, bs, false)
	assert.Equal(t, bs.BasePeakMemory, ps.PeakMemory)
	assert.Equal(t, bs.MemoryAudit, ps.MemoryAudit)
	assert.Equal(t, bs.BaseServiceResponseTimes, ps.ServiceResponseTimes)
	assert.Equal(t, bs.ModifiedDate, bs.GenerationDate)
}

func TestValidateResponseBody(t *testing.T) {
	b := []byte("content")
	empty := []byte{}

	assert.True(t, ValidateResponseBody(b, "test"))
	assert.False(t, ValidateResponseBody(empty, "test"))
}

func TestValidateResponseStatusCode(t *testing.T) {
	assert.True(t, ValidateResponseStatusCode(http.StatusOK, http.StatusOK, "test"))
	assert.False(t, ValidateResponseStatusCode(http.StatusOK, http.StatusInternalServerError, "test"))
}

func TestValidateServiceResponseTime(t *testing.T) {
	assert.True(t, ValidateServiceResponseTime(10, "test"))
	assert.False(t, ValidateServiceResponseTime(-10, "test"))
	assert.False(t, ValidateServiceResponseTime(0, "test"))
}

func TestValidatePeakMemoryVariance(t *testing.T) {
	assert.True(t, ValidatePeakMemoryVariance(15, 0.1))
	assert.False(t, ValidatePeakMemoryVariance(15, 16.5))
	assert.True(t, ValidatePeakMemoryVariance(15, 15))
}

func TestValidateAverageServiceResponeTimeVariance(t *testing.T) {
	assert.True(t, ValidateAverageServiceResponeTimeVariance(15, 10, "test"))
	assert.True(t, ValidateAverageServiceResponeTimeVariance(15, 15, "test"))
	assert.False(t, ValidateAverageServiceResponeTimeVariance(15, 16, "test"))
}

func TestGenerateEnvBasePerfOutputFile(t *testing.T) {
	ps := &PerfStats{
		TestDate:   time.Now(),
		PeakMemory: 10e6,
	}

	bs := &BasePerfStats{
		BaseServiceResponseTimes: make(map[string]int64),
	}

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

	willCallOsExit := false
	exit := func(i int) { willCallOsExit = true }

	GenerateEnvBasePerfOutputFile(ps, bs, &Config{ReBaseMemory: true, BaseStatsOutputDir: "env", TargetHost: "localhost"}, exit, mockedFs)
	assert.False(t, willCallOsExit)
}

func TestGenerateEnvBasePerfOutputFileFailCreate(t *testing.T) {
	ps := &PerfStats{
		TestDate:   time.Now(),
		PeakMemory: 10e6,
	}

	bs := &BasePerfStats{
		BaseServiceResponseTimes: make(map[string]int64),
	}

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

	willCallOsExit := false
	exit := func(i int) { willCallOsExit = true }

	GenerateEnvBasePerfOutputFile(ps, bs, &Config{ReBaseMemory: true, BaseStatsOutputDir: "env", ExecutionHost: "FAIL"}, exit, mockedFs)
	assert.True(t, willCallOsExit)
}
