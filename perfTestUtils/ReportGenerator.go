package perfTestUtils

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

func GenerateReport(basePerfstats *BasePerfStats, perfStats *PerfStats, configurationSettings *Config) {

	fileContent, fileErr := ioutil.ReadFile("./report/template.html")
	if fileErr != nil {
		fmt.Println(fileErr)
	}
	stringContents := string(fileContent)

	//Add Test date to report
	stringContents = strings.Replace(stringContents, "###testDate###", time.Now().Format(time.RFC850), 1)

	//Add Allowed Variace values
	stringContents = strings.Replace(stringContents, "###allowedMemoryVariance###", fmt.Sprintf("%4.2f", configurationSettings.AllowablePeakMemoryVariance), 1)
	stringContents = strings.Replace(stringContents, "###allowedServiceRespTimeVariance###", fmt.Sprintf("%4.2f", configurationSettings.AllowableServiceResponseTimeVariance), 1)

	//Add memory stats
	peakMemoryVariancePercentage := float64(CalcPeakMemoryVariancePercentage(basePerfstats.BasePeakMemory, perfStats.PeakMemory))
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

		baseTimeMillis := float64(float32(baseResponseTime) / float32(1000000))
		serviceResponseTimesBase = append(serviceResponseTimesBase, []byte(strconv.FormatFloat(baseTimeMillis, 'f', 3, 64))...)
		serviceResponseTimesBase = append(serviceResponseTimesBase, []byte(",")...)

		if averageServiceResponseTime == 0 {

			serviceResponseTimesTest = append(serviceResponseTimesTest, []byte(strconv.FormatFloat(0.0, 'f', 3, 64))...)
			serviceResponseTimesTest = append(serviceResponseTimesTest, []byte(",")...)

			serviceNames = append(serviceNames, []byte(serviceName+" ("+fmt.Sprintf("%3.2f", 0.0)+" %)")...)
			serviceNames = append(serviceNames, []byte("','")...)

			varianceError := "style=\"color:red\""
			servicesPass = false

			serviceResponseTimesTable = append(serviceResponseTimesTable, []byte("<tr height=10px><td>"+serviceName+"</td><td>"+fmt.Sprintf("%4.3f", baseTimeMillis)+"</td><td>"+"FAILED"+"</td><td "+varianceError+">"+"FAILED"+"</td></tr>\n")...)

		} else {
			responseTimeVariancePercentage := CalcAverageResponseVariancePercentage(averageServiceResponseTime, baseResponseTime)

			testTimeMillis := float64(float32(averageServiceResponseTime) / float32(1000000))

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
	file, err := os.Create(configurationSettings.ReportOutputDir + "/PerformanceReport.html")
	if err != nil {
		defer file.Close()
	}
	file.Write([]byte(stringContents))
}
