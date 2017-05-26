This directory stored definitions of test cases. Test case configuration can be defined in XML of TOML specification.
Below is a example of a test case definition.

        <testDefinition>
            <!--Test Name, This can be any name of your choosing-->
            <testName></testName>
            <!--Http method assciated with this request-->
            <httpMethod></httpMethod>
            <!--BaseURi of the request, excluding host and port. Path parameters if any should be placed here-->
            <baseUri></baseUri>
            <!--Request body, This can be Json or xml data. XML payload should be wrapped in cdata tags-->
            <payload></payload>
            <!--Indicated to the test, what is the expected http response code. This value is asserted during the test.-->
            <responseStatusCode></responseStatusCode>
            <!--request headers-->
            <!-- keys will be converted in the canonical format of the MIME header key. For example, the canonical key for "accept-encoding" is "Accept-Encoding". -->
            <headers>
                <header key=""></header>
                <header key=""></header>
            </headers>

            <!--
                Values to be extracted from the response for use in future
                requests.
                  - Use JMESPath to query JSON output.
                  - The "extractionKey" directive sets the variable name for
                    use in subsequent testCase entries.
                  - Valid results are single values or arrays. Note: if an
                    array is returned, a random index will be used for the
                    result in the current iteration.
            -->
            <responseProperties>
            	<!-- Set a variable with a single value: -->
                <value extractionKey="itemType">data.itemType</value>
                <!-- Set a variable with a random value from the array: -->
                <value extractionKey="wi_num">data.items[].workItemNumber</value>
            </responseProperties>
        </testDefinition>