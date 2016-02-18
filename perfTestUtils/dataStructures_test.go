package perfTestUtils

import (
	"bytes"
	"encoding/xml"
	"github.com/BurntSushi/toml"
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	xmlTestDefinition = `<testDefinition>
    <testName>XiwsLoginLTPASuccess</testName>
    <httpMethod>POST</httpMethod>
    <baseUri>/xiws/webservices/LoginLTPA</baseUri>
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
)

func TestSetDefaults(t *testing.T) {
	c := &Config{}
	c.SetDefaults()
	assert.Equal(t, defaultAPIName, c.APIName)
	assert.Equal(t, defaultTargetHost, c.TargetHost)
	assert.Equal(t, defaultTargetPort, c.TargetPort)
	assert.Equal(t, defaultNumIterations, c.NumIterations)
	assert.Equal(t, defaultAllowablePeakMemoryVariance, c.AllowablePeakMemoryVariance)
	assert.Equal(t, defaultAllowableServiceResponseTimeVariance, c.AllowableServiceResponseTimeVariance)
	assert.Equal(t, defaultTestDefinitionsDir, c.TestDefinitionsDir)
	assert.Equal(t, defaultBaseStatsOutputDir, c.BaseStatsOutputDir)
	assert.Equal(t, defaultReportOutputDir, c.ReportOutputDir)
	assert.Equal(t, defaultConcurrentUsers, c.ConcurrentUsers)
	assert.Equal(t, false, c.GBS)
	assert.Equal(t, false, c.ReBaseMemory)
	assert.Equal(t, false, c.ReBaseAll)
}

func TestPrintAndValidateConfig(t *testing.T) {
	willCallOsExit := false
	exit := func(i int) { willCallOsExit = true }
	c := &Config{}
	c.SetDefaults()
	c.PrintAndValidateConfig(exit)
	assert.False(t, willCallOsExit)
}

func TestPrintAndValidateConfigErr(t *testing.T) {
	willCallOsExit := false
	exit := func(i int) { willCallOsExit = true }
	c := &Config{}
	c.PrintAndValidateConfig(exit)
	assert.True(t, willCallOsExit)
}

func TestMarshalTomlTestDefinition(t *testing.T) {
	td := &TestDefinition{}
	err := xml.Unmarshal([]byte(xmlTestDefinition), td)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(td.Headers))
	t.Logf("%+v\n", *td)

	//toml
	buf := new(bytes.Buffer)
	err = toml.NewEncoder(buf).Encode(td)
	assert.Nil(t, err)
	t.Logf("%s\n", buf.Bytes())
}

func TestMarshalTomlMultipartTestDefinition(t *testing.T) {
	td := &TestDefinition{}
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
	td := &TomlTestDefinition{}
	err := toml.Unmarshal([]byte(tomlTestDefinition), td)
	assert.Nil(t, err)
	t.Logf("%+v\n", td)
}

func TestUnmarshalTomlMultipartTestDefinition(t *testing.T) {
	td := &TomlTestDefinition{}
	err := toml.Unmarshal([]byte(tomlMultipartTestDefinition), td)
	assert.Nil(t, err)
	t.Logf("%+v\n", td)
}
