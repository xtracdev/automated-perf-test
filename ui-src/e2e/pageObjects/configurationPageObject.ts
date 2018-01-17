import { browser, element, by, By, $, $$, ExpectedConditions } from "protractor"
import { read } from "fs";
var path = require('path');
const configFileLocation = "../../../config/"
class configurationPageObject {

    configFilePath = element(by.name('configFilePath'));
    applicationName = element(by.name('apiName'));
    targetHost = element(by.name('targetHost'));
    targetPort = element(by.name('targetPort'));
    memoryEndpoint = element(by.name('memoryEndpoint'));
    submitBtn = element(by.tagName('submit-widget'));
    numIterations = element(by.name('numIterations'));
    concurrentUsers = element(by.name('concurrentUsers'));
    memoryVariance = element(by.name('allowablePeakMemoryVariance'));
    serviceVariance = element(by.name('allowableServiceResponseTimeVariance'));
    testSuite = element(by.name('testSuite'));
    requestDelay = element(by.name('requestDelay'));
    tpsFreq = element(by.name('TPSFreq'));
    rampUsers = element(by.name('rampUsers'));
    rampDelay = element(by.name('rampDelay'));
    testCaseDir = element(by.name('testCaseDir'));
    testSuiteDir = element(by.name('testSuiteDir'));
    baseStatsDir = element(by.name('baseStatsOutputDir'));
    reportsDir = element(by.name('reportOutputDir'));

    labels = $('json-schema-form').$$('label');
    required = $$('p');

    absolutePath = path.resolve(__dirname, configFileLocation);




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
    setRampDelay() {
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

    addData() {
        // this.setConfigPath();
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

    checkRequired() {
        this.applicationName.sendKeys("x"); 
        this.targetHost.sendKeys("x");
        this.targetPort.sendKeys(1);
        this.numIterations.sendKeys(1);
        this.concurrentUsers.sendKeys(1);
        this.memoryVariance.sendKeys(1);
        this.serviceVariance.sendKeys(1);
        this.requestDelay.sendKeys(1);
        this.tpsFreq.sendKeys(1);
        this.rampUsers.sendKeys(1);
        this.rampDelay.sendKeys(1);
        this.testCaseDir.sendKeys("x");
        this.testSuiteDir.sendKeys("x");
        this.baseStatsDir.sendKeys("x");
        this.reportsDir.sendKeys("x");

    }

}
export = configurationPageObject;