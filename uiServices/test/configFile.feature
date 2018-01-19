Feature: Create Configuration File
  As an API user
  I want to be able to create a configuration file
  So that I can test my application using custom metrics

  Scenario: Successful creation of config file
    Given the automated performance ui server is available
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
    And the header configsDirPath is "/uiServices/test/GodogConfig.xml"
    Then the response code should be 201
    And the response body should be empty
    And the config file was created at location defined by configsPathDir

  Scenario: Try to create config file with "PUT" request
    When I send "PUT" request to "/configs"
    Then the response code should be 405

  Scenario: Try to create config file with "GET" request
    When I send "GET" request to "/configs"
    Then the response code should be 405

  Scenario: Try to create config file with "DELETE" request
    When I send "DELETE" request to "/configs"
    Then the response code should be 405

  Scenario: Unsuccessful creation of config file (invalid URL)
    When I send "POST" request to "/xxx"
    Then the response code should be 404
