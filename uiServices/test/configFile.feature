Feature: Create Configuration File
  As an API user
  I want to be able to create a configuration file
  So that I can test my application using custom metrics

  Scenario: Successful creation of config file
    Given the automated performance ui server is available
    And the header configsDirPath is "/uiServices/test/GodogConfig.xml"
    When I send "POST" request to "/configs" with a body:
         """
      {
       "apiName": "GodogConfig",
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
    And the config file was created at location defined by configsPathDir

  Scenario: Try to create config file with "PUT" request
    When I send "PUT" request to "/configs"
    Then the response code should be 405

  Scenario: Try to retrieve config file with valid "GET" request
    Given the config file "GodogConfig.xml" exists at "/uiServices/test/"
    Given the automated performance ui server is available
    And the header configsDirPath is "/uiServices/test/"
    And the file name is "GodogConfig.xml"
    When I send a "GET" request to "/configs/GodogConfig"
    Then the response code should be 200
    And the response body should match json:
    """
      {
       "apiName": "GodogConfig",
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

  Scenario: Unsuccessful retrieval of config file (File Not Found)
    Given the automated performance ui server is available
    Given the config file "GodogConfig.xml" exists at "/uiServices/test/"
    Given the automated performance ui server is available
    And the header configsDirPath is "/uiServices/test/"
    When I send a "GET" request to "/configs/xxx"
    Then the response code should be 404


  Scenario: Unsuccessful retrieval of config file (No Header)
    Given the automated performance ui server is available
    And the header configsDirPath is ""
    When I send a "GET" request to "/configs/GodogConfig"
    Then the response code should be 400

  Scenario: Try to create config file with "DELETE" request
    When I send "DELETE" request to "/configs"
    Then the response code should be 405

  Scenario: Unsuccessful retrieval of config file (invalid URL)
    When I send "GET" request to "/xxx"
    Then the response code should be 404

  Scenario: Unsuccessful creation of config file (Missing required field)
    Given the automated performance ui server is available
    When I send "POST" request to "/configs" with a body:
         """
      {
       "apiName": "GodogConfig2",
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
    And the header configsDirPath is "/uiServices/test/GodogConfig.xml"
    Then the response code should be 400


  Scenario: Unsuccessful creation of config file (Missing Header)
    Given the automated performance ui server is available
    And the header configsDirPath is ""
    When I send "POST" request to "/configs" with a body:
         """
      {
       "apiName": "GodogConfig3",
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