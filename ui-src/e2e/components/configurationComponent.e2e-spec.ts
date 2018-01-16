import { browser, element, by, By, $, $$, ExpectedConditions, protractor, Key } from "protractor"
import { read } from "fs";
import { Input } from "@angular/core/src/metadata/directives";

import configurationPageObject = require('../pageObjects/configurationPageObject');
var configPo: configurationPageObject = new configurationPageObject();
var path = require('path');







describe('configuration component -e2e testing', () => {
    beforeEach(() => {
        browser.get("http://localhost:49155/configurations");

    });


    //this test is failing cannot click button 
    // it('should succesfully add all values inputted', () => {

    //     configPo.addData();

    //     browser.wait(protractor.ExpectedConditions.alertIsPresent(), 5000);
    //     configPo.submitBtn.click();

    //     expect(alert.getText()).toEqual("hi Colm");
    //     alert.accept();
    //     expect(configPo.configFilePath.getAttribute('value')).toEqual('C:/Users/A586754/go/src/github.com/xtracdev/automated-perf-test/config');


    // });

    // it('should check that all labels are correct', () => {

    //     configPo.addData();
    //     expect(configPo.label.get(0).getText()).toEqual('Config File Path' + ' *')
    //     expect(configPo.label.get(1).getText()).toEqual('Application Name' + ' *')
    //     expect(configPo.label.get(2).getText()).toEqual('Target Host' + ' *')
    //     expect(configPo.label.get(3).getText()).toEqual('Target Port' + ' *')
    //     expect(configPo.label.get(4).getText()).toEqual('Num Iterations' + ' *')
    //     expect(configPo.label.get(5).getText()).toEqual('Memory Variance' + ' *')
    //     expect(configPo.label.get(6).getText()).toEqual('Service Variance' + ' *')
    //     expect(configPo.label.get(7).getText()).toEqual('Test Case Directory' + ' *')
    //     expect(configPo.label.get(8).getText()).toEqual('Test Suite Directory' + ' *')
    //     expect(configPo.label.get(9).getText()).toEqual('Base Stats Directory' + ' *')
    //     expect(configPo.label.get(10).getText()).toEqual('Report Directory' + ' *')
    //     expect(configPo.label.get(11).getText()).toEqual('Concurrent Users' + ' *')
    //     expect(configPo.label.get(12).getText()).toEqual('Test Suite' + ' *')
    //     expect(configPo.label.get(13).getText()).toEqual('Memory Endpoint')
    //     expect(configPo.label.get(14).getText()).toEqual('Request Delay' + ' *')
    //     expect(configPo.label.get(15).getText()).toEqual('Tps Frequency' + ' *')
    //     expect(configPo.label.get(16).getText()).toEqual('Ramp Users' + ' *')
    //     expect(configPo.label.get(17).getText()).toEqual('Ramp Delay' + ' *')



    // });


    // it('should check required label appears when required input is blank(configFilePath)', () => {

    //     configPo.configFilePath.sendKeys("x");
    //     configPo.configFilePath.sendKeys(Key.BACK_SPACE);
    //     expect(configPo.required.getText()).toContain('This field is required.');



    // });

    // it('should check required label appears when required input is blank(appliceationName)', () => {

    //     configPo.applicationName.sendKeys('x')
    //     configPo.applicationName.sendKeys(Key.BACK_SPACE);
    //     expect(configPo.required.getText()).toContain('This field is required.');



    // });

    // it('should check required label appears when required input is blank(targetHost)', () => {

    //     configPo.targetHost.sendKeys("x");
    //     configPo.targetHost.sendKeys(Key.BACK_SPACE);
    //     expect(configPo.required.getText()).toContain('This field is required.');



    // });

    // it('should check required label appears when required input is blank(targetPort)', () => {

    //     configPo.targetPort.sendKeys("x");
    //     configPo.targetPort.sendKeys(Key.BACK_SPACE);
    //     expect(configPo.required.getText()).toContain('This field is required.');



    // });
    // it('should check required label appears when required input is blank(numIterations)', () => {

    //     configPo.numIterations.sendKeys(1);
    //     configPo.numIterations.sendKeys(Key.BACK_SPACE);
    //     expect(configPo.required.getText()).toContain('This field is required.');



    // });

    // it('should check required label appears when required input is blank(memoryVariance)', () => {

    //     configPo.memoryVariance.sendKeys(1);
    //     configPo.memoryVariance.sendKeys(Key.BACK_SPACE);
    //     expect(configPo.required.getText()).toContain('This field is required.');



    // });


    // it('should check required label appears when required input is blank(serviceVariance)', () => {

    //     configPo.serviceVariance.sendKeys(1);
    //     configPo.serviceVariance.sendKeys(Key.BACK_SPACE);
    //     expect(configPo.required.getText()).toContain('This field is required.');



    // });


    // it('should check required label appears when required input is blank(serviceVariance)', () => {

    //     configPo.serviceVariance.sendKeys(1);
    //     configPo.serviceVariance.sendKeys(Key.BACK_SPACE);
    //     expect(configPo.required.getText()).toContain('This field is required.');



    // });

    // it('should check required label appears when required input is blank(testCaseDirectory)', () => {

    //     configPo.testCaseDir.sendKeys(1);
    //     configPo.testCaseDir.sendKeys(Key.BACK_SPACE);
    //     expect(configPo.required.getText()).toContain('This field is required.');



    // });

    // it('should check required label appears when required input is blank(testSuiteDirectory)', () => {

    //     configPo.testSuiteDir.sendKeys(1);
    //     configPo.testSuiteDir.sendKeys(Key.BACK_SPACE);
    //     expect(configPo.required.getText()).toContain('This field is required.');



    // });


    // it('should check required label appears when required input is blank(baseStatsDirectory)', () => {

    //     configPo.baseStatsDir.sendKeys(1);
    //     configPo.baseStatsDir.sendKeys(Key.BACK_SPACE);
    //     expect(configPo.required.getText()).toContain('This field is required.');



    // });



    // it('should check required label appears when required input is blank(reportDirectory)', () => {

    //     configPo.reportsDir.sendKeys(1);
    //     configPo.reportsDir.sendKeys(Key.BACK_SPACE);
    //     expect(configPo.required.getText()).toContain('This field is required.');



    // });

    // it('should check required label appears when required input is blank(concurrentUsers)', () => {

    //     configPo.concurrentUsers.sendKeys(1);
    //     configPo.concurrentUsers.sendKeys(Key.BACK_SPACE);
    //     expect(configPo.required.getText()).toContain('This field is required.');



    // });



    // it('should check required label appears when required input is blank(testSuite)', () => {

    //     configPo.testSuite.sendKeys(1);
    //     configPo.testSuite.sendKeys(Key.BACK_SPACE);
    //     expect(configPo.required.getText()).toContain('This field is required.');



    // });



    // it('should check required label appears when required input is blank(requestDelay)', () => {

    //     configPo.requestDelay.sendKeys(1);
    //     configPo.requestDelay.sendKeys(Key.BACK_SPACE);
    //     expect(configPo.required.getText()).toContain('This field is required.');



    // });


    // it('should check required label appears when required input is blank(tpsFrequency)', () => {

    //     configPo.tpsFreq.sendKeys(1);
    //     configPo.tpsFreq.sendKeys(Key.BACK_SPACE);
    //     expect(configPo.required.getText()).toContain('This field is required.');



    // });


    // it('should check required label appears when required input is blank(rampUsers)', () => {

    //     configPo.rampUsers.sendKeys(1);
    //     configPo.rampUsers.sendKeys(Key.BACK_SPACE);
    //     expect(configPo.required.getText()).toContain('This field is required.');



    // });



    // it('should check required label appears when required input is blank(rampDelay)', () => {

    //     configPo.rampDelay.sendKeys(1);
    //     configPo.rampDelay.sendKeys(Key.BACK_SPACE);
    //     expect(configPo.required.getText()).toContain('This field is required.');



    // });


    // it('should show that string cannot be entered into a integer field(numIterations)', () => {

    //     configPo.numIterations.sendKeys("x");
    //     expect(configPo.numIterations.getAttribute('value')).toEqual('')



    // });


    // it('should show that string cannot be entered into a integer field(memoryVariance)', () => {

    //     configPo.memoryVariance.sendKeys("x");
    //     expect(configPo.memoryVariance.getAttribute('value')).toEqual('')



    // });

    // it('should show that string cannot be entered into a integer field(serviceVariance)', () => {

    //     configPo.serviceVariance.sendKeys("x");
    //     expect(configPo.serviceVariance.getAttribute('value')).toEqual('')



    // });


    // it('should show that string cannot be entered into a integer field(concurrentUsers)', () => {

    //     configPo.concurrentUsers.sendKeys("x");
    //     expect(configPo.concurrentUsers.getAttribute('value')).toEqual('')



    // });

    // it('should show that string cannot be entered into a integer field(requestDelay)', () => {

    //     configPo.requestDelay.sendKeys("x");
    //     expect(configPo.requestDelay.getAttribute('value')).toEqual('')



    // });


    // it('should show that string cannot be entered into a integer field(requestDelay)', () => {

    //     configPo.requestDelay.sendKeys("x");
    //     expect(configPo.requestDelay.getAttribute('value')).toEqual('')



    // });

    // it('should show that string cannot be entered into a integer field(tpsFreq)', () => {

    //     configPo.tpsFreq.sendKeys("x");
    //     expect(configPo.tpsFreq.getAttribute('value')).toEqual('')



    // });


    // it('should show that string cannot be entered into a integer field(rampUsers)', () => {

    //     configPo.tpsFreq.sendKeys("x");
    //     expect(configPo.tpsFreq.getAttribute('value')).toEqual('')



    // });


    // it('should show that string cannot be entered into a integer field(rampDelay)', () => {

    //     configPo.rampDelay.sendKeys("x");
    //     expect(configPo.rampDelay.getAttribute('value')).toEqual('')



    // });


    it('should check that e is not accepted in interger field(numIterations)', () => {

        configPo.numIterations.sendKeys('e');
        expect(configPo.required.getText()).toContain('This field is required.');



    });


    it('should check that e is not accepted in interger field(numIterations)', () => {

        configPo.numIterations.sendKeys('e');
        expect(configPo.required.getText()).toContain('This field is required.');



    });







});

