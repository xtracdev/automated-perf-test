package testStrategies

import (
	"encoding/xml"
	"github.com/stretchr/testify/assert"
	"github.com/xtracdev/automated-perf-test/perfTestUtils"
	"strings"
	"testing"
)

const (
	xmlTestDefinition = `<testDefinition>
    <testName>XiwsLoginLTPASuccess</testName>
    <httpMethod>POST</httpMethod>
    <baseUri>/xiws/webservices/LoginLTPA</baseUri>
    <overrideHost>localhost</overrideHost>
    <overridePort>9192</overridePort>
    <payload>
        <![CDATA[<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/"
				   xmlns:log="http://webservices.xtrac.com/loginltpa"
				   xmlns:elem="http://webservices.xtrac.com/elements">
            <soapenv:Header/>
            <soapenv:Body>
                <log:loginltpaRequest>
                    <log:loginLTPA>
                        <elem:operatorName>drscan</elem:operatorName>
                        <elem:password>exteaw0rkb5nch</elem:password>
                    </log:loginLTPA>
                </log:loginltpaRequest>
            </soapenv:Body>
        </soapenv:Envelope>]]>
    </payload>
    <responseStatusCode>200</responseStatusCode>
    <headers>
        <header key="scenario">ltpa-success1</header>
    </headers>
     <responseProperties>
        <value extractionKey="sessionToken">sessionToken</value>
    </responseProperties>
</testDefinition>`

	multipartXmlTestDefinition = `<testDefinition>
    <!--Test Name, This can be any name of your choosing-->
    <testName>XiwsAddContentSuccess</testName>
    <!--Http method assciated with this request-->
    <httpMethod>POST</httpMethod>
    <!--BaseURI of the request, excluding host and port. Path parameters if any should be placed here-->
    <baseUri>/xiws/camel/ciws/addContent</baseUri>
    <!-- Indicate whether the request is multipart or not -->
    <multipart>true</multipart>
    <!--REquest body, This can be Json or xml data. XML payload should be wrapped in cdata tags-->
    <!--payload></payload-->
    <!-- Please uncomment the following structure for muleipart request payload -->
    <multipartPayload>
        <multipartFormField>
            <fieldName>sessionToken</fieldName>
            <fieldValue>whateverTokenValue</fieldValue>
        </multipartFormField>
        <multipartFormField>
            <fieldName>folderpath</fieldName>
            <fieldValue>Customer/XWIS</fieldValue>
        </multipartFormField>
        <multipartFormField>
            <fieldName>workitemnumber</fieldName>
            <fieldValue>whatever value</fieldValue>
        </multipartFormField>
        <multipartFormField>
            <fieldName>metadatafields</fieldName>
            <fieldValue><![CDATA[<ns3:metaDataDescriptionList xmlns="http://webservices.xtrac.com/elements"
                         xmlns:elem="http://webservices.xtrac.com/elements"
                         xmlns:ns3="http://webservices.xtrac.com/types/document">
                            <ns3:metaDataDescription>
                                <elem:name>title</elem:name>
                                <elem:value>SOAPUI_UPLOAD</elem:value>
                            </ns3:metaDataDescription>
                            <ns3:metaDataDescription>
                                <elem:name>document_type</elem:name>
                                <elem:value></elem:value>
                            </ns3:metaDataDescription>
                            <ns3:metaDataDescription>
                                <elem:name>document_sub_type</elem:name>
                                <elem:value></elem:value>
                            </ns3:metaDataDescription>
                        </ns3:metaDataDescriptionList>]]></fieldValue>
        </multipartFormField>
        <multipartFormField>
            <fieldName>document</fieldName>
            <fileName>whateverName</fileName>
            <fileContent>whatever content</fileContent>
        </multipartFormField>
    </multipartPayload>
    <!--Indicated to the test, what is the expected http response code. This value is asserted during the test.-->
    <responseStatusCode>200</responseStatusCode>
    <headers>
        <header key="scenario">addcontent-valid</header>
    </headers>
</testDefinition>`

	xmlTestSuite = `<testSuite>
    <name>testSuite</name>
    <testStrategy>ServiceBased</testStrategy>
    <testCases>
        <testCase preThinkTime="10" postThinkTime="20" execWeight="Infrequent">xiws-loginLTPA-success.xml</testCase>
        <testCase>xiws-workitem-create-success.xml</testCase>
        <testCase>xiws-workitem-search-success.xml</testCase>
    </testCases>
</testSuite>`
)

func TestUnmarshalXmlTestSuite(t *testing.T) {
	ts := &TestSuite{}
	err := xml.Unmarshal([]byte(xmlTestSuite), ts)
	assert.Nil(t, err)
	assert.Equal(t, 3, len(ts.TestCases))
	t.Logf("%+v\n", *ts)
}

func TestLoadTestDefinitionXml(t *testing.T) {
	config := &perfTestUtils.Config{}
	config.SetDefaults()
	td, err := loadTestDefinition([]byte(xmlTestDefinition))
	assert.Nil(t, err)
	assert.NotNil(t, td)
	assert.Equal(t, "XiwsLoginLTPASuccess", td.TestName)
	assert.Equal(t, "POST", td.HTTPMethod)
	assert.Equal(t, "/xiws/webservices/LoginLTPA", td.BaseURI)
	assert.Equal(t, 200, td.ResponseStatusCode)
	assert.Equal(t, "ltpa-success1", td.Headers[0].Value)
	assert.Equal(t, "sessionToken", td.ResponseValues[0].Value)
	assert.True(t, strings.Contains(td.Payload, "<elem:operatorName>drscan</elem:operatorName>"))
}

func TestLoadTestDefinitionXmlErr(t *testing.T) {
	config := &perfTestUtils.Config{}
	config.SetDefaults()
	td, err := loadTestDefinition([]byte(`This is not XML.`))
	assert.NotNil(t, err)
	assert.Nil(t, td)
}

func TestLoadTestSuiteDefinitionXml(t *testing.T) {
	config := &perfTestUtils.Config{}
	config.SetDefaults()
	ts := new(TestSuite)
	err := ts.loadTestSuiteDefinition([]byte(xmlTestSuite))
	assert.Nil(t, err)
	assert.NotNil(t, ts)
	assert.Equal(t, "testSuite", ts.Name)
	assert.Equal(t, "ServiceBased", ts.TestStrategy)
	assert.Equal(t, 3, len(ts.TestCases))
	assert.Equal(t, "xiws-loginLTPA-success.xml", ts.TestCases[0].Name)
	assert.Equal(t, "xiws-workitem-create-success.xml", ts.TestCases[1].Name)
	assert.Equal(t, "xiws-workitem-search-success.xml", ts.TestCases[2].Name)
	assert.Equal(t, int64(10), ts.TestCases[0].PreThinkTime)
	assert.Equal(t, int64(20), ts.TestCases[0].PostThinkTime)
	assert.Equal(t, "Infrequent", ts.TestCases[0].ExecWeight)
	assert.Equal(t, int64(0), ts.TestCases[1].PreThinkTime)
}

func TestLoadTestSuiteDefinitionXmlErr(t *testing.T) {
	config := &perfTestUtils.Config{}
	config.SetDefaults()
	ts := new(TestSuite)
	err := ts.loadTestSuiteDefinition([]byte(`This is not XML.`))
	assert.NotNil(t, err)
}
