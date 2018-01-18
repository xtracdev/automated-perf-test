import {
  browser,
  element,
  by,
  By,
  $,
  $$,
  ExpectedConditions,
  protractor,
  Key
} from "protractor";
import { read } from "fs";
import { Input } from "@angular/core/src/metadata/directives";
import {ToastModule} from 'ng2-toastr/ng2-toastr';
var since = require("jasmine2-custom-message");

import configurationPageObject = require("../pages/configuration-PO");
var configPO: configurationPageObject = new configurationPageObject();

describe("configuration component", () => {
  beforeEach(() => {
    browser.get("/configurations");
  });

  it("should create xml file", () => {
    configPO.addData();
    configPO.submitBtn.click();
    expect(configPO.toastrMessage.getText()).toContain('Success!');
  });


  it("should show submit button is disabled when required data is blank", () => {
    //used to clear default data in this test
    configPO.checkRequired();
    expect(configPO.submitBtn.isEnabled()).toBe(false);
  });
  

  it('should check cancel button clears all fields', () =>{
    configPO.addData();
   
    expect(configPO.applicationName.getAttribute("value")).toEqual("");
    expect(configPO.targetHost.getAttribute("value")).toEqual("");
    expect(configPO.targetPort.getAttribute("value")).toEqual("");
    expect(configPO.memoryEndpoint.getAttribute("value")).toEqual("");
    expect(configPO.testSuite.getAttribute("value")).toEqual("");
    expect(configPO.testCaseDir.getAttribute("value")).toEqual("");
    expect(configPO.testSuiteDir.getAttribute("value")).toEqual("");
    expect(configPO.baseStatsDir.getAttribute("value")).toEqual("");
    expect(configPO.reportsDir.getAttribute("value")).toEqual("");
    expect(configPO.numIterations.getAttribute("value")).toEqual("");
    expect(configPO.memoryVariance.getAttribute("value")).toEqual("");
    expect(configPO.serviceVariance.getAttribute("value")).toEqual("");
    expect(configPO.concurrentUsers.getAttribute("value")).toEqual("");
    expect(configPO.requestDelay.getAttribute("value")).toEqual("");
    expect(configPO.tpsFreq.getAttribute("value")).toEqual("");
    expect(configPO.tpsFreq.getAttribute("value")).toEqual("");
    expect(configPO.rampDelay.getAttribute("value")).toEqual("");
    expect(configPO.rampUsers.getAttribute("value")).toEqual("");

  });

  it("should check that all text box names are correct", () => {
    configPO.addData();
    expect(configPO.labels.get(0).getText()).toContain("Api Name");
    expect(configPO.labels.get(1).getText()).toContain("Num Iterations");
    expect(configPO.labels.get(2).getText()).toContain("Request Delay (ms)");
    expect(configPO.labels.get(3).getText()).toContain("Target Host");
    expect(configPO.labels.get(4).getText()).toContain("Concurrent Users");
    expect(configPO.labels.get(5).getText()).toContain("TPS Frequency (s");
    expect(configPO.labels.get(6).getText()).toContain("Target Port");
    expect(configPO.labels.get(7).getText()).toContain("Memory Variance (%)");
    expect(configPO.labels.get(8).getText()).toContain("Ramp Users");
    expect(configPO.labels.get(9).getText()).toContain("Memory Endpoint");
    expect(configPO.labels.get(10).getText()).toContain("Service Variance (%)");
    expect(configPO.labels.get(11).getText()).toContain("Ramp Delay (s)");
    expect(configPO.labels.get(12).getText()).toContain("Test Suite");
    expect(configPO.labels.get(13).getText()).toContain("Test Case Directory");
    expect(configPO.labels.get(14).getText()).toContain("Test Suites Directory");
    expect(configPO.labels.get(15).getText()).toContain("Base Stats Output Directory");
    expect(configPO.labels.get(16).getText()).toContain("Report Output Director");
  });

  it('should check required warning appears when required input is blank', () => {

      configPO.checkRequired();
      since('(apiName) #{actual} =/= #{expected}').expect(configPO.required.get(0).getText()).toContain('This field is required.');
      since('(targetHost) #{actual} =/= #{expected}').expect(configPO.required.get(1).getText()).toContain('This field is required.');
      since('(targetPort) #{actual} =/= #{expected}').expect(configPO.required.get(2).getText()).toContain('This field is required.');
      since('(numIterations) #{actual} =/= #{expected}').expect(configPO.required.get(3).getText()).toContain('This field is required.');
      since('(concurrentUsers) #{actual} =/= #{expected}').expect(configPO.required.get(4).getText()).toContain('This field is required.');
      since('(memoryVariance) #{actual} =/= #{expected}').expect(configPO.required.get(5).getText()).toContain('This field is required.');
      since('(serviceVariance) #{actual} =/= #{expected}').expect(configPO.required.get(6).getText()).toContain('This field is required.');
      since('(requestDelay) #{actual} =/= #{expected}').expect(configPO.required.get(7).getText()).toContain('This field is required.');
      since('(tpsFreq) #{actual} =/= #{expected}').expect(configPO.required.get(8).getText()).toContain('This field is required.');
      since('(rampUsers) #{actual} =/= #{expected}').expect(configPO.required.get(9).getText()).toContain('This field is required.');
      since('(rampDelay) #{actual} =/= #{expected}').expect(configPO.required.get(10).getText()).toContain('This field is required.');
      since('(testCaseDir) #{actual} =/= #{expected}').expect(configPO.required.get(11).getText()).toContain('This field is required.');
      since('(testSuiteDir) #{actual} =/= #{expected}').expect(configPO.required.get(12).getText()).toContain('This field is required.');
      since('(baseStatsDir) #{actual} =/= #{expected}').expect(configPO.required.get(13).getText()).toContain('This field is required.');
      since('(reportsDir) #{actual} =/= #{expected}').expect(configPO.required.get(14).getText()).toContain('This field is required.');

  });

  it("should show that string cannot be entered into a integer field", () => {
    configPO.checkForStrings();
    expect(configPO.numIterations.getAttribute("value")).toEqual("");
    expect(configPO.memoryVariance.getAttribute("value")).toEqual("");
    expect(configPO.serviceVariance.getAttribute("value")).toEqual("");
    expect(configPO.concurrentUsers.getAttribute("value")).toEqual("");
    expect(configPO.requestDelay.getAttribute("value")).toEqual("");
    expect(configPO.tpsFreq.getAttribute("value")).toEqual("");
    expect(configPO.tpsFreq.getAttribute("value")).toEqual("");
    expect(configPO.rampDelay.getAttribute("value")).toEqual("");
    expect(configPO.rampUsers.getAttribute("value")).toEqual("");
  });

  it('should check that e is not accepted in interger field', () => {
      configPO.checkE();
      since('(numIterations) #{actual} =/= #{expected}').expect(configPO.required.get(0).getText()).toContain('This field is required.');
      since('(concurrentUsers) #{actual} =/= #{expected}').expect(configPO.required.get(1).getText()).toContain('This field is required.');
      since('(memoryVariance) #{actual} =/= #{expected}').expect(configPO.required.get(2).getText()).toContain('This field is required.');
      since('(serviceVariance) #{actual} =/= #{expected}').expect(configPO.required.get(3).getText()).toContain('This field is required.');
      since('(requestDelay) #{actual} =/= #{expected}').expect(configPO.required.get(4).getText()).toContain('This field is required.');
      since('(tpsFreq) #{actual} =/= #{expected}').expect(configPO.required.get(5).getText()).toContain('This field is required.');
      since('(rampUsers) #{actual} =/= #{expected}').expect(configPO.required.get(6).getText()).toContain('This field is required.');
      since('(rampDelay) #{actual} =/= #{expected}').expect(configPO.required.get(7).getText()).toContain('This field is required.');

  });

  it('should check that warning appears if negative number is enter to integer field', () => {
      configPO.checkNegativeValues();   
      since('(numIterations) #{actual} =/= #{expected}').expect(configPO.required.get(0).getText()).toContain('Must be 1 or more');
      since('(requestDelay) #{actual} =/= #{expected}').expect(configPO.required.get(1).getText()).toContain('Must be 1 or more');
      since('(concurrentUsers) #{actual} =/= #{expected}').expect(configPO.required.get(2).getText()).toContain('Must be 1 or more');
      since('(tpsFreq) #{actual} =/= #{expected}').expect(configPO.required.get(3).getText()).toContain('Must be 1 or more');
      since('(memoryVariance) #{actual} =/= #{expected}').expect(configPO.required.get(4).getText()).toContain('Must be 0 or more');
      since('(rampUsers) #{actual} =/= #{expected}').expect(configPO.required.get(5).getText()).toContain('Must be 1 or more');
      since('(serviceVariance) #{actual} =/= #{expected}').expect(configPO.required.get(6).getText()).toContain('Must be 0 or more');
      since('(rampDelay) #{actual} =/= #{expected}').expect(configPO.required.get(7).getText()).toContain('Must be 0 or more');

  });

  it('should check that warning appears if value exceeds maximum', () => {
      configPO.memoryVariance.sendKeys(101);
      configPO.serviceVariance.sendKeys(101);
      since('(memoryVariance) #{actual} =/= #{expected}').expect(configPO.required.get(0).getText()).toEqual('Must be 100 or less');
      since('(serviceVariance) #{actual} =/= #{expected}').expect(configPO.required.get(1).getText()).toEqual('Must be 100 or less');

  });
});
