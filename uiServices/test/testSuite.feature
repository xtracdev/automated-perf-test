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
      "execWeight": "Infrequent"
    },
    {
      "name":"file2.xml",
      "preThinkTime": 1,
      "postThinkTime": 10,
      "execWeight": "Sparse"
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
      "execWeight": "Infrequent"
    },
    {
      "name":"file2.xml",
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
    And the header configsDirPath is "/uiServices/test/GodogTestSuite.xml"
    When I send "POST" request to "/test-suites" with a body:
         """
      {
  "testCases": [
    {
      "name":"file1.xml",
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
      "execWeight": "Infrequent"
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
      "execWeight": "Infrequent"
    }
  ]
}
      """
    Then the response code should be 405

  Scenario: Unsuccessful creation of test Suite(No Name defined)
    Given the automated performance ui server is available
    And the header configsDirPath is "/uiServices/test/GodogTestSuite.xml"
    When I send "POST" request to "/test-suites" with a body:
         """
      {
  "name": "",
  "testStrategy": "SuiteBased",
  "testCases": [
    {
      "name":"file1.xml",
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
    Given the test-suite file "GodogTestSuite.xml" exists at "/uiServices/test/"
    Given the automated performance ui server is available
    And the header configsDirPath is ""
    When I send "PUT" request to "/test-suites/GodogTestSuite.xml" with body:

      """
       {
  "name": "GodogTestSuite",
  "testStrategy": "SuiteBased",
  "testCases": [
    {
       "name":"file1.xml",
       "preThinkTime": 1000,
       "postThinkTime": 2000,
       "execWeight": "Infrequent"
    },
        {
          "name":"file2.xml",
          "preThinkTime": 1,
          "postThinkTime": 10,
          "execWeight": "Sparse"
        }
    ]
  }
       """
   Then the response code should be 400


  Scenario: Unsuccessful update of test-suite file with PUT request (Incorrect File Name)
    Given the config file "GodogTestSuite.xml" exists at "/uiServices/test/"
    Given the automated performance ui server is available
    And the header configsDirPath is "/uiServices/test/"
    When I send "PUT" request to "/test-suites/xxx" with body:
        """
       {
  "name": "GodogTestSuite",
  "testStrategy": "SuiteBased",
  "testCases": [
    {
       "name":"file1.xml",
       "preThinkTime": 1000,
       "postThinkTime": 2000,
       "execWeight": "Infrequent"
    },
        {
          "name":"file2.xml",
          "preThinkTime": 1,
          "postThinkTime": 10,
          "execWeight": "Sparse"
        }
    ]
  }
       """
    Then the response code should be 404


  Scenario: Unsuccessful update of test-suite file with PUT request (No File Name)
    Given the config file "GodogTestSuite.xml" exists at "/uiServices/test/"
    Given the automated performance ui server is available
    And the header configsDirPath is "/uiServices/test/"
    When I send "PUT" request to "/test-suites/" with body:
        """
       {
  "name": "GodogTestSuite",
  "testStrategy": "SuiteBased",
  "testCases": [
    {
       "name":"file1.xml",
       "preThinkTime": 1000,
       "postThinkTime": 2000,
       "execWeight": "Infrequent"
    },
        {
          "name":"file2.xml",
          "preThinkTime": 1,
          "postThinkTime": 10,
          "execWeight": "Sparse"
        }
    ]
  }
       """
    Then the response code should be 404

  Scenario: Unsuccessful update of test-suite file with PUT request (Missing Required Fields)
    Given the config file "GodogTestSuite.xml" exists at "/uiServices/test/"
    Given the automated performance ui server is available
    And the header configsDirPath is "/uiServices/test/"
    When I send "PUT" request to "/test-suites/GodogTestSuite.xml" with body:
           """
       {
  "name": "GodogTestSuite",
  "testStrategy": "",
  "testCases": [
    {
       "name":"",
       "preThinkTime": ,
       "postThinkTime": ,
       "execWeight": "Infrequent"
    },
        {
          "name":"file2.xml",
          "preThinkTime": ,
          "postThinkTime": 10,
          "execWeight": ""
        }
    ]
  }
       """
    Then the response code should be 400


  Scenario: Successful update of test-suite file with PUT request
    Given the config file "GodogTestSuite.xml" exists at "/uiServices/test/"
    Given the automated performance ui server is available
    And the header configsDirPath is "/uiServices/test/"
    When I send "PUT" request to "/test-suites/GodogTestSuite.xml" with body:
    """
       {
  "name": "GodogTestSuite",
  "testStrategy": "SuiteBased",
  "testCases": [
    {
       "name":"file1.xml",
       "preThinkTime": 1000,
       "postThinkTime": 2000,
       "execWeight": "Infrequent"
    },
        {
          "name":"file2.xml",
          "preThinkTime": 1,
          "postThinkTime": 10,
          "execWeight": "Sparse"
        }
    ]
  }
       """
    Then the response code should be 204
    And the response body should be empty
    When I send a "GET" request to "/test-suites/GodogTestSuite.xml"
    And the updated file should match json:
     """
       {
  "name": "GodogTestSuite",
  "testStrategy": "SuiteBased",
  "testCases": [
    {
       "name":"file1.xml",
       "preThinkTime": 1000,
       "postThinkTime": 2000,
       "execWeight": "Infrequent"
    },
        {
          "name":"file2.xml",
          "preThinkTime": 1,
          "postThinkTime": 10,
          "execWeight": "Sparse"
        }
    ]
  }
       """

  Scenario: Successful update of test-suite file with PUT request (Update API Name to not match Filename)
    Given the config file "GodogTestSuite.xml" exists at "/uiServices/test/"
    Given the automated performance ui server is available
    And the header configsDirPath is "/uiServices/test/"
    When I send "PUT" request to "/test-suites/GodogTestSuite.xml" with body:
        """
       {
  "name": "GodogApi",
  "testStrategy": "SuiteBased",
  "testCases": [
    {
       "name":"file1.xml",
       "preThinkTime": 1000,
       "postThinkTime": 2000,
       "execWeight": "Infrequent"
    },
        {
          "name":"file2.xml",
          "preThinkTime": 1,
          "postThinkTime": 10,
          "execWeight": "Sparse"
        }
    ]
  }
       """
    Then the response code should be 204
    And the response body should be empty
    When I send a "GET" request to "/test-suites/GodogTestSuite.xml"
    And the updated file should match json:
        """
       {
  "name": "GodogApi",
  "testStrategy": "SuiteBased",
  "testCases": [
    {
       "name":"file1.xml",
       "preThinkTime": 1000,
       "postThinkTime": 2000,
       "execWeight": "Infrequent"
    },
        {
          "name":"file2.xml",
          "preThinkTime": 1,
          "postThinkTime": 10,
          "execWeight": "Sparse"
        }
    ]
  }
       """