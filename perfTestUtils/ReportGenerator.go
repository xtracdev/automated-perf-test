package perfTestUtils

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"html/template"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
)

type perfStatsModel struct {
	BasePerfStats        *BasePerfStats
	PerfStats            *PerfStats
	Config               *Config
	JSONTimeServiceNames template.JS
	TestStrategy         string
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
	for k := range p.BasePerfStats.BaseServiceResponseTimes {
		if !p.IsServiceTimePass(k) {
			return false
		}
	}
	return true
}

func (p *perfStatsModel) PeakMemoryVariancePercentage() float64 {
	return float64(CalcPeakMemoryVariancePercentage(p.BasePerfStats.BasePeakMemory, p.PerfStats.PeakMemory))
}

func memoryMB(pm uint64) float64 {
	return float64((float32(pm) / float32(1024)) / float32(1024))
}

func formatMemory(m float64) string {
	return strconv.FormatFloat(m, 'f', 3, 64)
}

func div(num int64, den int64) float64 {
	return float64(float32(num) / float32(den))
}

func jsonMemoryArray(name string, array []uint64) template.JS {
	jsonMemoryAudit := []byte("['" + name + "',")
	for _, memValue := range array {
		jsonMemoryAudit = append(jsonMemoryAudit, []byte(strconv.FormatFloat(float64((float32(memValue)/float32(1024))), 'f', 3, 64))...)
		jsonMemoryAudit = append(jsonMemoryAudit, []byte(",")...)
	}
	jsonMemoryAudit = append(jsonMemoryAudit, []byte("]")...)
	return template.JS(jsonMemoryAudit)
}

func (p *perfStatsModel) JSONTestPartitions() template.JS {
	testpatritions := []byte("")
	for _, testPartition := range p.PerfStats.TestPartitions {
		testpatritions = append(testpatritions, []byte("{value: "+fmt.Sprint(testPartition.Count)+" , text: '"+testPartition.TestName+"'},")...)
	}
	return template.JS(testpatritions)
}

// JSONTimeArray returns series data for the chart in json format suitable
// to be inserted as javascript within <script> tags. The array is alpha
// sorted by service name.
func (p *perfStatsModel) JSONTimeArray() template.JS {
	serviceResponseTimesBase := []byte("['Base',")
	serviceResponseTimesTest := []byte("['Test',")
	serviceNames := []byte("['")

	// Sort the keys to get consistent charts between data runs for
	// easier comparison.
	sortedkeys := make([]string, len(p.PerfStats.ServiceResponseTimes))
	i := 0
	for k := range p.PerfStats.ServiceResponseTimes {
		sortedkeys[i] = k
		i++
	}
	sort.Strings(sortedkeys)

	// Format the chart data.
	for _, name := range sortedkeys {
		averageServiceResponseTime := p.PerfStats.ServiceResponseTimes[name]
		baseTimeMillis := float64(float32(p.BasePerfStats.BaseServiceResponseTimes[name]) / float32(1000000))
		serviceResponseTimesBase = append(serviceResponseTimesBase, []byte(strconv.FormatFloat(baseTimeMillis, 'f', 3, 64))...)
		serviceResponseTimesBase = append(serviceResponseTimesBase, []byte(",")...)

		testTimeMillis := float64(float32(averageServiceResponseTime) / float32(1000000))

		serviceResponseTimesTest = append(serviceResponseTimesTest, []byte(strconv.FormatFloat(testTimeMillis, 'f', 3, 64))...)
		serviceResponseTimesTest = append(serviceResponseTimesTest, []byte(",")...)

		if averageServiceResponseTime != 0 {
			responseTimeVariancePercentage := CalcAverageResponseVariancePercentage(averageServiceResponseTime, p.BasePerfStats.BaseServiceResponseTimes[name])
			serviceNames = append(serviceNames, []byte(name+" ("+fmt.Sprintf("%3.2f", responseTimeVariancePercentage)+" %)")...)
		} else {
			serviceNames = append(serviceNames, []byte(name+" ("+fmt.Sprintf("%3.2f", 0.0)+" %)")...)
		}
		serviceNames = append(serviceNames, []byte("','")...)
	}
	serviceResponseTimesBase = append(serviceResponseTimesBase, []byte("],")...)
	serviceResponseTimesTest = append(serviceResponseTimesTest, []byte("]")...)
	serviceNames = append(serviceNames, []byte("']")...)
	p.JSONTimeServiceNames = template.JS(serviceNames)

	serviceResponseTimesBase = append(serviceResponseTimesBase, serviceResponseTimesTest...)
	return template.JS(serviceResponseTimesBase)
}

// GenerateTemplateReport wraps the generateTemplate() function that creates
// the final performance report html.
func GenerateTemplateReport(basePerfstats *BasePerfStats, perfStats *PerfStats, configurationSettings *Config, fs FileSystem, testSuiteName string, testStrategy string) {
	// Check for existence of output dir and create if needed.
	err := os.MkdirAll(configurationSettings.ReportOutputDir, os.ModePerm)
	if err != nil {
		// Non-fatal error. Don't exit.
		log.Errorf("Failed to create path: [%s]. Error: %s\n", configurationSettings.ReportOutputDir, err)
	}

	filename := "PerformanceReport-" + configurationSettings.APIName + "-" + testSuiteName + ".html"

	file, err := fs.Create(configurationSettings.ReportOutputDir + "/" + filename)
	if err != nil {
		log.Errorf("Error creating report file: %v", err)
	}
	if file != nil {
		defer file.Close()
	} else {
		log.Errorf("No file was created, falling back to stdout: %v", err)
		file = os.Stdout
	}
	tf := configurationSettings.ReportTemplateFile
	err = generateTemplate(basePerfstats, perfStats, configurationSettings, file, tf, testStrategy)
	if err != nil {
		log.Errorf("Error generating template report: %v", err)
	}
}

func generateTemplate(bstats *BasePerfStats, pstats *PerfStats, configurationSettings *Config, wr io.Writer, templFile string, testStrategy string) error {
	ps := &perfStatsModel{BasePerfStats: bstats, PerfStats: pstats, Config: configurationSettings, TestStrategy: testStrategy}
	s1 := template.New("main")
	var err error
	s1 = s1.Funcs(template.FuncMap{"memToMB": memoryMB, "formatMem": formatMemory, "jsonMem": jsonMemoryArray, "div": div, "avgVar": CalcAverageResponseVariancePercentage})
	if templFile != "" {
		s1, err = s1.ParseFiles(templFile)
		if err != nil {
			return fmt.Errorf("Error loading template files: %v", err)
		}
		err = s1.ExecuteTemplate(wr, filepath.Base(templFile), ps)
		if err != nil {
			return fmt.Errorf("Error executing template: %v", err)
		}
	} else {
		//use builtin report
		for _, tname := range []string{"report/header.tmpl", "report/content.tmpl", "report/footer.tmpl"} {
			header, err := Asset(tname)
			if err != nil {
				return fmt.Errorf("Error asset not found: %v", err)
			}
			_, err = s1.New(tname).Parse(string(header))
			if err != nil {
				return fmt.Errorf("Error parsing template %v: %v", tname, err)
			}
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
	}
	return nil
}
