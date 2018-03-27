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


  describe('Test Cases ', () => {
      beforeEach(()  => {
        browser.get("/");
        browser.executeScript("window.onbeforeunload = function(e){};");
        browser.driver.manage().window().maximize();
      });

      fit('should create a test case file', () => {
      //click tab button
      testCasePO.setTestCaseData();
        
      //set test data
        testCasePO.btnAdd.click();
      //click submit

     // browser.sleep(520000000000000000)
      //confirm message toaster
      expect(testCasePO.toastrMessage.getText()).toEqual(
        "Success!\nYour data has been saved!"
      );
    });


  });