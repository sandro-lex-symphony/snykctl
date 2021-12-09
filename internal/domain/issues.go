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

func (p *Projects) GetProjectIssues(prj_id string) (string, error) {
	path := fmt.Sprintf(issuesPath, p.Org.Id, prj_id)

	resp := p.client.RequestPost(path, nil)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("getProjectIssues failed: %s ", resp.Status)
	}

	var result ProjectIssuesResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	var out string
	for _, issue := range result.Issues {
		out += fmt.Sprintf("%-38s%-30s%s%t\n", issue.Id, issue.PkgName, issue.IssueData.Severity, issue.IsIgnored)
	}
	return out, nil
}
