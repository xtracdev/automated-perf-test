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



//----- Globals ------------------------------------------------------------------
var configurationSettings *perfTestUtils.Config
var osFileSystem = perfTestUtils.OsFS{}

// Command line arguments:
var configFilePath string
var checkTestReadiness bool
var verbose bool
var debug bool

const (
	TRAINING_MODE = 1
	TESTING_MODE  = 2
)




//----- main ------------------------------------------------------------------
func main() {
	log.Debugf("[START]")

	initConfig(os.Args[1:], osFileSystem, os.Exit)

	//Validate config()
	configurationSettings.PrintAndValidateConfig()

	//Generate a test suite based on configuration settings
	testSuite := new(testStrategies.TestSuite)
	testSuite.BuildTestSuite(configurationSettings)
	numTestCases := len(testSuite.TestCases) //convenience variable

	if checkTestReadiness {
		readyForTest, _ := perfTestUtils.IsReadyForTest(configurationSettings, testSuite.Name, numTestCases)
		if !readyForTest {
			log.Warn("System is not ready for testing.")
			os.Exit(1)
		} else {
			log.Warn("System is ready for testing.")
			os.Exit(0)
		}
	}

	//Determine testing mode.
	if configurationSettings.GBS || configurationSettings.ReBaseAll {
		if configurationSettings.ReBaseAll {
			runInTrainingMode(configurationSettings.ExecutionHost, true, testSuite)
		} else {
			readyForTest, _ := perfTestUtils.IsReadyForTest(configurationSettings, testSuite.Name, numTestCases)
			if !readyForTest {
				runInTrainingMode(configurationSettings.ExecutionHost, false, testSuite)
			} else {
				log.Warn("System is ready for testing. Training is not required.")
			}
		}
	} else {
		readyForTest, basePerfStats := perfTestUtils.IsReadyForTest(configurationSettings, testSuite.Name, numTestCases)
		if readyForTest {
			runInTestingMode(basePerfStats, configurationSettings.ExecutionHost, perfTestUtils.GenerateTemplateReport, testSuite)
		} else {
			log.Warn("System is not ready for testing. Attempting to run training mode....")
			runInTrainingMode(configurationSettings.ExecutionHost, false, testSuite)
			readyForTest, basePerfStats = perfTestUtils.IsReadyForTest(configurationSettings, testSuite.Name, numTestCases)
			if readyForTest {
				runInTestingMode(basePerfStats, configurationSettings.ExecutionHost, perfTestUtils.GenerateTemplateReport, testSuite)
			} else {
				log.Error("System is not ready for testing. Failed to run in training mode. Check service logs for more details.")
				os.Exit(1)
			}
		}
	}
}




//----- initConfig ------------------------------------------------------------
func initConfig(args []string, fs perfTestUtils.FileSystem, exit func(code int)) {
	//----- Initialize config data structure and set defaults.
	// Note: Defaults will be overridden as needed. The user can ignore
	// unnecessary parameters in config file and command prompt.
	configurationSettings = new(perfTestUtils.Config)
	configurationSettings.SetDefaults()

	//----- Get Hostname for this machine.
	host, err := os.Hostname()
	if err != nil {
		log.Error("Failed to resolve host name. Error:", err)
		exit(1)
	}
	configurationSettings.ExecutionHost = host

	//----- Process command line args.
	// Global controls outside of Config struct:
	flag.StringVar(&configFilePath, "configFilePath", "", "The location of the configuration file.")
	flag.BoolVar(&checkTestReadiness, "checkTestReadiness", false, "Simple check to see if system requires training.")

	// Log level simplified for the end user.
	flag.BoolVar(&verbose, "v", false, "Set logging verbosity to 'info' from default of 'warn'. Use -vv for debug.")
	flag.BoolVar(&debug, "vv", false, "Set verbosity to debug.")

	// Args that override default options in Config struct:
	flag.BoolVar(&configurationSettings.GBS, "gbs", false, "Generate 'Base Statistics' for this server")
	flag.BoolVar(&configurationSettings.ReBaseMemory, "reBaseMemory", false, "Generate new base peak memory for this server")
	flag.BoolVar(&configurationSettings.ReBaseAll, "reBaseAll", false, "Generate new base for memory and service resposne times for this server")
	flag.StringVar(&configurationSettings.ConfigFileFormat, "configFileFormat", "xml", "The format of the configuration file {xlm, toml}")
	flag.StringVar(&configurationSettings.TestFileFormat, "testFileFormat", "xml", "The format of the test definition file {xlm, toml}")
	flag.CommandLine.Parse(args)

	setLogLevel( verbose, debug )

	//----- Parse the config file.
	if configFilePath == "" {
		log.Warn("No config file found. - Using default values.")
		return
	}

	cf, err := fs.Open(configFilePath)
	if cf != nil {
		defer cf.Close()
	}
	if err != nil {
		log.Error("No config file found at path: ", configFilePath, " - Using default values.")
	} else {
		fileContent, fileErr := ioutil.ReadAll(cf)
		if fileErr != nil {
			log.Error("No readable config file found at path: ", configFilePath, " - Using default values.")
		} else {
			switch configurationSettings.ConfigFileFormat {
			case "toml":
				err := toml.Unmarshal(fileContent, &configurationSettings)
				if err != nil {
					log.Error("Failed to parse config file ", configFilePath, ". Error:", err, " - Using default values.")
				}
			default:
				xmlError := xml.Unmarshal(fileContent, &configurationSettings)
				if xmlError != nil {
					log.Error("Failed to parse config file ", configFilePath, ". Error:", xmlError, " - Using default values.")
				}
			}
		}
	}
}




//----- setLogLevel -----------------------------------------------------------
// Set log level using a simplified interface for end user. See "Process
// command line args" in initConfig().
func setLogLevel( verbose, debug bool ) {
	// Set default to WarnLevel
	log.SetLevel( log.WarnLevel )

	// Increase verbosity as set by user at command line.
	if verbose {
		log.SetLevel( log.InfoLevel )
	}
	if debug {
		log.SetLevel( log.DebugLevel )
	}
}




//----- runInTrainingMode -----------------------------------------------------
func runInTrainingMode(host string, reBaseAll bool, testSuite *testStrategies.TestSuite) {
	log.Info("Running performance test in Training mode for host ", host)

	//Start Test Timer
	executionStartTime := time.Now()

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
		f, _ := os.Open(configurationSettings.BaseStatsOutputDir + "/" + host + "-" + testSuite.Name + "-perfBaseStats")
		basePerfstats, _ = perfTestUtils.ReadBasePerfFile(f)
	}

	//initilize Performance statistics struct for this test run
	perfStatsForTest := &perfTestUtils.PerfStats{ServiceResponseTimes: make(map[string]int64), ServiceTps: make(map[string]float64)}

	//Run the test
	runTests(perfStatsForTest, TRAINING_MODE, testSuite)
	durExecRunTime := time.Since(executionStartTime)

	//Generate base statistics output file for this training run.
	perfTestUtils.GenerateEnvBasePerfOutputFile(perfStatsForTest, basePerfstats, configurationSettings, os.Exit, osFileSystem, testSuite.Name)

	log.Info("Training mode completed successfully. ")
	log.Info( "Execution Run Time [", perfTestUtils.GetExecutionTimeDisplay(durExecRunTime), "]" )
}




//----- runInTestingMode ------------------------------------------------------
func runInTestingMode(
		basePerfstats *perfTestUtils.BasePerfStats,
		host string,
		frg func(*perfTestUtils.BasePerfStats, *perfTestUtils.PerfStats, *perfTestUtils.Config, perfTestUtils.FileSystem, string),
		testSuite *testStrategies.TestSuite,
) {
	log.Info("Running Performance test in Testing mode for host ", host)

	//initilize Performance statistics struct for this test run
	perfStatsForTest := &perfTestUtils.PerfStats{ServiceResponseTimes: make(map[string]int64), TestDate: time.Now(), ServiceTps: make(map[string]float64)}

	//Start Test Timer
	executionStartTime := time.Now()
	//Run the test
	runTests(perfStatsForTest, TESTING_MODE, testSuite)
	durExecRunTime := time.Since(executionStartTime)

	//Validate test results
	assertionFailures := runAssertions(basePerfstats, perfStatsForTest)

	//Generate performance test report
	frg(basePerfstats, perfStatsForTest, configurationSettings, osFileSystem, testSuite.Name)

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

	totalOps := (configurationSettings.NumIterations * configurationSettings.ConcurrentUsers * len(testSuite.TestCases))
	log.Infof( "Total Ops: [%d]", totalOps )
	log.Infof( "Run Time:  [%s]", perfTestUtils.GetExecutionTimeDisplay(durExecRunTime) )
	log.Infof( "TPS:       [%f]", perfStatsForTest.OverAllTPS)
	log.Info("=====================================================")

	if len(assertionFailures) > 0 {
		os.Exit(1)
	}
}




//----- runTests --------------------------------------------------------------
//This function does two things,
//1 Start a go routine to periodically grab the memory foot print and set the peak memory value
//2 Run all test using mock servers and gather performance stats
func runTests(perfStatsForTest *perfTestUtils.PerfStats, mode int, testSuite *testStrategies.TestSuite) {
	// Start the overall timer.
	testStartTime := time.Now()
	var testExecutionTime time.Duration

	//Initialize Memory analysis
	var peakMemoryAllocation = new(uint64)
	memoryAudit := make([]uint64, 0)
	testPartitions := make([]perfTestUtils.TestPartition, 0)
	counter := 0
	testPartitions = append(testPartitions, perfTestUtils.TestPartition{Count: counter, TestName: "StartUp"})

	//Start go routine to grab memory in use
	//Peak memory is stored in peakMemoryAllocation variable.
	quit := make(chan bool)
	go func() {
		for {
			select {
			case <-quit:
				return
			default:

				memoryStatsUrl := "http://" + configurationSettings.TargetHost + ":" + configurationSettings.TargetPort + configurationSettings.MemoryEndpoint
				resp, err := http.Get(memoryStatsUrl)
				if err != nil {
					log.Error("Memory analysis unavailable. Failed to retrieve memory Statistics from endpoint ", memoryStatsUrl, ". Error: ", err)
					quit <- true
				} else {

					body, _ := ioutil.ReadAll(resp.Body)

					defer resp.Body.Close()

					m := new(perfTestUtils.Entry)
					unmarshalErr := json.Unmarshal(body, m)
					if unmarshalErr != nil {
						log.Error("Memory analysis unavailable. Failed to unmarshal memory statistics from endpoint: ", memoryStatsUrl, ". UnmarsahlErr: ", unmarshalErr)
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

	//Check the test strategy
	if testSuite.TestStrategy == testStrategies.SERVICE_BASED_TESTING {
		log.Info("Running Service Based Testing Strategy")

		//Determine load per concurrent user
		loadPerUser := int(configurationSettings.NumIterations / configurationSettings.ConcurrentUsers)
		remainder := configurationSettings.NumIterations % configurationSettings.ConcurrentUsers

		var index int
		var testDefinition *testStrategies.TestDefinition
		for index, testDefinition = range testSuite.TestCases {
			log.Infof("Running Test case [%d] [Name:%s]", index, testDefinition.TestName)
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
		// Stop overall test timer
		testExecutionTime = time.Since(testStartTime)

		// Calc overall TPS.
		perfStatsForTest.OverAllTPS = perfTestUtils.CalcTps(
			(configurationSettings.NumIterations * configurationSettings.ConcurrentUsers),
			testExecutionTime,
		)

		// Calc individual service TPS.
		for index, testDefinition = range testSuite.TestCases {
			perfStatsForTest.ServiceTps[testDefinition.TestName] =
				perfTestUtils.CalcTps(configurationSettings.NumIterations, testExecutionTime)
		}
	} else if testSuite.TestStrategy == testStrategies.SUITE_BASED_TESTING {
		log.Info("Running Suite Based Testing Strategy. Suite Name: [", testSuite.Name, "]")

		allServicesResponseTimesMap := testStrategies.ExecuteTestSuiteWrapper(testSuite, configurationSettings, testStartTime)

		// Stop overall test timer
		testExecutionTime = time.Since(testStartTime)

		// Calc overall TPS
		perfStatsForTest.OverAllTPS = perfTestUtils.CalcTps(
				(configurationSettings.NumIterations * configurationSettings.ConcurrentUsers * len(testSuite.TestCases)),
				testExecutionTime,
		)

		for serviceName, serviceResponseTimes := range allServicesResponseTimesMap {
			if len(serviceResponseTimes) == (configurationSettings.NumIterations * configurationSettings.ConcurrentUsers) {
				averageResponseTime := perfTestUtils.CalcAverageResponseTime(serviceResponseTimes, configurationSettings.NumIterations)
				if averageResponseTime > 0 {
					perfStatsForTest.ServiceResponseTimes[serviceName] = averageResponseTime
					perfStatsForTest.ServiceTps[serviceName] =
						perfTestUtils.CalcTps(len(serviceResponseTimes), testExecutionTime)
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




//----- runAssertions ---------------------------------------------------------
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
		varianceOk := perfTestUtils.ValidateAverageServiceResponseTimeVariance(configurationSettings.AllowableServiceResponseTimeVariance, responseTimeVariancePercentage, serviceName)
		if !varianceOk {
			assertionFailures = append(assertionFailures, fmt.Sprintf("Service Failure: Service test %-60s response time variance exceeded by %3.2f %1s", serviceName, responseTimeVariancePercentage, "%"))
		}
	}
	return assertionFailures
}

