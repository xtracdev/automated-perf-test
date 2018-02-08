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


                                ###################################
                                #######    PUT REQUESTS ###########
                                ###################################

  Scenario: Unsuccessful update of test-suite file with PUT request (No File Path)
    Given the file "GodogTestSuite.xml" exists at "/uiServices/test/"
    Given the automated performance ui server is available
    And the header "testSuitePathDir" is ""
    When I send "PUT" request to "/test-suites/GodogTestSuite" with body:
    """
      {
       "name": "GodogTestSuite2",
       "testStrategy": "SuiteBased",
       "description": "ServiceDesc",
       "testCases":[
         {
         "name":"file1.xml",
          "preThinkTime": 2000,
          "postThinkTime": 5000,
          "execWeight": "Infrequent",
           "description": "Desc1"
         }
        ]
      }
    """
   Then the response code should be 400


  Scenario: Unsuccessful update of test-suite file with PUT request (Incorrect File Name)
    Given the file "GodogTestSuite.xml" exists at "/uiServices/test/"
    Given the automated performance ui server is available
    And the header "testSuitePathDir" is "/uiServices/test/"
    When I send "PUT" request to "/test-suites/xxx" with body:
   """
      {
       "name": "GodogTestSuite2",
       "testStrategy": "SuiteBased",
       "description": "ServiceDesc",
       "testCases":[
         {
         "name":"file1.xml",
          "preThinkTime": 2000,
          "postThinkTime": 5000,
          "execWeight": "Infrequent",
           "description": "Desc1"
         }
        ]
      }
    """
    Then the response code should be 404


  Scenario: Unsuccessful update of test-suite file with PUT request (No File Name)
    Given the file "GodogTestSuite.xml" exists at "/uiServices/test/"
    Given the automated performance ui server is available
    And the header "testSuitePathDir" is "/uiServices/test/"
    When I send "PUT" request to "/test-suites/" with body:
    """
      {
       "name": "GodogTestSuite2",
       "testStrategy": "SuiteBased",
       "description": "ServiceDesc",
       "testCases":[
         {
         "name":"file1.xml",
          "preThinkTime": 2000,
          "postThinkTime": 5000,
          "execWeight": "Infrequent",
           "description": "Desc1"
         }
        ]
      }
    """
    Then the response code should be 404

  Scenario: Unsuccessful update of test-suite file with PUT request (Missing Required Fields)
    Given the file "GodogTestSuite.xml" exists at "/uiServices/test/"
    Given the automated performance ui server is available
    And the header "testSuitePathDir" is "/uiServices/test/"
    When I send "PUT" request to "/test-suites/GodogTestSuite" with body:
  """
      {
       "testCases":[
         {
         "name":"file1.xml",
          "preThinkTime": 2000,
          "postThinkTime": 5000,
          "execWeight": "Infrequent",
           "description": "Desc1"
         }
        ]
      }
    """
    Then the response code should be 400


  Scenario: Successful update of test-suite file with PUT request
    Given the file "GodogTestSuite.xml" exists at "/uiServices/test/"
    Given the automated performance ui server is available
    And the header "testSuitePathDir" is "/uiServices/test/"
    When I send "PUT" request to "/test-suites/GodogTestSuite" with body:
    """
      {
       "name": "GodogTestSuite2",
       "testStrategy": "SuiteBased",
       "description": "ServiceDesc",
       "testCases":[
         {
         "name":"file1.xml",
          "preThinkTime": 2000,
          "postThinkTime": 5000,
          "execWeight": "Infrequent",
           "description": "Desc1"
         }
        ]
      }
    """
    Then the response code should be 204
    And the response body should be empty
    When I send a "GET" request to "/test-suites/GodogTestSuite"
    And the updated file should match json:
"""
      {
       "name": "GodogTestSuite2",
       "testStrategy": "SuiteBased",
       "description": "ServiceDesc",
       "testCases":[
         {
         "name":"file1.xml",
          "preThinkTime": 2000,
          "postThinkTime": 5000,
          "execWeight": "Infrequent",
           "description": "Desc1"
         }
        ]
      }
    """
  Scenario: Successful update of test-suite file with PUT request (Update API Name to not match Filename)
    Given the file "GodogTestSuite.xml" exists at "/uiServices/test/"
    Given the automated performance ui server is available
    And the header "testSuitePathDir" is "/uiServices/test/"
    When I send "PUT" request to "/test-suites/GodogTestSuite" with body:
    """
      {
       "name": "GodogTestSuite2",
       "testStrategy": "SuiteBased",
        "description": "ServiceDesc",
       "testCases":[
         {
         "name":"file1.xml",
          "preThinkTime": 2000,
          "postThinkTime": 5000,
          "execWeight": "Infrequent",
           "description": "Desc1"
         }
        ]
      }
    """
    Then the response code should be 204
    And the response body should be empty
    When I send a "GET" request to "/test-suites/GodogTestSuite"
    And the updated file should match json:
    """
      {
       "name": "GodogTestSuite2",
       "testStrategy": "SuiteBased",
        "description": "ServiceDesc",
       "testCases":[
         {
         "name":"file1.xml",
          "preThinkTime": 2000,
          "postThinkTime": 5000,
          "execWeight": "Infrequent",
           "description": "Desc1"
         }
        ]
      }
    """