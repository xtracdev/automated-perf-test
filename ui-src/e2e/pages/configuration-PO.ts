import {
  browser,
  element,
  by,
  By,
  $,
  $$,
  ExpectedConditions,
  Key
} from "protractor";
import { read } from "fs";
const path = require("path");
const configFileLocation = "../../../config/";

class ConfigurationPageObject {
  configFilePath = element(by.id("config-file-path"));
  btnUpdate = element(by.id("update-config-file-btn"));
  xmlFileName = element(by.id("xml-file-name"));
  applicationName = element(by.name("apiName"));
  targetHost = element(by.name("targetHost"));
  targetPort = element(by.name("targetPort"));
  memoryEndpoint = element(by.name("memoryEndpoint"));
  submitBtn = element(by.className("btn"));
  cancelBtn = element(by.id("cancel-btn"));

  numIterations = element(by.name("numIterations"));
  concurrentUsers = element(by.name("concurrentUsers"));
  memoryVariance = element(by.name("allowablePeakMemoryVariance"));
  serviceVariance = element(by.name("allowableServiceResponseTimeVariance"));
  testSuite = element(by.name("testSuite")).$$("option");
  requestDelay = element(by.name("requestDelay"));
  tpsFreq = element(by.name("TPSFreq"));
  rampUsers = element(by.name("rampUsers"));
  rampDelay = element(by.name("rampDelay"));
  testCaseDir = element(by.name("testCaseDir"));
  testSuiteDir = element(by.name("testSuiteDir"));
  baseStatsDir = element(by.name("baseStatsOutputDir"));
  reportsDir = element(by.name("reportOutputDir"));

  toastrMessage = element(by.className("toast-message"));
  labels = $("json-schema-form").$$("label");
  requiredFields = $$("p");

  absolutePath = path.resolve(__dirname, configFileLocation);

  setConfigPath() {
    return this.configFilePath.sendKeys(this.absolutePath);
  }

  setApplicationName() {
    return this.applicationName.sendKeys("Xtrac API");
  }
  setTargetHost() {
    return this.targetHost.sendKeys("localhost");
  }
  setTargetPort() {
    return this.targetPort.sendKeys("9191");
  }
  setMemoryEndpoint() {
    return this.memoryEndpoint.sendKeys("/alt/debug/vars");
  }
  setConcurrentUsers() {
    return this.concurrentUsers.sendKeys(10);
  }
  setNumberIterations() {
    return this.numIterations.sendKeys(1000);
  }

  setTestSuite() {
    return this.testSuite.get(0).click();
  }

  setRequestDelay() {
    return this.requestDelay.sendKeys(300);
  }

  setRampUsers() {
    return this.rampUsers.sendKeys(2);
  }
  setRampDelay() {
    return this.rampDelay.sendKeys(60);
  }
  setTPSfreq() {
    return this.tpsFreq.sendKeys(60);
  }
  setTestCaseDir() {
    return this.testCaseDir.sendKeys("./definitions/testCases");
  }
  setTestSuiteDir() {
    return this.testSuiteDir.sendKeys("./definitions/testSuites");
  }
  setBaseStatsDir() {
    return this.baseStatsDir.sendKeys("./envStats");
  }
  setReportDir() {
    return this.reportsDir.sendKeys("./report");
  }

  setConfigData() {
    this.setConfigPath();
    this.setApplicationName();
    this.setTargetHost();
    this.setTargetPort();
    this.setMemoryEndpoint();
    this.setConcurrentUsers();
    this.setNumberIterations();
    // used to activate dropdown
    this.testSuite.click();
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

  checkRequiredFields() {
    this.applicationName.sendKeys("x");
    this.targetHost.sendKeys("x");
    this.targetPort.sendKeys(1);
    this.numIterations.sendKeys(1);
    this.concurrentUsers.sendKeys(1);
    this.requestDelay.sendKeys(1);
    this.tpsFreq.sendKeys(1);
    this.rampUsers.sendKeys(1);
    this.rampDelay.sendKeys(1);
    this.testCaseDir.sendKeys("x");
    this.testSuiteDir.sendKeys("x");
    this.baseStatsDir.sendKeys("x");
    this.reportsDir.sendKeys("x");

    this.applicationName.sendKeys(Key.BACK_SPACE);
    this.numIterations.sendKeys(Key.BACK_SPACE);
    this.concurrentUsers.sendKeys(Key.BACK_SPACE);
    this.targetHost.sendKeys(Key.BACK_SPACE);
    this.targetPort.sendKeys(Key.BACK_SPACE);
    // Clear default data in these fields
    this.memoryVariance.sendKeys(Key.BACK_SPACE);
    this.memoryVariance.sendKeys(Key.BACK_SPACE);
    this.serviceVariance.sendKeys(Key.BACK_SPACE);
    this.serviceVariance.sendKeys(Key.BACK_SPACE);

    this.requestDelay.sendKeys(Key.BACK_SPACE);
    this.tpsFreq.sendKeys(Key.BACK_SPACE);
    this.rampUsers.sendKeys(Key.BACK_SPACE);
    this.rampDelay.sendKeys(Key.BACK_SPACE);
    this.testCaseDir.sendKeys(Key.BACK_SPACE);
    this.testSuiteDir.sendKeys(Key.BACK_SPACE);
    this.baseStatsDir.sendKeys(Key.BACK_SPACE);
    this.reportsDir.sendKeys(Key.BACK_SPACE);
  }
  checkEisNotAccepted() {
    this.numIterations.sendKeys("e");
    this.concurrentUsers.sendKeys("e");
    this.memoryVariance.sendKeys("e");
    this.serviceVariance.sendKeys("e");
    this.requestDelay.sendKeys("e");
    this.tpsFreq.sendKeys("e");
    this.rampUsers.sendKeys("e");
    this.rampDelay.sendKeys("e");
  }
  checkNegativeValues() {
    this.numIterations.sendKeys(-1);
    this.concurrentUsers.sendKeys(-1);
    // Clear default data in these fields
    this.memoryVariance.sendKeys(Key.BACK_SPACE);
    this.memoryVariance.sendKeys(Key.BACK_SPACE);
    this.serviceVariance.sendKeys(Key.BACK_SPACE);
    this.serviceVariance.sendKeys(Key.BACK_SPACE);
    this.memoryVariance.sendKeys(-1);
    this.serviceVariance.sendKeys(-1);
    this.requestDelay.sendKeys(-1);
    this.tpsFreq.sendKeys(-1);
    this.rampUsers.sendKeys(-1);
    this.rampDelay.sendKeys(-1);
  }

  checkForStrings() {
    this.targetPort.sendKeys("x");
    this.numIterations.sendKeys("x");
    this.concurrentUsers.sendKeys("x");
    this.memoryVariance.sendKeys(Key.BACK_SPACE);
    this.memoryVariance.sendKeys(Key.BACK_SPACE);
    this.serviceVariance.sendKeys(Key.BACK_SPACE);
    this.serviceVariance.sendKeys(Key.BACK_SPACE);
    this.requestDelay.sendKeys("x");
    this.tpsFreq.sendKeys("x");
    this.rampUsers.sendKeys("x");
    this.rampDelay.sendKeys("x");
  }
}
export = ConfigurationPageObject;
