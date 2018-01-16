import { browser, element, by, By, $, $$, ExpectedConditions } from "protractor"
import { read } from "fs";
var path = require('path');
class configurationPageObject {

    configFilePath = element(by.name('configFilePath'));
    applicationName = element(by.name('applicationName'));s
    targetHost = element(by.name('targetHost'));
    targetPort = element(by.name('targetPort'));
    memoryEndpoint = element(by.name('memoryEndpoint'));
    submitBtn = $('submit-widget');
    numIterations = element(by.name('numIterations'));
    concurrentUsers = element(by.name('concurrentUsers'));
    memoryVariance = element(by.name('memoryVariance'));
    serviceVariance = element(by.name('serviceVariance'));
    testSuite = element(by.name('testSuite'));
    requestDelay = element(by.name('requestDelay'));
    tpsFreq =element(by.name('tpsFrequency'));
    rampUsers = element(by.name('rampUsers'));
    rampDelay = element(by.name('rampDelay'));
    testCaseDir =element(by.name('testCaseDirectory'));
    testSuiteDir =element(by.name('testSuiteDirectory'));
    baseStatsDir = element(by.name('baseStatsDirectory'));
    reportsDir = element(by.name('reportDirectory'));
    
    label = $('json-schema-form').$$('label');
    required = element(by.className('help-block'));
   
    configFileLocation = "../../../config/"
    absolutePath = path.resolve(__dirname, this.configFileLocation);
    
   
    
    
    setConfigPath() {
         return this.configFilePath.sendKeys(this.absolutePath)

    }

    setApplicationName() {
        return this.applicationName.sendKeys("Xtrac API")
    }
    setTargetHost() {
        return this.targetHost.sendKeys("localhost")
    }
    setTargetPort() {
        return this.targetPort.sendKeys("9191")
    }
    setMemoryEndpoint() {
        return this.memoryEndpoint.sendKeys("/alt/debug/vars")

    }
    setConcurrentUsers() {
        return this.concurrentUsers.sendKeys(10)
    }
    setNumberIterations() {
        return this.numIterations.sendKeys(1000)
    }

    setServiceVariance() {
        return this.serviceVariance.sendKeys(15)
    }
    setMemoryVariance() {
        return this.memoryVariance.sendKeys(15)
    }
    setTestSuite() {
        return this.testSuite.sendKeys("suiteFileName.xml")
    }

    setRequestDelay() {
        return this.requestDelay.sendKeys(300)
    }

    setRampUsers() {
        return this.rampUsers.sendKeys(2)
    }
    setRampDelay(){
        return this.rampDelay.sendKeys(60)
    }
    setTPSfreq() {
        return this.tpsFreq.sendKeys(60)
    }
    setTestCaseDir() {
        return this.testCaseDir.sendKeys("./definitions/testCases")

    }
    setTestSuiteDir() {
        return this.testSuiteDir.sendKeys("./definitions/testSuites")
    }
    setBaseStatsDir() {
        return this.baseStatsDir.sendKeys("./envStats")
    }
    setReportDir() {
        return this.reportsDir.sendKeys("./report")
    }

    addData(){
        this.setConfigPath();
        this.setApplicationName();
        this.setTargetHost();
        this.setTargetPort();
        this.setMemoryEndpoint();
        this.setConcurrentUsers();
        this.setNumberIterations();
        this.setMemoryVariance();
        this.setServiceVariance();
        this.setTestSuite();
        this.setRequestDelay();
        this.setTPSfreq();
        this.setRampUsers();
        this.setRampDelay();
        this.setTestCaseDir();
        this.setTestSuiteDir();
        this.setBaseStatsDir();
        this.setReportDir();
    }

}
export = configurationPageObject;