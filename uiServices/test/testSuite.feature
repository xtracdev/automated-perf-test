Feature: Test Suite Creation
  As an API user
  I want to be able to create a test suite
  So that I can test my application using custom metrics

                                ###################################
                                #######    POST REQUESTS ##########
                                ###################################

  Scenario: Successful creation of test Suite
    Given there is no existing test file "GodogTestSuite.xml"
    Given the automated performance ui server is available
    And the header configsDirPath is "/uiServices/test/GodogTestSuite.xml"
    When I send "POST" request to "/test-suites" with a body:
         """
      {
  "name": "GodogTestSuite",
  "testStrategy": "SuiteBased",
  "testCases": [
    {
      "name":"file1.xml",
      "preThinkTime": 1000,
      "postThinkTime": 2000,
      "execWeight": "infrequent"
    },
    {
      "name":"file2.xml",
      "preThinkTime": 1,
      "postThinkTime": 10,
      "execWeight": "sparce"
    }
  ]
}
      """
    Then the response code should be 201
    And the response body should be empty
    And the config file was created at location defined by configsPathDir

  Scenario: Unsuccessful creation of test Suite (file already exists )
    Given the automated performance ui server is available
    And the header configsDirPath is "/uiServices/test/GodogTestSuite.xml"
    When I send "POST" request to "/test-suites" with a body:
         """
      {
  "name": "GodogTestSuite",
  "testStrategy": "SuiteBased",
  "testCases": [
    {
      "name":"file1.xml",
      "preThinkTime": 1000,
      "postThinkTime": 2000,
      "execWeight": "infrequent"
    },
    {
      "name":"file2.xml",
      "preThinkTime": 1,
      "postThinkTime": 10,
      "execWeight": "sparce"
    }
  ]
}
      """
    Then the response code should be 400


  Scenario: Unsuccessful creation of test Suite ( Missing Required Fields )
    Given the automated performance ui server is available
    And the header configsDirPath is "/uiServices/test/GodogTestSuite.xml"
    When I send "POST" request to "/test-suites" with a body:
         """
      {
  "testStrategy": "SuiteBased",
  "testCases": [
    {
      "name":"file1.xml",
      "preThinkTime": 1000,
      "postThinkTime": 2000,
      "execWeight": "infrequent"
    },
    {
      "name":"file2.xml",
      "preThinkTime": 1,
      "postThinkTime": 10,
      "execWeight": "sparce"
    }
  ]
}
      """
    Then the response code should be 400


  Scenario: Unsuccessful creation of test Suite ( No header defined )
    Given the automated performance ui server is available
    And the header configsDirPath is ""
    When I send "POST" request to "/test-suites" with a body:
         """
      {
  "name": "GodogTestSuite",
  "testStrategy": "testStrat",
  "testCases": [
    {
      "name":"file1.xml",
      "preThinkTime": 1000,
      "postThinkTime": 2000,
      "execWeight": "infrequent"
    },
    {
      "name":"file2.xml",
      "preThinkTime": 1,
      "postThinkTime": 10,
      "execWeight": "sparce"
    }
  ]
}
      """
    Then the response code should be 400


  Scenario: Unsuccessful creation of test Suite ( Incorrect URL )
    Given the automated performance ui server is available
    And the header configsDirPath is "/uiServices/test/GodogTestSuite.xml"
    When I send "POST" request to "/xxxx" with a body:
         """
      {
  "name": "GodogTestSuite",
  "testStrategy": "SuiteBased",
  "testCases": [
    {
      "name":"file1.xml",
      "preThinkTime": 1000,
      "postThinkTime": 2000,
      "execWeight": "infrequent"
    },
    {
      "name":"file2.xml",
      "preThinkTime": 1,
      "postThinkTime": 10,
      "execWeight": "sparce"
    }
  ]
}
      """
    Then the response code should be 405
