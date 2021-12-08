package domain

import (
	"net/http"
	"snykctl/internal/tools"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Project_Get_httpError(t *testing.T) {
	client := tools.NewMockClient()
	client.ResponseBody = ""
	client.StatusCode = http.StatusUnauthorized
	client.Status = "XXX"
	prjs := NewProjects(client, "xxx")

	err := prjs.Get()
	expectedErrorMsg := "GetProjects failed: XXX"
	assert.EqualErrorf(t, err, expectedErrorMsg, "Error should be: %v, got: %v", expectedErrorMsg, err)
	assert.Equal(t, false, prjs.IsSync())
}

func Test_Project_Get_badBody(t *testing.T) {
	client := tools.NewMockClient()
	client.ResponseBody = "filler"
	client.StatusCode = http.StatusOK
	client.Status = "XXX"
	prjs := NewProjects(client, "xxx")

	err := prjs.Get()
	expectedErrorMsg := "GetProjects failed:"
	assert.Containsf(t, err.Error(), expectedErrorMsg, "Error should be: %v, got: %v", expectedErrorMsg, err)
	assert.Equal(t, false, prjs.IsSync())
}

func Test_Project_Get_Ok(t *testing.T) {
	client := tools.NewMockClient()
	client.ResponseBody = `{"org":{"name":"Sandbox","id":"16df2e12-d4cb-400d-aaf2-547db9ff07e9"},"projects":[{"id":"5c8e7160-5b60-4f49-824f-c01c891ea29f","name":"com.symphony.is:zoom-frontend","created":"2021-11-22T10:03:05.435Z","origin":"cli","type":"maven","readOnly":false,"testFrequency":"daily","isMonitored":true,"totalDependencies":0,"issueCountsBySeverity":{"low":0,"high":0,"medium":0,"critical":0},"remoteRepoUrl":"http://github.com/SymphonyOSF/sms-zoom.git","imageTag":"1.0.0-SNAPSHOT","lastTestedDate":"2021-12-05T06:20:50.043Z","browseUrl":"https://app.snyk.io/org/sandbox-pie/project/5c8e7160-5b60-4f49-824f-c01c891ea29f","owner":null,"importingUser":{"id":"7261cefe-93f4-472d-b6cd-27d8f41fcc29","name":"sandbox","username":"sandbox","email":null},"tags":[],"attributes":{"criticality":[],"lifecycle":[],"environment":[]},"branch":null},{"id":"9931b808-9f92-4283-a8aa-d96289e363bd","name":"com.symphony:cmd-mock-conveyor","created":"2021-11-19T13:45:23.001Z","origin":"cli","type":"maven","readOnly":false,"testFrequency":"daily","isMonitored":true,"totalDependencies":152,"issueCountsBySeverity":{"low":0,"high":3,"medium":4,"critical":0},"remoteRepoUrl":"http://github.com/SymphonyOSF/command-middleware.git","imageTag":"0.0.1-SNAPSHOT","lastTestedDate":"2021-12-05T05:43:10.379Z","browseUrl":"https://app.snyk.io/org/sandbox-pie/project/9931b808-9f92-4283-a8aa-d96289e363bd","owner":null,"importingUser":{"id":"7261cefe-93f4-472d-b6cd-27d8f41fcc29","name":"sandbox","username":"sandbox","email":null},"tags":[],"attributes":{"criticality":[],"lifecycle":[],"environment":[]},"branch":null}]}"`
	client.StatusCode = http.StatusOK
	client.Status = "XXX"
	prjs := NewProjects(client, "org_id")

	err := prjs.Get()

	assert.Equal(t, nil, err)
	assert.Equal(t, true, prjs.IsSync())

	assert.Equal(t, 2, len(prjs.Projects))
	// id is read from json output
	assert.Equal(t, "16df2e12-d4cb-400d-aaf2-547db9ff07e9", prjs.Org.Id)

	var idFound bool
	for _, o := range prjs.Projects {
		if o.Id == "5c8e7160-5b60-4f49-824f-c01c891ea29f" {
			idFound = true
			assert.Equal(t, "com.symphony.is:zoom-frontend", o.Name)
		}
	}

	assert.True(t, idFound)
}

func Test_Project_Get_Ids(t *testing.T) {
	client := tools.NewMockClient()
	client.ResponseBody = `{"org":{"name":"Sandbox","id":"16df2e12-d4cb-400d-aaf2-547db9ff07e9"},"projects":[{"id":"5c8e7160-5b60-4f49-824f-c01c891ea29f","name":"com.symphony.is:zoom-frontend","created":"2021-11-22T10:03:05.435Z","origin":"cli","type":"maven","readOnly":false,"testFrequency":"daily","isMonitored":true,"totalDependencies":0,"issueCountsBySeverity":{"low":0,"high":0,"medium":0,"critical":0},"remoteRepoUrl":"http://github.com/SymphonyOSF/sms-zoom.git","imageTag":"1.0.0-SNAPSHOT","lastTestedDate":"2021-12-05T06:20:50.043Z","browseUrl":"https://app.snyk.io/org/sandbox-pie/project/5c8e7160-5b60-4f49-824f-c01c891ea29f","owner":null,"importingUser":{"id":"7261cefe-93f4-472d-b6cd-27d8f41fcc29","name":"sandbox","username":"sandbox","email":null},"tags":[],"attributes":{"criticality":[],"lifecycle":[],"environment":[]},"branch":null},{"id":"9931b808-9f92-4283-a8aa-d96289e363bd","name":"com.symphony:cmd-mock-conveyor","created":"2021-11-19T13:45:23.001Z","origin":"cli","type":"maven","readOnly":false,"testFrequency":"daily","isMonitored":true,"totalDependencies":152,"issueCountsBySeverity":{"low":0,"high":3,"medium":4,"critical":0},"remoteRepoUrl":"http://github.com/SymphonyOSF/command-middleware.git","imageTag":"0.0.1-SNAPSHOT","lastTestedDate":"2021-12-05T05:43:10.379Z","browseUrl":"https://app.snyk.io/org/sandbox-pie/project/9931b808-9f92-4283-a8aa-d96289e363bd","owner":null,"importingUser":{"id":"7261cefe-93f4-472d-b6cd-27d8f41fcc29","name":"sandbox","username":"sandbox","email":null},"tags":[],"attributes":{"criticality":[],"lifecycle":[],"environment":[]},"branch":null}]}"`
	client.StatusCode = http.StatusOK
	client.Status = "XXX"
	prjs := NewProjects(client, "org_id")

	err := prjs.Get()

	assert.Equal(t, nil, err)
	assert.Equal(t, true, prjs.IsSync())

	assert.Equal(t, 2, len(prjs.Projects))
	// id is read from json output
	assert.Equal(t, "16df2e12-d4cb-400d-aaf2-547db9ff07e9", prjs.Org.Id)

	expected := "5c8e7160-5b60-4f49-824f-c01c891ea29f\n9931b808-9f92-4283-a8aa-d96289e363bd\n"
	actual, err := prjs.Quiet()

	assert.Nil(t, err)
	assert.Equal(t, expected, actual)

}

func Test_Project_Get_Names(t *testing.T) {
	client := tools.NewMockClient()
	client.ResponseBody = `{"org":{"name":"Sandbox","id":"16df2e12-d4cb-400d-aaf2-547db9ff07e9"},"projects":[{"id":"5c8e7160-5b60-4f49-824f-c01c891ea29f","name":"com.symphony.is:zoom-frontend","created":"2021-11-22T10:03:05.435Z","origin":"cli","type":"maven","readOnly":false,"testFrequency":"daily","isMonitored":true,"totalDependencies":0,"issueCountsBySeverity":{"low":0,"high":0,"medium":0,"critical":0},"remoteRepoUrl":"http://github.com/SymphonyOSF/sms-zoom.git","imageTag":"1.0.0-SNAPSHOT","lastTestedDate":"2021-12-05T06:20:50.043Z","browseUrl":"https://app.snyk.io/org/sandbox-pie/project/5c8e7160-5b60-4f49-824f-c01c891ea29f","owner":null,"importingUser":{"id":"7261cefe-93f4-472d-b6cd-27d8f41fcc29","name":"sandbox","username":"sandbox","email":null},"tags":[],"attributes":{"criticality":[],"lifecycle":[],"environment":[]},"branch":null},{"id":"9931b808-9f92-4283-a8aa-d96289e363bd","name":"com.symphony:cmd-mock-conveyor","created":"2021-11-19T13:45:23.001Z","origin":"cli","type":"maven","readOnly":false,"testFrequency":"daily","isMonitored":true,"totalDependencies":152,"issueCountsBySeverity":{"low":0,"high":3,"medium":4,"critical":0},"remoteRepoUrl":"http://github.com/SymphonyOSF/command-middleware.git","imageTag":"0.0.1-SNAPSHOT","lastTestedDate":"2021-12-05T05:43:10.379Z","browseUrl":"https://app.snyk.io/org/sandbox-pie/project/9931b808-9f92-4283-a8aa-d96289e363bd","owner":null,"importingUser":{"id":"7261cefe-93f4-472d-b6cd-27d8f41fcc29","name":"sandbox","username":"sandbox","email":null},"tags":[],"attributes":{"criticality":[],"lifecycle":[],"environment":[]},"branch":null}]}"`
	client.StatusCode = http.StatusOK
	client.Status = "XXX"
	prjs := NewProjects(client, "org_id")

	err := prjs.Get()

	assert.Equal(t, nil, err)
	assert.Equal(t, true, prjs.IsSync())

	assert.Equal(t, 2, len(prjs.Projects))
	// id is read from json output
	assert.Equal(t, "16df2e12-d4cb-400d-aaf2-547db9ff07e9", prjs.Org.Id)

	expected := "com.symphony.is:zoom-frontend\ncom.symphony:cmd-mock-conveyor\n"
	actual, err := prjs.Names()

	assert.Nil(t, err)
	assert.Equal(t, expected, actual)

}

func Test_Project_Get_Raw(t *testing.T) {
	client := tools.NewMockClient()
	raw := `{"org":{"name":"Sandbox","id":"16df2e12-d4cb-400d-aaf2-547db9ff07e9"},"projects":[{"id":"5c8e7160-5b60-4f49-824f-c01c891ea29f","name":"com.symphony.is:zoom-frontend","created":"2021-11-22T10:03:05.435Z","origin":"cli","type":"maven","readOnly":false,"testFrequency":"daily","isMonitored":true,"totalDependencies":0,"issueCountsBySeverity":{"low":0,"high":0,"medium":0,"critical":0},"remoteRepoUrl":"http://github.com/SymphonyOSF/sms-zoom.git","imageTag":"1.0.0-SNAPSHOT","lastTestedDate":"2021-12-05T06:20:50.043Z","browseUrl":"https://app.snyk.io/org/sandbox-pie/project/5c8e7160-5b60-4f49-824f-c01c891ea29f","owner":null,"importingUser":{"id":"7261cefe-93f4-472d-b6cd-27d8f41fcc29","name":"sandbox","username":"sandbox","email":null},"tags":[],"attributes":{"criticality":[],"lifecycle":[],"environment":[]},"branch":null},{"id":"9931b808-9f92-4283-a8aa-d96289e363bd","name":"com.symphony:cmd-mock-conveyor","created":"2021-11-19T13:45:23.001Z","origin":"cli","type":"maven","readOnly":false,"testFrequency":"daily","isMonitored":true,"totalDependencies":152,"issueCountsBySeverity":{"low":0,"high":3,"medium":4,"critical":0},"remoteRepoUrl":"http://github.com/SymphonyOSF/command-middleware.git","imageTag":"0.0.1-SNAPSHOT","lastTestedDate":"2021-12-05T05:43:10.379Z","browseUrl":"https://app.snyk.io/org/sandbox-pie/project/9931b808-9f92-4283-a8aa-d96289e363bd","owner":null,"importingUser":{"id":"7261cefe-93f4-472d-b6cd-27d8f41fcc29","name":"sandbox","username":"sandbox","email":null},"tags":[],"attributes":{"criticality":[],"lifecycle":[],"environment":[]},"branch":null}]}"`
	client.ResponseBody = raw
	client.StatusCode = http.StatusOK
	client.Status = "XXX"
	prjs := NewProjects(client, "org_id")

	actual, err := prjs.GetRaw()

	assert.Equal(t, nil, err)
	assert.Equal(t, true, prjs.IsSync())

	assert.Equal(t, raw, actual)
}

func Test_AddTag_OK(t *testing.T) {
	client := tools.NewMockClient()
	client.ResponseBody = ``
	client.StatusCode = http.StatusOK
	client.Status = "XXX"
	prjs := NewProjects(client, "org_id")
	err := prjs.AddTag("org2", "k=v")
	assert.Nil(t, err)
}

func Test_AddTag_parseFailed(t *testing.T) {
	client := tools.NewMockClient()
	client.ResponseBody = ``
	client.StatusCode = http.StatusOK
	client.Status = "XXX"
	prjs := NewProjects(client, "org_id")
	err := prjs.AddTag("org2", "vvv")
	expectedErrorMsg := "invalid tag. Not a key=value format"
	assert.EqualErrorf(t, err, expectedErrorMsg, "Error should be: %v, got: %v", expectedErrorMsg, err)
}

func Test_AddTag_KO(t *testing.T) {
	client := tools.NewMockClient()
	client.ResponseBody = ``
	client.StatusCode = http.StatusUnauthorized
	client.Status = "XXX"
	prjs := NewProjects(client, "org_id")
	err := prjs.AddTag("org2", "k=v")
	expectedErrorMsg := "failed to add tag XXX"
	assert.EqualErrorf(t, err, expectedErrorMsg, "Error should be: %v, got: %v", expectedErrorMsg, err)
}

type testDataAttributes struct {
	env         string
	lifecycle   string
	criticality string
	errorMsg    string
}

func Test_ParseAttributes(t *testing.T) {
	tests := []testDataAttributes{
		testDataAttributes{env: "frontend", lifecycle: "", criticality: "", errorMsg: ""},
		testDataAttributes{env: "xxx", lifecycle: "", criticality: "", errorMsg: "invalid environment value: xxx\nValid values: [frontend backend internal external mobile saas on-prem hosted distributed]"},
		testDataAttributes{env: "", lifecycle: "production", criticality: "", errorMsg: ""},
		testDataAttributes{env: "", lifecycle: "xxx", criticality: "", errorMsg: "invalid lifecycle value: xxx\nValid values: [production development sandbox]"},
		testDataAttributes{env: "", lifecycle: "", criticality: "high", errorMsg: ""},
		testDataAttributes{env: "", lifecycle: "", criticality: "xxx", errorMsg: "invalid lifecycle value: xxx\nValid values: [critical high medium low]"},
		testDataAttributes{env: "frontend", lifecycle: "production", criticality: "medium", errorMsg: ""},
		testDataAttributes{env: "xxx", lifecycle: "xxx", criticality: "xxx", errorMsg: "invalid environment value: xxx\nValid values: [frontend backend internal external mobile saas on-prem hosted distributed]"},
	}
	// var err error
	for _, test := range tests {
		err := ParseAttributes(test.env, test.lifecycle, test.criticality)
		if err == nil {
			assert.Equal(t, test.errorMsg, "")
		} else {
			assert.EqualErrorf(t, err, test.errorMsg, "Error should be: %v, got: %v", test.errorMsg, err)
		}
	}
}

func Test_ParseTags(t *testing.T) {
	// one ok
	tag := []string{"key=value"}
	pTags, err := ParseTags(tag)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(pTags))
	assert.Equal(t, pTags["key"], "value")

	// two ok
	tag = append(tag, "key2=value2")
	pTags, err = ParseTags(tag)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(pTags))
	assert.Equal(t, pTags["key"], "value")
	assert.Equal(t, pTags["key2"], "value2")

	// ko
	tag = []string{"qweasd"}
	_, err = ParseTags(tag)
	expected := "invalid tag. Not a key=value format"
	assert.EqualErrorf(t, err, expected, "Error should be: %v, got: %v", expected, err)
}

type testAttributesBody struct {
	env         string
	lifecycle   string
	criticality string
	body        string
}

func Test_Build_AttributesBody(t *testing.T) {
	tests := []testAttributesBody{
		testAttributesBody{env: "", lifecycle: "", criticality: "", body: ""},
		testAttributesBody{env: "production", lifecycle: "", criticality: "", body: `{ "environment": ["production"] }`},
		testAttributesBody{env: "", lifecycle: "frontend", criticality: "", body: `{ "lifecycle": ["frontend"] }`},
		testAttributesBody{env: "", lifecycle: "", criticality: "high", body: `{ "criticality": ["high"] }`},
		testAttributesBody{env: "production", lifecycle: "frontend", criticality: "", body: `{ "environment": ["production"],"lifecycle": ["frontend"] }`},
		testAttributesBody{env: "production", lifecycle: "frontend", criticality: "medium", body: `{ "environment": ["production"],"lifecycle": ["frontend"],"criticality": ["medium"] }`},
	}

	for _, test := range tests {
		body := BuildAttributesBody(test.env, test.lifecycle, test.criticality)
		assert.Equal(t, test.body, body)
	}
}
