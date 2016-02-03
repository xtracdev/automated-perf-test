package perfTestUtils

import (
	"encoding/json"
	"fmt"
	//log "github.com/Sirupsen/logrus"
	"io/ioutil"
	"os"
	"sort"
	"time"
)

//=============================
//Testing run utility functions
//=============================
//This function reads a base perf file for this host and converts it to a base perf struct
func ReadBasePerfFile(host string, baseStatsOutputDir string) (*BasePerfStats, error) {
	basePerfstats := &BasePerfStats{
		BaseServiceResponseTimes: make(map[string]int64),
		MemoryAudit:              make([]uint64, 0),
	}
	var errorFound error

	fileContent, fileErr := ioutil.ReadFile(baseStatsOutputDir + "/" + host + "-perfBaseStats")
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

func ValidateTestDefinitionAmount(baselineAmount int, configurationSettings *Config) bool {

	d, err := os.Open(configurationSettings.TestDefinationsDir)
	if err != nil {
		//log.Error("Failed to open test definations directory. Error:", err)
		fmt.Println("Failed to open test definations directory. Error:", err)
		os.Exit(1)
	}
	defer d.Close()
	fi, err := d.Readdir(-1)
	if err != nil {
		//log.Error("Failed to read files in test definations directory. Error:", err)
		fmt.Println("Failed to read files in test definations directory. Error:", err)
		os.Exit(1)
	}

	definitionAmount := len(fi)

	fmt.Println("Number of defined test cases:", definitionAmount)
	fmt.Println("Number of base line test cases:", baselineAmount)
	if definitionAmount != baselineAmount {
		//log.Errorf("Amount of test definition: %d does not equal to baseline amount: %d.", definitionAmount, baselineAmount)
		fmt.Println(fmt.Sprintf("Amount of test definition: %d does not equal to baseline amount: %d.", definitionAmount, baselineAmount))
		return false
	}
	return true
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

	averageResponseTime = int64(float64(totalOfAllresponseTimes) / float64(numIterations-numberToRemove))

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

	if allowablePeakMemoryVariance >= peakMemoryVariancePercentage {
		return true
	} else {
		return false
	}
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

	if allowableServiceResponseTimeVariance >= serviceResponseTimeVariancePercentage {
		return true
	} else {
		return false
	}
}

//=====================================
//Response times sort functions
//=====================================
func (a RspTimes) Len() int           { return len(a) }
func (a RspTimes) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a RspTimes) Less(i, j int) bool { return a[i] < a[j] }

//==============================================
//Generate base environment stats file functions
//==============================================
func populateBasePerfStats(perfStatsForTest *PerfStats, basePerfstats *BasePerfStats, resetPeakMemory bool) {
	modified := false

	//Setting memory data
	if basePerfstats.BasePeakMemory == 0 || resetPeakMemory {
		basePerfstats.BasePeakMemory = perfStatsForTest.PeakMemory
		modified = true
	}
	if basePerfstats.MemoryAudit == nil || len(basePerfstats.MemoryAudit) == 0 || resetPeakMemory {
		basePerfstats.MemoryAudit = perfStatsForTest.MemoryAudit
		modified = true
	}

	//Setting service response time data
	for serviceName, responseTime := range perfStatsForTest.ServiceResponseTimes {
		serviceBaseResponseTime := basePerfstats.BaseServiceResponseTimes[serviceName]
		if serviceBaseResponseTime == 0 {
			basePerfstats.BaseServiceResponseTimes[serviceName] = responseTime
			modified = true
		}
	}

	//Setting time stamps
	currentTime := time.Now().Format(time.RFC850)
	if basePerfstats.GenerationDate == "" {
		basePerfstats.GenerationDate = currentTime
	}
	if modified {
		basePerfstats.ModifiedDate = currentTime
	}
}

func GenerateEnvBasePerfOutputFile(perfStatsForTest *PerfStats, basePerfstats *BasePerfStats, configurationSettings *Config) {

	//Set base performance based on training test run
	populateBasePerfStats(perfStatsForTest, basePerfstats, configurationSettings.ResetPeakMemory)

	//Convert base perf stat to Json and write out to file
	basePerfstatsJson, err := json.Marshal(basePerfstats)
	if err != nil {
		//log.Error("Failed to marshal to Json. Error:", err)
		fmt.Println("Failed to marshal to Json. Error:", err)
		os.Exit(1)
	}
	file, err := os.Create(configurationSettings.BaseStatsOutputDir + "/" + configurationSettings.ExecutionHost + "-perfBaseStats")
	if err != nil {
		//log.Error("Failed to create output file. Error:", err)
		fmt.Println("Failed to create output file. Error:", err)
		os.Exit(1)
	}
	defer file.Close()
	file.Write(basePerfstatsJson)
}
