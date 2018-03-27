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
    testCaseDir = element(by.id("testCase-directory"));
    btnAdd = element(by.id("btn-add"));
    btnSave = element(by.id("btn-save"));

    testName = element(by.name("testName"));
    describe = element(by.name("description"));
    baseUri = element(by.name("baseUri"));
    overrideHost = element(by.name("overrideHost"));
    overridePort = element(by.name("overRidePort"));
    payload = element(by.name("payload"));
    httpMethod = element(by.name("httpMethod"));
    header = element(by.name("header"));
    preThinkTime = element(by.name("preThinkTime"));
    postThinkTime = element(by.name("postThinkTime"));
    execWeight = element(by.name("execWeight"));
    responseStatusCode = element(by.name("responseStatusCode"));
    responseContentType = element(by.name("responseContentType"));
    multipartPayload = element(by.name("multipartPayload"));
    responseValue = element(by.name("responseValue")); 


    toastrMessage = element(by.id("toast-container"));
    labels = $("json-schema-form").$$("label");
    requiredFields = $$("p");

    absolutePath = path.resolve(__dirname, testCaseFileLocation);

    setTestCaseDir(){
    return this.testCaseDir.sendKeys(this.absolutePath);
    }
    
    setTestName(){
      return this.testName.sendKeys("Xtrac Test Case");
    }

    setDescribe(){
      return this.describe.sendKeys("test case is cool");
    }
     setBaseUri(){
       return this.baseUri.sendKeys("./baseURI/testCases ");
     }

     setOverrideHost(){
       return this.overrideHost.sendKeys("overrideHost is:");
     }

     setOverridePort(){
      return this.overridePort.sendKeys("overridePort is:");
    }

    setPayload(){
      return this.payload.sendKeys("payload is:");
    }

    setHttpMethod(){
      return this.httpMethod.get(0).click();
    }

    // setHeader(){
    //   return this.header.get(0);
    // }
   
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
    // setMultipartPayload(){
    //   return this.multipartPayload.get(0);
    // }

    // setResponseValue(){
    //   return this.responseValue.get(0);
    // }
     
    setTestCaseData(){
      this.setTestCaseDir();
      this.setTestName();
      this.setDescribe();
      this.setBaseUri();
      this.setOverrideHost();
      this.setOverridePort();
      this.setPayload();
      this.httpMethod.click();
      //this.setHeader();
      this.setPreThinkTime();
      this.setPostThinkTime();
      this.execWeight.click();
      this.setResponseStatusCode();
      this.setResponseContentType();
     // this.setMultipartPayload();
     // this.setResponseValue();

    }

    checkRequiredFields(){
      this.testCaseDir.sendKeys("x");
      this.testName.sendKeys("x");
      this.describe.sendKeys("x");
      this.baseUri.sendKeys("x");
      this.overrideHost.sendKeys("x");
      this.overridePort.sendKeys("x");
      this.payload.sendKeys("X");
      //this.header.get(0);
      this.preThinkTime.sendKeys(1);
      this.postThinkTime.sendKeys(1);
      this.responseStatusCode.sendKeys(1);
      this.responseContentType.sendKeys("x");
     // this.multipartPayload.get(1);
      //this.responseValue.get(1);

      this.testCaseDir.sendKeys(Key.BACK_SPACE);
      this.testName.sendKeys(Key.BACK_SPACE);
      this.describe.sendKeys(Key.BACK_SPACE);
      this.baseUri.sendKeys(Key.BACK_SPACE);
      this.overrideHost.sendKeys(Key.BACK_SPACE);
      this.overridePort.sendKeys(Key.BACK_SPACE);
      this.payload.sendKeys(Key.BACK_SPACE);
     // this.header.get(1);
      this.preThinkTime.sendKeys(Key.BACK_SPACE);
      this.postThinkTime.sendKeys(Key.BACK_SPACE);
      this.responseStatusCode.sendKeys(Key.BACK_SPACE);
      this.responseContentType.sendKeys(Key.BACK_SPACE);
     // this.multipartPayload.get(1);
      //this.responseValue.get(1);

    }
  }
export = TestCasePageObject;