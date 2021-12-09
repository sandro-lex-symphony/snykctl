package domain

import (
	"net/http"
	"snykctl/internal/tools"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Org_Get_httpError(t *testing.T) {
	client := tools.NewMockClient()
	client.ResponseBody = ""
	client.StatusCode = http.StatusUnauthorized
	client.Status = "XXX"
	orgs := NewOrgs(client)

	err := orgs.Get()
	expectedErrorMsg := "GetOrgs failed: XXX"
	assert.EqualErrorf(t, err, expectedErrorMsg, "Error should be: %v, got: %v", expectedErrorMsg, err)
	assert.Equal(t, false, orgs.Sync())
}

func Test_Org_Get_badBody(t *testing.T) {
	client := tools.NewMockClient()
	client.ResponseBody = "filler"
	client.StatusCode = http.StatusOK
	client.Status = "XXX"
	orgs := NewOrgs(client)

	err := orgs.Get()
	expectedErrorMsg := "GetOrgs failed:"
	assert.Containsf(t, err.Error(), expectedErrorMsg, "Error should be: %v, got: %v", expectedErrorMsg, err)
	assert.Equal(t, false, orgs.Sync())
}

func Test_Org_Get_Ok(t *testing.T) {
	client := tools.NewMockClient()
	client.ResponseBody = `{"orgs":[{"id":"f6910fd7-43a3-4e20-8327-6b621b7746b3","name":"JDC On Prem","slug":"jdc-on-prem","url":"https://app.snyk.io/org/jdc-on-prem","group":{"name":"Symphony","id":"25c3050c-d3c7-464c-8517-4181b4b12308"},"created":"2021-09-08T14:42:36.756Z"},{"id":"711c53b6-a85d-4a51-a34f-42552cc8572e","name":"Release - Current","slug":"release-current","url":"https://app.snyk.io/org/release-current","group":{"name":"Symphony","id":"25c3050c-d3c7-464c-8517-4181b4b12308"},"created":"2021-08-26T14:53:49.842Z"},{"id":"10fee9f9-c85c-470d-b9b7-4c9e20b09f07","name":"Directory","slug":"directory-yvz","url":"https://app.snyk.io/org/directory-yvz","group":{"name":"Symphony","id":"25c3050c-d3c7-464c-8517-4181b4b12308"},"created":"2021-03-30T15:33:07.854Z"}]}`
	client.StatusCode = http.StatusOK
	client.Status = "XXX"
	orgs := NewOrgs(client)

	err := orgs.Get()

	assert.Equal(t, nil, err)
	assert.Equal(t, true, orgs.Sync())

	assert.Equal(t, 3, len(orgs.Orgs))

	var idFound bool
	for _, o := range orgs.Orgs {
		if o.Id == "f6910fd7-43a3-4e20-8327-6b621b7746b3" {
			assert.Equal(t, "JDC On Prem", o.Name)
			idFound = true
		}
	}
	assert.True(t, idFound)
}

func Test_Org_Get_Quiet(t *testing.T) {
	client := tools.NewMockClient()
	client.ResponseBody = `{"orgs":[{"id":"f6910fd7-43a3-4e20-8327-6b621b7746b3","name":"JDC On Prem","slug":"jdc-on-prem","url":"https://app.snyk.io/org/jdc-on-prem","group":{"name":"Symphony","id":"25c3050c-d3c7-464c-8517-4181b4b12308"},"created":"2021-09-08T14:42:36.756Z"},{"id":"711c53b6-a85d-4a51-a34f-42552cc8572e","name":"Release - Current","slug":"release-current","url":"https://app.snyk.io/org/release-current","group":{"name":"Symphony","id":"25c3050c-d3c7-464c-8517-4181b4b12308"},"created":"2021-08-26T14:53:49.842Z"},{"id":"10fee9f9-c85c-470d-b9b7-4c9e20b09f07","name":"Directory","slug":"directory-yvz","url":"https://app.snyk.io/org/directory-yvz","group":{"name":"Symphony","id":"25c3050c-d3c7-464c-8517-4181b4b12308"},"created":"2021-03-30T15:33:07.854Z"}]}`
	client.StatusCode = http.StatusOK
	client.Status = "XXX"
	orgs := NewOrgs(client)

	err := orgs.Get()

	assert.Equal(t, nil, err)
	assert.Equal(t, true, orgs.Sync())

	assert.Equal(t, 3, len(orgs.Orgs))

	expected := "f6910fd7-43a3-4e20-8327-6b621b7746b3\n711c53b6-a85d-4a51-a34f-42552cc8572e\n10fee9f9-c85c-470d-b9b7-4c9e20b09f07\n"
	actual := orgs.Quiet()
	assert.Equal(t, expected, actual)
	assert.Equal(t, nil, err)
}

func Test_Org_Get_String(t *testing.T) {
	client := tools.NewMockClient()
	client.ResponseBody = `{"orgs":[{"id":"f6910fd7-43a3-4e20-8327-6b621b7746b3","name":"JDC On Prem","slug":"jdc-on-prem","url":"https://app.snyk.io/org/jdc-on-prem","group":{"name":"Symphony","id":"25c3050c-d3c7-464c-8517-4181b4b12308"},"created":"2021-09-08T14:42:36.756Z"},{"id":"711c53b6-a85d-4a51-a34f-42552cc8572e","name":"Release - Current","slug":"release-current","url":"https://app.snyk.io/org/release-current","group":{"name":"Symphony","id":"25c3050c-d3c7-464c-8517-4181b4b12308"},"created":"2021-08-26T14:53:49.842Z"},{"id":"10fee9f9-c85c-470d-b9b7-4c9e20b09f07","name":"Directory","slug":"directory-yvz","url":"https://app.snyk.io/org/directory-yvz","group":{"name":"Symphony","id":"25c3050c-d3c7-464c-8517-4181b4b12308"},"created":"2021-03-30T15:33:07.854Z"}]}`
	client.StatusCode = http.StatusOK
	client.Status = "XXX"
	orgs := NewOrgs(client)

	err := orgs.Get()

	assert.Equal(t, nil, err)
	assert.Equal(t, true, orgs.Sync())

	assert.Equal(t, 3, len(orgs.Orgs))

	expected := "f6910fd7-43a3-4e20-8327-6b621b7746b3   JDC On Prem\n711c53b6-a85d-4a51-a34f-42552cc8572e   Release - Current\n10fee9f9-c85c-470d-b9b7-4c9e20b09f07   Directory\n"
	actual := orgs.String()
	assert.Equal(t, expected, actual)
	assert.Equal(t, nil, err)
}

func Test_Org_Get_Names(t *testing.T) {
	client := tools.NewMockClient()
	client.ResponseBody = `{"orgs":[{"id":"f6910fd7-43a3-4e20-8327-6b621b7746b3","name":"JDC On Prem","slug":"jdc-on-prem","url":"https://app.snyk.io/org/jdc-on-prem","group":{"name":"Symphony","id":"25c3050c-d3c7-464c-8517-4181b4b12308"},"created":"2021-09-08T14:42:36.756Z"},{"id":"711c53b6-a85d-4a51-a34f-42552cc8572e","name":"Release - Current","slug":"release-current","url":"https://app.snyk.io/org/release-current","group":{"name":"Symphony","id":"25c3050c-d3c7-464c-8517-4181b4b12308"},"created":"2021-08-26T14:53:49.842Z"},{"id":"10fee9f9-c85c-470d-b9b7-4c9e20b09f07","name":"Directory","slug":"directory-yvz","url":"https://app.snyk.io/org/directory-yvz","group":{"name":"Symphony","id":"25c3050c-d3c7-464c-8517-4181b4b12308"},"created":"2021-03-30T15:33:07.854Z"}]}`
	client.StatusCode = http.StatusOK
	client.Status = "XXX"
	orgs := NewOrgs(client)

	err := orgs.Get()

	assert.Equal(t, nil, err)
	assert.Equal(t, true, orgs.Sync())

	assert.Equal(t, 3, len(orgs.Orgs))

	expected := "JDC On Prem\nRelease - Current\nDirectory\n"
	actual := orgs.Names()
	assert.Equal(t, expected, actual)
	assert.Equal(t, nil, err)
}

func Test_Org_Get_Raw(t *testing.T) {
	client := tools.NewMockClient()
	raw := `{"orgs":[{"id":"f6910fd7-43a3-4e20-8327-6b621b7746b3","name":"JDC On Prem","slug":"jdc-on-prem","url":"https://app.snyk.io/org/jdc-on-prem","group":{"name":"Symphony","id":"25c3050c-d3c7-464c-8517-4181b4b12308"},"created":"2021-09-08T14:42:36.756Z"},{"id":"711c53b6-a85d-4a51-a34f-42552cc8572e","name":"Release - Current","slug":"release-current","url":"https://app.snyk.io/org/release-current","group":{"name":"Symphony","id":"25c3050c-d3c7-464c-8517-4181b4b12308"},"created":"2021-08-26T14:53:49.842Z"},{"id":"10fee9f9-c85c-470d-b9b7-4c9e20b09f07","name":"Directory","slug":"directory-yvz","url":"https://app.snyk.io/org/directory-yvz","group":{"name":"Symphony","id":"25c3050c-d3c7-464c-8517-4181b4b12308"},"created":"2021-03-30T15:33:07.854Z"}]}`
	client.ResponseBody = raw
	client.StatusCode = http.StatusOK
	client.Status = "XXX"
	orgs := NewOrgs(client)

	out, err := orgs.GetRaw()

	assert.Equal(t, nil, err)
	assert.Equal(t, true, orgs.Sync())

	assert.Equal(t, raw, out)
}

func Test_Org_Get_Raw_Err(t *testing.T) {
	client := tools.NewMockClient()
	raw := `{"orgs":[{"id":"f6910fd7-43a3-4e20-8327-6b621b7746b3","name":"JDC On Prem","slug":"jdc-on-prem","url":"https://app.snyk.io/org/jdc-on-prem","group":{"name":"Symphony","id":"25c3050c-d3c7-464c-8517-4181b4b12308"},"created":"2021-09-08T14:42:36.756Z"},{"id":"711c53b6-a85d-4a51-a34f-42552cc8572e","name":"Release - Current","slug":"release-current","url":"https://app.snyk.io/org/release-current","group":{"name":"Symphony","id":"25c3050c-d3c7-464c-8517-4181b4b12308"},"created":"2021-08-26T14:53:49.842Z"},{"id":"10fee9f9-c85c-470d-b9b7-4c9e20b09f07","name":"Directory","slug":"directory-yvz","url":"https://app.snyk.io/org/directory-yvz","group":{"name":"Symphony","id":"25c3050c-d3c7-464c-8517-4181b4b12308"},"created":"2021-03-30T15:33:07.854Z"}]}`
	client.ResponseBody = raw
	client.StatusCode = http.StatusUnauthorized
	client.Status = "XXX"
	orgs := NewOrgs(client)

	out, err := orgs.GetRaw()
	assert.Equal(t, "", out)
	expectedErrorMsg := "GetOrgs failed: XXX"
	assert.EqualErrorf(t, err, expectedErrorMsg, "Error should be: %v, got: %v", expectedErrorMsg, err)

	assert.Equal(t, false, orgs.Sync())
}

func Test_Org_GetOrgName_OK(t *testing.T) {
	client := tools.NewMockClient()
	client.ResponseBody = `{"orgs":[{"id":"f6910fd7-43a3-4e20-8327-6b621b7746b3","name":"JDC On Prem","slug":"jdc-on-prem","url":"https://app.snyk.io/org/jdc-on-prem","group":{"name":"Symphony","id":"25c3050c-d3c7-464c-8517-4181b4b12308"},"created":"2021-09-08T14:42:36.756Z"},{"id":"711c53b6-a85d-4a51-a34f-42552cc8572e","name":"Release - Current","slug":"release-current","url":"https://app.snyk.io/org/release-current","group":{"name":"Symphony","id":"25c3050c-d3c7-464c-8517-4181b4b12308"},"created":"2021-08-26T14:53:49.842Z"},{"id":"10fee9f9-c85c-470d-b9b7-4c9e20b09f07","name":"Directory","slug":"directory-yvz","url":"https://app.snyk.io/org/directory-yvz","group":{"name":"Symphony","id":"25c3050c-d3c7-464c-8517-4181b4b12308"},"created":"2021-03-30T15:33:07.854Z"}]}`
	client.StatusCode = http.StatusOK
	client.Status = "XXX"
	orgs := NewOrgs(client)

	out, err := orgs.GetOrgName("f6910fd7-43a3-4e20-8327-6b621b7746b3")

	assert.Equal(t, nil, err)
	assert.Equal(t, true, orgs.Sync())
	assert.Equal(t, "JDC On Prem", out)
}

func Test_Org_GetOrgName_Notfound(t *testing.T) {
	client := tools.NewMockClient()
	client.ResponseBody = `{"orgs":[{"id":"f6910fd7-43a3-4e20-8327-6b621b7746b3","name":"JDC On Prem","slug":"jdc-on-prem","url":"https://app.snyk.io/org/jdc-on-prem","group":{"name":"Symphony","id":"25c3050c-d3c7-464c-8517-4181b4b12308"},"created":"2021-09-08T14:42:36.756Z"},{"id":"711c53b6-a85d-4a51-a34f-42552cc8572e","name":"Release - Current","slug":"release-current","url":"https://app.snyk.io/org/release-current","group":{"name":"Symphony","id":"25c3050c-d3c7-464c-8517-4181b4b12308"},"created":"2021-08-26T14:53:49.842Z"},{"id":"10fee9f9-c85c-470d-b9b7-4c9e20b09f07","name":"Directory","slug":"directory-yvz","url":"https://app.snyk.io/org/directory-yvz","group":{"name":"Symphony","id":"25c3050c-d3c7-464c-8517-4181b4b12308"},"created":"2021-03-30T15:33:07.854Z"}]}`
	client.StatusCode = http.StatusOK
	client.Status = "XXX"
	orgs := NewOrgs(client)

	out, err := orgs.GetOrgName("123")

	assert.Equal(t, out, "")
	expectedErrorMsg := "getOrgName: org not found 123"
	assert.EqualErrorf(t, err, expectedErrorMsg, "Error should be: %v, got: %v", expectedErrorMsg, err)
}
