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
    And the header "testSuitePathDir" is "/uiServices/test/"
    When I send "POST" request to "/test-suites" with a body:
    """
      {
       "name": "GodogTestSuite",
       "testStrategy": "SuiteBased",
       "description": "Services for XYZ",
       "testCases":[
         {
         "name":"file1",
          "preThinkTime": 1000,
          "postThinkTime": 2000,
          "execWeight": "Infrequent",
           "description": "Desc1"
         },
         {
          "name":"file2",
          "preThinkTime": 1,
          "postThinkTime": 10,
          "execWeight": "Sparse",
          "description": "Desc2"
         }
        ]
      }
    """
    Then the response code should be 201
    And the response body should be empty

  Scenario: Unsuccessful creation of test Suite (file already exists )
    Given the automated performance ui server is available
    And the header "testSuitePathDir" is "/uiServices/test/"
    When I send "POST" request to "/test-suites" with a body:
    """
      {
       "name": "GodogTestSuite",
       "testStrategy": "SuiteBased",
       "description": "Services for XYZ",
       "testCases":[
         {
         "name":"file1",
          "preThinkTime": 1000,
          "postThinkTime": 2000,
          "execWeight": "Infrequent",
           "description": "Desc1"
         },
         {
          "name":"file2",
          "preThinkTime": 1,
          "postThinkTime": 10,
          "execWeight": "Sparse",
          "description": "Desc2"
         }
        ]
      }
    """
    Then the response code should be 400


  Scenario: Unsuccessful creation of test Suite ( Missing Required Fields )
    Given the automated performance ui server is available
    And the header "testSuitePathDir" is "/uiServices/test/"
    When I send "POST" request to "/test-suites" with a body:
    """
      {
       "testCases": [
        {
        "name":"file1",
        "preThinkTime": 1000,
        "postThinkTime": 2000,
        "execWeight": "Infrequent"
        }
      ]
    }
    """
    Then the response code should be 400


  Scenario: Unsuccessful creation of test Suite ( No header defined )
    Given the automated performance ui server is available
    And the header "testSuitePathDir" is ""
    When I send "POST" request to "/test-suites" with a body:
      """
      {
       "name": "GodogTestSuite",
       "testStrategy": "SuiteBased",
       "description": "Services for XYZ",
       "testCases":[
         {
         "name":"file1",
          "preThinkTime": 1000,
          "postThinkTime": 2000,
          "execWeight": "Infrequent",
           "description": "Desc1"
         }
        ]
      }
    """
    Then the response code should be 400


  Scenario: Unsuccessful creation of test Suite ( Incorrect URL )
    Given the automated performance ui server is available
    And the header "testSuitePathDir" is "/uiServices/test/GodogTestSuite.xml"
    When I send "POST" request to "/xxxx" with a body:
      """
      {
       "name": "GodogTestSuite",
       "testStrategy": "SuiteBased",
       "description": "Services for XYZ",
       "testCases":[
         {
         "name":"file1",
          "preThinkTime": 1000,
          "postThinkTime": 2000,
          "execWeight": "Infrequent",
           "description": "Desc1"
         }
        ]
      }
    """
    Then the response code should be 405

  Scenario: Unsuccessful creation of test Suite(No Name defined)
    Given the automated performance ui server is available
    And the header "testSuitePathDir" is "/uiServices/test/GodogTestSuite.xml"
    When I send "POST" request to "/test-suites" with a body:
       """
      {
       "name": "",
       "testStrategy": "SuiteBased",
       "description": "Services for XYZ",
       "testCases":[
         {
         "name":"file1",
          "preThinkTime": 1000,
          "postThinkTime": 2000,
          "execWeight": "Infrequent",
           "description": "Desc1"
         }
        ]
      }
    """
    Then the response code should be 400
