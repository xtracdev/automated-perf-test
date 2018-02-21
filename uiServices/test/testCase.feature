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