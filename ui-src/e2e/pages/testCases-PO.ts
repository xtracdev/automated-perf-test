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
  const testCaseFileLocation = "../../../config/";
 
  
  
   class TestCasePageObject {

    testCaseTabBtn = element(by.id("test-cases"));
    
    testCaseDir = element(by.id("testCase-dir"));
    btnAdd = element(by.id("btn-add"));
    btnSave = element(by.id("btn-save"));
    btnCancel = element(by.id("btn-cancel"));

    testName = element(by.name("testname"));
    baseUri = element(by.name("baseURI"));
    overrideHost = element(by.name("overrideHost"));
    overridePort = element(by.name("overridePort"));
    multipart = element(by.name("multipart"))
    payload = element(by.name("payload"));
    httpMethod = element(by.name("httpMethod")).$$("option");
    header = element(by.name("header"));
    preThinkTime = element(by.name("preThinkTime"));
    postThinkTime = element(by.name("postThinkTime"));
    execWeight = element(by.name("execWeight")).$$("option");
    responseStatusCode = element(by.name("responseStatusCode"));
    responseContentType = element(by.name("responseContentType"));
    multipartPayload = element(by.name("multipartPayload"));
    responseValue = element(by.name("responseValue")); 


    toastrMessage = element(by.id("toast-container"));
    labels = $("json-schema-form").$$("label");
    requiredFields =  $$("p");

    absolutePath = path.resolve(__dirname, testCaseFileLocation);


    setToastrMessage(){
      return this.toastrMessage;
    }

    setTestCasesPage(){
      return this.testCaseTabBtn.click();
    }

    setTestCaseDir(){
    return this.testCaseDir.sendKeys(this.absolutePath);
    }
    
    setTestName(){
      return this.testName.sendKeys("Xtrac Test Case");
    }

    setMultipart(){
      return this.multipart.get(0).click();
    }
    
     setBaseUri(){
       return this.baseUri.sendKeys("./baseURI/testCases");
     }

     setOverrideHost(){
       return this.overrideHost.sendKeys("/etc/host");
     }

     setOverridePort(){
      return this.overridePort.sendKeys("80");
    }

    setPayload(){
      return this.payload.sendKeys("Xtrac Data");
    }

    setHttpMethod(){
      return this.httpMethod.get(0).click();
    }
    
    setPreThinkTime(){
      return this.preThinkTime.sendKeys(1);
    }
    setPostThinkTime(){
      return this.postThinkTime.sendKeys(5);
    }
    setExecWeight(){
      return this.execWeight.get(0).click();
    }
    setResponseStatusCode(){
      return this.responseStatusCode.sendKeys(500);
    }
    setResponseContentType(){
      return this.responseContentType.sendKeys("pass/fail");
    }
     
    setTestCaseData(){
      this.setTestCaseDir();
      this.setTestName();
      this.multipart.click();
      this.setBaseUri();
      this.setOverrideHost();
      this.setOverridePort();
      this.setPayload();
      this.httpMethod.click();
      this.setHttpMethod();
      this.setPreThinkTime();
      this.setPostThinkTime();
      this.execWeight.click();
      this.setExecWeight();
      this.setResponseStatusCode();
      this.setResponseContentType();
    
    }

    checkRequiredFields(){
      this.testName.sendKeys("x");
      this.baseUri.sendKeys("x");
      this.overrideHost.sendKeys("x");
      this.overridePort.sendKeys("x");
      this.payload.sendKeys("X");
      this.preThinkTime.sendKeys(1);
      this.postThinkTime.sendKeys(1);
      this.responseStatusCode.sendKeys(1);
      this.responseContentType.sendKeys("x");
      
      this.httpMethod.click();
      this.setHttpMethod();
      this.execWeight.click();
      this.setExecWeight();
     
      this.testName.sendKeys(Key.BACK_SPACE);
      this.baseUri.sendKeys(Key.BACK_SPACE);
      this.overrideHost.sendKeys(Key.BACK_SPACE);
      this.overridePort.sendKeys(Key.BACK_SPACE);
      this.payload.sendKeys(Key.BACK_SPACE);
      this.preThinkTime.sendKeys(Key.BACK_SPACE);
      this.postThinkTime.sendKeys(Key.BACK_SPACE);
      this.responseStatusCode.sendKeys(Key.BACK_SPACE);
      this.responseContentType.sendKeys(Key.BACK_SPACE);
     
    }

    checkForStrings() {
      this.responseStatusCode.sendKeys("x");
      this.postThinkTime.sendKeys("x");
      this.preThinkTime.sendKeys("x");
    }

    checkEisNotAccepted() {
      this.preThinkTime.sendKeys("e");
      this.postThinkTime.sendKeys("e");
      this.responseStatusCode.sendKeys("e");
    }

    checkNegativeValues(){
      this.preThinkTime.sendKeys(-1);
      this.postThinkTime.sendKeys(-1);
      this.responseStatusCode.sendKeys(-1);
    }
  }
export = TestCasePageObject;