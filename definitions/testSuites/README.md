This directory stores definitions of Test Suite. Test suites pull together a selection of test cases to be used in system load testing. Test Suite configuration can be defined in XML of TOML specification.

Below is a example of a test suite definition.

```xml
<testSuite>
    <!--Test suite name-->
    <name>Xtrac API QA</name>
    <!--Test Strategy to be used for this test - (ServiceBased / SuiteBased)-->
    <testStrategy>SuiteBased</testStrategy>
    <!--A list of predefined test case to be executed as part of this suite. testCase element should be populated
    with the name of the test case definition file.-->
    <testCases>
        <!--
            * The attributes "preThinkTime" and "postThinkTime" specify an amount of milliseconds to pause
              before or after executing the testCase. These are independent of and additive to the global
              configuration element <randomDelay>. Leaving out either attribute, as illustrated below, is
              equivalent to setting the value to 0.

            * The "execWeight" attribute controls whether a testCase runs full time or a fraction of the time.
              Default is to execute every iteration. Available settings are "Infrequent" at 20% execution and
              "Sparse" at 8% execution.
        -->
        <testCase preThinkTime="2500" postThinkTime="5000">testCase-definition1.xml</testCase>
        <testCase preThinkTime="3000" execWeight="Infrequent">testCase-definition2.xml</testCase>
        <testCase>testCase-definition2.xml</testCase>
    </testCases>
</testSuite>
```

