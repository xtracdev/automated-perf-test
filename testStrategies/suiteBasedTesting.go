package testStrategies

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/xtracdev/automated-perf-test/perfTestUtils"
	"sync"
	"time"
	"sync/atomic"
)




//----- ExecuteTestSuiteWrapper ----------------------------------------------
func ExecuteTestSuiteWrapper(
		testSuite *TestSuite,
		configurationSettings *perfTestUtils.Config,
		startTime time.Time,
) map[string][]int64 {
	allServicesResponseTimesMap := make(map[string][]int64, 0)
	testSuiteResponseTimesChan := make(chan []map[string]int64, 1)
	var suiteWaitGroup sync.WaitGroup

	// Variables to track total current TPS across all threads.
	// We'll use an unsigned integer to represent our (always-positive) counter,
	// and use a wait group from the sync package to track a concurrent counter.
	var curOps uint64 = 0
	quitShowTPSChan := make(chan bool)

	suiteWaitGroup.Add(configurationSettings.ConcurrentUsers)
	for i := 0; i < configurationSettings.ConcurrentUsers; i++ {
		go executeTestSuite(testSuiteResponseTimesChan, testSuite, configurationSettings, i, GlobalsLockCounter, &curOps)
		go aggregateSuiteResponseTimes(testSuiteResponseTimesChan, allServicesResponseTimesMap, &suiteWaitGroup)
	}

	// Display the ongoing TPS to log.Info based on period specified in configurationSettings.TPSFreq:
	go showCurrentTPS(quitShowTPSChan, &curOps, startTime, configurationSettings)

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
		curOps *uint64,
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

			// Increment the global ops counter
			atomic.AddUint64( curOps, 1 )
		}

		allSuiteResponseTimes = append(allSuiteResponseTimes, testSuiteResponseTimes)

		globalsMap.Lock()
		globalsMap.m[uniqueTestRunId] = nil
		globalsMap.Unlock()
	}
	testSuiteResponseTimesChan <- allSuiteResponseTimes
	log.Infof("Test Suite [%s::%s] Finished", testSuite.Name, uniqueTestRunId)
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
func showCurrentTPS( chQuit chan bool, curOps *uint64, startTime time.Time, configurationSettings *perfTestUtils.Config ) {
	for {
		// Concurrent controls:
		select {
		case <-chQuit:
			return
		default:
			// We only want one output line during any given second.
			// This effectively sets the lower bound for TPSFreq.
			time.Sleep(time.Second)

			// No need to display until at least one operation has completed.
			if (*curOps < uint64(1)) {
				break
			}

			elapsedTime := time.Since(startTime)
			// No need to display if not within the period set in config:
			if (int64(elapsedTime.Seconds()) % int64(configurationSettings.TPSFreq) != 0) {
				break
			}

			// Print the display.
			num_ops := atomic.LoadUint64(curOps)
			tps := 0.0
			if (int(elapsedTime.Seconds()) > 0) {
				tps = (float64(num_ops) / elapsedTime.Seconds())
			}
			log.Infof("[showCurrentTPS] {\"curOps\":\"%d\",\"elapsedTime\":\"%f\",\"TPS\":\"%f\"}",
				num_ops,
				elapsedTime.Seconds(),
				tps,
			)
		}
	}
}
