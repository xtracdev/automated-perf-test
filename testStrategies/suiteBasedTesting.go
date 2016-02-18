package testStrategies

import (
	"fmt"
	"github.com/xtracdev/automated-perf-test/perfTestUtils"
	"sync"
)

func ExecuteTestSuiteWrapper(testSuite *TestSuite, configurationSettings *perfTestUtils.Config) map[string][]int64 {
	allServicesResponseTimesMap := make(map[string][]int64, 0)

	testSuiteResponseTimesChan := make(chan []map[string]int64, 1)

	var suiteWaitGroup sync.WaitGroup
	suiteWaitGroup.Add(configurationSettings.ConcurrentUsers)
	for i := 0; i < configurationSettings.ConcurrentUsers; i++ {
		go executeTestSuite(testSuiteResponseTimesChan, testSuite, configurationSettings)
		go aggregateSuiteResponseTimes(testSuiteResponseTimesChan, allServicesResponseTimesMap, &suiteWaitGroup)
	}
	suiteWaitGroup.Wait()

	fmt.Println("Waitgrouop done")
	fmt.Println("allServicesResponseTimesMap", allServicesResponseTimesMap)
	return allServicesResponseTimesMap
}

func executeTestSuite(testSuiteResponseTimesChan chan []map[string]int64, testSuite *TestSuite, configurationSettings *perfTestUtils.Config) {
	allSuiteResponseTimes := make([]map[string]int64, 0)
	for i := 0; i < configurationSettings.NumIterations; i++ {
		testSuiteResponseTimes := make(map[string]int64)
		for _, testDefinition := range testSuite.TestCases {
			responseTime := testDefinition.BuildAndSendRequest(configurationSettings.TargetHost, configurationSettings.TargetPort)
			testSuiteResponseTimes[testDefinition.TestName] = responseTime
		}
		allSuiteResponseTimes = append(allSuiteResponseTimes, testSuiteResponseTimes)
	}
	testSuiteResponseTimesChan <- allSuiteResponseTimes
}

func aggregateSuiteResponseTimes(testSuiteResponseTimesChan chan []map[string]int64, allServicesResponseTimesMap map[string][]int64, suiteWaitGroup *sync.WaitGroup) {
	perUserSuiteResponseTimes := <-testSuiteResponseTimesChan
	fmt.Println("allServicesResponseTimesMap:", allServicesResponseTimesMap)
	for _, singleSuiteRunResponseTimes := range perUserSuiteResponseTimes {
		for serviceName, serviceResponseTime := range singleSuiteRunResponseTimes {
			if allServicesResponseTimesMap[serviceName] == nil {
				serviceResponseSlice := make([]int64, 0)
				allServicesResponseTimesMap[serviceName] = serviceResponseSlice
			}
			allServicesResponseTimesMap[serviceName] = append(allServicesResponseTimesMap[serviceName], serviceResponseTime)
		}
	}
	fmt.Println("allServicesResponseTimesMap:", allServicesResponseTimesMap)
	suiteWaitGroup.Done()
}
