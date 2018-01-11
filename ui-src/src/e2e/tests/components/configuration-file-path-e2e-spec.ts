import { browser, element, by, By, $, $$, ExpectedConditions } from "protractor"

import configurationPageObject = require('../pageObjects/config');
var configPo: configurationPageObject = new configurationPageObject();
var path = require('path');

describe('configuration component -e2e testing', () => {


  beforeEach(() => {
    browser.get('http://localhost:9191/configs');
  });





  it('should create xml file through UI', () => {
    configPo.addData();

    //Click Save
    configPo.saveBtn.click();

    browser.refresh();
    configPo.configPathDir.sendKeys(this.absolutePath + "Xtrac API.xml");

    //Test expected values
    expect(configPo.applicationName.getText()).toEqual('Xtrac API');
    expect(configPo.numIterations.getText()).toEqual(1000);
    expect(configPo.configPathDir.getText()).toEqual(this.absolutePath + "Xtrac API.xml");



  });

  it('should update xml file that already exists through UI', () => {
    //Populate Application Properties

    configPo.addData();
    configPo.configPathDir.sendKeys(this.absolutePath + "config.xml");
    //Click Save
    configPo.saveBtn.click();

    browser.refresh();
    configPo.configPathDir.sendKeys(this.absolutePath + "config.xml");

    //Test expected values
    expect(configPo.applicationName.getText()).toEqual('Xtrac API');
    expect(configPo.numIterations.getText()).toEqual(1100);

    expect(configPo.configPathDir.getText()).toEqual(this.absolutePath + "config.xml");



  });

  it('should not enable save button if all req fields are not filled (applicationName) ', () => {
    configPo.addData();
    configPo.applicationName.sendKeys(null);


    //Test
    expect(configPo.saveBtn.isEnabled()).toBe(false);


  });

  it('should not enable save button if all req fields are not filled (targetPort) ', () => {

    configPo.addData();
    configPo.targetPort.sendKeys(null);


    //Test
    expect(configPo.saveBtn.isEnabled()).toBe(false);


  });
  it('should not enable save button if all req fields are not filled (targetHost) ', () => {
    configPo.addData();
    configPo.targetHost.sendKeys(null);

    //Test
    expect(configPo.saveBtn.isEnabled()).toBe(false);


  });


  it('should not enable save button if all req fields are not filled (configPathDir) ', () => {

    configPo.addData();
    configPo.configPathDir.sendKeys(null);


    //Test
    expect(configPo.saveBtn.isEnabled()).toBe(false);





  });

  it('should not enable save button if all req fields are not filled (numIterations) ', () => {
    configPo.addData();
    configPo.numIterations.sendKeys(null);


    //Test
    expect(configPo.saveBtn.isEnabled()).toBe(false);


  });

  it('should not enable save button if all req fields are not filled (concurrentUsers) ', () => {
    configPo.addData();
    configPo.numIterations.sendKeys(1000);


    //Test
    expect(configPo.saveBtn.isEnabled()).toBe(false);


  });

  it('should not enable save button if all req fields are not filled (memoryVariance) ', () => {
    configPo.addData();
    configPo.memoryVariance.sendKeys(null);


    //Test
    expect(configPo.saveBtn.isEnabled()).toBe(false);


  });
  it('should not enable save button if all req fields are not filled (serviceVariance) ', () => {
    configPo.addData();
    configPo.serviceVariance.clear;


    //Test
    expect(configPo.saveBtn.isEnabled()).toBe(false);


  });

  it('should not enable save button if all req fields are not filled (testSuite) ', () => {
    configPo.addData();
    configPo.testSuite.sendKeys(null);


    //Test
    expect(configPo.saveBtn.isEnabled()).toBe(false);


  });

  it('should not enable save button if all req fields are not filled (requestDelay) ', () => {
    configPo.addData();
    configPo.requestDelay.sendKeys(null);


    //Test
    expect(configPo.saveBtn.isEnabled()).toBe(false);


  });



  it('should not enable save button if all req fields are not filled (tpsFreq) ', () => {
    configPo.addData();
    configPo.tpsFreq.sendKeys(null);


    //Test
    expect(configPo.saveBtn.isEnabled()).toBe(false);


  });


  it('should not enable save button if all req fields are not filled (tpsFreq) ', () => {
    configPo.addData();
    configPo.tpsFreq.sendKeys(null);


    //Test
    expect(configPo.saveBtn.isEnabled()).toBe(false);


  });


  it('should not enable save button if all req fields are not filled (rampUsers) ', () => {
    configPo.addData();
    configPo.rampUsers.sendKeys(null);


    //Test
    expect(configPo.saveBtn.isEnabled()).toBe(false);


  });


  it('should not enable save button if all req fields are not filled (rampDelay) ', () => {
    configPo.addData();
    configPo.rampDelay.sendKeys(null);

    //Test
    expect(configPo.saveBtn.isEnabled()).toBe(false);


  });



  it('should not enable save button if all req fields are not filled (testCaseDir) ', () => {
    configPo.addData();

    //Populate Output paths
    configPo.testCaseDir.sendKeys(null);


    //Test
    expect(configPo.saveBtn.isEnabled()).toBe(false);


  });


  it('should not enable save button if all req fields are not filled (testSuiteDirTxt) ', () => {
    configPo.addData();
    configPo.testSuiteDir.sendKeys(null);


    //Test
    expect(configPo.saveBtn.isEnabled()).toBe(false);


  });



  it('should not enable save button if all req fields are not filled (baseStatsDir) ', () => {
    configPo.addData();
    configPo.baseStatsDir.sendKeys(null);


    //Test
    expect(configPo.saveBtn.isEnabled()).toBe(false);


  });



  it('should not enable save button if all req fields are not filled (reportsDirTxt) ', () => {
    configPo.addData();

    configPo.reportsDir.sendKeys(null);

    //Test
    expect(configPo.saveBtn.isEnabled()).toBe(false);

  });
  it('shoul    d create file when memoryEndppoint field is empty ', () => {

    configPo.addData();
    //Left Blank
    configPo.memoryEndPoint.sendKeys(null);




    //Click Save
    configPo.saveBtn.click();

    browser.refresh();
    configPo.configPathDir.sendKeys(this.absolutePath + "Xtrac API.xml");

    //Test expected values
    expect(configPo.applicationName.getText()).toEqual('Xtrac API');
    expect(configPo.numIterations.getText()).toEqual(1000);
    expect(configPo.configPathDir.getText()).toEqual(this.absolutePath + "XtracTest.xml");


  });

});