package testStrategies

import (
	"bytes"
	"encoding/xml"
	"github.com/BurntSushi/toml"
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
    <!--BaseURi of the request, excluding host and port. Path parameters if any should be placed here-->
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

	tomlTestDefinition = `testName = "createSession"
httpMethod = "POST"
baseUri = "/xtrac/api/v1/sessions"
multipart = false
payload = "\n        {\n        \"credentials\": {\n        \"operatorName\": \"DRSCAN\",\n        \"password\": \"Tester01\"\n        }\n        }\n    "
responseStatusCode = 200

[[responseProperties]]
	extractionKey = "extractionKey"
	value = "sessionToken"

[headers]
  Scenario = ["success1"]
  Xtractoken = ["509760429188261892213816064522998760"]`

	tomlMultipartTestDefinition = `testName = "XiwsAddContentSuccess"
httpMethod = "POST"
baseUri = "/xiws/camel/ciws/addContent"
multipart = true
payload = ""
responseStatusCode = 200

[[multipartFormField]]
  fieldName = "sessionToken"
  fieldValue = "whateverTokenValue"
  fileName = ""

[[multipartFormField]]
  fieldName = "folderpath"
  fieldValue = "Customer/XWIS"
  fileName = ""

[[multipartFormField]]
  fieldName = "workitemnumber"
  fieldValue = "whatever value"
  fileName = ""

[[multipartFormField]]
  fieldName = "metadatafields"
  fieldValue = "<ns3:metaDataDescriptionList xmlns=\"http://webservices.xtrac.com/elements\"\n                         xmlns:elem=\"http://webservices.xtrac.com/elements\"\n                         xmlns:ns3=\"http://webservices.xtrac.com/types/document\">\n                            <ns3:metaDataDescription>\n                                <elem:name>title</elem:name>\n                                <elem:value>SOAPUI_UPLOAD</elem:value>\n                            </ns3:metaDataDescription>\n                            <ns3:metaDataDescription>\n                                <elem:name>document_type</elem:name>\n                                <elem:value></elem:value>\n                            </ns3:metaDataDescription>\n                            <ns3:metaDataDescription>\n                                <elem:name>document_sub_type</elem:name>\n                                <elem:value></elem:value>\n                            </ns3:metaDataDescription>\n                        </ns3:metaDataDescriptionList>"
  fileName = ""

[[multipartFormField]]
  fieldName = "document"
  fieldValue = ""
  fileName = "whateverName"
  fileContent = [119, 104, 97, 116, 101, 118, 101, 114, 32, 99, 111, 110, 116, 101, 110, 116]

[headers]
  Scenario = ["addcontent-valid"]`

	xmlTestSuite = `<testSuite>
    <name>testSuite</name>
    <testStrategy>ServiceBased</testStrategy>
    <testCases>
        <testCase>xiws-loginLTPA-success.xml</testCase>
        <testCase>xiws-workitem-create-success.xml</testCase>
        <testCase>xiws-workitem-search-success.xml</testCase>
    </testCases>
</testSuite>`

	tomlTestSuite = `name = "testSuite"
testStrategy = "SuiteBased"
testCases = ["xiws-loginLTPA-success.toml", "xiws-workitem-create-success.toml", "xiws-workitem-search-success.toml"]`
)

func TestMarshalTomlTestDefinition(t *testing.T) {
	td := &XmlTestDefinition{}
	err := xml.Unmarshal([]byte(xmlTestDefinition), td)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(td.Headers))
	assert.Equal(t, 1, len(td.ResponseValues))
	assert.Equal(t, "9192", td.OverridePort)
	t.Logf("%+v\n", *td)

	//toml
	buf := new(bytes.Buffer)
	err = toml.NewEncoder(buf).Encode(td)
	assert.Nil(t, err)
	t.Logf("%s\n", buf.Bytes())
}

func TestMarshalTomlMultipartTestDefinition(t *testing.T) {
	td := &XmlTestDefinition{}
	err := xml.Unmarshal([]byte(multipartXmlTestDefinition), td)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(td.Headers))
	t.Logf("%+v\n", *td)

	//toml
	buf := new(bytes.Buffer)
	err = toml.NewEncoder(buf).Encode(td)
	assert.Nil(t, err)
	t.Logf("%s\n", buf.Bytes())
}

func TestUnmarshalTomlTestDefinition(t *testing.T) {
	td := &TestDefinition{}
	err := toml.Unmarshal([]byte(tomlTestDefinition), td)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(td.ResponseValues))
	t.Logf("%+v\n", td)
}

func TestUnmarshalTomlMultipartTestDefinition(t *testing.T) {
	td := &TestDefinition{}
	err := toml.Unmarshal([]byte(tomlMultipartTestDefinition), td)
	assert.Nil(t, err)
	t.Logf("%+v\n", td)
}

func TestUnmarshalXmlTestSuite(t *testing.T) {
	ts := &TestSuiteDefinition{}
	err := xml.Unmarshal([]byte(xmlTestSuite), ts)
	assert.Nil(t, err)
	assert.Equal(t, 3, len(ts.TestCases))
	t.Logf("%+v\n", *ts)

	//toml
	buf := new(bytes.Buffer)
	err = toml.NewEncoder(buf).Encode(ts)
	assert.Nil(t, err)
	t.Logf("%s\n", buf.Bytes())
}

func TestTomlHeaders(t *testing.T) {
	headers := []Header{Header{Value: "val", Key: "key"}, Header{Value: "hello", Key: "header"}}
	h := tomlHeaders(headers)
	assert.NotNil(t, h)
	assert.Equal(t, 2, len(h))
	assert.Equal(t, "val", h.Get("key"))
	assert.Equal(t, "hello", h.Get("header"))
}

func TestLoadTestDefinitionXml(t *testing.T) {
	config := &perfTestUtils.Config{}
	config.SetDefaults()
	td, err := loadTestDefinition([]byte(xmlTestDefinition), config)
	assert.Nil(t, err)
	assert.NotNil(t, td)
	assert.Equal(t, "XiwsLoginLTPASuccess", td.TestName)
	assert.Equal(t, "POST", td.HttpMethod)
	assert.Equal(t, "/xiws/webservices/LoginLTPA", td.BaseUri)
	assert.Equal(t, 200, td.ResponseStatusCode)
	assert.Equal(t, "ltpa-success1", td.Headers.Get("scenario"))
	assert.Equal(t, "sessionToken", td.ResponseValues[0].Value)
	assert.True(t, strings.Contains(td.Payload, "<elem:operatorName>drscan</elem:operatorName>"))
}

func TestLoadTestDefinitionXmlErr(t *testing.T) {
	config := &perfTestUtils.Config{}
	config.SetDefaults()
	td, err := loadTestDefinition([]byte(tomlMultipartTestDefinition), config)
	assert.NotNil(t, err)
	assert.Nil(t, td)
}

func TestLoadTestDefinitionToml(t *testing.T) {
	config := &perfTestUtils.Config{}
	config.SetDefaults()
	config.TestFileFormat = "toml"
	td, err := loadTestDefinition([]byte(tomlTestDefinition), config)
	assert.Nil(t, err)
	assert.NotNil(t, td)
	assert.Equal(t, "createSession", td.TestName)
	assert.Equal(t, "POST", td.HttpMethod)
	assert.Equal(t, "/xtrac/api/v1/sessions", td.BaseUri)
	assert.Equal(t, 200, td.ResponseStatusCode)
	assert.Equal(t, "success1", td.Headers.Get("scenario"))
	assert.Equal(t, "509760429188261892213816064522998760", td.Headers.Get("Xtractoken"))
	assert.Equal(t, "sessionToken", td.ResponseValues[0].Value)
	assert.True(t, strings.Contains(td.Payload, "operatorName"))
}

func TestLoadTestDefinitionTomlErr(t *testing.T) {
	config := &perfTestUtils.Config{}
	config.SetDefaults()
	config.TestFileFormat = "toml"
	td, err := loadTestDefinition([]byte(xmlTestDefinition), config)
	assert.NotNil(t, err)
	assert.Nil(t, td)
}

func TestLoadTestSuiteDefinitionXml(t *testing.T) {
	config := &perfTestUtils.Config{}
	config.SetDefaults()
	ts, err := loadTestSuiteDefinition([]byte(xmlTestSuite), config)
	assert.Nil(t, err)
	assert.NotNil(t, ts)
	assert.Equal(t, "testSuite", ts.Name)
	assert.Equal(t, "ServiceBased", ts.TestStrategy)
	assert.Equal(t, 3, len(ts.TestCases))
	assert.Equal(t, "xiws-loginLTPA-success.xml", ts.TestCases[0])
	assert.Equal(t, "xiws-workitem-create-success.xml", ts.TestCases[1])
	assert.Equal(t, "xiws-workitem-search-success.xml", ts.TestCases[2])
}

func TestLoadTestSuiteDefinitionXmlErr(t *testing.T) {
	config := &perfTestUtils.Config{}
	config.SetDefaults()
	ts, err := loadTestSuiteDefinition([]byte(tomlTestSuite), config)
	assert.Nil(t, ts)
	assert.NotNil(t, err)
}

func TestLoadTestSuiteDefinitionToml(t *testing.T) {
	config := &perfTestUtils.Config{}
	config.SetDefaults()
	config.TestFileFormat = "toml"
	ts, err := loadTestSuiteDefinition([]byte(tomlTestSuite), config)
	assert.Nil(t, err)
	assert.NotNil(t, ts)
	assert.Equal(t, "testSuite", ts.Name)
	assert.Equal(t, "SuiteBased", ts.TestStrategy)
	assert.Equal(t, 3, len(ts.TestCases))
	assert.Equal(t, "xiws-loginLTPA-success.toml", ts.TestCases[0])
	assert.Equal(t, "xiws-workitem-create-success.toml", ts.TestCases[1])
	assert.Equal(t, "xiws-workitem-search-success.toml", ts.TestCases[2])
}

func TestLoadTestSuiteDefinitionTomlErr(t *testing.T) {
	config := &perfTestUtils.Config{}
	config.SetDefaults()
	config.TestFileFormat = "toml"
	ts, err := loadTestSuiteDefinition([]byte(xmlTestSuite), config)
	assert.Nil(t, ts)
	assert.NotNil(t, err)
}
