package testStrategies

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/xtracdev/automated-perf-test/perfTestUtils"
	"sync"
)

func ExecuteTestSuiteWrapper(testSuite *TestSuite, configurationSettings *perfTestUtils.Config) map[string][]int64 {
	allServicesResponseTimesMap := make(map[string][]int64, 0)

	testSuiteResponseTimesChan := make(chan []map[string]int64, 1)

	var suiteWaitGroup sync.WaitGroup
	suiteWaitGroup.Add(configurationSettings.ConcurrentUsers)
	for i := 0; i < configurationSettings.ConcurrentUsers; i++ {
		go executeTestSuite(testSuiteResponseTimesChan, testSuite, configurationSettings, i, GlobalsLockCounter)
		go aggregateSuiteResponseTimes(testSuiteResponseTimesChan, allServicesResponseTimesMap, &suiteWaitGroup)
	}
	suiteWaitGroup.Wait()
	return allServicesResponseTimesMap
}

func executeTestSuite(testSuiteResponseTimesChan chan []map[string]int64, testSuite *TestSuite, configurationSettings *perfTestUtils.Config, userId int, globalsMap GlobalsMaps) {
	log.Info("Test Suite started")

	allSuiteResponseTimes := make([]map[string]int64, 0)
	for i := 0; i < configurationSettings.NumIterations; i++ {
		uniqueTestRunId := fmt.Sprintf("User%dIter%d", userId, i)
		testSuiteResponseTimes := make(map[string]int64)
		for _, testDefinition := range testSuite.TestCases {
			log.Info("Test case: [", testDefinition.TestName, "] UniqueRunID: [", uniqueTestRunId, "]")

			targetHost, targetPort := determineHostandPortforRequest(testDefinition, configurationSettings)

			responseTime := testDefinition.BuildAndSendRequest(configurationSettings.RequestDelay, targetHost, targetPort, uniqueTestRunId, globalsMap)
			testSuiteResponseTimes[testDefinition.TestName] = responseTime
		}
		allSuiteResponseTimes = append(allSuiteResponseTimes, testSuiteResponseTimes)

		globalsMap.Lock()
		globalsMap.m[uniqueTestRunId] = nil
		globalsMap.Unlock()
	}
	testSuiteResponseTimesChan <- allSuiteResponseTimes
	log.Infof("Test Suite [%s] Finished", testSuite.Name)
}

func aggregateSuiteResponseTimes(testSuiteResponseTimesChan chan []map[string]int64, allServicesResponseTimesMap map[string][]int64, suiteWaitGroup *sync.WaitGroup) {
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
