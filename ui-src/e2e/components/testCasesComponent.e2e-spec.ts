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
      testCasePO.setTestCasesPage();
      testCasePO.setTestCaseData();
      testCasePO.btnSave.click();
      expect(testCasePO.toastrMessage.getText()).toEqual(
        "Success!\nYour data has been saved!"
      );
    });

    fit("should check that all text box names are correct", () => {
      testCasePO.setTestCasesPage();
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
    

    fit("should throw error when a test case file path does not exist", () => {
      testCasePO.setTestCasesPage();
      testCasePO.setTestCaseData();
      testCasePO.testCaseDir.sendKeys("/path/to/bad/location");
      testCasePO.btnSave.click();
      expect(testCasePO.toastrMessage.getText()).toContain(
        "Some of the fields do not conform to the schema"
      );
    });

   fit("should check requiredFields warning appears when requiredFields input is blank", () => {
      testCasePO.setTestCasesPage();
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