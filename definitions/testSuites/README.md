This directory stored definitions of Test Suite. Test suites pull together a selection of test cases to be used in system load testing. 
Test Suite configuration can be defined in XML of TOML specification. 
Below is a example of a test suite definition.

        <testSuite>
            <!--Test suite name-->
            <name>Xtrac API QA</name>
             <!--Test Strategy to be used for this test - (ServiceBased / SuiteBased)-->
            <testStrategy>SuiteBased</testStrategy>
            <!--A list of predefined test case to be executed as part of this suite. testCase element should be populated 
            with the name of the test case definition file.-->
            <testCases>
                <testCase>testCase-definition.xml</testCase>
            </testCases>
        </testSuite>