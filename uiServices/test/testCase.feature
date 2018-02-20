Feature: Test Case Scenarios
  As an API user
  I want to be able use various requests for test cases
  So that I can test my application using custom metrics


                                ###################################
                                #######    GET ALL REQUESTS #######
                                ###################################

  Scenario: Susscessful retrieval all test cases with valid "GET" request
    Given the automated performance ui server is available
    And the header "testSuitePathDir" is "/uiServices/test/"
    When I send a "GET" request to "/test-cases"
    Then the response code should be 200
    And the test case collection response body should match json:
      """
        [
          {
          "name": "GodogTestCase,
          "description": "Case Desc",
          "httpMethod": "GET"
          }
        ]
    """


  Scenario: Unsuccessful retrieval of test-cases (No Header)
    Given the automated performance ui server is available
    And the header "testSuitePathDir" is ""
    When I send a "GET" request to "/test-cases"
    Then the response code should be 400