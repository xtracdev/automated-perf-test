package testStrategies

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/xtracdev/automated-perf-test/perfTestUtils"
	"sync/atomic"
	"time"
	"sync"
)




//----- ExecuteTestSuiteWrapper ----------------------------------------------
func ExecuteTestSuiteWrapper(
		testSuite *TestSuite,
		configSettings *perfTestUtils.Config,
		perfStatsForTest *perfTestUtils.PerfStats,
		scenarioTimeStart time.Time,
) (map[string][]int64) {
	allServicesResponseTimesMap := make(map[string][]int64, 0)
	testSuiteResponseTimesChan := make(chan []map[string]int64, 1)
	var suiteWaitGroup sync.WaitGroup

	// Set concurrency control:
	suiteWaitGroup.Add(configSettings.ConcurrentUsers)

	for i := 0; i < configSettings.ConcurrentUsers; i++ {
		go executeTestSuite(testSuiteResponseTimesChan, testSuite, configSettings, i, GlobalsLockCounter, perfStatsForTest)
		go aggregateSuiteResponseTimes(testSuiteResponseTimesChan, allServicesResponseTimesMap, &suiteWaitGroup)
	}

	// Display the ongoing TPS to log.Info based on period specified in configurationSettings.TPSFreq:
	quitShowTPSChan := make(chan bool)
	go showCurrentTPS(quitShowTPSChan, configSettings, scenarioTimeStart, &perfStatsForTest.OverAllTransCount)

	suiteWaitGroup.Wait()
	quitShowTPSChan <- true

	return allServicesResponseTimesMap
}




//----- executeTestSuite ------------------------------------------------------
func executeTestSuite(
		testSuiteResponseTimesChan chan []map[string]int64,
		testSuite *TestSuite,
		configurationSettings *perfTestUtils.Config,
		userId int,
		globalsMap GlobalsMaps,
		perfStatsForTest *perfTestUtils.PerfStats,
) {
	log.Info("Test Suite started")


	allSuiteResponseTimes := make([]map[string]int64, 0)
	uniqueTestRunId := ""

	for i := 0; i < configurationSettings.NumIterations; i++ {
		// Run all services of the test suite NumIterations of times.
		uniqueTestRunId = fmt.Sprintf("User%dIter%d", userId, i)
		testSuiteResponseTimes := make(map[string]int64)

		for _, testDefinition := range testSuite.TestCases {
			log.Info("Test case: [", testDefinition.TestName, "] UniqueRunID: [", uniqueTestRunId, "]")

			targetHost, targetPort := determineHostandPortforRequest(testDefinition, configurationSettings)
			responseTime := testDefinition.BuildAndSendRequest(configurationSettings.RequestDelay, targetHost, targetPort, uniqueTestRunId, globalsMap)
			testSuiteResponseTimes[testDefinition.TestName] = responseTime

			// Update the concurrent counters.
			// Overall counter:
			atomic.AddUint64(&perfStatsForTest.OverAllTransCount, 1)

			// Service-level counters:
			// Create the counters on the fly and increment atomically.
			if perfStatsForTest.ServiceTransCount[testDefinition.TestName] == nil {
				perfStatsForTest.ServiceTransCount[testDefinition.TestName] = new(uint64)
				atomic.StoreUint64(
					perfStatsForTest.ServiceTransCount[testDefinition.TestName],
					0,
				)
			}
			atomic.AddUint64(
				perfStatsForTest.ServiceTransCount[testDefinition.TestName],
				1,
			)
		}

		allSuiteResponseTimes = append(allSuiteResponseTimes, testSuiteResponseTimes)

		globalsMap.Lock()
		globalsMap.m[uniqueTestRunId] = nil
		globalsMap.Unlock()
	}

	testSuiteResponseTimesChan <- allSuiteResponseTimes
}




//----- aggregateSuiteResponseTimes -------------------------------------------
func aggregateSuiteResponseTimes(
		testSuiteResponseTimesChan chan []map[string]int64,
		allServicesResponseTimesMap map[string][]int64,
		suiteWaitGroup *sync.WaitGroup,
) {
	perUserSuiteResponseTimes := <-testSuiteResponseTimesChan
	for _, singleSuiteRunResponseTimes := range perUserSuiteResponseTimes {
		for serviceName, serviceResponseTime := range singleSuiteRunResponseTimes {
			if allServicesResponseTimesMap[serviceName] == nil {
				serviceResponseSlice := make([]int64, 0)
				allServicesResponseTimesMap[serviceName] = serviceResponseSlice
			}
			allServicesResponseTimesMap[serviceName] = append(allServicesResponseTimesMap[serviceName], serviceResponseTime)
		}
	}
	suiteWaitGroup.Done()
}




//----- showCurrentTPS -------------------------------------------------------------------------------------------------
// Print current TPS progress every period of time defined by configurationSettings.TPSFREQ.
func showCurrentTPS(
		chQuit chan bool,
		confgSettings *perfTestUtils.Config,
		scenarioStartTime time.Time,
		nNumberOfTrans *uint64,
) {
	for {
		// Concurrent controls:
		select {
		case <-chQuit:
			return
		default:
			// Set variables for convenience.
			durElapsedTime := time.Since(scenarioStartTime)
			num_ops := atomic.LoadUint64(nNumberOfTrans)

			// We only want one output line during any given second. This
			// effectively sets the lower bound for TPSFreq to one second.
			time.Sleep(time.Second)

			// No need to display until at least one operation has completed.
			if (num_ops < uint64(1)) {
				break
			}

			// No need to display if not within the period set in config:
			if (int64(durElapsedTime.Seconds()) % int64(confgSettings.TPSFreq) != 0) {
				break
			}

			// Print the display.
			tps := 0.0
			if (int(durElapsedTime.Seconds()) > 0) {
				tps = (float64(num_ops) / durElapsedTime.Seconds())
			}

			log.Infof("[showCurrentTPS] {\"TransCount\":\"%d\",\"ElapsedTime\":\"%v\",\"TPS\":\"%f\"}",
				num_ops,
				durElapsedTime,
				tps,
			)
		}
	}
}
