package perfTestUtils

import (
	"math/rand"
	"os"
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

	err := generateTemplate(bs, ps, c, os.Stdout, gopath+`/src/github.com/xtracdev/automated-perf-test/report`)
	if err != nil {
		t.Errorf("Expected to be nil: %v", err)
	}

}
