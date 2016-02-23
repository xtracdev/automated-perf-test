package main

import (
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	"github.com/BurntSushi/toml"
	log "github.com/Sirupsen/logrus"
	"github.com/xtracdev/automated-perf-test/perfTestUtils"
	"github.com/xtracdev/automated-perf-test/testStrategies"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

var configurationSettings *perfTestUtils.Config
var checkTestReadyness bool

const (
	TRAINING_MODE = 1
	TESTING_MODE  = 2
)

var osFileSystem = perfTestUtils.OsFS{}

func initConfig(args []string, fs perfTestUtils.FileSystem, exit func(code int)) {
	//Command line ags
	var gbs bool
	var reBaseMemory bool
	var reBaseAll bool
	var configFilePath string
	var configFileFormat string
	var testFileFormat string

	//Process command line arugments.
	flag.BoolVar(&gbs, "gbs", false, "Genertate Base Performance Staticists for this server")
	flag.BoolVar(&reBaseMemory, "reBaseMemory", false, "Generate new base peak memory for this server")
	flag.BoolVar(&reBaseAll, "reBaseAll", false, "Generate new base for memory and service resposne times for this server")
	flag.BoolVar(&checkTestReadyness, "checkTestReadyness", false, "Simple check to see if system requires training.")
	flag.StringVar(&configFilePath, "configFilePath", "", "The location of the configuration file.")
	flag.StringVar(&configFileFormat, "configFileFormat", "xml", "The format of the configuration file {xlm, toml}")
	flag.StringVar(&testFileFormat, "testFileFormat", "xml", "The format of the test definition file {xlm, toml}")
	flag.CommandLine.Parse(args)

	//Read and paser config file if present.
	configurationSettings = new(perfTestUtils.Config)
	if configFilePath != "" {
		cf, err := fs.Open(configFilePath)
		if cf != nil {
			defer cf.Close()
		}
		if err != nil {
			log.Error("No config file found at path: ", configFilePath, " - Using default values.")
			configurationSettings.SetDefaults()
		} else {
			fileContent, fileErr := ioutil.ReadAll(cf)
			if fileErr != nil {
				log.Error("No readable config file found at path: ", configFilePath, " - Using default values.")
				configurationSettings.SetDefaults()
			} else {
				switch configFileFormat {
				case "toml":
					err := toml.Unmarshal(fileContent, &configurationSettings)
					if err != nil {
						fmt.Println("Failed to parse config file ", configFilePath, ". Error:", err)
						configurationSettings.SetDefaults()
					}
				default:
					xmlError := xml.Unmarshal(fileContent, &configurationSettings)
					if xmlError != nil {
						log.Error("Failed to parse config file ", configFilePath, ". Error:", xmlError, " - Using default values.")
						configurationSettings.SetDefaults()
					}
				}
			}
		}
	}

	//Get Hostname for this machine.
	host, err := os.Hostname()
	if err != nil {
		log.Error("Failed to resolve host name. Error:", err)
		exit(1)
	}
	configurationSettings.ExecutionHost = host
	configurationSettings.GBS = gbs
	configurationSettings.ReBaseMemory = reBaseMemory
	configurationSettings.ReBaseAll = reBaseAll

	configurationSettings.ConfigFileFormat = configFileFormat
	configurationSettings.TestFileFormat = testFileFormat
}

//Main Test Method
func main() {
	log.Debugf("[START]")
	initConfig(os.Args[1:], osFileSystem, os.Exit)

	if checkTestReadyness {
		readyForTest, _ := perfTestUtils.IsReadyForTest(configurationSettings, osFileSystem)
		if !readyForTest {
			log.Info("System is not ready for testing.")
			os.Exit(1)
		} else {
			log.Info("System is ready for testing.")
			os.Exit(0)
		}
	}

	//Validate config()
	configurationSettings.PrintAndValidateConfig(os.Exit)

	//Determine testing mode.
	if configurationSettings.GBS || configurationSettings.ReBaseAll {
		if configurationSettings.ReBaseAll {
			runInTrainingMode(configurationSettings.ExecutionHost, true)
		} else {
			readyForTest, _ := perfTestUtils.IsReadyForTest(configurationSettings, osFileSystem)
			if !readyForTest {
				runInTrainingMode(configurationSettings.ExecutionHost, false)
			} else {
				log.Info("System is ready for testing. Training is not required.")
			}
		}
	} else {
		readyForTest, basePerfStats := perfTestUtils.IsReadyForTest(configurationSettings, osFileSystem)
		if readyForTest {
			runInTestingMode(basePerfStats, configurationSettings.ExecutionHost, perfTestUtils.GenerateTemplateReport)
		} else {
			log.Info("System is not ready for testing. Attempting to run training mode....")
			runInTrainingMode(configurationSettings.ExecutionHost, false)
			readyForTest, basePerfStats = perfTestUtils.IsReadyForTest(configurationSettings, osFileSystem)
			if readyForTest {
				runInTestingMode(basePerfStats, configurationSettings.ExecutionHost, perfTestUtils.GenerateTemplateReport)
			} else {
				log.Info("System is not ready for testing. Attempting to run training failed. Check logs for more details.")
				os.Exit(1)
			}
		}
	}
}

func runInTrainingMode(host string, reBaseAll bool) {
	log.Info("Running performance test in Training mode for host ", host)
	testStratTime := time.Now().UnixNano()

	var basePerfstats *perfTestUtils.BasePerfStats
	if reBaseAll {
		log.Info("Performing full rebase of performance statistics for host ", host)
		basePerfstats = &perfTestUtils.BasePerfStats{
			BaseServiceResponseTimes: make(map[string]int64),
			MemoryAudit:              make([]uint64, 0),
		}
	} else {
		//Check to see if this server already has a base perf file defined.
		//If so, only values not previously populated will be set.
		//if not, a default base perf struct is created with nil values for all fields
		f, _ := os.Open(configurationSettings.BaseStatsOutputDir + "/" + host + "-perfBaseStats")
		basePerfstats, _ = perfTestUtils.ReadBasePerfFile(f)
	}

	//initilize Performance statistics struct for this test run
	perfStatsForTest := &perfTestUtils.PerfStats{ServiceResponseTimes: make(map[string]int64)}

	//Run the test
	runTests(perfStatsForTest, TRAINING_MODE)

	//Generate base statistics output file for this training run.
	perfTestUtils.GenerateEnvBasePerfOutputFile(perfStatsForTest, basePerfstats, configurationSettings, os.Exit, osFileSystem)

	testRunTime := time.Now().UnixNano() - testStratTime
	log.Info("Training mode completed successfully. ")
	log.Info("Execution Run Time :", perfTestUtils.GetExecutionTimeDisplay(testRunTime))
}

func runInTestingMode(basePerfstats *perfTestUtils.BasePerfStats, host string, frg func(*perfTestUtils.BasePerfStats, *perfTestUtils.PerfStats, *perfTestUtils.Config, perfTestUtils.FileSystem)) {
	log.Info("Running Performance test in Testing mode for host ", host)
	testStratTime := time.Now().UnixNano()

	//initilize Performance statistics struct for this test run
	perfStatsForTest := &perfTestUtils.PerfStats{ServiceResponseTimes: make(map[string]int64), TestDate: time.Now()}

	//Run the test
	runTests(perfStatsForTest, TESTING_MODE)

	//Validate test results
	assertionFailures := runAssertions(basePerfstats, perfStatsForTest)

	//Generate performance test report
	frg(basePerfstats, perfStatsForTest, configurationSettings, osFileSystem)

	//Print test results to std out
	log.Info("=================== TEST RESULTS ===================")
	if len(assertionFailures) > 0 {
		log.Info("Number of Failures : ", len(assertionFailures))
		for _, failure := range assertionFailures {
			log.Info(failure)
		}
	} else {
		log.Info("Testing mode completed successfully")
	}

	testRunTime := time.Now().UnixNano() - testStratTime
	log.Info("Execution Run Time :", perfTestUtils.GetExecutionTimeDisplay(testRunTime))
	log.Info("=====================================================")

	if len(assertionFailures) > 0 {
		os.Exit(1)
	}
}

//This function does two thing,
//1 Start a go routine to preiodically grab the memory foot print and set the peak memory value
//2 Run all test using mock servers and gather performance stats
func runTests(perfStatsForTest *perfTestUtils.PerfStats, mode int) {

	//Initilize Memory analysis
	var peakMemoryAllocation = new(uint64)
	memoryAudit := make([]uint64, 0)
	testPartitions := make([]perfTestUtils.TestPartition, 0)
	counter := 0
	testPartitions = append(testPartitions, perfTestUtils.TestPartition{Count: counter, TestName: "StartUp"})

	//Start go routine to grab memory in use
	//Peak memory is stored in peakMemoryAlocation variable.
	quit := make(chan bool)
	go func() {
		for {
			select {
			case <-quit:
				return
			default:

				memoryStatsUrl := "http://" + configurationSettings.TargetHost + ":" + configurationSettings.TargetPort + "/debug/vars"
				resp, err := http.Get(memoryStatsUrl)
				if err != nil {
					log.Error("Memory analysis unavailable. Failed to retrieve memory Statistics from endpoint ", memoryStatsUrl, ". Error:", err)
					quit <- true
				} else {

					body, _ := ioutil.ReadAll(resp.Body)

					defer resp.Body.Close()

					m := new(perfTestUtils.Entry)
					unmarshalErr := json.Unmarshal(body, m)
					if unmarshalErr != nil {
						log.Error("Memory analysis unavailable. Failed to unmarshal memory statistics. ", unmarshalErr)
						quit <- true
					} else {
						if m.Memstats.Alloc > *peakMemoryAllocation {
							*peakMemoryAllocation = m.Memstats.Alloc
						}
						memoryAudit = append(memoryAudit, m.Memstats.Alloc)
						counter++
						time.Sleep(time.Millisecond * 200)
					}
				}
			}
		}
	}()

	//Add a 1 second delay before running test case to allow the graph get some initial memory data before test cases are executed.
	time.Sleep(time.Second * 1)

	//Generate a test suite based on configuration settings
	testSuite := new(testStrategies.TestSuite)
	testSuite.BuildTestSuite(configurationSettings)

	//Check the test strategy
	if testSuite.TestStrategy == testStrategies.SERVICE_BASED_TESTING {

		log.Info("Running Service Based Testing Strategy")

		//Determine load per concurrent user
		loadPerUser := int(configurationSettings.NumIterations / configurationSettings.ConcurrentUsers)
		remainder := configurationSettings.NumIterations % configurationSettings.ConcurrentUsers

		for index, testDefinition := range testSuite.TestCases {
			log.Info("Running Test case ", index, " [Name:", testDefinition.TestName, "]")
			testPartitions = append(testPartitions, perfTestUtils.TestPartition{Count: counter, TestName: testDefinition.TestName})
			averageResponseTime := testStrategies.ExecuteServiceTest(testDefinition, loadPerUser, remainder, configurationSettings)
			if averageResponseTime > 0 {
				perfStatsForTest.ServiceResponseTimes[testDefinition.TestName] = averageResponseTime
			} else {
				if mode == TRAINING_MODE {
					//Fail fast on training mode if any requests fail. If training fails we cannot guarantee the results.
					log.Error("Training mode failed due to invalid response on service [Name:", testDefinition.TestName, "]")
					os.Exit(1)
				}
			}
		}
	} else if testSuite.TestStrategy == testStrategies.SUITE_BASED_TESTING {

		log.Info("Running Suite Based Testing Strategy. Suite:", testSuite.Name)
		allServicesResponseTimesMap := testStrategies.ExecuteTestSuiteWrapper(testSuite, configurationSettings)

		for serviceName, serviceResponseTimes := range allServicesResponseTimesMap {
			if len(serviceResponseTimes) == (configurationSettings.NumIterations * configurationSettings.ConcurrentUsers) {
				averageResponseTime := perfTestUtils.CalcAverageResponseTime(serviceResponseTimes, configurationSettings.NumIterations)
				if averageResponseTime > 0 {
					perfStatsForTest.ServiceResponseTimes[serviceName] = averageResponseTime
				} else {
					if mode == TRAINING_MODE {
						//Fail fast on training mode if any requests fail. If training fails we cannot guarantee the results.
						log.Error("Training mode failed due to invalid response on service [Name:", serviceName, "]")
						os.Exit(1)
					}
				}
			}
		}
	}

	time.Sleep(time.Second * 1)
	perfStatsForTest.PeakMemory = *peakMemoryAllocation
	perfStatsForTest.MemoryAudit = memoryAudit
	perfStatsForTest.TestPartitions = testPartitions

}

//This function runs the assertions to ensure memory and service have not deviated past the allowed variance
func runAssertions(basePerfstats *perfTestUtils.BasePerfStats, perfStats *perfTestUtils.PerfStats) []string {

	assertionFailures := make([]string, 0)

	//Asserts Peak memory growth has not exceeded the allowable variance
	peakMemoryVariancePercentage := perfTestUtils.CalcPeakMemoryVariancePercentage(basePerfstats.BasePeakMemory, perfStats.PeakMemory)
	varianceOk := perfTestUtils.ValidatePeakMemoryVariance(configurationSettings.AllowablePeakMemoryVariance, peakMemoryVariancePercentage)
	if !varianceOk {
		assertionFailures = append(assertionFailures, fmt.Sprintf("Memory Failure: Peak variance exceeded by %3.2f %1s", peakMemoryVariancePercentage, "%"))
	}

	//Asserts service response times have not exceeded the allowable variance
	for serviceName, baseResponseTime := range basePerfstats.BaseServiceResponseTimes {
		averageServiceResponseTime := perfStats.ServiceResponseTimes[serviceName]
		if averageServiceResponseTime == 0 {
			assertionFailures = append(assertionFailures, fmt.Sprintf("Service Failure: Service test %-60s did not execute correctly. See logs for more details.", serviceName))
		}

		responseTimeVariancePercentage := perfTestUtils.CalcAverageResponseVariancePercentage(averageServiceResponseTime, baseResponseTime)
		varianceOk := perfTestUtils.ValidateAverageServiceResponeTimeVariance(configurationSettings.AllowableServiceResponseTimeVariance, responseTimeVariancePercentage, serviceName)
		if !varianceOk {
			assertionFailures = append(assertionFailures, fmt.Sprintf("Service Failure: Service test %-60s response time variance exceeded by %3.2f %1s", serviceName, responseTimeVariancePercentage, "%"))
		}
	}
	return assertionFailures
}
