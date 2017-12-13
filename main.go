package main

import (
	"github.com/xtracdev/automated-perf-test/services"
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
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
var boolVerbose bool
var boolDebug bool
var configOverrides *perfTestUtils.Config
var uiMode bool

const (
	trainingMode = 1
	testingMode = 2
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
	numTestCases := len(testSuite.TestDefinitions) //convenience variable

	if checkTestReadiness {
		readyForTest, _ := perfTestUtils.IsReadyForTest(configurationSettings, numTestCases)
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
			readyForTest, _ := perfTestUtils.IsReadyForTest(configurationSettings, numTestCases)
			if !readyForTest {
				runInTrainingMode(configurationSettings.ExecutionHost, false, testSuite)
			} else {
				log.Warn("System is ready for testing. Training is not required.")
			}
		}
	} else {
		readyForTest, basePerfStats := perfTestUtils.IsReadyForTest(configurationSettings, numTestCases)
		if readyForTest {
			runInTestingMode(basePerfStats, configurationSettings.ExecutionHost, perfTestUtils.GenerateTemplateReport, testSuite)
		} else {
			log.Warn("System is not ready for testing. Attempting to run training mode....")
			runInTrainingMode(configurationSettings.ExecutionHost, false, testSuite)
			readyForTest, basePerfStats = perfTestUtils.IsReadyForTest(configurationSettings, numTestCases)
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

	//Retrieve the URL from command line for using the user interface
	flag.BoolVar(&uiMode, "ui", false, "Get URL for User Interface")

	flag.StringVar(&configFilePath, "configFilePath", "", "The location of the configuration file.")
	flag.BoolVar(&checkTestReadiness, "checkTestReadiness", false, "Simple check to see if system requires training.")

	// Log level simplified for the end user.
	flag.BoolVar(&boolVerbose, "v", false, "Set logging verbosity to 'info' from default of 'warn'. Use -vv for debug.")
	flag.BoolVar(&boolDebug, "vv", false, "Set verbosity to debug.")

	// Controls that reside in Config struct:
	flag.BoolVar(&configurationSettings.GBS, "gbs", false, "Generate 'Base Statistics' for this server")
	flag.BoolVar(&configurationSettings.ReBaseMemory, "reBaseMemory", false, "Generate new base peak memory for this server")
	flag.BoolVar(&configurationSettings.ReBaseAll, "reBaseAll", false, "Generate new base for memory and service response times for this server")

	// Config options that can be overridden from command line args.
	configOverrides = new(perfTestUtils.Config)
	flag.StringVar(&configOverrides.APIName, "apiName", "", "Uniqe name of the test scenario. Do not use spaces. (Default_API_NAME)")
	flag.StringVar(&configOverrides.TargetHost, "host", "", "Remote target host for service calls. (localhost)")
	flag.StringVar(&configOverrides.TargetPort, "port", "", "Remote target port for service calls. (8080)")
	flag.IntVar(&configOverrides.NumIterations, "i", 0, "Number of iterations. (1000)")
	flag.Float64Var(&configOverrides.AllowablePeakMemoryVariance, "allowedMemVar", 0.0, "Allowed peak memory variance percent. (15)")
	flag.Float64Var(&configOverrides.AllowableServiceResponseTimeVariance, "allowedTimeVar", 0.0, "Allowed response time variance percent. (15)")
	flag.StringVar(&configOverrides.TestCaseDir, "tcDir", "", "Path to the testCase definition XML files. (./definitions/testCases)")
	flag.StringVar(&configOverrides.TestSuiteDir, "tsDir", "", "Path to the testSuite defnition XML files. [Optional. Leave off for service-based test runs.] (./definitions/testSuites)")
	flag.StringVar(&configOverrides.BaseStatsOutputDir, "statsDir", "", "Path to base stats files. Filename based on apiName. (./envStats)")
	flag.StringVar(&configOverrides.ReportOutputDir, "rDir", "", "Path to report output files. Filename based on apiName. (./reports)")
	flag.IntVar(&configOverrides.ConcurrentUsers, "u", 0, "Concurrent users, or number of threads. (1)")
	flag.StringVar(&configOverrides.TestSuite, "ts", "", "Name of testSuite definition file. [Optional. Leave off for service-based test runs.]")
	flag.StringVar(&configOverrides.MemoryEndpoint, "mem", "", "Override endpoint in URL for memory metrics. (/debug/vars)")
	flag.IntVar(&configOverrides.RequestDelay, "d", 0, "Delay between calls in ms. (1)")
	flag.IntVar(&configOverrides.TPSFreq, "tps", 0, "Delay between TPS reporting log events in sec. (30)")
	flag.IntVar(&configOverrides.RampUsers, "ru", 0, "Number of users/threads to batch for ramp up. (0)")
	flag.IntVar(&configOverrides.RampDelay, "rd", 0, "Seconds between user/thread batches for ramp up. (15)")

	// Parse the args!
	flag.CommandLine.Parse(args)

	// Start server for the User Interface Mode
	if uiMode {
		log.Print("http:\\localhost:9191")
		services.StartUiMode()
	}

	setLogLevel(boolVerbose, boolDebug)

	// Override defaults with args.
	overrideConfigOpts()





	//----- Parse the config file if specified.
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
			xmlError := xml.Unmarshal(fileContent, &configurationSettings)
			if xmlError != nil {
				log.Error("Failed to parse config file ", configFilePath, ". Error:", xmlError, " - Using default values.")
			}
		}
	}

	// If a config file was loaded, override with args.
	overrideConfigOpts()
}

//----- setLogLevel -----------------------------------------------------------
// Set log level using a simplified interface for end user. See "Process
// command line args" in initConfig().
func setLogLevel(verbose, debug bool) {
	// Set default to WarnLevel
	log.SetLevel(log.WarnLevel)

	// Increase verbosity as set by user at command line.
	if verbose {
		log.SetLevel(log.InfoLevel)
	}
	if debug {
		log.SetLevel(log.DebugLevel)
	}
}

//----- overrideConfigOpts ----------------------------------------------------
// Now that configs have been loaded from file (if specified), override any
// config options that have been specified on command line.
func overrideConfigOpts() {
	if configOverrides.APIName != "" {
		configurationSettings.APIName = configOverrides.APIName
	}
	if configOverrides.TargetHost != "" {
		configurationSettings.TargetHost = configOverrides.TargetHost
	}
	if configOverrides.TargetPort != "" {
		configurationSettings.TargetPort = configOverrides.TargetPort
	}
	if configOverrides.NumIterations != 0 {
		configurationSettings.NumIterations = configOverrides.NumIterations
	}
	if configOverrides.AllowablePeakMemoryVariance != 0.0 {
		configurationSettings.AllowablePeakMemoryVariance = configOverrides.AllowablePeakMemoryVariance
	}
	if configOverrides.AllowableServiceResponseTimeVariance != 0.0 {
		configurationSettings.AllowableServiceResponseTimeVariance = configOverrides.AllowableServiceResponseTimeVariance
	}
	if configOverrides.TestCaseDir != "" {
		configurationSettings.TestCaseDir = configOverrides.TestCaseDir
	}
	if configOverrides.TestSuiteDir != "" {
		configurationSettings.TestSuiteDir = configOverrides.TestSuiteDir
	}
	if configOverrides.BaseStatsOutputDir != "" {
		configurationSettings.BaseStatsOutputDir = configOverrides.BaseStatsOutputDir
	}
	if configOverrides.ReportOutputDir != "" {
		configurationSettings.ReportOutputDir = configOverrides.ReportOutputDir
	}
	if configOverrides.ConcurrentUsers != 0 {
		configurationSettings.ConcurrentUsers = configOverrides.ConcurrentUsers
	}
	if configOverrides.TestSuite != "" {
		configurationSettings.TestSuite = configOverrides.TestSuite
	}
	if configOverrides.MemoryEndpoint != "" {
		configurationSettings.MemoryEndpoint = configOverrides.MemoryEndpoint
	}
	if configOverrides.RequestDelay != 0 {
		configurationSettings.RequestDelay = configOverrides.RequestDelay
	}
	if configOverrides.TPSFreq != 0 {
		configurationSettings.TPSFreq = configOverrides.TPSFreq
	}
	if configOverrides.RampUsers != 0 {
		configurationSettings.RampUsers = configOverrides.RampUsers
	}
	if configOverrides.RampDelay != 0 {
		configurationSettings.RampDelay = configOverrides.RampDelay
	}
}

//----- runInTrainingMode -----------------------------------------------------
func runInTrainingMode(host string, reBaseAll bool, testSuite *testStrategies.TestSuite) {
	log.Info("Running performance test in Training mode for host ", host)

	// Start test timer.
	scenarioTimeStart := time.Now()

	// Initialize the performance statistics struct.
	perfStatsForTest := &perfTestUtils.PerfStats{
		TestTimeStart:        scenarioTimeStart,
		ServiceResponseTimes: make(map[string]int64),
		ServiceTransCount:    make(map[string]*uint64),
		ServiceErrorCount:    make(map[string]*uint64),
		ServiceTPS:           make(map[string]float64),
	}

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
		f, _ := os.Open(configurationSettings.BaseStatsOutputDir + "/" + host + "-" + configurationSettings.APIName + "-perfBaseStats")
		basePerfstats, _ = perfTestUtils.ReadBasePerfFile(f)
	}

	//Run the test
	runTests(perfStatsForTest, trainingMode, testSuite, scenarioTimeStart)
	scenarioTimeElapsed := time.Since(scenarioTimeStart)
	perfStatsForTest.TestTimeEnd = time.Now()

	//Generate base statistics output file for this training run.
	perfTestUtils.GenerateEnvBasePerfOutputFile(perfStatsForTest, basePerfstats, configurationSettings, os.Exit, osFileSystem)

	log.Info("Training mode completed successfully. ")
	log.Infof("Execution Run Time [%v]", scenarioTimeElapsed)
}

//----- runInTestingMode ------------------------------------------------------
func runInTestingMode(
	basePerfstats *perfTestUtils.BasePerfStats,
	host string,
	frg func(*perfTestUtils.BasePerfStats, *perfTestUtils.PerfStats, *perfTestUtils.Config, perfTestUtils.FileSystem, string, string),
	testSuite *testStrategies.TestSuite,
) {
	log.Info("Running Performance test in Testing mode for host ", host)

	// Start test timer. This will give us a basis for all TPS calculations,
	// and will enable the engineer to:
	//     o  Adjust config.NumIterations to control the overall length of the
	//        test run.
	//     o  Set config.ConcurrentUsers to adjust load (see documentation).
	scenarioTimeStart := time.Now()

	// Initialize performance statistics struct.
	perfStatsForTest := &perfTestUtils.PerfStats{
		TestTimeStart:        scenarioTimeStart,
		ServiceResponseTimes: make(map[string]int64),
		ServiceTransCount:    make(map[string]*uint64),
		ServiceErrorCount:    make(map[string]*uint64),
		ServiceTPS:           make(map[string]float64),
	}

	// Run the test.
	runTests(perfStatsForTest, testingMode, testSuite, scenarioTimeStart)

	// Stop the timer. See comment on scenarioTimeStart above.
	scenarioTimeElapsed := time.Since(scenarioTimeStart)
	perfStatsForTest.TestTimeEnd = time.Now()

	// Save overall TPS.
	perfStatsForTest.OverAllTPS = perfTestUtils.CalcTps(perfStatsForTest.OverAllTransCount, scenarioTimeElapsed)

	// Save per-service TPS.
	for key, val := range perfStatsForTest.ServiceTransCount {
		perfStatsForTest.ServiceTPS[key] = perfTestUtils.CalcTps(*val, scenarioTimeElapsed)
		log.Debugf("ServiceName[%v]=%v TPS=%v",
			key,
			*val,
			perfTestUtils.CalcTps(*val, scenarioTimeElapsed),
		)
	}

	// Validate test results
	assertionFailures := runAssertions(basePerfstats, perfStatsForTest)

	// Generate performance test report
	frg(basePerfstats, perfStatsForTest, configurationSettings, osFileSystem, testSuite.Name, testSuite.TestStrategy)

	// Print test results to std out at log level "INFO".
	log.Info("=================== TEST RESULTS ===================")
	if len(assertionFailures) > 0 {
		log.Info("Number of Failures : ", len(assertionFailures))
		for _, failure := range assertionFailures {
			log.Info(failure)
		}
	} else {
		log.Info("Testing mode completed successfully")
	}

	log.Infof("Scenario Time:   [%v]", scenarioTimeElapsed)
	log.Infof("Overall Trans:   [%d]", perfStatsForTest.OverAllTransCount)
	log.Infof("Overall TPS:     [%f]", perfStatsForTest.OverAllTPS)
	log.Info("=====================================================")

	if len(assertionFailures) > 0 {
		os.Exit(1)
	}
}

//----- runTests --------------------------------------------------------------
// This function does two things,
// 1. Start a go routine to periodically grab the memory foot print and set the
//    peak memory value.
// 2. Run all test cases depending on Service-based or Suite-based strategy.
func runTests(perfStatsForTest *perfTestUtils.PerfStats, mode int, testSuite *testStrategies.TestSuite, scenarioTimeStart time.Time) {
	// Initialize Memory analysis.
	var peakMemoryAllocation = new(uint64)
	memoryAudit := make([]uint64, 0)
	testPartitions := make([]perfTestUtils.TestPartition, 0)
	counter := 0
	testPartitions = append(testPartitions, perfTestUtils.TestPartition{Count: counter, TestName: "StartUp"})

	// 1. Start go routine to grab memory in use.
	// Peak memory is stored in peakMemoryAllocation variable.
	chanQuitPkMem := make(chan bool)
	go func() {
		for {
			select {
			case <-chanQuitPkMem:
				return
			default:
				memoryStatsURL := "http://" + configurationSettings.TargetHost + ":" + configurationSettings.TargetPort + configurationSettings.MemoryEndpoint
				resp, err := http.Get(memoryStatsURL)
				if err != nil {
					log.Error("Memory analysis unavailable. Failed to retrieve memory Statistics from endpoint ", memoryStatsURL, ". Error: ", err)
					return
				} else {
					body, _ := ioutil.ReadAll(resp.Body)
					resp.Body.Close()
					m := new(perfTestUtils.Entry)
					unmarshalErr := json.Unmarshal(body, m)
					if unmarshalErr != nil {
						log.Error("Memory analysis unavailable. Failed to unmarshal memory statistics from endpoint: ", memoryStatsURL, ". UnmarsahlErr: ", unmarshalErr)
						return
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

	// Add a 1 second delay before running test case to allow the graph to get
	// some initial memory data before test cases are executed.
	time.Sleep(time.Second * 1)

	// 2. Execute tests based on strategy defaulting to ServiceBasedTesting.
	if testSuite.TestStrategy == testStrategies.SuiteBasedTesting {
		// SuiteBasedTesting strategy runs service requests in the order
		// they are defined in the suite config file. Each full suite
		// definition is run concurrently across the number of threads defined
		// by the config.ConcurrentUsers value. Each thread runs the scenario
		// for config.NumIterations number of times. Usually used for capacity
		// and longevity test runs against a live back end.
		log.Info("Running Suite Based Testing Strategy. Suite Name: [", testSuite.Name, "]")

		// Execute the suite.
		allServicesResponseTimesMap := testStrategies.ExecuteTestSuiteWrapper(
			testSuite,
			configurationSettings,
			perfStatsForTest,
			scenarioTimeStart,
		)

		// Collate the service-level response time data.
		for serviceName, serviceResponseTimes := range allServicesResponseTimesMap {
			averageResponseTime := perfTestUtils.CalcAverageResponseTime(serviceResponseTimes, mode)
			if averageResponseTime == 0 && mode == trainingMode {
				// If all response times average to zero, all attempts to call the
				// service failed. In training mode, abort so the problem can be
				// remedied. In testing mode, continue, but record the zero.
				log.Error("Training mode failed due to invalid response on service [Name:", serviceName, "]")
				os.Exit(1)
			}
			perfStatsForTest.ServiceResponseTimes[serviceName] = averageResponseTime
		}
	} else {
		// ServiceBasedTesting strategy runs sequentially through all test
		// cases in the config.TestCaseDir folder for config.NumIterations
		// number of times in parallel across the number of threads defined by
		// the config.ConcurrentUsers value. Usually used with mock calls.
		log.Info("Running Service Based Testing Strategy")

		// Determine load per concurrent user.
		loadPerUser := int(configurationSettings.NumIterations / configurationSettings.ConcurrentUsers)
		remainder := configurationSettings.NumIterations % configurationSettings.ConcurrentUsers

		// Set the overall TransCount, which will subsequently be used to
		// calculate OverallTPS (see runInTestingMode() above).
		perfStatsForTest.OverAllTransCount = uint64(len(testSuite.TestCases) * configurationSettings.NumIterations)

		log.Infof("ServiceBasedTesting loadPerUser=[%d] remainder=[%d]", loadPerUser, remainder)

		var index int
		var testDefinition *testStrategies.TestDefinition
		for index, testDefinition = range testSuite.TestDefinitions {
			log.Infof("Running Test case [%d] [Name:%s]", index, testDefinition.TestName)
			testPartitions = append(testPartitions, perfTestUtils.TestPartition{Count: counter, TestName: testDefinition.TestName})
			averageResponseTime := testStrategies.ExecuteServiceTest(testDefinition, loadPerUser, remainder, configurationSettings, mode)

			if averageResponseTime > 0 {
				perfStatsForTest.ServiceResponseTimes[testDefinition.TestName] = averageResponseTime
			} else {
				if mode == trainingMode {
					//Fail fast on training mode if any requests fail. If training fails we cannot guarantee the results.
					log.Error("Training mode failed due to invalid response on service [Name:", testDefinition.TestName, "]")
					os.Exit(1)
				}
			}
		}
	}

	// Kill the peak memory thread to avoid race condition when saving metrics.
	close(chanQuitPkMem)

	// Save the peak memory metrics:
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
		varianceOk := perfTestUtils.ValidateAverageServiceResponseTimeVariance(configurationSettings.AllowableServiceResponseTimeVariance, responseTimeVariancePercentage)
		if !varianceOk {
			assertionFailures = append(assertionFailures, fmt.Sprintf("Service Failure: Service test %-60s response time variance exceeded by %3.2f %1s", serviceName, responseTimeVariancePercentage, "%"))
		}
	}
	return assertionFailures
}
