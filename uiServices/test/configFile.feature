Feature: Create Configuration File
As an API user
I want to be able to create a configuration file
So that I can test my application using custom metrics


  ###################################
  #######    POST REQUESTS ##########
  ###################################

  Scenario: Successful creation of config file
    Given there is no existing test file "PostConfigSAMPLE.xml"
    Given the automated performance ui server is available
    And the header "path" is "/uiServices/test"
    When I send "POST" request to "/configs" with a body:
      """
      {
      "apiName": "PostConfigSAMPLE",
      "targetHost": "localhost",
      "targetPort":"9191",
      "memoryEndpoint": "/alt/debug/vars",
      "numIterations": 1000,
      "allowablePeakMemoryVariance": 30,
      "allowableServiceResponseTimeVariance": 30,
      "testCaseDir": "./definitions/testCases",
      "testSuiteDir": "./definitions/testSuites",
      "baseStatsOutputDir": "./envStats",
      "reportOutputDir": "./report",
      "concurrentUsers": 50,
      "testSuite": "Default-3",
      "requestDelay": 5000,
      "TPSFreq": 30,
      "rampUsers": 5,
      "rampDelay": 15
      }
      """
    Then the response code should be 201
    And the response body should be empty

  Scenario: Unsuccessful creation of config file (file already exists)
    Given the automated performance ui server is available
    Given there is no existing test file "DUPConfigSAMPLE.xml"
    And the header "path" is "/uiServices/test"
    When I send "POST" request to "/configs" with a body:
      """
      {
      "apiName": "DUPConfigSAMPLE",
      "targetHost": "localhost",
      "targetPort":"9191",
      "memoryEndpoint": "/alt/debug/vars",
      "numIterations": 1000,
      "allowablePeakMemoryVariance": 30,
      "allowableServiceResponseTimeVariance": 30,
      "testCaseDir": "./definitions/testCases",
      "testSuiteDir": "./definitions/testSuites",
      "baseStatsOutputDir": "./envStats",
      "reportOutputDir": "./report",
      "concurrentUsers": 50,
      "testSuite": "Default-3",
      "requestDelay": 5000,
      "TPSFreq": 30,
      "rampUsers": 5,
      "rampDelay": 15
      }
      """
    Then the response code should be 201
    And the header "path" is "/uiServices/test"
    When I send "POST" request to "/configs" with a body:
      """
      {
      "apiName": "DUPConfigSAMPLE",
      "targetHost": "localhost",
      "targetPort":"9191",
      "memoryEndpoint": "/alt/debug/vars",
      "numIterations": 1000,
      "allowablePeakMemoryVariance": 30,
      "allowableServiceResponseTimeVariance": 30,
      "testCaseDir": "./definitions/testCases",
      "testSuiteDir": "./definitions/testSuites",
      "baseStatsOutputDir": "./envStats",
      "reportOutputDir": "./report",
      "concurrentUsers": 50,
      "testSuite": "Default-3",
      "requestDelay": 5000,
      "TPSFreq": 30,
      "rampUsers": 5,
      "rampDelay": 15
      }
      """
    Then the response code should be 400

  Scenario: Unsuccessful creation of config file (Missing required field)
    Given the automated performance ui server is available
    When I send "POST" request to "/configs" with a body:
      """
      {
      "apiName": "ConfigSAMPLE2",
      "targetHost": "localhost",
      "targetPort":"9191",
      "memoryEndpoint": "/alt/debug/vars",
      "numIterations": 1000,
      "allowablePeakMemoryVariance": 30,
      "allowableServiceResponseTimeVariance": 30,
      "testCaseDir": "./definitions/testCases",
      "testSuiteDir": "./definitions/testSuites",
      "baseStatsOutputDir": "./envStats",
      "reportOutputDir": "./report",
      "concurrentUsers": 50,
      "testSuite": "",
      "requestDelay": 5000,
      "TPSFreq": 30,
      "rampUsers": 5,
      "rampDelay": 15
      }
      """
    And the header "path" is "/uiServices/test"
    Then the response code should be 400


  Scenario: Unsuccessful creation of config file (Missing Header)
    Given the automated performance ui server is available
    And the header "path" is ""
    When I send "POST" request to "/configs" with a body:
      """
      {
      "apiName": "ConfigSAMPLE3",
      "targetHost": "localhost",
      "targetPort":"9191",
      "memoryEndpoint": "/alt/debug/vars",
      "numIterations": 1000,
      "allowablePeakMemoryVariance": 30,
      "allowableServiceResponseTimeVariance": 30,
      "testCaseDir": "./definitions/testCases",
      "testSuiteDir": "./definitions/testSuites",
      "baseStatsOutputDir": "./envStats",
      "reportOutputDir": "./report",
      "concurrentUsers": 50,
      "testSuite": "Default-1",
      "requestDelay": 5000,
      "TPSFreq": 30,
      "rampUsers": 5,
      "rampDelay": 15
      }
      """
    Then the response code should be 400


  ###################################
  #######    GET REQUESTS ########
  ###################################


  Scenario: Try to retrieve config file with valid "GET" request
    Given the file "testConfig.xml" exists at "/uiServices/test/samples"
    And the header "path" is "/uiServices/test/samples"
    And the file name is "testConfig.xml"
    When I send a "GET" request to "/configs/testConfig"
    Then the response code should be 200
    And the response body should match json:
      """
      {
        "apiName": "testConfig",
        "targetHost": "localhost",
        "targetPort": "8080",
        "numIterations": 55,
        "allowablePeakMemoryVariance": 12,
        "allowableServiceResponseTimeVariance": 15,
        "testCaseDir": "./definitions/testCases",
        "testSuiteDir": "./definitions/testSuites",
        "baseStatsOutputDir": "./envStats",
        "reportOutputDir": ".efreefef",
        "concurrentUsers": 50,
        "testSuite": "Default-1",
        "memoryEndpoint": "/alt/debug/vars",
        "requestDelay": 5000,
        "TPSFreq": 30,
        "rampUsers": 15,
        "rampDelay": 15,
        "GBS": false,
        "ReBaseMemory": false,
        "ReBaseAll": false,
        "ExecutionHost": "",
        "ReportTemplateFile": ""
      }
      """

  Scenario: Unsuccessful retrieval of config file (File Not Found)
    Given the automated performance ui server is available
    Given the file "ConfigSAMPLE.xml" exists at "/uiServices/test/"
    Given the automated performance ui server is available
    And the header "path" is "/uiServices/test/"
    When I send a "GET" request to "/configs/xxx"
    Then the response code should be 404



  ##################################
  ######    PUT REQUESTS ###########
  ##################################

  Scenario: Unsuccessful update of config file with PUT request (No File Path)
    Given the file "PUTConfigSAMPLE.xml" exists at "/uiServices/test/"
    Given the automated performance ui server is available
    And the header "path" is ""
    When I send "PUT" request to "/configs/PUTConfigSAMPLE" with body:
      """
      {
      "apiName": "PUTConfigSAMPLE",
      "targetHost": "localhost2",
      "targetPort":"1001",
      "memoryEndpoint": "/alt/debug/vars",
      "numIterations": 4000,
      "allowablePeakMemoryVariance": 50,
      "allowableServiceResponseTimeVariance": 50,
      "testCaseDir": "./definitions/testCases",
      "testSuiteDir": "./definitions/testSuites",
      "baseStatsOutputDir": "./envStats",
      "reportOutputDir": "./report",
      "concurrentUsers": 50,
      "testSuite": "Default-3",
      "requestDelay": 1000,
      "TPSFreq": 10,
      "rampUsers": 10,
      "rampDelay": 10
      }
      """
    Then the response code should be 400

  Scenario: Unsuccessful update of config file with PUT request (Incorrect File Name)
    Given the file "PUTConfigSAMPLE-1.xml" exists at "/uiServices/test/"
    Given the automated performance ui server is available
    And the header "path" is "/uiServices/test/"
    When I send "PUT" request to "/configs/xxx" with body:
      """
      {
      "apiName": "PUTConfigSAMPLE-1",
      "targetHost": "localhost",
      "targetPort":"1001",
      "memoryEndpoint": "/alt/debug/vars",
      "numIterations": 4000,
      "allowablePeakMemoryVariance": 50,
      "allowableServiceResponseTimeVariance": 50,
      "testCaseDir": "./definitions/testCases",
      "testSuiteDir": "./definitions/testSuites",
      "baseStatsOutputDir": "./envStats",
      "reportOutputDir": "./report",
      "concurrentUsers": 50,
      "testSuite": "Default-3",
      "requestDelay": 1000,
      "TPSFreq": 10,
      "rampUsers": 10,
      "rampDelay": 10
      }
      """
    Then the response code should be 404

  Scenario: Unsuccessful update of config file with PUT request (No File Name)
    Given the file "ErrorConfigSAMPLE.xml" exists at "/uiServices/test/"
    Given the automated performance ui server is available
    And the header "path" is "/uiServices/test/"
    When I send "PUT" request to "/configs/" with body:
      """
      {
      "apiName": "ErrorConfigSAMPLE",
      "targetHost": "localhost",
      "targetPort":"1001",
      "memoryEndpoint": "/alt/debug/vars",
      "numIterations": 4000,
      "allowablePeakMemoryVariance": 50,
      "allowableServiceResponseTimeVariance": 50,
      "testCaseDir": "./definitions/testCases",
      "testSuiteDir": "./definitions/testSuites",
      "baseStatsOutputDir": "./envStats",
      "reportOutputDir": "./report",
      "concurrentUsers": 50,
      "testSuite": "Default-3",
      "requestDelay": 1000,
      "TPSFreq": 10,
      "rampUsers": 10,
      "rampDelay": 10
      }
      """
    Then the response code should be 404

    
