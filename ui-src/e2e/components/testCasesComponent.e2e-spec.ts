import {
    browser,
    element,
    by,
    By,
    $,
    $$,
    ExpectedConditions,
    protractor,
    WebDriver,
    Key
  } from "protractor";
  import { read } from "fs";
  import { Input } from "@angular/core/src/metadata/directives";
  import { ToastModule } from "ng2-toastr/ng2-toastr";
  const since = require("jasmine2-custom-message");

  import TestCasePageObject = require("../pages/testCases-PO");
  const testCasePO: TestCasePageObject = new TestCasePageObject();


  describe("test cases component", () => {
      beforeEach(()  => {
        browser.get("/");
        browser.executeScript("window.onbeforeunload = function(e){};");
        browser.driver
        .manage()
        .window()
        .maximize();
      });


      fit('should create a new test case file', () => {
      testCasePO.setTestNavigateToTestCasesPage();
      testCasePO.setTestCaseData();
      testCasePO.btnSave.click();
      expect(testCasePO.toastrMessage.getText()).toEqual(
        "Success!\nYour data has been saved!"
      );
    });

    it("should check that all text box names are correct", () => {
      testCasePO.setTestNavigateToTestCasesPage();
      testCasePO.setTestCaseData();
      expect(testCasePO.labels.get(0).getText()).toContain("Testname");
      expect(testCasePO.labels.get(1).getText()).toContain("Description");
      expect(testCasePO.labels.get(2).getText()).toContain("Base URI");
      expect(testCasePO.labels.get(3).getText()).toContain("Override Host");
      expect(testCasePO.labels.get(4).getText()).toContain("Multipart");
      expect(testCasePO.labels.get(5).getText()).toContain("Override Port");
      expect(testCasePO.labels.get(6).getText()).toContain("Payload");
      expect(testCasePO.labels.get(7).getText()).toContain("Http Method");
      expect(testCasePO.labels.get(8).getText()).toContain("Headers");
      expect(testCasePO.labels.get(11).getText()).toContain("Pre Think Time");
      expect(testCasePO.labels.get(12).getText()).toContain("Post Think Time");
      expect(testCasePO.labels.get(13).getText()).toContain("Exec Weight");
      expect(testCasePO.labels.get(14).getText()).toContain("Response Status Code");
      expect(testCasePO.labels.get(15).getText()).toContain("Response Content Type");
      expect(testCasePO.labels.get(16).getText()).toContain("Multipart Payload");

     
    });
    
    it("should check values of existing test case file are as expected", () => {
      testCasePO.setTestNavigateToTestCasesPage();
      testCasePO.setTestCaseData();
      testCasePO.testCaseDir.sendKeys(testCasePO.absolutePath);
      testCasePO.testName.sendKeys("suites");
      testCasePO.btnCancel.click();
      expect(testCasePO.testName.getAttribute("value")).toEqual("Xtrac Test Case");
      expect(testCasePO.baseUri.getAttribute("value")).toEqual("./baseURI/testCases");
      expect(testCasePO.overrideHost.getAttribute("value")).toEqual("overrideHost is:");
      expect(testCasePO.overridePort.getAttribute("value")).toEqual("overridePort is:");
      expect(testCasePO.payload.getAttribute("value")).toEqual("payload is:");
      expect(testCasePO.preThinkTime.getAttribute("value")).toEqual("1");
      expect(testCasePO.postThinkTime.getAttribute("value")).toEqual("5");
      expect(testCasePO.responseStatusCode.getAttribute("value")).toEqual("500");
      expect(testCasePO.responseContentType.getAttribute("value")).toEqual("pass/fail");
      
    });

    it("should throw error when a test case file path does not exist", () => {
      testCasePO.setTestNavigateToTestCasesPage();
      testCasePO.setTestCaseData();
      testCasePO.testCaseDir.sendKeys("/path/to/bad/location");
      testCasePO.btnSave.click();
      expect(testCasePO.toastrMessage.getText()).toContain(
        "Some of the fields do not conform to the schema"
      );
    });

    fit("should check requiredFields warning appears when requiredFields input is blank", () => {
      testCasePO.setTestNavigateToTestCasesPage();
      testCasePO.checkRequiredFields();
      expect(testCasePO.requiredFields.get(0).getText()).toContain("This field is required.")
      expect(testCasePO.requiredFields.get(1).getText()).toContain("This field is required.");
      expect(testCasePO.requiredFields.get(2).getText()).toContain("This field is required.");
      expect(testCasePO.requiredFields.get(3).getText()).toContain("This field is required.");
      expect(testCasePO.requiredFields.get(4).getText()).toContain("This field is required.");
      expect(testCasePO.requiredFields.get(5).getText()).toContain("This field is required.");
      expect(testCasePO.requiredFields.get(6).getText()).toContain("This field is required.");
      expect(testCasePO.requiredFields.get(7).getText()).toContain("This field is required.");
      expect(testCasePO.requiredFields.get(8).getText()).toContain("This field is required.");
      
    });
  
  });