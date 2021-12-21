package domain

import (
	"encoding/json"
	"fmt"
	"net/http"
	"snykctl/internal/tools"
)

const issuesPathPath = "/org/%s/project/%s/issue/%s/paths"
const issueCountPath = "/reporting/counts/issues/latest?groupBy=project,severity"

type ProjectIssuesResult struct {
	Issues []*Issue
}

type Issue struct {
	Id                string
	PkgName           string
	PkgVersion        string
	IssueData         IssueData
	IsIgnored         bool
	IssueType         string
	IntroducedThrough string
}

type IssuePathsResult struct {
	SnapshotId string
	Paths      [][]IssuePath
}

type IssuePath struct {
	Name    string
	Version string
}

type IssueCountResults struct {
	Results []IssueCountResult
}

type IssueCountResult struct {
	Day     string
	Results []ProjectIssueCount
}

type ProjectIssueCount struct {
	ProjectId string
	Count     int
	Severity  Severity
}

type Severity struct {
	Critical int
	High     int
	Medium   int
	Low      int
}

type IssueData struct {
	Id              string
	Title           string
	Severity        string
	ExploitMaturity string
	CVSSv3          string
	CvssScore       float32
}

func (p *Projects) GetIssues(prj_id string, issueType string) (ProjectIssuesResult, error) {
	var out ProjectIssuesResult

	if err := CheckIssueType(issueType); err != nil {
		return out, err
	}

	path := fmt.Sprintf(issuesPath, p.Org.Id, prj_id)

	body := BuildIssueTypeFilter(issueType)
	var jsonStr = []byte(body)
	resp := p.client.RequestPost(path, jsonStr)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return out, fmt.Errorf("getIssues failed: %s", resp.Status)
	}

	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return out, err
	}

	return out, nil
}

func FormatProjectIssues(in ProjectIssuesResult, prj_id string) string {
	var out string
	for _, issue := range in.Issues {
		if prj_id == "" {
			out += fmt.Sprintf("%-38s%-30s%-15s%-10s%t\n", issue.Id, issue.PkgName, issue.IssueData.Severity, issue.IssueType, issue.IsIgnored)
		} else {
			out += fmt.Sprintf("%-38s%-38s%-30s%-15s%-10s%t\n", prj_id, issue.Id, issue.PkgName, issue.IssueData.Severity, issue.IssueType, issue.IsIgnored)
		}
	}
	return out
}

func CheckIssueType(t string) error {
	if t != "" {
		if t != "license" && t != "vuln" {
			return fmt.Errorf("invalid type. (license | vuln)")
		}
	}
	return nil
}

func BuildIssueTypeFilter(issueType string) string {
	var content string
	if issueType == "vuln" {
		content = `"types": [ "vuln" ]`
	} else if issueType == "license" {
		content = `"types": [ "license" ]`
	} else {
		content = `"types": [ "vuln", "license" ]`
	}
	return fmt.Sprintf(`{"filters": {%s}}`, content)
}

func GetProjectIssuePaths(client tools.HttpClient, org_id string, prj_id string, issue_id string) (*IssuePathsResult, error) {
	path := fmt.Sprintf(issuesPathPath, org_id, prj_id, issue_id)
	resp := client.RequestGet(path)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GetProjectIssues failed %s", resp.Status)
	}

	var result IssuePathsResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func IssueCount(client tools.HttpClient, org_id string) (IssueCountResults, error) {
	var result IssueCountResults
	str := fmt.Sprintf("{\"filters\": {\"orgs\": [\"%s\"] } }", org_id)

	var jsonStr = []byte(str)
	resp := client.RequestPost(issueCountPath, jsonStr)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return result, fmt.Errorf("Issue count failed: %s ", resp.Status)
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return result, err
	}

	return result, nil
}

func FormatIssueCountResults(r IssueCountResults, output string) string {
	var out string
	var count, c, h, m, l int

	for _, results := range r.Results {
		for _, projectResult := range results.Results {
			count += projectResult.Count
			c += projectResult.Severity.Critical
			h += projectResult.Severity.High
			m += projectResult.Severity.Medium
			l += projectResult.Severity.Low
		}
	}
	if output == "html" {
		header := "<tabsle><tr><td>Total</td><td>Critical</td><td>High</td><td>Medium</td><td>Low</td></tr>"
		body := fmt.Sprintf("<tr><td>%d</td><td>%d</td><td>%d</td><td>%d</td><td>%d</td></tr>", count, c, h, m, l)
		footer := "</table>"
		out = header + body + footer
	} else {
		out = fmt.Sprintf("COUNT: %d\nCRITICAL: %d\nHIGH: %d\nMEDIUM: %d\nLOW: %d\n", count, c, h, m, l)
	}
	return out
}
