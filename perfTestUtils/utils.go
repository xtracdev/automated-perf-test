package perfTestUtils

import (
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"io"
	"io/ioutil"
	"os"
	"sort"
	"time"
)

//FileSystem is an interface to access os filesystem or mock it
type FileSystem interface {
	Open(name string) (File, error)
	Create(name string) (File, error)
}

//File is an interface to access os.File or mock it
type File interface {
	Readdir(n int) (fi []os.FileInfo, err error)
	io.WriteCloser
	Read(p []byte) (n int, err error)
}

// OsFS implements fileSystem using the local disk.
type OsFS struct{}

//Open calls os function
func (OsFS) Open(name string) (File, error) { return os.Open(name) }

//Create calls os function
func (OsFS) Create(name string) (File, error) { return os.Create(name) }

//=============================
//Test run utility functions
//=============================
//This function reads a base perf and converts it to a base perf struct
func ReadBasePerfFile(r io.Reader) (*BasePerfStats, error) {
	basePerfstats := &BasePerfStats{
		BaseServiceResponseTimes: make(map[string]int64),
		MemoryAudit:              make([]uint64, 0),
	}
	var errorFound error

	content, err := ioutil.ReadAll(r)
	if err != nil {
		errorFound = err
	} else {
		jsonError := json.Unmarshal(content, basePerfstats)
		if jsonError != nil {
			errorFound = jsonError
		}
	}
	return basePerfstats, errorFound
}

func GetExecutionTimeDisplay(durExecTime time.Duration) string {
	return fmt.Sprintf( "%s (%.9fs)", durExecTime.String(), durExecTime.Seconds() )
}

func IsReadyForTest(configurationSettings *Config, testSuiteName string, numTestCases int) (bool, *BasePerfStats) {
	//1) read in perf base stats
	f, err := os.Open(configurationSettings.BaseStatsOutputDir + "/" + configurationSettings.ExecutionHost + "-" + testSuiteName + "-perfBaseStats")
	if err != nil {
		log.Errorf("Failed to open env stats for %v. Error: %v.", configurationSettings.ExecutionHost, err)
		return false, nil
	}
	basePerfstats, err := ReadBasePerfFile(f)
	if err != nil {
		log.Error("Failed to read env stats for " + configurationSettings.ExecutionHost + ". Error:" + err.Error() + ".")
		return false, nil
	}

	//2) validate content  of base stats file
	isBasePerfStatsValid := validateBasePerfStat(basePerfstats)
	if !isBasePerfStatsValid {
		log.Error("Base Perf stats are not fully populated for  " + configurationSettings.ExecutionHost + ".")
		return false, nil
	}

	//3) Verify the number of base test cases is equal to the number of service test cases.
	baselineAmount := len(basePerfstats.BaseServiceResponseTimes)
	log.Info("Number of defined test cases:", numTestCases)
	log.Info("Number of base line test cases:", baselineAmount)

	if baselineAmount != numTestCases {
		log.Errorf(
			"The number of test definitions [%d] does not equal the number of baseline metrics [%d].",
			numTestCases,
			baselineAmount,
		)
		return false, nil
	}

	return true, basePerfstats
}

func validateBasePerfStat(basePerfstats *BasePerfStats) bool {
	isBasePerfStatsValid := true

	if basePerfstats.BasePeakMemory <= 0 {
		isBasePerfStatsValid = false
	}
	if basePerfstats.GenerationDate == "" {
		isBasePerfStatsValid = false
	}
	if basePerfstats.ModifiedDate == "" {
		isBasePerfStatsValid = false
	}
	if len(basePerfstats.MemoryAudit) <= 0 {
		isBasePerfStatsValid = false
	}
	if basePerfstats.BaseServiceResponseTimes != nil {
		for _, baseResponseTime := range basePerfstats.BaseServiceResponseTimes {
			if baseResponseTime <= 0 {
				isBasePerfStatsValid = false
				break
			}
		}
	} else {
		isBasePerfStatsValid = false
	}
	return isBasePerfStatsValid
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

func CalcTpsForService(averageRespTimeForServiceinNanoSeconds int64) float64 {
	timeInMilliSeconds := float64(float64(averageRespTimeForServiceinNanoSeconds) / float64(1000000))
	return float64(float64(1000) / float64(timeInMilliSeconds))
}

func CalcTpsOverAllBasedOnAverageserviceTPS(ServiceTps map[string]float64) float64 {
	totalServiceTPS := float64(0)
	for _, v := range ServiceTps {
		totalServiceTPS = totalServiceTPS + v
	}
	averageServiceTPSTPS := float64(float64(totalServiceTPS) / float64(len(ServiceTps)))

	return averageServiceTPSTPS
}

func CalcTpsOverAllBasedOnAverageServiceResponseTimes(serviceResponseTimes map[string]int64) float64 {
	totalServiceResponseTimes := int64(0)
	for _, v := range serviceResponseTimes {
		totalServiceResponseTimes = totalServiceResponseTimes + v
	}
	averageServiceResponseTimesTPS := float64(float64(totalServiceResponseTimes) / float64(len(serviceResponseTimes)))

	averageinMilieseonds := float64(float64(averageServiceResponseTimesTPS) / float64(1000000))
	return float64(float64(1000) / averageinMilieseonds)
}

//============================
//Calc Response time functions
//============================
func CalcAverageResponseTime(responseTimes RspTimes, numIterations int, testMode int) int64 {

	averageResponseTime := int64(0)
	numberToRemove := 0
	sort.Sort(responseTimes)

	if testMode == 2 {
		//It in testing mode, remove the highest 10% to take out anomolies and outliers
		numberToRemove = int(float32(numIterations) * float32(0.1))
		responseTimes = responseTimes[0 : len(responseTimes)-numberToRemove]
	}
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
		delta := uint64(averageResponseTime) - uint64(baseResponseTime)
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
func ValidateResponseStatusCode(responseStatusCode int, expectedStatusCode int, testName string) bool {

	isResponseStatusCodeValid := false
	if responseStatusCode == expectedStatusCode {
		isResponseStatusCodeValid = true
	} else {
		log.Errorf("Incorrect status code of %d returned for service %s. %d expected", responseStatusCode, testName, expectedStatusCode)
	}
	return isResponseStatusCodeValid
}

func ValidateServiceResponseTime(responseTime int64, testName string) bool {

	isResponseTimeValid := false
	if responseTime > 0 {
		isResponseTimeValid = true
	} else {
		log.Error(fmt.Sprintf("Time taken to complete request %s was 0 nanoseconds", testName))
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

func ValidateAverageServiceResponseTimeVariance(allowableServiceResponseTimeVariance float64, serviceResponseTimeVariancePercentage float64, serviceName string) bool {
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
func populateBasePerfStats(perfStatsForTest *PerfStats, basePerfstats *BasePerfStats, reBaseMemory bool) {
	modified := false

	//Setting memory data
	if basePerfstats.BasePeakMemory == 0 || reBaseMemory {
		basePerfstats.BasePeakMemory = perfStatsForTest.PeakMemory
		modified = true
	}
	if basePerfstats.MemoryAudit == nil || len(basePerfstats.MemoryAudit) == 0 || reBaseMemory {
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

func GenerateEnvBasePerfOutputFile(perfStatsForTest *PerfStats, basePerfstats *BasePerfStats, configurationSettings *Config, exit func(code int), fs FileSystem, testSuiteName string) {

	//Set base performance based on training test run
	populateBasePerfStats(perfStatsForTest, basePerfstats, configurationSettings.ReBaseMemory)

	//Convert base perf stat to Json
	basePerfstatsJson, err := json.Marshal(basePerfstats)
	if err != nil {
		log.Error("Failed to marshal to Json. Error:", err)
		exit(1)
	}

	// Check for existence of output dir and create if needed.
	if os.MkdirAll(configurationSettings.BaseStatsOutputDir, os.ModePerm); err != nil {
		log.Errorf("Failed to create path: [%s]. Error: %s\n", configurationSettings.BaseStatsOutputDir, err)
		exit(1)
	}

	// Write base perf stat to file.
	file_name := configurationSettings.ExecutionHost + "-" + testSuiteName + "-perfBaseStats"
	file, err := fs.Create(configurationSettings.BaseStatsOutputDir + "/" + file_name)
	if err != nil {
		log.Error("Failed to create output file. Error:", err)
		exit(1)
	}
	if file != nil {
		defer file.Close()
		file.Write(basePerfstatsJson)
	}
}
