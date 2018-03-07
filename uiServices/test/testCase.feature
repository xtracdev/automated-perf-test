Feature: Test Case Scenarios
  As an API user
  I want to be able use various requests for test cases
  So that I can test my application using custom metrics


                                ###################################
                                #######    POST REQUESTS    #######
                                ###################################

  Scenario: Successful creation of Test Case
    Given there is no existing test file "GodogTestCase.xml"
    Given the automated performance ui server is available
    And the header "testCasePathDir" is "/uiServices/test/"
    When I send "POST" request to "/test-cases" with a body:
    """
      {
       "testname":"GodogTestCase",
       "description":"desc",
       "overrideHost":"host",
       "overridePort":"9191",
       "httpMethod":"GET",
       "baseURI": "path/to/URI",
       "multipart":false,
       "payload": "payload",
       "responseStatusCode":200,
       "responseContentType": "JSON" ,
       "preThinkTime": 1000,
       "postThinkTime":2000,
       "execWeight": "Sparse",
       "headers":[{
   	     "key": "Authorization",
         "value" :"Header-Value"
        }],
      "responseValues":[{
         "value":"Res-Value",
         "extractionKey": "Res-Key"
       }],
      "multipartPayload":[{
         "fieldName": "F-Name",
         "fieldValue":"PayloadName",
         "fileName": "file-name"
       }]
      }
    """
    Then the response code should be 201
    And the response body should be empty

  Scenario: Unsuccessful creation of Test Case (file already exists )
    Given the automated performance ui server is available
    And the header "testSuitePathDir" is "/uiServices/test/"
    When I send "POST" request to "/test-cases" with a body:
    """
      {
       "testname":"GodogTestCase",
       "description":"desc",
       "overrideHost":"host",
       "overridePort":"9191",
       "httpMethod":"GET",
       "baseURI": "path/to/URI",
       "multipart":false,
       "payload": "payload",
       "responseStatusCode":200,
       "responseContentType": "JSON" ,
       "preThinkTime": 1000,
       "postThinkTime":2000,
       "execWeight": "Sparse",
       "headers":[{
   	     "key": "Authorization",
         "value" :"Header-Value"
        }],
      "responseValues":[{
         "value":"Res-Value",
         "extractionKey": "Res-Key"
       }],
      "multipartPayload":[{
         "fieldName": "F-Name",
         "fieldValue":"PayloadName",
         "fileName": "file-name"
       }]
      }
    """
    Then the response code should be 400


  Scenario: Unsuccessful creation of test Case ( Missing Required Fields )
    Given the automated performance ui server is available
    And the header "testSuitePathDir" is "/uiServices/test/"
    When I send "POST" request to "/test-suites" with a body:
       """
      {
       "testname":"GodogTestCase",
       "description":"desc",
       "overrideHost":"host",
       "overridePort":"9191",
       "HttpMethod":"GET",
       "BaseURI": "path/to/URI",
       "multipart":false,
       "payload": "payload",
       "responseStatusCode":200,
       "responseContentType": "JSON" ,
       "preThinkTime": 1000,
       "postThinkTime":2000,
       "execWeight": "Sparse",
       "Headers":[{
   	     "Key": "Authorization",
         "Value" :"Header-Value"
        }],
      "ResponseValues":[{
         "Value":"Res-Value",
         "ExtractionKey": "Res-Key"
       }],
      "MultipartPayload":[{
         "FieldName": "F-Name",
         "FieldValue":"PayloadName",
         "FileName": "file-name"
       }]
      }
    """
    Then the response code should be 400


  Scenario: Unsuccessful creation of Test Case ( No header defined )
    Given the automated performance ui server is available
    And the header "testCasePathDir" is ""
    When I send "POST" request to "/test-cases" with a body:
       """
      {
       "testname":"GodogTestCase",
       "description":"desc",
       "overrideHost":"host",
       "overridePort":"9191",
       "HttpMethod":"GET",
       "BaseURI": "path/to/URI",
       "multipart":false,
       "payload": "payload",
       "responseStatusCode":200,
       "responseContentType": "JSON" ,
       "preThinkTime": 1000,
       "postThinkTime":2000,
       "execWeight": "Sparse",
       "Headers":[{
   	     "Key": "Authorization",
         "Value" :"Header-Value"
        }],
      "ResponseValues":[{
         "Value":"Res-Value",
         "ExtractionKey": "Res-Key"
       }],
      "MultipartPayload":[{
         "FieldName": "F-Name",
         "FieldValue":"PayloadName",
         "FileName": "file-name"
       }]
      }
    """
    Then the response code should be 400



                                ###################################
                                #######    PUT REQUESTS    #######
                                ###################################

  Scenario: Unsuccessful update of test-case file with PUT request (No File Path)
    Given the file "GodogTestCase.xml" exists at "/uiServices/test/"
    Given the automated performance ui server is available
    And the header "testCasePathDir" is ""
    When I send "PUT" request to "/test-suites/GodogTestSuite" with body:
  """
     {
   "testname":"GodogTestCase",
   "description":"desc",
   "overrideHost":"host",
   "overridePort":"9191",
   "HttpMethod":"GET",
   "BaseURI": "path/to/URI",
   "multipart":false,
   "payload": "payload",
   "responseStatusCode":200,
   "responseContentType": "JSON" ,
   "preThinkTime": 1000,
   "postThinkTime":2000,
   "execWeight": "Sparse",
   "Headers":[{
   	 "Key": "Authorization",
     "Value" :"Header-Value"
   }],
  "ResponseValues":[{
     "Value":"Res-Value",
     "ExtractionKey": "Res-Key"
  }],
  "MultipartPayload":[{
     "fieldName": "F-Name",
   	 "FieldValue":"PayloadName",
     "FileName": "file-name"
  }]

}
    """
    Then the response code should be 400


  Scenario: Unsuccessful update of test-case file with PUT request (Incorrect File Name)
    Given the file "GodogTestCase.xml" exists at "/uiServices/test/"
    Given the automated performance ui server is available
    And the header "testCasePathDir" is "/uiServices/test/"
    When I send "PUT" request to "/test-cases/xxx" with body:
    """
     {
   "testname":"GodogTestCase",
   "description":"desc",
   "overrideHost":"host",
   "overridePort":"9191",
   "HttpMethod":"GET",
   "BaseURI": "path/to/URI",
   "multipart":false,
   "payload": "payload",
   "responseStatusCode":200,
   "responseContentType": "JSON" ,
   "preThinkTime": 1000,
   "postThinkTime":2000,
   "execWeight": "Sparse",
   "Headers":[{
   	 "Key": "Authorization",
     "Value" :"Header-Value"
   }],
  "ResponseValues":[{
     "Value":"Res-Value",
     "ExtractionKey": "Res-Key"
  }],
  "MultipartPayload":[{
     "fieldName": "F-Name",
   	 "FieldValue":"PayloadName",
     "FileName": "file-name"
  }]

}
    """
    Then the response code should be 404


  Scenario: Unsuccessful update of test-case file with PUT request (No File Name)
    Given the file "GodogTestCase.xml" exists at "/uiServices/test/"
    Given the automated performance ui server is available
    And the header "testCasePathDir" is "/uiServices/test/"
    When I send "PUT" request to "/test-suites/" with body:
    """
     {
   "testname":"GodogTestCase",
   "description":"desc",
   "overrideHost":"host",
   "overridePort":"9191",
   "HttpMethod":"GET",
   "BaseURI": "path/to/URI",
   "multipart":false,
   "payload": "payload",
   "responseStatusCode":200,
   "responseContentType": "JSON" ,
   "preThinkTime": 1000,
   "postThinkTime":2000,
   "execWeight": "Sparse",
   "Headers":[{
   	 "Key": "Authorization",
     "Value" :"Header-Value"
   }],
  "ResponseValues":[{
     "Value":"Res-Value",
     "ExtractionKey": "Res-Key"
  }],
  "MultipartPayload":[{
     "fieldName": "F-Name",
   	 "FieldValue":"PayloadName",
     "FileName": "file-name"
  }]

}
    """
    Then the response code should be 404

  Scenario: Unsuccessful update of test-suite file with PUT request (Missing Required Fields)
    Given the file "GodogTestCase.xml" exists at "/uiServices/test/"
    Given the automated performance ui server is available
    And the header "testCasePathDir" is "/uiServices/test/"
    When I send "PUT" request to "/test-suites/GodogTestCase" with body:
    """
     {
   "testname":"GodogTestCase",
   "description":"",
   "overridePort":"",
   "HttpMethod":"",
   "BaseURI": "path/to/URI",
   "multipart":false,
   "payload": "payload",
   "responseStatusCode":200,
   "responseContentType": "JSON" ,
   "preThinkTime": 1000,
   "postThinkTime":2000,
   "execWeight": "Sparse",
   "Headers":[{
   	 "Key": "Authorization",
     "Value" :"Header-Value"
   }],
  "ResponseValues":[{
     "Value":"Res-Value",
     "ExtractionKey": "Res-Key"
  }],
  "MultipartPayload":[{
     "fieldName": "F-Name",
   	 "FieldValue":"PayloadName",
     "FileName": "file-name"
  }]

}
    """
    Then the response code should be 400


  Scenario: Successful update of test-case file with PUT request
    Given the file "GodogTestCase.xml" exists at "/uiServices/test/"
    Given the automated performance ui server is available
    And the header "testCasePathDir" is "/uiServices/test/"
    When I send "PUT" request to "/test-cases/GodogTestCase" with body:
     """
     {
   "testname":"GodogTestCase",
   "description":"desc",
   "overrideHost":"host",
   "overridePort":"1001",
   "HttpMethod":"POST",
   "BaseURI": "path/to/URI",
   "multipart":false,
   "payload": "payload",
   "responseStatusCode":200,
   "responseContentType": "JSON" ,
   "preThinkTime": 1000,
   "postThinkTime":2000,
   "execWeight": "Sparse",
   "Headers":[{
   	 "Key": "Authorization",
     "Value" :"Header-Value"
   }],
  "ResponseValues":[{
     "Value":"Res-Value",
     "ExtractionKey": "Res-Key"
  }],
  "MultipartPayload":[{
     "fieldName": "F-Name",
   	 "FieldValue":"PayloadName",
     "FileName": "file-name"
  }]

}
    """
    Then the response code should be 204
    And the response body should be empty


                                ###################################
                                #######    GET ALL REQUESTS #######
                                ###################################
  Scenario: Successful retrieval all test cases with valid "GET" request
    ##Add additional file first so there are multiple files to GET
    Given there is no existing test file "GodogTestCase2.xml"
    Given the automated performance ui server is available
    And the header "testCasePathDir" is "/uiServices/test/"
    When I send "POST" request to "/test-cases" with a body:
    """
      {
       "testname":"GodogTestCase2",
       "description":"desc2",
       "overrideHost":"host",
       "overridePort":"9191",
       "httpMethod":"PUT",
       "baseURI": "path/to/URI",
       "multipart":false,
       "payload": "payload",
       "responseStatusCode":200,
       "responseContentType": "JSON" ,
       "preThinkTime": 1000,
       "postThinkTime":2000,
       "execWeight": "Sparse",
       "headers":[{
   	     "key": "Authorization",
         "value" :"Header-Value"
        }],
      "responseValues":[{
         "value":"Res-Value",
         "extractionKey": "Res-Key"
       }],
      "multipartPayload":[{
         "fieldName": "F-Name",
         "fieldValue":"PayloadName",
         "fileName": "file-name"
       }]
      }
    """
    Then the response code should be 201
    And the response body should be empty
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
          },
          {
          "name": "GodogTestCase2,
          "description": "Case Desc2",
          "httpMethod": "PUT"
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
    Given the file "GodogTestCase2.xml" exists at "/uiServices/test/"
    Given the automated performance ui server is available
    And the header "testCasePathDir" is "/uiServices/test/"
    And the file name is "Case1.xml"
    When I send a "GET" request to "/test-cases/GodogTestCase2"
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
                                #######    DELETE REQUESTS #######
                                ###################################

  Scenario:  Fail to remove Test Case file with "DELETE" request (File Not Found)
    Given the file "GodogTestCase.xml" exists at "/uiServices/test/"
    Given the automated performance ui server is available
    And the header "testCasePathDir" is "/uiServices/test/"
    When I send a "DELETE" request to "/test-cases/xxxx"
    Then the response code should be 404

  Scenario:  Fail to remove Test Case file with "DELETE" request (No Header)
    Given the file "GodogTestCase.xml" exists at "/uiServices/test/"
    Given the automated performance ui server is available
    And the header "testCasePathDir" is ""
    When I send a "DELETE" request to "/test-cases/xxxx"
    Then the response code should be 400

  Scenario:  Successful removal Test Case file with "DELETE" request
    #create file to delete
    Given there is no existing test file "GodogTestCase3.xml"
    Given the automated performance ui server is available
    And the header "testCasePathDir" is "/uiServices/test/"
    When I send "POST" request to "/test-cases" with a body:
    """
      {
       "testname":"GodogTestCase3",
       "description":"desc",
       "overrideHost":"host",
       "overridePort":"9191",
       "httpMethod":"GET",
       "baseURI": "path/to/URI",
       "multipart":false,
       "payload": "payload",
       "responseStatusCode":200,
       "responseContentType": "JSON" ,
       "preThinkTime": 1000,
       "postThinkTime":2000,
       "execWeight": "Sparse",
       "headers":[{
   	     "key": "Authorization",
         "value" :"Header-Value"
        }],
      "responseValues":[{
         "value":"Res-Value",
         "extractionKey": "Res-Key"
       }],
      "multipartPayload":[{
         "fieldName": "F-Name",
         "fieldValue":"PayloadName",
         "fileName": "file-name"
       }]
      }
    """
    Then the response code should be 201
    #Delete
    Given the file "GodogTestCase3.xml" exists at "/uiServices/test/"
    Given the automated performance ui server is available
    And the header "testCasePathDir" is "/uiServices/test/"
    When I send a "DELETE" request to "/test-cases/GodogTestCase3"
    Then the response code should be 204

