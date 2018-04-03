Feature: Test Suite Creation
  As an API user
  I want to be able to create a test suite
  So that I can test my application using custom metrics

                                ###################################
                                #######    POST REQUESTS ##########
                                ###################################

  Scenario: Successful creation of test Suite
    Given there is no existing test file "TestSuiteSAMPLE.xml"
    Given the automated performance ui server is available
    And the header "path" is "/uiServices/test/"
    When I send "POST" request to "/test-suites" with a body:
    """
      {
       "name": "TestSuiteSAMPLE",
       "testStrategy": "SuiteBased",
       "description": "Services for XYZ",
       "testCases":[
         {
         "name":"file1",
          "preThinkTime": 1000,
          "postThinkTime": 2000,
          "execWeight": "Infrequent"
         }
        ]
      }
    """
    Then the response code should be 201
    And the response body should be empty

  Scenario: Unsuccessful creation of test Suite (file already exists )
    Given the automated performance ui server is available
    Given there is no existing test file "TestSuiteSAMPLE.xml"
    And the header "path" is "/uiServices/test/"
    When I send "POST" request to "/test-suites" with a body:
    """
      {
       "name": "TestSuiteSAMPLE",
       "testStrategy": "SuiteBased",
       "description": "Services for XYZ",
       "testCases":[
         {
         "name":"file1",
          "preThinkTime": 1000,
          "postThinkTime": 2000,
          "execWeight": "Infrequent"
         }
        ]
      }
    """
    Then the response code should be 201
    And the header "path" is "/uiServices/test/"
    When I send "POST" request to "/test-suites" with a body:
    """
      {
       "name": "TestSuiteSAMPLE",
       "testStrategy": "SuiteBased",
       "description": "Services for XYZ",
       "testCases":[
         {
         "name":"file1",
          "preThinkTime": 1000,
          "postThinkTime": 2000,
          "execWeight": "Infrequent"
         },
         {
          "name":"file2",
          "preThinkTime": 1,
          "postThinkTime": 10,
          "execWeight": "Sparse"
         }
        ]
      }
    """
    Then the response code should be 400


  Scenario: Unsuccessful creation of test Suite ( Missing Required Fields )
    Given the automated performance ui server is available
    And the header "path" is "/uiServices/test/"
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
    And the header "path" is ""
    When I send "POST" request to "/test-suites" with a body:
      """
      {
       "name": "TestSuiteSAMPLE",
       "testStrategy": "SuiteBased",
       "description": "Services for XYZ",
       "testCases":[
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


  Scenario: Unsuccessful creation of test Suite ( Incorrect URL )
    Given the automated performance ui server is available
    And the header "path" is "/uiServices/test/TestSuiteSAMPLE.xml"
    When I send "POST" request to "/xxxx" with a body:
      """
      {
       "name": "TestSuiteSAMPLE",
       "testStrategy": "SuiteBased",
       "description": "Services for XYZ",
       "testCases":[
         {
         "name":"file1",
          "preThinkTime": 1000,
          "postThinkTime": 2000,
          "execWeight": "Infrequent"
         }
        ]
      }
    """
    Then the response code should be 405

  Scenario: Unsuccessful creation of test Suite(No Name defined)
    Given the automated performance ui server is available
    And the header "path" is "/uiServices/test/TestSuiteSAMPLE.xml"
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
          "execWeight": "Infrequent"
         }
        ]
      }
    """
    Then the response code should be 400


                                ###################################
                                #######    PUT REQUESTS ###########
                                ###################################

  Scenario: Unsuccessful update of test-suite file with PUT request (No File Path)
    Given the file "TestSuiteSAMPLE.xml" exists at "/uiServices/test/"
    Given the automated performance ui server is available
    And the header "path" is ""
    When I send "PUT" request to "/test-suites/TestSuiteSAMPLE" with body:
    """
      {
       "name": "TestSuiteSAMPLE",
       "testStrategy": "SuiteBased",
       "description": "ServiceDesc",
       "testCases":[
         {
         "name":"file1.xml",
          "preThinkTime": 2000,
          "postThinkTime": 5000,
          "execWeight": "Infrequent"
         }
        ]
      }
    """
   Then the response code should be 400


  Scenario: Unsuccessful update of test-suite file with PUT request (Incorrect File Name)
    Given the file "TestSuiteSAMPLE.xml" exists at "/uiServices/test/"
    Given the automated performance ui server is available
    And the header "path" is "/uiServices/test/"
    When I send "PUT" request to "/test-suites/xxx" with body:
   """
      {
       "name": "TestSuiteSAMPLE2",
       "testStrategy": "SuiteBased",
       "description": "ServiceDesc",
       "testCases":[
         {
         "name":"file1.xml",
          "preThinkTime": 2000,
          "postThinkTime": 5000,
          "execWeight": "Infrequent"
         }
        ]
      }
    """
    Then the response code should be 404


  Scenario: Unsuccessful update of test-suite file with PUT request (No File Name)
    Given the file "TestSuiteSAMPLE.xml" exists at "/uiServices/test/"
    Given the automated performance ui server is available
    And the header "path" is "/uiServices/test/"
    When I send "PUT" request to "/test-suites/" with body:
    """
      {
       "name": "TestSuiteSAMPLE2",
       "testStrategy": "SuiteBased",
       "description": "ServiceDesc",
       "testCases":[
         {
         "name":"file1.xml",
          "preThinkTime": 2000,
          "postThinkTime": 5000,
          "execWeight": "Infrequent"
         }
        ]
      }
    """
    Then the response code should be 404

  Scenario: Unsuccessful update of test-suite file with PUT request (Missing Required Fields)
    Given the file "TestSuiteSAMPLE.xml" exists at "/uiServices/test/"
    Given the automated performance ui server is available
    And the header "path" is "/uiServices/test/"
    When I send "PUT" request to "/test-suites/TestSuiteSAMPLE" with body:
  """
      {
       "testCases":[
         {
         "name":"file1.xml",
          "preThinkTime": 2000,
          "postThinkTime": 5000,
          "execWeight": "Infrequent"
         }
        ]
      }
    """
    Then the response code should be 400


  Scenario: Successful update of test-suite file with PUT request
    Given the file "TestSuiteSAMPLE.xml" exists at "/uiServices/test/"
    Given the automated performance ui server is available
    And the header "path" is "/uiServices/test/"
    When I send "PUT" request to "/test-suites/TestSuiteSAMPLE" with body:
    """
      {
       "name": "TestSuiteSAMPLE2",
       "testStrategy": "SuiteBased",
       "description": "ServiceDesc",
       "testCases":[
         {
         "name":"file1.xml",
          "preThinkTime": 2000,
          "postThinkTime": 5000,
          "execWeight": "Infrequent"
         }
        ]
      }
    """
    Then the response code should be 204
    And the response body should be empty
    When I send a "GET" request to "/test-suites/TestSuiteSAMPLE"
    And the updated file should match json:
"""
    {
       "name": "TestSuiteSAMPLE2",
       "testStrategy": "SuiteBased",
       "description": "ServiceDesc",
       "testCases":[
         {
         "name":"file1.xml",
          "preThinkTime": 2000,
          "postThinkTime": 5000,
          "execWeight": "Infrequent"
         }
        ]
      }
    """
  Scenario: Successful update of test-suite file with PUT request (Update API Name to not match Filename)
    Given the file "TestSuiteSAMPLE.xml" exists at "/uiServices/test/"
    Given the automated performance ui server is available
    And the header "path" is "/uiServices/test/"
    When I send "PUT" request to "/test-suites/TestSuiteSAMPLE" with body:
    """
      {
       "name": "TestSuiteSAMPLE2",
       "testStrategy": "SuiteBased",
        "description": "ServiceDesc",
       "testCases":[
         {
         "name":"file1",
          "preThinkTime": 2000,
          "postThinkTime": 5000,
          "execWeight": "Infrequent"
         }
        ]
      }
    """
    Then the response code should be 204
    And the response body should be empty
    When I send a "GET" request to "/test-suites/TestSuiteSAMPLE"
    And the updated file should match json:
    """
      {
       "name": "TestSuiteSAMPLE2",
       "testStrategy": "SuiteBased",
        "description": "ServiceDesc",
       "testCases":[
         {
         "name":"file1.xml",
          "preThinkTime": 2000,
          "postThinkTime": 5000,
          "execWeight": "Infrequent"
         }
        ]
      }
    """



                                ###################################
                                #######    GET REQUESTS ########
                                ###################################

  Scenario: Try to retrieve test-suite file with valid "GET" request
    Given the automated performance ui server is available
    Given the file "testSuite.xml" exists at "/uiServices/test/samples"
    And the header "path" is "/uiServices/test/samples"
    And the file name is "testSuite.xml"
    When I send a "GET" request to "/test-suites/testSuite"
    Then the response code should be 200
    And the test suite response body should match json:
    """
      {
          "XMLName": {
          "Space": "",
          "Local": "testSuite"
          },
           "name": "testSuite",
           "description": "Services for XYZrrr",
           "testStrategy": "SuiteBased"
      }
    """


  Scenario: Unsuccessful retrieval of test-suites file (File Not Found)
    Given the automated performance ui server is available
    Given the file "TestSuiteSAMPLE.xml" exists at "/uiServices/test/"
    Given the automated performance ui server is available
    And the header "path" is "/uiServices/test/"
    When I send a "GET" request to "/test-suites/xxx"
    Then the response code should be 404


  Scenario: Unsuccessful retrieval of test-suites file (No Header)
    Given the automated performance ui server is available
    And the header "path" is ""
    When I send a "GET" request to "/test-suites/TestSuiteSAMPLE"
    Then the response code should be 400


                                ###################################
                                #######    GET ALL REQUESTS #######
                                ###################################

  Scenario: Try to retrieve all test suites with valid "GET" request
    Given the automated performance ui server is available
    And the header "path" is "/uiServices/test/"
    When I send a "GET" request to "/test-suites"
    Then the response code should be 200
    And the test suite collection response body should match json:
     """
        [
          {
          "file": "TestSuiteSAMPLE.xml",
          "name": "TestSuiteSAMPLE2,
          "description": "ServiceDesc",
          }
        ]
    """


  Scenario: Unsuccessful retrieval of test-suites (No Header)
    Given the automated performance ui server is available
    And the header "path" is ""
    When I send a "GET" request to "/test-suites"
    Then the response code should be 400