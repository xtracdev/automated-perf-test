## Automated performance test framework

Automated performance test framework allows API developers to test the performance of their APIs as part of a CI / CD pipeline. The framework is intended to be run after code is committed and validates that performance has not
degraded since the last commit. It analyzes the memory footprint and service response times of the API services and reports on success and failure scenarios.

Features include:
* **Training mode** This mode executes all defined test cases and outputs the results to a JSON file. This JSON file represents the base memory profile of the API in relation to the target environment.
* **Testing mode** This mode executes all defined test cases and compares the results with the pre-generated base profile.
* **Peak memory analysis.** If the peak memory usage during testing mode exceeds an allowed threshold (configurable) the test will fail.
* **Service response time analysis** If the response time of any service test case exceeds an allowed threshold (configurable) the test will fail.
* **Simulated concurrent users** The system can be configured to allow multiple concurrent hitting the API at once, dividing the test load across these users.
* **Report generation** In both successful and failed test runs, a report will be generated indicating the test finding and targeting areas where performance issue have been found.
* **Easily defined test cases** To define a new test case, all that is required is a request definition which describes the service request.
* **Integration into Build Pipeline** Can be configured to run as part of a build pipeline to validate performance as code is checked in.


### Usage
#### Configuration file
To use the framework, users should provide a configuration file which define the setting specific to the API under test. A sample configuration file can be in the config directory.
The configuration file parameters are described in the table below.

| Property                                | Description                                                                                                                                 |
|-----------------------------------------|:--------------------------------------------------------------------------------------------------------------------------------------------|
| \<apiName>                              | Provide a name to the API under test to be used for report generation.                                                                      |
| \<targetHost>                           | Target host of the API under test.                                                                                                          |
| \<targetPort>                           | Target port of the API under test.                                                                                                          |
| \<allowablePeakMemoryVariance>          | The percentage by which the peak memory can exceed during the test without the test without the test being considered a failure scenario.   |
| \<allowableServiceResponseTimeVariance> | The percentage by which a service test case response time can exceed during the test without the test being considered a failure scenario.  |
| \<numIterations>                        | The number of times each test case will be executed. Controls the run time length of the test run.                                          |
| \<concurrentUsers>                      | Specify the number of threads across which to spread the load. Controls overall TPS.                                                        |
| \<testCasesDir>                         | Directory location for test cases.                                                                                                          |
| \<testSuitesDir>                        | Directory location for test suites.                                                                                                         |
| \<testSuite>                            | Name of the test suite XML file. Remove this property for service-base strategy.                                                            |
| \<baseStatsOutputDir>                   | Directory location of the output file for the base performance statistics JSON file.                                                        |
| \<reportOutputDir>                      | Directory location of the output report file (HTML).                                                                                        |
| \<memoryEndpoint>                       | Override "/debug/vars" as the endpoint call for memory metrics.                                                                             |
| \<requestDelay>                         | Add a random delay between all requests specified in milliseconds.                                                                          |
| \<tpsFreq>                              | Specify the number of seconds between displaying the overall TPS message in log.info output.                                                |
| \<rampUsers>                            | Specify the number of user threads to start in a batch during ramp up. Eg. Start 5 threads every 15 seconds.                                |
| \<rampDelay>                            | Specify number of seconds between starting user threads batched during ramp up.                                                             |

#### Command line arguments
In addition the configuration parameters, command line arguments can the passed in to control specifics of each individual test run. The command line arguments are described in the table below.

| Argument          | Description                                                                                                    |
|-------------------|----------------------------------------------------------------------------------------------------------------|
| -configFilePath   | The location of a configuration file which describes the test behavior for the API under test.                 |
| -configFileFormat | The format of the configuration file, the supported formats are XML and TOML (default XML).                    |
| -gbs              | Generate Base Statistics. Indicates to the framework that this test run should be considered a training run.   |
| -reBaseMemory     | Run a training run which will overwrite the memory statistics only of previous training on the execution host. |
| -reBaseAll        | Run a training run which will overwrite the all statistics of previous training on the execution host.         |
| -reBaseAll        | Run a training run which will overwrite the all statistics of previous training on the execution host.         |
| -testFileFormat   | The format of the test definition files, the supported formats are XML and TOML (default XML).
<<<<<<< HEAD
| -ui               | Runs local server on localhost:9191 to support CRUD on configuration and test definitions in UI mode
=======
| -ui               | Runs local server on localhost:9191 to support CRUD on configuration and test definitions in UI mode.
>>>>>>> f4024b89511db1b26d0815870b1901009f410f5d

#### Testing Strategies
The framework supports two type of testing strategies, ServiceBased and SuiteBased. These testing strategies allow flexibility when performing performance
test under different conditions, for example Build pipeline mock vs Live back end load test.
##### ServiceBased
This is the default testing strategy and will be used if no test suite is defined in the configuration file. In this scenario, all files in the test case dir will
for an informal test suite. Service Based testing focuses on each service  independently of others. Memory and service response time data is gathered during the test and analysis is performed once the test is complete. Service based testing is very
appropriate when used in conjunction with a build pipeline and mock back end. These test should run quickly to ensure fast overall run time of the pipeline. This type of testing
divides the load across concurrent users. Eg. For 1000 iterations per test case with 10 concurrent users, each user will perform 100 requests concurrently per test case.

##### SuiteBased
Suite based testing is designed to simulate real load testing hitting a live back-end. Data can be passed between requests so response data from one request can be used
in the request of another. Memory and service response time data is gathered during the test and analysis is performed once the test is complete.
In suite based testing, the number of iteration controls the number of time the suite is run per concurrent user. Thus adding more concurrent user will increase the
testing load.

### Contributing

To contribute, you must certify you agree with the [Developer Certificate of Origin](http://developercertificate.org/)
by signing your commits via `git -s`. To create a signature, configure your user name and email address in git.
Sign with your real name, do not use pseudonyms or submit anonymous commits.


In terms of workflow:

0. For significant changes or improvement, create an issue before commencing work.
1. Fork the repository, and create a branch for your edits.
2. Add tests that cover your changes, unit tests for smaller changes, acceptance test
for more significant functionality.
3. Run gofmt on each file you change before committing your changes.
4. Run golint on each file you change before committing your changes.
5. Make sure all the tests pass before committing your changes.
6. Commit your changes and issue a pull request.
