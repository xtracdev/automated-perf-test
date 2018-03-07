Feature: Test Case Scenarios
  As an API user
  I want to be able use various requests for test cases
  So that I can test my application using custom metrics


                                ###################################
                                #######    GET ALL REQUESTS #######
                                ###################################

  Scenario: Successful retrieval all test cases with valid "GET" request
    Given the automated performance ui server is available
    And the header "testCasePathDir" is "/uiServices/test/"
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

                                ###################################
                                #######    GET REQUESTS #######
                                ###################################

  Scenario: Unsuccessful retrieval of test-Cases file (File Not Found)
    Given the automated performance ui server is available
    And the header "testCasePathDir" is "/uiServices/test/"
    When I send a "GET" request to "/test-cases/xxx"
    Then the response code should be 404


  Scenario: Unsuccessful retrieval of test-suites file (No Header)
    Given the automated performance ui server is available
    And the header "testCasePathDir" is ""
    When I send a "GET" request to "/test-suites/GodogTestCase"
    Then the response code should be 400

  Scenario: Retrieve Test Case file with valid "GET" request
    Given the file "Case3.xml" exists at "/uiServices/test/"
    Given the automated performance ui server is available
    And the header "testCasePathDir" is "/uiServices/test/"
    And the file name is "Case1.xml"
    When I send a "GET" request to "/test-cases/Case1"
    Then the response code should be 200
    And the test case response body should match json:
    """
   {
    "XMLName": {
    "Space": "",
    "Local": "testDefinition"
      },
      "TestName": GodogTestCase",
      "OverrideHost": "host",
      "OverridePort": "9191",
      "HTTPMethod": "GET",
      "Description": "desc",
      "BaseURI": "",
      "Multipart": false,
      "Payload": "",
      "MultipartPayload": null,
      "ResponseStatusCode": 0,
      "ResponseContentType": "",
      "Headers": null,
      "ResponseValues": null,
      "PreThinkTime": 0,
      "PostThinkTime": 0,
      "ExecWeight": ""
    }
    """

                                ###################################
                                #######    DELETE REQUESTS ########
                                ###################################

Scenario: Unsuccessful deleting of test-case (No Header)
  Given the automated performance ui server is available
  And the header "testCasePathDir" is ""
  When I send a "DELETE" request to "/test-cases/all"
  Then the response code should be 400


Scenario: Unsuccessful deleting of test-case file (Empty Directory)
  Given the automated performance ui server is available
  And the header "testCasePathDir" is "/uiServices/test/cases"
  When I send a "DELETE" request to "/test-cases/"
  Then the response code should be 404
   

Scenario: Successful deleting of test-case file with DELETE request
  Given the "deleteTest.xml" has been created at "/uiServices/test/cases"
  Given the automated performance ui server is available
  And the header "testCasePathDir" is "/uiServices/test/cases"
  When I send a "DELETE" request to "/test-cases/all"
  Then the response code should be 200
