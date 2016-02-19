package testStrategies

import (
	"github.com/stretchr/testify/assert"
	"github.com/xtracdev/automated-perf-test/perfTestUtils"
	"sync"
	"testing"
)

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
