## Automated performance test framework

Automated performance test framework allows API developers to test the performance of their APIs os part of a CI / CD pipeline. The frame work is intended to be run after code is committed and validates that performance has not 
degraded since the last commit. It analisizes the memory footprint of the service response times of the API service and reports on success and failure scenarios. 

Featrues include:
* ** Training mode ** This mode executes all defined test cases and outputs the results to a json file. This json file represents the base memory profile of the API in relation to the target environment.
* ** Testing mode ** This mode executes all defined test cases and compares the results with the pre-generated base profile.
* ** Peak memory analysis. ** If the peak memory usage during testing mode exceeds an allowed threshold (configurable) the test will fail.
* ** Servcie response time analysis ** If the response time of any service test case exceeds an allowed threshold (configurable) the test will fail.
* ** Simulated concurrent users. ** The system can be configured to allow multiple concurrent hitting the API at once, dividing teh test load across these users.
* ** Report generation. ** In both successful and failed test runs, a report will be generated indicating the test finding and trageting ares where performance issue have been found.
* ** Easily defined test cases. ** To define a new test case, all that is required is a request definition which describes the service request.
* ** Integration into Build Pipeline. Can be configured to run as part of a build pipeline to validate performance as code is checked in. 


### Usage 
To use the framework, users should provide a configuration file which define the setting specific to the API under test. A sample configuration file can be in the config directory.
The configuration file parameters are described in the table below. 

| Property                             	|                                                                                                                                                        Description 	|
|--------------------------------------	|-------------------------------------------------------------------------------------------------------------------------------------------------------------------:	|
| apiName                              	| This property is used to provide a name to the API under test. The property is used for report generation.                                                         	|
| targetHost                           	| Target host of the API under test                                                                                                                                  	|
| targetPort                           	| Target port of the API under test                                                                                                                                  	|
| allowablePeakMemoryVariance          	| This is the percentage by which the peak memory can exceed during the test without the test without the test being considered a failure scenario                   	|
| allowableServiceResponseTimeVariance 	| This is the percentage by which a service test case response time can exceed during the test without the test without the test being considered a failure scenario 	|
| numIterations                        	| This property represents the number a request that will be executed for each test case.                                                                            	|
| concurrentUsers                      	| This property allow for the API to be tested with multiple concurrent users. The load specified above will be equally distributed across these users.              	|
| testDefinitionsDir                   	| This is the directory location for test cases                                                                                                                      	|
| baseStatsOutputDir                   	| This is the directory location of the output file for the base performance statistics json file.                                                                   	|
| reportOutputDir                      	| This is the directory location of the output report file (HTML)                                                                                                    	|

