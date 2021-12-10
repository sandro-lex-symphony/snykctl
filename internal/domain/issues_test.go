package domain

import (
	"net/http"
	"snykctl/internal/tools"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_CheckIssuesType(t *testing.T) {
	assert.Nil(t, CheckIssueType("license"))
	assert.Nil(t, CheckIssueType("vuln"))
	err := CheckIssueType("xxx")
	expectedErrorMsg := "invalid type. (license | vuln)"
	assert.EqualErrorf(t, err, expectedErrorMsg, "Error should be: %v, got: %v", expectedErrorMsg, err)
}

type testIssueData struct {
	in  string
	out string
}

func Test_BuildIssueTypeFilter(t *testing.T) {
	tests := []testIssueData{
		{"license", `{"filters": {"types": [ "license" ]}}`},
		{"vuln", `{"filters": {"types": [ "vuln" ]}}`},
		{"xxx", `{"filters": {"types": [ "vuln", "license" ]}}`},
		{"", `{"filters": {"types": [ "vuln", "license" ]}}`},
	}

	for _, test := range tests {
		assert.Equal(t, BuildIssueTypeFilter(test.in), test.out)
	}
}

func Test_FormatProjectIssues(t *testing.T) {
	id1 := IssueData{Id: "id1", Title: "t1", Severity: "high", ExploitMaturity: "m1", CVSSv3: "c1", CvssScore: 1.0}
	i1 := Issue{Id: "i1", PkgName: "p1", PkgVersion: "v1", IssueData: id1, IsIgnored: false, IssueType: "t1", IntroducedThrough: "i1"}
	var issues []*Issue
	issues = append(issues, &i1)
	res := ProjectIssuesResult{Issues: issues}

	out := FormatProjectIssues(res, "")
	expected := "i1                                    p1                            high           t1        false\n"
	assert.Equal(t, expected, out)

	out = FormatProjectIssues(res, "p1")
	expected = "p1                                    i1                                    p1                            high           t1        false\n"
	assert.Equal(t, expected, out)
}

func Test_GetIssues_BadBody(t *testing.T) {
	client := tools.NewMockClient()
	client.ResponseBody = `xxx`
	client.StatusCode = http.StatusOK
	client.Status = "XXX"

	prjs := NewProjects(client, "org")

	_, err := prjs.GetIssues("prj", "")
	expectedErrorMsg := "invalid character 'x' looking for beginning of value"
	assert.EqualErrorf(t, err, expectedErrorMsg, "Error should be: %v, got: %v", expectedErrorMsg, err)
}

func Test_GetIssues_KO(t *testing.T) {
	client := tools.NewMockClient()
	client.ResponseBody = `xxx`
	client.StatusCode = http.StatusUnauthorized
	client.Status = "XXX"

	prjs := NewProjects(client, "org")

	_, err := prjs.GetIssues("prj", "")
	expectedErrorMsg := "getIssues failed: XXX"
	assert.EqualErrorf(t, err, expectedErrorMsg, "Error should be: %v, got: %v", expectedErrorMsg, err)
}

func Test_GetIssues_OK(t *testing.T) {
	raw := `{"issues":[{"id":"SNYK-JAVA-COMNIMBUSDS-536068","issueType":"vuln","pkgName":"com.nimbusds:nimbus-jose-jwt","pkgVersions":["5.1"],"priorityScore":422,"priority":{"score":422,"factors":[{"name":"exploitMaturity","description":"Proof of Concept exploit"},{"name":"cvssScore","description":"CVSS 6.3"}]},"issueData":{"id":"SNYK-JAVA-COMNIMBUSDS-536068","title":"Improper Check for Unusual or Exceptional Conditions","severity":"medium","url":"https://snyk.io/vuln/SNYK-JAVA-COMNIMBUSDS-536068","identifiers":{"CVE":["CVE-2019-17195"],"CWE":["CWE-754"],"GHSA":["GHSA-f6vf-pq8c-69m4"]},"credit":["Unknown"],"exploitMaturity":"proof-of-concept","semver":{"vulnerable":["[,7.8.1)"]},"publicationTime":"2019-11-26T10:16:06Z","disclosureTime":"2019-10-15T14:15:00Z","CVSSv3":"CVSS:3.1/AV:N/AC:L/PR:N/UI:R/S:U/C:L/I:L/A:L/E:P/RL:O/RC:R","cvssScore":6.3,"language":"java","patches":[],"nearestFixedInVersion":"","isMaliciousPackage":false},"isPatched":false,"isIgnored":false,"fixInfo":{"isUpgradable":false,"isPinnable":false,"isPatchable":false,"isFixable":false,"isPartiallyFixable":false,"nearestFixedInVersion":"","fixedIn":["7.8.1"]},"links":{"paths":"https://app.snyk.io/api/v1/org/4c961058-b36e-4510-894f-99d3f39d3498/project/200b6d0a-e0d5-4d93-aaf8-d93f6edf5260/history/83bec2e3-ae8c-4757-b97b-2382a85f301d/issue/SNYK-JAVA-COMNIMBUSDS-536068/paths"}},{"id":"SNYK-JAVA-NETMINIDEV-1298655","issueType":"vuln","pkgName":"net.minidev:json-smart","pkgVersions":["2.3"],"priorityScore":265,"priority":{"score":265,"factors":[{"name":"cvssScore","description":"CVSS 5.3"}]},"issueData":{"id":"SNYK-JAVA-NETMINIDEV-1298655","title":"Denial of Service (DoS)","severity":"medium","url":"https://snyk.io/vuln/SNYK-JAVA-NETMINIDEV-1298655","identifiers":{"CVE":["CVE-2021-31684"],"CWE":["CWE-400"]},"credit":["@pcy190"],"exploitMaturity":"no-known-exploit","semver":{"vulnerable":["[0,1.3.3)","[2.0.0, 2.4.5)"]},"publicationTime":"2021-06-02T11:13:38.628733Z","disclosureTime":"2021-06-02T10:02:59Z","CVSSv3":"CVSS:3.1/AV:N/AC:L/PR:N/UI:N/S:U/C:N/I:N/A:L","cvssScore":5.3,"language":"java","patches":[],"nearestFixedInVersion":"","isMaliciousPackage":false},"isPatched":false,"isIgnored":false,"fixInfo":{"isUpgradable":false,"isPinnable":false,"isPatchable":false,"isFixable":false,"isPartiallyFixable":false,"nearestFixedInVersion":"","fixedIn":["1.3.3","2.4.5"]},"links":{"paths":"https://app.snyk.io/api/v1/org/4c961058-b36e-4510-894f-99d3f39d3498/project/200b6d0a-e0d5-4d93-aaf8-d93f6edf5260/history/83bec2e3-ae8c-4757-b97b-2382a85f301d/issue/SNYK-JAVA-NETMINIDEV-1298655/paths"}},{"id":"SNYK-JAVA-NETMINIDEV-1078499","issueType":"vuln","pkgName":"net.minidev:json-smart","pkgVersions":["2.3"],"priorityScore":265,"priority":{"score":265,"factors":[{"name":"cvssScore","description":"CVSS 5.3"}]},"issueData":{"id":"SNYK-JAVA-NETMINIDEV-1078499","title":"Denial of Service (DoS)","severity":"medium","url":"https://snyk.io/vuln/SNYK-JAVA-NETMINIDEV-1078499","identifiers":{"CVE":["CVE-2021-27568"],"CWE":["CWE-400"],"GHSA":["GHSA-v528-7hrm-frqp"]},"credit":["Tobias Mayer"],"exploitMaturity":"no-known-exploit","semver":{"vulnerable":["[0,1.3.2)","[2.0.0, 2.4.1)"]},"publicationTime":"2021-02-23T16:42:39Z","disclosureTime":"2021-02-23T11:31:54Z","CVSSv3":"CVSS:3.1/AV:N/AC:L/PR:N/UI:N/S:U/C:N/I:N/A:L","cvssScore":5.3,"language":"java","patches":[],"nearestFixedInVersion":"","isMaliciousPackage":false},"isPatched":false,"isIgnored":false,"fixInfo":{"isUpgradable":false,"isPinnable":false,"isPatchable":false,"isFixable":false,"isPartiallyFixable":false,"nearestFixedInVersion":"","fixedIn":["1.3.2","2.4.1"]},"links":{"paths":"https://app.snyk.io/api/v1/org/4c961058-b36e-4510-894f-99d3f39d3498/project/200b6d0a-e0d5-4d93-aaf8-d93f6edf5260/history/83bec2e3-ae8c-4757-b97b-2382a85f301d/issue/SNYK-JAVA-NETMINIDEV-1078499/paths"}}]}`
	client := tools.NewMockClient()
	client.ResponseBody = raw
	client.StatusCode = http.StatusOK
	client.Status = "XXX"

	prjs := NewProjects(client, "org")

	out, err := prjs.GetIssues("prj", "")
	assert.Nil(t, err)
	assert.Equal(t, 3, len(out.Issues))
}

func Test_GetIssues_KO_BadType(t *testing.T) {
	client := tools.NewMockClient()
	client.ResponseBody = ""
	client.StatusCode = http.StatusOK
	client.Status = "XXX"

	prjs := NewProjects(client, "org")

	_, err := prjs.GetIssues("prj", "xxx")
	expectedErrorMsg := "invalid type. (license | vuln)"
	assert.EqualErrorf(t, err, expectedErrorMsg, "Error should be: %v, got: %v", expectedErrorMsg, err)
}
