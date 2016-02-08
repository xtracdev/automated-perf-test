package perfTestUtils

import (
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

type perfStatsModel struct {
	BasePerfStats        *BasePerfStats
	PerfStats            *PerfStats
	Config               *Config
	JsonTimeServiceNames template.JS
}

func (p *perfStatsModel) IsMemoryPass() bool {
	return p.PeakMemoryVariancePercentage() < float64(p.Config.AllowablePeakMemoryVariance)
}

func (p *perfStatsModel) IsServiceTimePass(s string) bool {
	if p.PerfStats.ServiceResponseTimes[s] == 0 {
		return false
	}
	if p.Config.AllowableServiceResponseTimeVariance < CalcAverageResponseVariancePercentage(p.PerfStats.ServiceResponseTimes[s], p.BasePerfStats.BaseServiceResponseTimes[s]) {
		return false
	}
	return true
}

func (p *perfStatsModel) IsTimePass() bool {
	for k, _ := range p.BasePerfStats.BaseServiceResponseTimes {
		if !p.IsServiceTimePass(k) {
			return false
		}
	}
	return true
}

func (p *perfStatsModel) PeakMemoryVariancePercentage() float64 {
	return float64(CalcPeakMemoryVariancePercentage(p.BasePerfStats.BasePeakMemory, p.PerfStats.PeakMemory))
}

func MemoryMB(pm uint64) float64 {
	return float64((float32(pm) / float32(1024)) / float32(1024))
}

func MemoryKB(pm uint64) float64 {
	return float64((float32(pm) / float32(1024)))
}

func FormatMemory(m float64) string {
	return strconv.FormatFloat(m, 'f', 3, 64)
}

func Div(num int64, den int64) float64 {
	return float64(float32(num) / float32(den))
}

func JsonMemoryArray(name string, array []uint64) template.JS {
	jsonMemoryAudit := []byte("['" + name + "',")
	for _, memValue := range array {
		jsonMemoryAudit = append(jsonMemoryAudit, []byte(strconv.FormatFloat(float64((float32(memValue)/float32(1024))), 'f', 3, 64))...)
		jsonMemoryAudit = append(jsonMemoryAudit, []byte(",")...)
	}
	jsonMemoryAudit = append(jsonMemoryAudit, []byte("]")...)
	return template.JS(jsonMemoryAudit)
}

func (p *perfStatsModel) JsonTestPartitions() template.JS {
	testpatritions := []byte("")
	for _, testPartition := range p.PerfStats.TestPartitions {
		testpatritions = append(testpatritions, []byte("{value: "+fmt.Sprint(testPartition.Count)+" , text: '"+testPartition.TestName+"'},")...)
	}
	return template.JS(testpatritions)
}

func (p *perfStatsModel) JsonTimeArray() template.JS {
	serviceResponseTimesBase := []byte("['Base',")
	serviceResponseTimesTest := []byte("['Test',")
	serviceNames := []byte("['")
	for i := 1; i < len(p.PerfStats.TestPartitions); i++ {
		tp := p.PerfStats.TestPartitions[i]
		averageServiceResponseTime := p.PerfStats.ServiceResponseTimes[tp.TestName]

		baseTimeMillis := float64(float32(p.BasePerfStats.BaseServiceResponseTimes[tp.TestName]) / float32(1000000))
		serviceResponseTimesBase = append(serviceResponseTimesBase, []byte(strconv.FormatFloat(baseTimeMillis, 'f', 3, 64))...)
		serviceResponseTimesBase = append(serviceResponseTimesBase, []byte(",")...)

		testTimeMillis := float64(float32(averageServiceResponseTime) / float32(1000000))

		serviceResponseTimesTest = append(serviceResponseTimesTest, []byte(strconv.FormatFloat(testTimeMillis, 'f', 3, 64))...)
		serviceResponseTimesTest = append(serviceResponseTimesTest, []byte(",")...)

		if averageServiceResponseTime != 0 {
			responseTimeVariancePercentage := CalcAverageResponseVariancePercentage(averageServiceResponseTime, p.BasePerfStats.BaseServiceResponseTimes[tp.TestName])
			serviceNames = append(serviceNames, []byte(tp.TestName+" ("+fmt.Sprintf("%3.2f", responseTimeVariancePercentage)+" %)")...)
		} else {
			serviceNames = append(serviceNames, []byte(tp.TestName+" ("+fmt.Sprintf("%3.2f", 0.0)+" %)")...)
		}
		serviceNames = append(serviceNames, []byte("','")...)
	}
	serviceResponseTimesBase = append(serviceResponseTimesBase, []byte("],")...)
	serviceResponseTimesTest = append(serviceResponseTimesTest, []byte("]")...)
	serviceNames = append(serviceNames, []byte("']")...)
	p.JsonTimeServiceNames = template.JS(serviceNames)

	serviceResponseTimesBase = append(serviceResponseTimesBase, serviceResponseTimesTest...)
	return template.JS(serviceResponseTimesBase)
}

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

func GenerateTemplateReport(basePerfstats *BasePerfStats, perfStats *PerfStats, configurationSettings *Config) {
	file, err := os.Create(configurationSettings.ReportOutputDir + "/PerformanceTemplateReport.html")
	if err != nil {
		defer file.Close()
	}
	generateTemplate(basePerfstats, perfStats, configurationSettings, file, ".report")
}

func generateTemplate(bstats *BasePerfStats, pstats *PerfStats, configurationSettings *Config, wr io.Writer, templDir string) error {
	ps := &perfStatsModel{BasePerfStats: bstats, PerfStats: pstats, Config: configurationSettings}
	s1 := template.New("main")
	fmt.Printf("template: %v\n", s1)
	s1 = s1.Funcs(template.FuncMap{"memToMB": MemoryMB, "memToKb": MemoryKB, "formatMem": FormatMemory, "jsonMem": JsonMemoryArray, "div": Div, "avgVar": CalcAverageResponseVariancePercentage})
	s1, err := s1.ParseFiles(templDir+"/header.tmpl", templDir+"/content.tmpl", templDir+"/footer.tmpl")
	if err != nil {
		return fmt.Errorf("Error loading template files: %v", err)
	}
	err = s1.ExecuteTemplate(wr, "header", ps)
	if err != nil {
		return fmt.Errorf("Error executing template: %v", err)
	}
	err = s1.ExecuteTemplate(wr, "content", ps)
	if err != nil {
		return fmt.Errorf("Error executing template: %v", err)
	}
	err = s1.ExecuteTemplate(wr, "footer", nil)
	if err != nil {
		return fmt.Errorf("Error executing template: %v", err)
	}
	return nil
}
