package domain

import (
	"encoding/json"
	"fmt"
	"net/http"
)

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
		return out, fmt.Errorf("getIssues failed: %s ", resp.Status)
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
		if t != "license" && t != "vuln" && t != "vulnerability" {
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
