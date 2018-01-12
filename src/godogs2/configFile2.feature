Feature: Create Configuration File
  As an API user
  I want to be able to create a configuration file
  So that I can test my application using custom metrics

  Scenario: Successful creation of config file
    When I send "POST" request to "/configs"
    Then the response code should be 200
    And the response should match json:
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
       "testSuite": "suiteFileName.xml",
       "requestDelay": 5000,
       "TPSFreq": 30,
       "rampUsers": 5,
       "rampDelay": 15

      }
      """

  Scenario: Try to create config file with "PUT" request
    When I send "PUT" request to "/configs"
    Then the response code should be 405
    And the response should match json:
      """
      {
       "error": "Method not allowed"
      }
      """

  Scenario: Unsuccessful creation of config file (invalid URL)
    When I send "POST" request to "/xxx"
    Then the response code should be 404
    And the response should match json:
      """
      {
       "error": "URL Not Found"
      }
      """
