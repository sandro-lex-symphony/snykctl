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
}

func Test_Project_Get_Ok(t *testing.T) {
	client := tools.NewMockClient()
	client.ResponseBody = `{"org":{"name":"org-test","id":"16df2e12-d4cb-4111-aaf2-547db9ff07e9"},"projects":[{"id":"5c8e7160-5b60-4f49-824f-c01c111ea29f","name":"prj1:front","created":"2021-11-22T10:03:05.435Z","origin":"cli","type":"maven","readOnly":false,"testFrequency":"daily","isMonitored":true,"totalDependencies":0,"issueCountsBySeverity":{"low":0,"high":0,"medium":0,"critical":0},"remoteRepoUrl":"http://example.com/repot/prj1.git","imageTag":"1.0.0-SNAPSHOT","lastTestedDate":"2021-12-05T06:20:50.043Z","browseUrl":"https://app.snyk.io/org/org-test-pie/project/5c8e7160-5b60-4f49-824f-c01c111ea29f","owner":null,"importingUser":{"id":"7261cefe-93f4-472d-b6cd-27d8f41f1111","name":"org-test","username":"org-test","email":null},"tags":[],"attributes":{"criticality":[],"lifecycle":[],"environment":[]},"branch":null},{"id":"9931b808-9f92-4283-a8aa-d96289e11111","name":"com.example:cmd-mock-conxxx","created":"2021-11-19T13:45:23.001Z","origin":"cli","type":"maven","readOnly":false,"testFrequency":"daily","isMonitored":true,"totalDependencies":152,"issueCountsBySeverity":{"low":0,"high":3,"medium":4,"critical":0},"remoteRepoUrl":"http://example.com/repo/prj2.git","imageTag":"0.0.1-SNAPSHOT","lastTestedDate":"2021-12-05T05:43:10.379Z","browseUrl":"https://app.snyk.io/org/org-test-pie/project/9931b808-9f92-4283-a8aa-d96289e11111","owner":null,"importingUser":{"id":"7261cefe-93f4-472d-b6cd-27d8f41f1111","name":"org-test","username":"org-test","email":null},"tags":[],"attributes":{"criticality":[],"lifecycle":[],"environment":[]},"branch":null}]}"`
	client.StatusCode = http.StatusOK
	client.Status = "XXX"
	prjs := NewProjects(client, "org_id")

	err := prjs.Get()

	assert.Equal(t, nil, err)

	assert.Equal(t, 2, len(prjs.Projects))
	// id is read from json output
	assert.Equal(t, "16df2e12-d4cb-4111-aaf2-547db9ff07e9", prjs.Org.Id)

	var idFound bool
	for _, o := range prjs.Projects {
		if o.Id == "5c8e7160-5b60-4f49-824f-c01c111ea29f" {
			idFound = true
			assert.Equal(t, "prj1:front", o.Name)
		}
	}

	assert.True(t, idFound)
}

func Test_Project_Get_Ids(t *testing.T) {
	client := tools.NewMockClient()
	client.ResponseBody = `{"org":{"name":"org-test","id":"16df2e12-d4cb-4111-aaf2-547db9ff07e9"},"projects":[{"id":"5c8e7160-5b60-4f49-824f-c01c111ea29f","name":"prj1:front","created":"2021-11-22T10:03:05.435Z","origin":"cli","type":"maven","readOnly":false,"testFrequency":"daily","isMonitored":true,"totalDependencies":0,"issueCountsBySeverity":{"low":0,"high":0,"medium":0,"critical":0},"remoteRepoUrl":"http://example.com/repot/prj1.git","imageTag":"1.0.0-SNAPSHOT","lastTestedDate":"2021-12-05T06:20:50.043Z","browseUrl":"https://app.snyk.io/org/org-test-pie/project/5c8e7160-5b60-4f49-824f-c01c111ea29f","owner":null,"importingUser":{"id":"7261cefe-93f4-472d-b6cd-27d8f41f1111","name":"org-test","username":"org-test","email":null},"tags":[],"attributes":{"criticality":[],"lifecycle":[],"environment":[]},"branch":null},{"id":"9931b808-9f92-4283-a8aa-d96289e11111","name":"com.example:cmd-mock-conxxx","created":"2021-11-19T13:45:23.001Z","origin":"cli","type":"maven","readOnly":false,"testFrequency":"daily","isMonitored":true,"totalDependencies":152,"issueCountsBySeverity":{"low":0,"high":3,"medium":4,"critical":0},"remoteRepoUrl":"http://example.com/repo/prj2.git","imageTag":"0.0.1-SNAPSHOT","lastTestedDate":"2021-12-05T05:43:10.379Z","browseUrl":"https://app.snyk.io/org/org-test-pie/project/9931b808-9f92-4283-a8aa-d96289e11111","owner":null,"importingUser":{"id":"7261cefe-93f4-472d-b6cd-27d8f41f1111","name":"org-test","username":"org-test","email":null},"tags":[],"attributes":{"criticality":[],"lifecycle":[],"environment":[]},"branch":null}]}"`
	client.StatusCode = http.StatusOK
	client.Status = "XXX"
	prjs := NewProjects(client, "org_id")

	err := prjs.Get()

	assert.Equal(t, nil, err)

	assert.Equal(t, 2, len(prjs.Projects))
	// id is read from json output
	assert.Equal(t, "16df2e12-d4cb-4111-aaf2-547db9ff07e9", prjs.Org.Id)

	expected := "5c8e7160-5b60-4f49-824f-c01c111ea29f\n9931b808-9f92-4283-a8aa-d96289e11111\n"
	actual := prjs.Quiet()

	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}

func Test_Project_Get_Names(t *testing.T) {
	client := tools.NewMockClient()
	client.ResponseBody = `{"org":{"name":"org-test","id":"16df2e12-d4cb-4111-aaf2-547db9ff07e9"},"projects":[{"id":"5c8e7160-5b60-4f49-824f-c01c111ea29f","name":"prj1:front","created":"2021-11-22T10:03:05.435Z","origin":"cli","type":"maven","readOnly":false,"testFrequency":"daily","isMonitored":true,"totalDependencies":0,"issueCountsBySeverity":{"low":0,"high":0,"medium":0,"critical":0},"remoteRepoUrl":"http://example.com/repot/prj1.git","imageTag":"1.0.0-SNAPSHOT","lastTestedDate":"2021-12-05T06:20:50.043Z","browseUrl":"https://app.snyk.io/org/org-test-pie/project/5c8e7160-5b60-4f49-824f-c01c111ea29f","owner":null,"importingUser":{"id":"7261cefe-93f4-472d-b6cd-27d8f41f1111","name":"org-test","username":"org-test","email":null},"tags":[],"attributes":{"criticality":[],"lifecycle":[],"environment":[]},"branch":null},{"id":"9931b808-9f92-4283-a8aa-d96289e11111","name":"com.example:cmd-mock-conxxx","created":"2021-11-19T13:45:23.001Z","origin":"cli","type":"maven","readOnly":false,"testFrequency":"daily","isMonitored":true,"totalDependencies":152,"issueCountsBySeverity":{"low":0,"high":3,"medium":4,"critical":0},"remoteRepoUrl":"http://example.com/repo/prj2.git","imageTag":"0.0.1-SNAPSHOT","lastTestedDate":"2021-12-05T05:43:10.379Z","browseUrl":"https://app.snyk.io/org/org-test-pie/project/9931b808-9f92-4283-a8aa-d96289e11111","owner":null,"importingUser":{"id":"7261cefe-93f4-472d-b6cd-27d8f41f1111","name":"org-test","username":"org-test","email":null},"tags":[],"attributes":{"criticality":[],"lifecycle":[],"environment":[]},"branch":null}]}"`
	client.StatusCode = http.StatusOK
	client.Status = "XXX"
	prjs := NewProjects(client, "org_id")

	err := prjs.Get()

	assert.Equal(t, nil, err)

	assert.Equal(t, 2, len(prjs.Projects))
	// id is read from json output
	assert.Equal(t, "16df2e12-d4cb-4111-aaf2-547db9ff07e9", prjs.Org.Id)

	expected := "prj1:front\ncom.example:cmd-mock-conxxx\n"
	actual := prjs.Names()

	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}

func Test_Project_Get_String(t *testing.T) {
	client := tools.NewMockClient()
	client.ResponseBody = `{"org":{"name":"org-test","id":"16df2e12-d4cb-4111-aaf2-547db9ff07e9"},"projects":[{"id":"5c8e7160-5b60-4f49-824f-c01c111ea29f","name":"prj1:front","created":"2021-11-22T10:03:05.435Z","origin":"cli","type":"maven","readOnly":false,"testFrequency":"daily","isMonitored":true,"totalDependencies":0,"issueCountsBySeverity":{"low":0,"high":0,"medium":0,"critical":0},"remoteRepoUrl":"http://example.com/repot/prj1.git","imageTag":"1.0.0-SNAPSHOT","lastTestedDate":"2021-12-05T06:20:50.043Z","browseUrl":"https://app.snyk.io/org/org-test-pie/project/5c8e7160-5b60-4f49-824f-c01c111ea29f","owner":null,"importingUser":{"id":"7261cefe-93f4-472d-b6cd-27d8f41f1111","name":"org-test","username":"org-test","email":null},"tags":[],"attributes":{"criticality":[],"lifecycle":[],"environment":[]},"branch":null},{"id":"9931b808-9f92-4283-a8aa-d96289e11111","name":"com.example:cmd-mock-conxxx","created":"2021-11-19T13:45:23.001Z","origin":"cli","type":"maven","readOnly":false,"testFrequency":"daily","isMonitored":true,"totalDependencies":152,"issueCountsBySeverity":{"low":0,"high":3,"medium":4,"critical":0},"remoteRepoUrl":"http://example.com/repo/prj2.git","imageTag":"0.0.1-SNAPSHOT","lastTestedDate":"2021-12-05T05:43:10.379Z","browseUrl":"https://app.snyk.io/org/org-test-pie/project/9931b808-9f92-4283-a8aa-d96289e11111","owner":null,"importingUser":{"id":"7261cefe-93f4-472d-b6cd-27d8f41f1111","name":"org-test","username":"org-test","email":null},"tags":[],"attributes":{"criticality":[],"lifecycle":[],"environment":[]},"branch":null}]}"`
	client.StatusCode = http.StatusOK
	client.Status = "XXX"
	prjs := NewProjects(client, "org_id")

	err := prjs.Get()

	assert.Equal(t, nil, err)

	assert.Equal(t, 2, len(prjs.Projects))
	// id is read from json output
	assert.Equal(t, "16df2e12-d4cb-4111-aaf2-547db9ff07e9", prjs.Org.Id)

	expected := "5c8e7160-5b60-4f49-824f-c01c111ea29f   prj1:front\n9931b808-9f92-4283-a8aa-d96289e11111   com.example:cmd-mock-conxxx\n"
	actual := prjs.String()

	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}

func Test_Project_Get_Verbose(t *testing.T) {
	client := tools.NewMockClient()
	client.ResponseBody = `{"org":{"name":"org-test","id":"16df2e12-d4cb-4111-aaf2-547db9ff07e9"},"projects":[{"id":"5c8e7160-5b60-4f49-824f-c01c111ea29f","name":"prj1:front","created":"2021-11-22T10:03:05.435Z","origin":"cli","type":"maven","readOnly":false,"testFrequency":"daily","isMonitored":true,"totalDependencies":0,"issueCountsBySeverity":{"low":0,"high":0,"medium":0,"critical":0},"remoteRepoUrl":"http://example.com/repot/prj1.git","imageTag":"1.0.0-SNAPSHOT","lastTestedDate":"2021-12-05T06:20:50.043Z","browseUrl":"https://app.snyk.io/org/org-test-pie/project/5c8e7160-5b60-4f49-824f-c01c111ea29f","owner":null,"importingUser":{"id":"7261cefe-93f4-472d-b6cd-27d8f41f1111","name":"org-test","username":"org-test","email":null},"tags":[{"key": "k1", "value": "v1"}, {"key":"k2", "value":"v2"}],"attributes":{"criticality":["high"],"lifecycle":["production"],"environment":["frontend"]},"branch":null}]}"`
	client.StatusCode = http.StatusOK
	client.Status = "XXX"
	prjs := NewProjects(client, "org_id")

	err := prjs.Get()
	assert.Nil(t, nil, err)
	expected := "5c8e7160-5b60-4f49-824f-c01c111ea29f   prj1:front                                        high,frontend,production,[k1=v1,k2=v2]\n"
	actual := prjs.Verbose()
	assert.Nil(t, nil, err)
	assert.Equal(t, expected, actual)
}

func Test_Project_Get_Raw(t *testing.T) {
	client := tools.NewMockClient()
	raw := `{"org":{"name":"org-test","id":"16df2e12-d4cb-4111-aaf2-547db9ff07e9"},"projects":[{"id":"5c8e7160-5b60-4f49-824f-c01c111ea29f","name":"prj1:front","created":"2021-11-22T10:03:05.435Z","origin":"cli","type":"maven","readOnly":false,"testFrequency":"daily","isMonitored":true,"totalDependencies":0,"issueCountsBySeverity":{"low":0,"high":0,"medium":0,"critical":0},"remoteRepoUrl":"http://example.com/repot/prj1.git","imageTag":"1.0.0-SNAPSHOT","lastTestedDate":"2021-12-05T06:20:50.043Z","browseUrl":"https://app.snyk.io/org/org-test-pie/project/5c8e7160-5b60-4f49-824f-c01c111ea29f","owner":null,"importingUser":{"id":"7261cefe-93f4-472d-b6cd-27d8f41f1111","name":"org-test","username":"org-test","email":null},"tags":[],"attributes":{"criticality":[],"lifecycle":[],"environment":[]},"branch":null},{"id":"9931b808-9f92-4283-a8aa-d96289e11111","name":"com.example:cmd-mock-conxxx","created":"2021-11-19T13:45:23.001Z","origin":"cli","type":"maven","readOnly":false,"testFrequency":"daily","isMonitored":true,"totalDependencies":152,"issueCountsBySeverity":{"low":0,"high":3,"medium":4,"critical":0},"remoteRepoUrl":"http://example.com/repo/prj2.git","imageTag":"0.0.1-SNAPSHOT","lastTestedDate":"2021-12-05T05:43:10.379Z","browseUrl":"https://app.snyk.io/org/org-test-pie/project/9931b808-9f92-4283-a8aa-d96289e11111","owner":null,"importingUser":{"id":"7261cefe-93f4-472d-b6cd-27d8f41f1111","name":"org-test","username":"org-test","email":null},"tags":[],"attributes":{"criticality":[],"lifecycle":[],"environment":[]},"branch":null}]}"`
	client.ResponseBody = raw
	client.StatusCode = http.StatusOK
	client.Status = "XXX"
	prjs := NewProjects(client, "org_id")

	actual, err := prjs.GetRaw()

	assert.Equal(t, nil, err)
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

func Test_GetProject_OK(t *testing.T) {
	client := tools.NewMockClient()
	raw := `{"name":"prj1:front","id":"5c8e7160-5b60-4f49-824f-c01c111ea29f","created":"2021-11-22T10:03:05.435Z","origin":"cli","type":"maven","readOnly":false,"testFrequency":"daily","totalDependencies":0,"issueCountsBySeverity":{"low":0,"high":0,"medium":0,"critical":0},"remoteRepoUrl":"http://example.com/repot/prj1.git","imageTag":"1.0.0-SNAPSHOT","hostname":"slex-mbp","lastTestedDate":"2021-12-08T10:13:19.059Z","browseUrl":"https://app.snyk.io/org/org-test-pie/project/5c8e7160-5b60-4f49-824f-c01c111ea29f","owner":null,"importingUser":{"id":"7261cefe-93f4-472d-b6cd-27d8f41f1111","name":"org-test","username":"org-test","email":null},"isMonitored":true,"tags":[{"key":"key2","value":"value4"},{"key":"key2","value":"value2"},{"key":"key","value":"value"}],"attributes":{"criticality":[],"lifecycle":[],"environment":[]},"remediation":{"pin":{},"patch":{},"upgrade":{}},"branch":null}`
	client.ResponseBody = raw
	client.StatusCode = http.StatusOK
	client.Status = "XXX"
	prjs := NewProjects(client, "org_id")

	out, err := prjs.GetRawProject("prj_id")

	assert.Nil(t, err)
	assert.Equal(t, raw, out)
}

func Test_AddAttributes_OK(t *testing.T) {
	client := tools.NewMockClient()
	client.ResponseBody = ``
	client.StatusCode = http.StatusOK
	client.Status = "XXX"
	prjs := NewProjects(client, "org_id")
	err := prjs.AddAttributes("prj_id", "frontend", "", "")
	assert.Nil(t, err)
}

func Test_AddAttributes_Parsefailed(t *testing.T) {
	client := tools.NewMockClient()
	client.ResponseBody = ``
	client.StatusCode = http.StatusOK
	client.Status = "XXX"
	prjs := NewProjects(client, "org_id")
	err := prjs.AddAttributes("prj_id", "xxx", "", "")
	expectedErrorMsg := "invalid environment value: xxx\nValid values: [frontend backend internal external mobile saas on-prem hosted distributed]"
	assert.EqualErrorf(t, err, expectedErrorMsg, "Error should be: %v, got: %v", expectedErrorMsg, err)
}

func Test_AddAttributes_KO(t *testing.T) {
	client := tools.NewMockClient()
	client.ResponseBody = ``
	client.StatusCode = http.StatusUnauthorized
	client.Status = "XXX"
	prjs := NewProjects(client, "org_id")
	err := prjs.AddAttributes("prj_id", "frontend", "", "")
	expectedErrorMsg := "failed to add attribute XXX"
	assert.EqualErrorf(t, err, expectedErrorMsg, "Error should be: %v, got: %v", expectedErrorMsg, err)
}

func Test_DeletePrj_OK(t *testing.T) {
	client := tools.NewMockClient()
	client.ResponseBody = ``
	client.StatusCode = http.StatusOK
	client.Status = "XXX"
	prjs := NewProjects(client, "org_id")
	err := prjs.DeleteProject("prj_id")
	assert.Nil(t, err)
}

func Test_DeleteProject_KO(t *testing.T) {
	client := tools.NewMockClient()
	client.ResponseBody = ``
	client.StatusCode = http.StatusUnauthorized
	client.Status = "XXX"
	prjs := NewProjects(client, "org_id")
	err := prjs.DeleteProject("prj_id")
	expectedErrorMsg := "deleteProject failed: XXX"
	assert.EqualErrorf(t, err, expectedErrorMsg, "Error should be: %v, got: %v", expectedErrorMsg, err)
}

func Test_GetFiltered_OK(t *testing.T) {
	client := tools.NewMockClient()
	raw := `{"org":{"name":"Sandbox","id":"16df2e12-d4cb-400d-aaf2-547db9ff1111"},"projects":[{"id":"2ab89519-82c1-45c6-86c3-ffe024c31111","name":"@example/app","created":"2021-11-22T10:03:04.501Z","origin":"cli","type":"yarn","readOnly":false,"testFrequency":"daily","isMonitored":true,"totalDependencies":89,"issueCountsBySeverity":{"low":0,"high":0,"medium":0,"critical":0},"remoteRepoUrl":"http://example.com/repo/app.git","imageTag":"1.0.0-SNAPSHOT","lastTestedDate":"2021-12-08T09:35:21.565Z","browseUrl":"https://app.snyk.io/org/sandbox-pie/project/2ab89519-82c1-45c6-86c3-ffe024c31111","owner":null,"importingUser":{"id":"7261cefe-93f4-472d-b6cd-27d8f41f1111","name":"org","username":"org","email":null},"tags":[],"attributes":{"criticality":["medium"],"lifecycle":["development"],"environment":["frontend"]},"branch":null}]}`
	client.ResponseBody = raw
	client.StatusCode = http.StatusOK
	client.Status = "XXX"
	prjs := NewProjects(client, "org_id")
	var mTags map[string]string
	err := prjs.GetFiltered("", "", "medium", mTags)
	assert.Nil(t, err)
}

func Test_GetFiltered_KO(t *testing.T) {
	client := tools.NewMockClient()
	raw := ""
	client.ResponseBody = raw
	client.StatusCode = http.StatusUnauthorized
	client.Status = "XXX"
	prjs := NewProjects(client, "org_id")
	var mTags map[string]string
	err := prjs.GetFiltered("", "", "medium", mTags)
	expectedErrorMsg := "Get filtered projects list failed: XXX "
	assert.EqualErrorf(t, err, expectedErrorMsg, "Error should be: %v, got: %v", expectedErrorMsg, err)
}

func Test_GetRawFiltered_OK(t *testing.T) {
	client := tools.NewMockClient()
	raw := `{"org":{"name":"Sandbox","id":"16df2e12-d4cb-400d-aaf2-547db9ff1111"},"projects":[{"id":"2ab89519-82c1-45c6-86c3-ffe024c31111","name":"@example/app","created":"2021-11-22T10:03:04.501Z","origin":"cli","type":"yarn","readOnly":false,"testFrequency":"daily","isMonitored":true,"totalDependencies":89,"issueCountsBySeverity":{"low":0,"high":0,"medium":0,"critical":0},"remoteRepoUrl":"http://example.com/repo/app.git","imageTag":"1.0.0-SNAPSHOT","lastTestedDate":"2021-12-08T09:35:21.565Z","browseUrl":"https://app.snyk.io/org/sandbox-pie/project/2ab89519-82c1-45c6-86c3-ffe024c31111","owner":null,"importingUser":{"id":"7261cefe-93f4-472d-b6cd-27d8f41f1111","name":"org","username":"org","email":null},"tags":[],"attributes":{"criticality":["medium"],"lifecycle":["development"],"environment":["frontend"]},"branch":null}]}`
	client.ResponseBody = raw
	client.StatusCode = http.StatusOK
	client.Status = "XXX"
	prjs := NewProjects(client, "org_id")
	var mTags map[string]string
	out, err := prjs.GetRawFiltered("", "", "medium", mTags)
	assert.Nil(t, err)
	assert.Equal(t, raw, out)
}

func Test_GetRawFiltered_KO(t *testing.T) {
	client := tools.NewMockClient()
	client.ResponseBody = ""
	client.StatusCode = http.StatusUnauthorized
	client.Status = "XXX"
	prjs := NewProjects(client, "org_id")
	var mTags map[string]string
	out, err := prjs.GetRawFiltered("", "", "medium", mTags)
	expectedErrorMsg := "Get filtered projects list failed: XXX "
	assert.EqualErrorf(t, err, expectedErrorMsg, "Error should be: %v, got: %v", expectedErrorMsg, err)
	assert.Equal(t, "", out)
}

func Test_DeleteAll_KO(t *testing.T) {
	client := tools.NewMockClient()
	client.ResponseBody = ""
	client.StatusCode = http.StatusUnauthorized
	client.Status = "XXX"
	prjs := NewProjects(client, "org_id")
	prj := Project{Id: "id", Name: "name"}
	prjs.Projects = append(prjs.Projects, &prj)
	out, err := prjs.DeleteAllProjects()

	expectedErrorMsg := "deleteProject failed: XXX"
	assert.EqualErrorf(t, err, expectedErrorMsg, "Error should be: %v, got: %v", expectedErrorMsg, err)
	assert.Equal(t, "", out)
}

func Test_DeleteAll_OK(t *testing.T) {
	client := tools.NewMockClient()
	client.ResponseBody = ""
	client.StatusCode = http.StatusOK
	client.Status = "XXX"
	prjs := NewProjects(client, "org_id")
	prj := Project{Id: "id", Name: "name"}
	prjs.Projects = append(prjs.Projects, &prj)
	out, err := prjs.DeleteAllProjects()
	assert.Nil(t, err)
	expected := "id                                    DELETED\n"

	assert.Equal(t, expected, out)
}
