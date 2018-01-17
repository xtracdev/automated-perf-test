import { browser, element, by, By, $, $$, ExpectedConditions, protractor, Key } from "protractor"
import { read } from "fs";
import { Input } from "@angular/core/src/metadata/directives";
var since = require('jasmine2-custom-message');


import configurationPageObject = require('../pageObjects/configurationPageObject');
var configPO: configurationPageObject = new configurationPageObject();

describe('configuration component', () => {
    beforeEach(() => {
        browser.get("/configurations", 1000);

    });


    it('should check that all text box names are correct', () => {

        configPO.addData();
        configPO.submitBtn.click();



    });

    it('should check that all text box names are correct', () => {

        configPO.addData();
        expect(configPO.labels.get(0).getText()).toContain('Api Name')
        expect(configPO.labels.get(1).getText()).toContain('Target Host')
        expect(configPO.labels.get(2).getText()).toContain('Target Port')
        expect(configPO.labels.get(3).getText()).toContain('Memory Endpoint')
        expect(configPO.labels.get(4).getText()).toContain('Num Iterations')
        expect(configPO.labels.get(5).getText()).toContain('Concurrent Users')
        expect(configPO.labels.get(6).getText()).toContain('Allowable Peak Memory Variance')
        expect(configPO.labels.get(7).getText()).toContain('Allowable Service Response Time Variance')
        expect(configPO.labels.get(8).getText()).toContain('Test Suite')
        expect(configPO.labels.get(9).getText()).toContain('Request Delay')
        expect(configPO.labels.get(10).getText()).toContain('TPSFreq')
        expect(configPO.labels.get(11).getText()).toContain('Ramp Users')
        expect(configPO.labels.get(12).getText()).toContain('Ramp Delay')
        expect(configPO.labels.get(13).getText()).toContain('Test Case Dir')
        expect(configPO.labels.get(14).getText()).toContain('Test Suite Dir')
        expect(configPO.labels.get(15).getText()).toContain('Base Stats Output Dir')
        expect(configPO.labels.get(16).getText()).toContain('Report Output Dir')



    });


    it('should check required warning appears when required input is blank', () => {

        configPO.checkRequired();
        configPO.applicationName.sendKeys(Key.BACK_SPACE);
        configPO.targetHost.sendKeys(Key.BACK_SPACE);
        configPO.targetPort.sendKeys(Key.BACK_SPACE);
        configPO.numIterations.sendKeys(Key.BACK_SPACE);
        configPO.concurrentUsers.sendKeys(Key.BACK_SPACE);
        configPO.memoryVariance.sendKeys(Key.BACK_SPACE);
        configPO.serviceVariance.sendKeys(Key.BACK_SPACE);
        configPO.requestDelay.sendKeys(Key.BACK_SPACE);
        configPO.tpsFreq.sendKeys(Key.BACK_SPACE);
        configPO.rampUsers.sendKeys(Key.BACK_SPACE);
        configPO.rampDelay.sendKeys(Key.BACK_SPACE);
        configPO.testCaseDir.sendKeys(Key.BACK_SPACE);
        configPO.testSuiteDir.sendKeys(Key.BACK_SPACE);
        configPO.baseStatsDir.sendKeys(Key.BACK_SPACE);
        configPO.backSpaceByField('reportsDir');

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




    it('should show that string cannot be entered into a integer field', () => {
        configPO.targetPort.sendKeys("x");
        configPO.numIterations.sendKeys("x");
        configPO.concurrentUsers.sendKeys("x");
        configPO.memoryVariance.sendKeys("x");
        configPO.serviceVariance.sendKeys("x");
        configPO.requestDelay.sendKeys("x");
        configPO.tpsFreq.sendKeys("x");
        configPO.rampUsers.sendKeys("x");
        configPO.rampDelay.sendKeys("x");

        expect(configPO.numIterations.getAttribute('value')).toEqual('');
        expect(configPO.memoryVariance.getAttribute('value')).toEqual('');
        expect(configPO.serviceVariance.getAttribute('value')).toEqual('');
        expect(configPO.concurrentUsers.getAttribute('value')).toEqual('');
        expect(configPO.requestDelay.getAttribute('value')).toEqual('');
        expect(configPO.tpsFreq.getAttribute('value')).toEqual('');
        expect(configPO.tpsFreq.getAttribute('value')).toEqual('');
        expect(configPO.rampDelay.getAttribute('value')).toEqual('');
        expect(configPO.rampUsers.getAttribute('value')).toEqual('');

    });

    it('should check that e is not accepted in interger field', () => {
        configPO.targetPort.sendKeys("e");
        configPO.numIterations.sendKeys("e");
        configPO.concurrentUsers.sendKeys("e");
        configPO.memoryVariance.sendKeys("e");
        configPO.serviceVariance.sendKeys("e");
        configPO.requestDelay.sendKeys("e");
        configPO.tpsFreq.sendKeys("e");
        configPO.rampUsers.sendKeys("e");
        configPO.rampDelay.sendKeys("e");


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
        configPO.numIterations.sendKeys(-1);
        configPO.concurrentUsers.sendKeys(-1);
        configPO.memoryVariance.sendKeys(-1);
        configPO.serviceVariance.sendKeys(-1);
        configPO.requestDelay.sendKeys(-1);
        configPO.tpsFreq.sendKeys(-1);
        configPO.rampUsers.sendKeys(-1);
        configPO.rampDelay.sendKeys(-1);


        since('(numIterations) #{actual} =/= #{expected}').expect(configPO.required.get(0).getText()).toContain('Must be 0 or more');
        since('(concurrentUsers) #{actual} =/= #{expected}').expect(configPO.required.get(1).getText()).toContain('Must be 0 or more');
        since('(memoryVariance) #{actual} =/= #{expected}').expect(configPO.required.get(2).getText()).toContain('Must be 0 or more');
        since('(serviceVariance) #{actual} =/= #{expected}').expect(configPO.required.get(3).getText()).toContain('Must be 0 or more');
        since('(requestDelay) #{actual} =/= #{expected}').expect(configPO.required.get(4).getText()).toContain('Must be 0 or more');
        since('(tpsFreq) #{actual} =/= #{expected}').expect(configPO.required.get(5).getText()).toContain('Must be 0 or more');
        since('(rampUsers) #{actual} =/= #{expected}').expect(configPO.required.get(6).getText()).toContain('Must be 0 or more');
        since('(rampDelay) #{actual} =/= #{expected}').expect(configPO.required.get(7).getText()).toContain('Must be 0 or more');




    });


    it('should check that warning appears if value exceeds maximum', () => {
        configPO.memoryVariance.sendKeys(101);
        configPO.serviceVariance.sendKeys(101);
        since('(memoryVariance) #{actual} =/= #{expected}').expect(configPO.required.get(0).getText()).toEqual('Must be 100 or less');
        since('(serviceVariance) #{actual} =/= #{expected}').expect(configPO.required.get(1).getText()).toEqual('Must be 100 or less');



    });
});





