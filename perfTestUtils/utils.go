package perfTestUtils

import (
	"encoding/json"
	"fmt"
	//log "github.com/Sirupsen/logrus"
	"io/ioutil"
	"sort"
)

//This function reads a base perf file for this host and converts it to a base perf struct
func ReadBasePerfFile(host string) (*BasePerfStats, error) {
	basePerfstats := &BasePerfStats{
		BaseServiceResponseTimes: make(map[string]int64),
		MemoryAudit:              make([]uint64, 0),
	}
	var errorFound error

	fileContent, fileErr := ioutil.ReadFile("./envStats/" + host + "-perfBaseStats")
	if fileErr != nil {
		errorFound = fileErr
	} else {
		jsonError := json.Unmarshal(fileContent, basePerfstats)
		if jsonError != nil {
			errorFound = jsonError
		}
	}
	return basePerfstats, errorFound
}

//=====================
//Calc Memory functions
//=====================
func CalcPeakMemoryVariancePercentage(basePeakMemory uint64, peakMemory uint64) float64 {

	peakMemoryVariancePercentage := float64(0)

	if basePeakMemory < peakMemory {
		peakMemoryDelta := peakMemory - basePeakMemory
		temp := float64(float64(peakMemoryDelta) / float64(basePeakMemory))
		peakMemoryVariancePercentage = temp * 100
	} else {
		peakMemoryDelta := basePeakMemory - peakMemory
		temp := float64(float64(peakMemoryDelta) / float64(basePeakMemory))
		peakMemoryVariancePercentage = (temp * 100) * -1
	}

	return peakMemoryVariancePercentage
}

//============================
//Calc Response time functions
//============================
func CalcAverageResponseTime(responseTimes RspTimes, numIterations int) int64 {

	averageResponseTime := int64(0)

	//Remove the highest 5% to take out anomolies
	sort.Sort(responseTimes)
	numberToRemove := int(float32(numIterations) * float32(0.05))
	responseTimes = responseTimes[0 : len(responseTimes)-numberToRemove]

	totalOfAllresponseTimes := int64(0)
	for _, val := range responseTimes {
		totalOfAllresponseTimes = totalOfAllresponseTimes + val
	}

	averageResponseTime = int64(float64(totalOfAllresponseTimes) / float64(numIterations))

	return averageResponseTime
}

func CalcAverageResponseVariancePercentage(averageResponseTime int64, baseResponseTime int64) float64 {

	responseTimeVariancePercentage := float64(0)

	if baseResponseTime < averageResponseTime {
		delta := averageResponseTime - baseResponseTime
		temp := float64(float64(delta) / float64(baseResponseTime))
		responseTimeVariancePercentage = temp * 100
	} else {
		delta := baseResponseTime - averageResponseTime
		temp := float64(float64(delta) / float64(baseResponseTime))
		responseTimeVariancePercentage = (temp * 100) * -1
	}

	return responseTimeVariancePercentage
}

//=====================================
//Service response validation functions
//=====================================
func ValidateResponseBody(body []byte, testName string) bool {

	isResponseBodyValid := false
	if len(body) > 0 {
		isResponseBodyValid = true
	} else {
		//log.Error(fmt.Sprintf("Incorrect Content lenght (%d) returned for service %s", len(body), testName))
		fmt.Println(fmt.Sprintf("Incorrect Content lenght (%d) returned for service %s", len(body), testName))
	}
	return isResponseBodyValid
}

func ValidateResponseStatusCode(responseStatusCode int, expectedStatusCode int, testName string) bool {

	isResponseStatusCodeValid := false
	if responseStatusCode == expectedStatusCode {
		isResponseStatusCodeValid = true
	} else {
		//log.Error(fmt.Sprintf("Incorrect status code of %d retruned for service %s. %d expected", responseStatusCode, testName, expectedStatusCode))
		fmt.Println(fmt.Sprintf("Incorrect status code of %d retruned for service %s. %d expected", responseStatusCode, testName, expectedStatusCode))
	}
	return isResponseStatusCodeValid
}

func ValidateServiceResponseTime(responseTime int64, testName string) bool {

	isResponseTimeValid := false
	if responseTime > 0 {
		isResponseTimeValid = true
	} else {
		//log.Error(fmt.Sprintf("Time taken to complete request %s was 0 nanoseconds", testName))
		fmt.Println(fmt.Sprintf("Time taken to complete request %s was 0 nanoseconds", testName))
	}
	return isResponseTimeValid
}

//=====================================
//Test Assertion functions
//=====================================
func ValidatePeakMemoryVariance(allowablePeakMemoryVariance float64, peakMemoryVariancePercentage float64) bool {

	isPeakMemoryVarianceValid := false
	if allowablePeakMemoryVariance < peakMemoryVariancePercentage {
		isPeakMemoryVarianceValid = true
	} else {
		//log.Error(fmt.Sprintf("Peak Memory Variance value of %f%% exceeded. Max allowable variance is %f%%", peakMemoryVariancePercentage, allowablePeakMemoryVariance))
		fmt.Println(fmt.Sprintf("Peak Memory Variance value of %f%% exceeded. Max allowable variance is %f%%", peakMemoryVariancePercentage, allowablePeakMemoryVariance))
	}
	return isPeakMemoryVarianceValid
}

func ValidateTestCaseCount(baseTestCaseCount int, testTestCaseCount int) bool {

	isTestCaseCountValid := false
	if baseTestCaseCount == testTestCaseCount {
		isTestCaseCountValid = true
	} else {
		//log.Error(fmt.Sprintf("Number of service tests in base is differnet to the number of services for this test run."))
		fmt.Println(fmt.Sprintf("Number of service tests in base is differnet to the number of services for this test run."))
	}
	return isTestCaseCountValid
}

func ValidateAverageServiceResponeTimeVariance(allowableServiceResponseTimeVariance float64, serviceResponseTimeVariancePercentage float64, serviceName string) bool {

	isAverageServiceResponeTimeVarianceValid := false
	if allowableServiceResponseTimeVariance < serviceResponseTimeVariancePercentage {
		isAverageServiceResponeTimeVarianceValid = true
	} else {
		//log.Error(fmt.Sprintf("%s Response Time Variance value of %f%% exceeded. Max allowable variance is %f%%", serviceName, serviceResponseTimeVariancePercentage, allowableServiceResponseTimeVariance))
		fmt.Println(fmt.Sprintf("%s Response Time Variance value of %f%% exceeded. Max allowable variance is %f%%", serviceName, serviceResponseTimeVariancePercentage, allowableServiceResponseTimeVariance))
	}
	return isAverageServiceResponeTimeVarianceValid
}

func (a RspTimes) Len() int           { return len(a) }
func (a RspTimes) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a RspTimes) Less(i, j int) bool { return a[i] < a[j] }
