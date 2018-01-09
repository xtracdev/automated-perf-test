@CreateConfigFile
Feature: Create Config File

  Scenario: User successfully creates a configuration file
    Given the automated performance ui server is available
    When the user makes a request for POST http://localhost:9191/configs with payload
    """
       {
         "apiName": "GucumberConfig",
           "targetHost": "localhost",
           "targetPort": "9191",
           "numIterations": 1000,
           "allowablePeakMemoryVariance": 30,
           "allowableServiceResponseTimeVariance": 30,
           "testCaseDir": "./definitions/testCases",
           "testSuiteDir": "./definitions/testSuites",
            "baseStatsOutputDir": "./envStats",
           "reportOutputDir": "./report",
           "concurrentUsers": 50,
           "testSuite": "suiteFileName.xml",
           "memoryEndpoint": "/alt/debug/vars",
           "requestDelay": 5000,
           "TPSFreq": 30,
           "rampUsers": 5,
           "rampDelay": 15
       }
    """
    Then the POST configuration service returns 201 HTTP status