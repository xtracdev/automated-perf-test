package main

import (
	"automated-perf-test/perfTestUtils"
	"bytes"
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	//log "github.com/Sirupsen/logrus"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
	"xtrac/api/common/test/backendServerMock"
)

var configurationSettings *perfTestUtils.Config

func init() {

	//Command line ags
	var gbs bool
	var resetPeakMemory bool
	var configFilePath string

	//Process command line arugments.
	flag.BoolVar(&gbs, "gbs", false, "Genertate Base Performance Staticists for this server")
	flag.BoolVar(&resetPeakMemory, "resetPeakMemory", false, "Generate new base peak memory for this server")
	flag.StringVar(&configFilePath, "configFilePath", "", "The location of the configuration file.")
	flag.Parse()

	//Read and paser config file if present.
	configurationSettings = new(perfTestUtils.Config)
	if configFilePath != "" {
		fileContent, fileErr := ioutil.ReadFile(configFilePath)
		if fileErr != nil {
			//log.Info("No config file found. ")
			fmt.Println("No config file found. ")
			os.Exit(1)
		} else {
			xmlError := xml.Unmarshal(fileContent, &configurationSettings)
			if xmlError != nil {
				//log.Info("Failed to parse config file ", configFilePath, ". Error:", xmlError)
				fmt.Println("Failed to parse config file ", configFilePath, ". Error:", xmlError)
				os.Exit(1)
			}
		}

	}

	//Get Hostname for this machine.
	host, err := os.Hostname()
	if err != nil {
		//log.Error("Failed to resolve host name. Error:", err)
		fmt.Println("Failed to resolve host name. Error:", err)
		os.Exit(1)
	}
	configurationSettings.ExecutionHost = host
	configurationSettings.GBS = gbs
	configurationSettings.ResetPeakMemory = resetPeakMemory
}

//Main Test Method
func main() {

	//Validate config()
	configurationSettings.PrintAndValidateConfig()

	//initilize Performance statistics struct for this test run
	perfStatsForTest := &perfTestUtils.PerfStats{ServiceResponseTimes: make(map[string]int64)}

	if configurationSettings.GBS {
		runInTrainingMode(perfStatsForTest, configurationSettings.ExecutionHost)
	} else {
		runInTestingMode(perfStatsForTest, configurationSettings.ExecutionHost)
	}
}

func runInTrainingMode(perfStatsForTest *perfTestUtils.PerfStats, host string) {
	//log.Info("Running Perf test in Training mode for host ", host)
	fmt.Println("Running Perf test in Training mode for host ", host)

	//Check to see if this server already has a base perf file defined.
	//If so, only values not previously populated will be set.
	//if not, a default base perf struct is created with nil values for all fields
	basePerfstats, _ := perfTestUtils.ReadBasePerfFile(host)

	//Run the test
	runTests(perfStatsForTest)

	//Set base performance based on this test run
	populateBasePerfStats(perfStatsForTest, basePerfstats)

	//Convert base perf stat to Json and write out to file
	basePerfstatsJson, err := json.Marshal(basePerfstats)
	if err != nil {
		//log.Error("Failed to marshal to Json. Error:", err)
		fmt.Println("Failed to marshal to Json. Error:", err)
		os.Exit(1)
	}
	file, err := os.Create("./envStats/" + host + "-perfBaseStats")
	if err != nil {
		//log.Error("Failed to create output file. Error:", err)
		fmt.Println("Failed to create output file. Error:", err)
		os.Exit(1)
	}
	defer file.Close()
	file.Write(basePerfstatsJson)
}

func runInTestingMode(perfStatsForTest *perfTestUtils.PerfStats, host string) {
	//log.Info("Running Perf test in Testing mode for host ", host)
	fmt.Println("Running Perf test in Testing mode for host ", host)
	//read in perf base stats
	basePerfstats, err := perfTestUtils.ReadBasePerfFile(host)
	if err != nil {
		//log.Error("Failed to read env stats for " + host + ". Error:" + err.Error() + ". Run go test -gbs to generate base performance statistics for this server.")
		fmt.Println("Failed to read env stats for " + host + ". Error:" + err.Error() + ". Run go test -gbs to generate base performance statistics for this server.")
		os.Exit(1)
	}

	validateTestDefinitionAmount(len(basePerfstats.BaseServiceResponseTimes))
	runTests(perfStatsForTest)
	runAssertions(basePerfstats, perfStatsForTest)
	generateReport(basePerfstats, perfStatsForTest)

}

func populateBasePerfStats(perfStatsForTest *perfTestUtils.PerfStats, basePerfstats *perfTestUtils.BasePerfStats) {
	modified := false
	if basePerfstats.BasePeakMemory == 0 || configurationSettings.ResetPeakMemory {
		basePerfstats.BasePeakMemory = perfStatsForTest.PeakMemory
		modified = true
	}
	for serviceName, responseTime := range perfStatsForTest.ServiceResponseTimes {
		serviceBaseResponseTime := basePerfstats.BaseServiceResponseTimes[serviceName]
		if serviceBaseResponseTime == 0 {
			basePerfstats.BaseServiceResponseTimes[serviceName] = responseTime
			modified = true
		}
	}
	if basePerfstats.MemoryAudit == nil || len(basePerfstats.MemoryAudit) == 0 {
		basePerfstats.MemoryAudit = perfStatsForTest.MemoryAudit
		modified = true
	}
	currentTime := time.Now().Format(time.RFC850)
	if basePerfstats.GenerationDate == "" {
		basePerfstats.GenerationDate = currentTime
	}

	if modified {
		basePerfstats.ModifiedDate = currentTime
	}
}

//This function does two thing,
//1 Start a go routine to preiodically grab the memory foot print and set the peak memory value
//2 Run all test using mock servers and gather performance stats
func runTests(perfStatsForTest *perfTestUtils.PerfStats) {

	var peakMemoryAllocation = new(uint64)
	var lastServiceName = "StartUp"
	var currentServiceName = "StartUp"

	memoryAudit := make([]uint64, 0)
	testPartitions := make([]perfTestUtils.TestPartition, 0)
	counter := 0
	testPartitions = append(testPartitions, perfTestUtils.TestPartition{Count: counter, TestName: currentServiceName})

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
					//log.Error("Memory analysis unavailable. Failed to retrieve memory Statistics from endpoint ", memoryStatsUrl)
					fmt.Println("Memory analysis unavailable. Failed to retrieve memory Statistics from endpoint ", memoryStatsUrl)
					quit <- true
				} else {
					body, _ := ioutil.ReadAll(resp.Body)

					defer resp.Body.Close()

					m := new(perfTestUtils.Entry)
					unmarshalErr := json.Unmarshal(body, m)
					if unmarshalErr != nil {
						//log.Error("Memory analysis unavailable. Failed to unmarshal memory statistics. ", unmarshalErr)
						fmt.Println("Memory analysis unavailable. Failed to unmarshal memory statistics. ", unmarshalErr)
						quit <- true
					} else {
						if m.Memstats.Alloc > *peakMemoryAllocation {
							*peakMemoryAllocation = m.Memstats.Alloc
						}
						memoryAudit = append(memoryAudit, m.Memstats.Alloc)

						if lastServiceName != currentServiceName {
							testPartitions = append(testPartitions, perfTestUtils.TestPartition{Count: counter, TestName: currentServiceName})
							lastServiceName = currentServiceName
						}

						counter++
						time.Sleep(time.Millisecond * 200)
					}
				}
			}
		}
	}()

	//Read test case files from test defination directory
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
	if len(fi) == 0 {
		//log.Error("No test case files found in specified directory ", configurationSettings.TestDefinationsDir)
		fmt.Println("No test case files found in specified directory ", configurationSettings.TestDefinationsDir)
		os.Exit(1)
	}

	//Add a 1 second delay before running test case to allow the graph get some initial memory data before test cases are executed.
	time.Sleep(time.Second * 1)
	for _, fi := range fi {
		bs, err := ioutil.ReadFile(configurationSettings.TestDefinationsDir + "/" + fi.Name())
		if err != nil {
			//log.Error("Failed to read test file. Filename: ", fi.Name(), err)
			fmt.Println("Failed to read test file. Filename: ", fi.Name(), err)
			continue
		}

		testDefination := new(perfTestUtils.TestDefination)
		xml.Unmarshal(bs, &testDefination)

		//log.Info("Running Test case [Name:", testDefination.TestName, ", File name:", fi.Name(), "]")
		fmt.Println("Running Test case [Name:", testDefination.TestName, ", File name:", fi.Name(), "]")
		currentServiceName = testDefination.TestName
		perfStatsForTest.ServiceResponseTimes[testDefination.TestName] = executeServiceTest(testDefination)
		time.Sleep(time.Millisecond * 200)
	}

	time.Sleep(time.Second * 1)
	perfStatsForTest.PeakMemory = *peakMemoryAllocation
	perfStatsForTest.MemoryAudit = memoryAudit
	perfStatsForTest.TestPartitions = testPartitions
}

//Single execution function for all service test.
//Runs multiple invocations of the test based on num iterations parameter
func executeServiceTest(testDefination *perfTestUtils.TestDefination) int64 {

	averageResponseTime := int64(0)
	loopExecutedToCompletion := true
	responseTimes := make(perfTestUtils.RspTimes, configurationSettings.NumIterations)
	//Execute the test in a loop
	for i := 0; i < configurationSettings.NumIterations; i++ {

		var req *http.Request

		if !testDefination.Multipart {
			if testDefination.Payload != "" {
				req, _ = http.NewRequest(testDefination.HttpMethod, "http://"+configurationSettings.TargetHost+":"+configurationSettings.TargetPort+testDefination.BaseUri, strings.NewReader(testDefination.Payload))
			} else {
				req, _ = http.NewRequest(testDefination.HttpMethod, "http://"+configurationSettings.TargetHost+":"+configurationSettings.TargetPort+testDefination.BaseUri, nil)
			}
		} else {
			if testDefination.HttpMethod != "POST" {
				//log.Fatal("Multipart request has to be 'POST' method.")
				fmt.Println("Multipart request has to be 'POST' method.")
			} else {
				body := new(bytes.Buffer)
				writer := multipart.NewWriter(body)
				for _, field := range testDefination.MultipartPayload {
					//log.Debugf("field: %s\n", field)
					fmt.Println(fmt.Sprintf("field: %s\n", field))
					if field.FileName == "" {
						writer.WriteField(field.FieldName, field.FieldValue)
					} else {
						part, _ := writer.CreateFormFile(field.FieldName, field.FileName)
						io.Copy(part, bytes.NewReader(field.FileContent))
					}
				}
				writer.Close()
				req, _ = http.NewRequest(testDefination.HttpMethod, "http://"+configurationSettings.TargetHost+":"+configurationSettings.TargetPort+testDefination.BaseUri, body)
				req.Header.Set("Content-Type", writer.FormDataContentType())
			}
		}

		req.Header.Add(backendServerMock.SCENARIO_HEADER_KEY, testDefination.Scenario)

		//add headers
		for _, v := range testDefination.Headers {
			req.Header.Add(v.Key, v.Value)
		}
		req.Header.Set("xtracToken", testDefination.XtracToken)
		startTime := time.Now()
		if resp, err := (&http.Client{}).Do(req); err != nil {
			//log.Error("Error by firing request: ", req, "Error:", err)
			fmt.Println("Error by firing request: ", req, "Error:", err)
		} else {

			timeTaken := time.Since(startTime)

			body, _ := ioutil.ReadAll(resp.Body)

			//Validate service response
			contentLengthOk := perfTestUtils.ValidateResponseBody(body, testDefination.TestName)
			responseCodeOk := perfTestUtils.ValidateResponseStatusCode(resp.StatusCode, testDefination.ResponseStatusCode, testDefination.TestName)
			responseTimeOK := perfTestUtils.ValidateServiceResponseTime(timeTaken.Nanoseconds(), testDefination.TestName)

			if contentLengthOk && responseCodeOk && responseTimeOK {
				responseTimes[i] = timeTaken.Nanoseconds()
			} else {
				loopExecutedToCompletion = false
				break
			}
		}
	}

	if loopExecutedToCompletion {
		averageResponseTime = perfTestUtils.CalcAverageResponseTime(responseTimes, configurationSettings.NumIterations)
	}
	return averageResponseTime
}

//This function runs the assertions to ensure memory and service have not deviated past the allowed variance
func runAssertions(basePerfstats *perfTestUtils.BasePerfStats, perfStats *perfTestUtils.PerfStats) {

	fmt.Println("\n===================================== Test Run statistics ====================================================================")
	fmt.Println("Memory")

	//Asserts Peak memory growth has not exceeded the allowable variance
	peakMemoryVariancePercentage := perfTestUtils.CalcPeakMemoryVariancePercentage(basePerfstats.BasePeakMemory, perfStats.PeakMemory)
	varianceOk := perfTestUtils.ValidatePeakMemoryVariance(configurationSettings.AllowablePeakMemoryVariance, peakMemoryVariancePercentage)
	fmt.Printf("%4s %-30s %10d %5s %4.2f %5s", "\t", "Base Memory", basePerfstats.BasePeakMemory, "B   (", (float32(basePerfstats.BasePeakMemory)/float32(1024))/float32(1024), "MB)\n")
	fmt.Printf("%4s %-30s %10d %5s %4.2f %5s", "\t", "Peak Memory", perfStats.PeakMemory, "B   (", (float32(perfStats.PeakMemory)/float32(1024))/float32(1024), "MB)\n")
	if !varianceOk {
		fmt.Printf("\x1b[31;1m")
	}
	fmt.Printf("%4s %-30s %9.2f %1s", "\t", "Peak Memory Variance", peakMemoryVariancePercentage, "%\n")
	fmt.Printf("\x1b[0m")
	fmt.Println("")

	//Assert Each service response time has not exceeded the allowable variance
	serviceCountsok := perfTestUtils.ValidateTestCaseCount(len(basePerfstats.BaseServiceResponseTimes), len(perfStats.ServiceResponseTimes))
	if serviceCountsok {
		fmt.Println("Services Resposne Times")
		fmt.Printf("\t ------------------------------------------------------------------------------------------------------------\n")
		fmt.Printf("%4s | %-40s | %20s | %20s | %15s | %1s", "\t", "TestName", "BaseTime (Milli)", "TestTime (Milli)", "%variance", "\n")
		fmt.Printf("\t ------------------------------------------------------------------------------------------------------------\n")
		for serviceName, baseResponseTime := range basePerfstats.BaseServiceResponseTimes {
			averageServiceResponseTime := perfStats.ServiceResponseTimes[serviceName]
			//assert.True(t, averageServiceResponseTime > 0, "Average Time taken to complete request %s was 0 nanoseconds", serviceName)

			responseTimeVariancePercentage := perfTestUtils.CalcAverageResponseVariancePercentage(averageServiceResponseTime, baseResponseTime)
			varianceOk := perfTestUtils.ValidateAverageServiceResponeTimeVariance(configurationSettings.AllowableServiceResponseTimeVariance, responseTimeVariancePercentage, serviceName)
			if !varianceOk {
				fmt.Printf("\x1b[31;1m")
			}
			fmt.Printf("%4s | %-40s | %20.2f | %20.2f | %15.2f | %1s", "\t", serviceName, float32(float32(baseResponseTime)/float32(1000000)), float32(float32(averageServiceResponseTime)/float32(1000000)), responseTimeVariancePercentage, " \n")
			fmt.Printf("\x1b[0m")
		}
		fmt.Printf("\t ------------------------------------------------------------------------------------------------------------\n")
	}
	fmt.Println("===============================================================================================================================")

}

func generateReport(basePerfstats *perfTestUtils.BasePerfStats, perfStats *perfTestUtils.PerfStats) {

	fileContent, fileErr := ioutil.ReadFile("./report/template.html")
	if fileErr != nil {
		fmt.Println(fileErr)
	}
	stringContents := string(fileContent)

	//Add Test date to report
	stringContents = strings.Replace(stringContents, "###testDate###", time.Now().Format(time.RFC850), 1)

	//Add Allowed Variace values
	stringContents = strings.Replace(stringContents, "###testDate###", time.Now().Format(time.RFC850), 1)
	stringContents = strings.Replace(stringContents, "###allowedMemoryVariance###", fmt.Sprintf("%4.2f", configurationSettings.AllowablePeakMemoryVariance), 1)
	stringContents = strings.Replace(stringContents, "###allowedServiceRespTimeVariance###", fmt.Sprintf("%4.2f", configurationSettings.AllowableServiceResponseTimeVariance), 1)

	//Add memory stats
	peakMemoryVariancePercentage := float64(perfTestUtils.CalcPeakMemoryVariancePercentage(basePerfstats.BasePeakMemory, perfStats.PeakMemory))
	stringContents = strings.Replace(stringContents, "###basePeakMemory###", strconv.FormatFloat(float64((float32(basePerfstats.BasePeakMemory)/float32(1024))/float32(1024)), 'f', 3, 64), 1)
	stringContents = strings.Replace(stringContents, "###testPeakMemory###", strconv.FormatFloat(float64((float32(perfStats.PeakMemory)/float32(1024))/float32(1024)), 'f', 3, 64), 1)
	stringContents = strings.Replace(stringContents, "###memoryVariance###", strconv.FormatFloat(peakMemoryVariancePercentage, 'f', 3, 64), 1)

	//Set memory Error style if variance is above allowed variance
	if float64(configurationSettings.AllowablePeakMemoryVariance) < peakMemoryVariancePercentage {
		stringContents = strings.Replace(stringContents, "###errorStyle###", "style=\"color:red\"", 1)
		stringContents = strings.Replace(stringContents, "###memoryPassFail###", "<font color=\"red\">FAIL</font>", 1)
	} else {
		stringContents = strings.Replace(stringContents, "###memoryPassFail###", "<font color=\"green\">PASS</font>", 1)
	}

	//Add data to  build memory chart
	baseMemoryAudit := []byte("['Base',")
	for _, memValue := range basePerfstats.MemoryAudit {
		baseMemoryAudit = append(baseMemoryAudit, []byte(strconv.FormatFloat(float64((float32(memValue)/float32(1024))), 'f', 3, 64))...)
		baseMemoryAudit = append(baseMemoryAudit, []byte(",")...)
	}
	baseMemoryAudit = append(baseMemoryAudit, []byte("]")...)

	testMemoryAudit := []byte("['Test',")
	for _, memValue := range perfStats.MemoryAudit {
		testMemoryAudit = append(testMemoryAudit, []byte(strconv.FormatFloat(float64((float32(memValue)/float32(1024))), 'f', 3, 64))...)
		testMemoryAudit = append(testMemoryAudit, []byte(",")...)
	}
	testMemoryAudit = append(testMemoryAudit, []byte("]")...)

	stringContents = strings.Replace(stringContents, "###baseMemoryArray###", string(baseMemoryAudit), 1)
	stringContents = strings.Replace(stringContents, "###testMemoryArray###", string(testMemoryAudit), 1)

	//Define partitions on chart.
	testpatritions := []byte("")
	for _, testPartition := range perfStats.TestPartitions {
		testpatritions = append(testpatritions, []byte("{value: "+fmt.Sprint(testPartition.Count)+" , text: '"+testPartition.TestName+"'},")...)
	}
	stringContents = strings.Replace(stringContents, "###testPartitions###", string(testpatritions), 1)

	servicesPass := true
	//Build service response time chart
	serviceResponseTimesBase := []byte("['Base',")
	serviceResponseTimesTest := []byte("['Test',")
	serviceNames := []byte("['")
	serviceResponseTimesTable := []byte("")
	for serviceName, baseResponseTime := range basePerfstats.BaseServiceResponseTimes {
		averageServiceResponseTime := perfStats.ServiceResponseTimes[serviceName]
		responseTimeVariancePercentage := perfTestUtils.CalcAverageResponseVariancePercentage(averageServiceResponseTime, baseResponseTime)
		baseTimeMillis := float64(float32(baseResponseTime) / float32(1000000))
		testTimeMillis := float64(float32(averageServiceResponseTime) / float32(1000000))

		serviceResponseTimesBase = append(serviceResponseTimesBase, []byte(strconv.FormatFloat(baseTimeMillis, 'f', 3, 64))...)
		serviceResponseTimesBase = append(serviceResponseTimesBase, []byte(",")...)

		serviceResponseTimesTest = append(serviceResponseTimesTest, []byte(strconv.FormatFloat(testTimeMillis, 'f', 3, 64))...)
		serviceResponseTimesTest = append(serviceResponseTimesTest, []byte(",")...)

		serviceNames = append(serviceNames, []byte(serviceName+" ("+fmt.Sprintf("%3.2f", responseTimeVariancePercentage)+" %)")...)
		serviceNames = append(serviceNames, []byte("','")...)

		varianceError := ""
		if configurationSettings.AllowableServiceResponseTimeVariance < responseTimeVariancePercentage {
			varianceError = "style=\"color:red\""
			servicesPass = false
		}
		serviceResponseTimesTable = append(serviceResponseTimesTable, []byte("<tr height=10px><td>"+serviceName+"</td><td>"+fmt.Sprintf("%4.3f", baseTimeMillis)+"</td><td>"+fmt.Sprintf("%4.3f", testTimeMillis)+"</td><td "+varianceError+">"+fmt.Sprintf("%4.2f", responseTimeVariancePercentage)+"</td></tr>\n")...)

	}
	serviceResponseTimesBase = append(serviceResponseTimesBase, []byte("]")...)
	serviceResponseTimesTest = append(serviceResponseTimesTest, []byte("]")...)
	serviceNames = append(serviceNames, []byte("']")...)

	stringContents = strings.Replace(stringContents, "###servcieResponseTimesBase###", string(serviceResponseTimesBase), 1)
	stringContents = strings.Replace(stringContents, "###servcieResponseTimesTest###", string(serviceResponseTimesTest), 1)
	stringContents = strings.Replace(stringContents, "###servcieNames###", string(serviceNames), 1)
	stringContents = strings.Replace(stringContents, "###serviceRespTimesTable###", string(serviceResponseTimesTable), 1)

	if servicesPass {
		stringContents = strings.Replace(stringContents, "###servicePassFail###", "<font color=\"green\">PASS</font>", 1)
	} else {
		stringContents = strings.Replace(stringContents, "###servicePassFail###", "<font color=\"red\">FAIL</font>", 1)
	}

	//Write out the file
	file, err := os.Create("./report/PerformanceReport.html")
	if err != nil {
		defer file.Close()
	}
	file.Write([]byte(stringContents))
}

func validateTestDefinitionAmount(baselineAmount int) {
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

	if definitionAmount != baselineAmount {
		//log.Errorf("Amount of test definition: %d does not equal to baseline amount: %d.", definitionAmount, baselineAmount)
		fmt.Println("Amount of test definition: %d does not equal to baseline amount: %d.", definitionAmount, baselineAmount)
		os.Exit(1)
	}
}
