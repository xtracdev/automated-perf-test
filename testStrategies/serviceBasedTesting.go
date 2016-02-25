package testStrategies

import (
	"github.com/xtracdev/automated-perf-test/perfTestUtils"
	"sync"
)

//Single execution function for all service test.
//Runs multiple invocations of the test based on num iterations parameter
func ExecuteServiceTest(testDefinition *TestDefinition, loadPerUser int, remainder int, configurationSettings *perfTestUtils.Config) int64 {

	averageResponseTime := int64(0)
	responseTimes := make([]int64, 0)
	subsetOfResponseTimesChan := make(chan perfTestUtils.RspTimes, 1)

	var wg sync.WaitGroup
	wg.Add(configurationSettings.ConcurrentUsers)
	for i := 0; i < configurationSettings.ConcurrentUsers; i++ {
		go buildAndSendUserRequests(subsetOfResponseTimesChan, loadPerUser, testDefinition, configurationSettings.TargetHost, configurationSettings.TargetPort, globals)
		go aggregateResponseTimes(&responseTimes, subsetOfResponseTimesChan, &wg)
	}
	if remainder > 0 {
		wg.Add(1)
		go buildAndSendUserRequests(subsetOfResponseTimesChan, remainder, testDefinition, configurationSettings.TargetHost, configurationSettings.TargetPort, globals)
		go aggregateResponseTimes(&responseTimes, subsetOfResponseTimesChan, &wg)
	}

	wg.Wait()

	if len(responseTimes) == configurationSettings.NumIterations {
		averageResponseTime = perfTestUtils.CalcAverageResponseTime(responseTimes, configurationSettings.NumIterations)
	}
	return averageResponseTime
}

func buildAndSendUserRequests(subsetOfResponseTimesChan chan perfTestUtils.RspTimes, loadPerUser int, testDefinition *TestDefinition, targetHost string, targetPort string, globalsMap map[string]map[string]string) {
	responseTimes := make(perfTestUtils.RspTimes, loadPerUser)
	loopExecutedToCompletion := true

	for i := 0; i < loadPerUser; i++ {
		responseTime := testDefinition.BuildAndSendRequest(targetHost, targetPort, "", globalsMap)

		if responseTime > 0 {
			responseTimes[i] = responseTime
		} else {
			loopExecutedToCompletion = false
			break
		}
	}

	if loopExecutedToCompletion {
		subsetOfResponseTimesChan <- responseTimes
	} else {
		subsetOfResponseTimesChan <- nil
	}
}

func aggregateResponseTimes(responseTimes *[]int64, subsetOfResponseTimesChan chan perfTestUtils.RspTimes, wg *sync.WaitGroup) {
	subsetOfResponseTimes := <-subsetOfResponseTimesChan
	if subsetOfResponseTimes != nil {
		*responseTimes = append(*responseTimes, subsetOfResponseTimes...)
	}
	wg.Done()
}
