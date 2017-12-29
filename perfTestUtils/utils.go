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

// FileSystem is an interface to access os filesystem or mock it
type FileSystem interface {
	Open(name string) (File, error)
	Create(name string) (File, error)
}

// File is an interface to access os.File or mock it
type File interface {
	Readdir(n int) (fi []os.FileInfo, err error)
	io.WriteCloser
	Read(p []byte) (n int, err error)
}

// OsFS implements fileSystem using the local disk.
type OsFS struct{}

// Open calls os function
func (OsFS) Open(name string) (File, error) { return os.Open(name) }

// Create calls os function
func (OsFS) Create(name string) (File, error) { return os.Create(name) }

//=============================
//Test run utility functions
//=============================

// ReadBasePerfFile reads a base perf and converts it to a base perf struct
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

// IsReadyForTest validates the basePerfStats file content, and verifies the
// number of base test cases equals the number of configured test cases.
func IsReadyForTest(configurationSettings *Config, numTestCases int) (bool, *BasePerfStats) {
	//1) read in perf base stats
	f, err := os.Open(configurationSettings.BaseStatsOutputDir + "/" + configurationSettings.ExecutionHost + "-" + configurationSettings.APIName + "-perfBaseStats")
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
	isBasePerfStatsValid := validateBasePerfStat(basePerfstats, configurationSettings)
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

func validateBasePerfStat(basePerfstats *BasePerfStats, configurationSettings *Config) bool {
	isBasePerfStatsValid := true

	if ! configurationSettings.SkipMemCheck {
		if basePerfstats.BasePeakMemory <= 0 {
			isBasePerfStatsValid = false
		}
		if len(basePerfstats.MemoryAudit) <= 0 {
			isBasePerfStatsValid = false
		}
	}
	if basePerfstats.GenerationDate == "" {
		isBasePerfStatsValid = false
	}
	if basePerfstats.ModifiedDate == "" {
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

// CalcPeakMemoryVariancePercentage calculates the variance percentage between
// the base peak memory and the recorded peak memory.
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

// CalcTps returns the transaction per second given a time duratiuon and
// number of iterations.
func CalcTps(numIterations uint64, testRunTime time.Duration) float64 {
	return float64(float64(numIterations) / testRunTime.Seconds())
}

//============================
//Calc Response time functions
//============================

// CalcAverageResponseTime returns the average response time of a set of
// recorded response time values.
func CalcAverageResponseTime(responseTimes RspTimes, testMode int) int64 {
	averageResponseTime := int64(0)

	// Remove the highest =10% outliers.
	numberToRemove := 0

	sort.Sort(responseTimes)

	if testMode == 2 {
		// If in testing mode, remove the highest 10% outliers.
		numberToRemove = int(float32(len(responseTimes)) * float32(0.1))
		responseTimes = responseTimes[0 : len(responseTimes)-numberToRemove]
	}

	totalOfAllresponseTimes := int64(0)
	for _, val := range responseTimes {
		totalOfAllresponseTimes = totalOfAllresponseTimes + val
	}
	averageResponseTime = int64(float64(totalOfAllresponseTimes) / float64(len(responseTimes)))

	return averageResponseTime
}

// CalcAverageResponseVariancePercentage returns the variance percentage between
// the base peak response average and the response average of a set of
// recorded resonse time values.
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

// ValidateResponseStatusCode returns true if the HTTP response code matches
// the response code contained in the test case definition. Otherwise, false.
func ValidateResponseStatusCode(responseStatusCode int, expectedStatusCode int, testName string) bool {
	isResponseStatusCodeValid := false
	if responseStatusCode == expectedStatusCode {
		isResponseStatusCodeValid = true
	} else {
		log.Errorf("Incorrect status code of %d returned for service %s. %d expected", responseStatusCode, testName, expectedStatusCode)
	}
	return isResponseStatusCodeValid
}

// ValidateServiceResponseTime returns true if the given response time is
// greater than zero. Otherwise, log an error and return false.
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

// ValidatePeakMemoryVariance returns true if the given percentage is less
// than or equal to the allowable variance contained in the config.
func ValidatePeakMemoryVariance(allowablePeakMemoryVariance float64, peakMemoryVariancePercentage float64) bool {
	if allowablePeakMemoryVariance >= peakMemoryVariancePercentage {
		return true
	}
	return false
}

// ValidateAverageServiceResponseTimeVariance returns true if the given
// percentage is less than or equal to the allowable variance contained in
// the config.
func ValidateAverageServiceResponseTimeVariance(allowableServiceResponseTimeVariance float64, serviceResponseTimeVariancePercentage float64) bool {
	if allowableServiceResponseTimeVariance >= serviceResponseTimeVariancePercentage {
		return true
	}
	return false
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

// GenerateEnvBasePerfOutputFile writes the basePerfStats file.
func GenerateEnvBasePerfOutputFile(perfStatsForTest *PerfStats, basePerfstats *BasePerfStats, configurationSettings *Config, exit func(code int), fs FileSystem) {
	//Set base performance based on training test run
	populateBasePerfStats(perfStatsForTest, basePerfstats, configurationSettings.ReBaseMemory)

	//Convert base perf stat to Json
	basePerfstatsJSON, err := json.Marshal(basePerfstats)
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
	fileName := configurationSettings.ExecutionHost + "-" + configurationSettings.APIName + "-perfBaseStats"
	file, err := fs.Create(configurationSettings.BaseStatsOutputDir + "/" + fileName)
	if err != nil {
		log.Error("Failed to create output file. Error:", err)
		exit(1)
	}
	if file != nil {
		defer file.Close()
		file.Write(basePerfstatsJSON)
	}
}
