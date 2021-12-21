package domain

import (
	"fmt"
	"snykctl/internal/tools"
)

type RepOrg struct {
	Id     string
	Name   string
	Prjs   []*RepProject
	client tools.HttpClient
}

type RepProject struct {
	Id     string
	Name   string
	Issues []*Issue
}

const projectIssuePath = "/org/%s/project/%s/aggregated-issues"
const projectIssuePathPath = "/org/%s/project/%s/issue/%s/paths"

func NewReport(client tools.HttpClient, org_id string) (*RepOrg, error) {
	var err error
	rep := new(RepOrg)
	rep.Id = org_id
	rep.client = client

	orgs := NewOrgs(client)

	rep.Name, err = orgs.GetOrgName(org_id)
	if err != nil {
		return rep, err
	}

	prjs := NewProjects(client, org_id)
	if err := prjs.Get(); err != nil {
		return rep, err
	}

	for _, v := range prjs.Projects {
		rep.addProject(v.Id, v.Name)
	}

	return rep, nil
}

func (r *RepOrg) addProject(prj_id, prj_name string) {
	p := new(RepProject)
	p.Id = prj_id
	p.Name = prj_name
	r.Prjs = append(r.Prjs, p)
}

func (r *RepOrg) GetIssues(issueType string) error {

	prjs := NewProjects(r.client, r.Id)

	for _, prj := range r.Prjs {
		issues, err := prjs.GetIssues(prj.Id, issueType)
		if err != nil {
			return err
		}
		for _, issue := range issues.Issues {
			paths, err := GetProjectIssuePaths(r.client, r.Id, prj.Id, issue.Id)
			if err != nil {
				return err
			}
			for _, path := range paths.Paths {
				issue.IntroducedThrough += path[0].Name + "@" + path[0].Version + " "
			}
			prj.Issues = append(prj.Issues, issue)
		}
	}

	return nil
}

func (r *RepOrg) ToString() string {
	var str string
	str = "ORG: " + r.Name + "\n"
	str += "PRJS:\n"
	for _, prj := range r.Prjs {
		str += "---- " + prj.Name + "\n"
		str += "--------ISSUES:\n"
		for _, issue := range prj.Issues {
			str += "======== " + issue.Id + "\n"
			str += "======== INTRODUCED " + issue.IntroducedThrough + "\n"
		}
	}
	return str
}

func (r *RepOrg) ToCsv() string {
	var str string
	// header
	str = "Org Name,Prj Name,Issue Id, Severity, IsIgnored, Introduced\n"
	// body
	for _, prj := range r.Prjs {
		for _, issue := range prj.Issues {
			ignoredStr := "false"
			if issue.IsIgnored {
				ignoredStr = "true"
			}
			str += r.Name + "," + prj.Name + "," + issue.Id + "," + issue.IssueData.Severity + "," + ignoredStr + "," + issue.IntroducedThrough + "\n"
		}
	}

	return str
}

func (r *RepOrg) ToHtmlTable() string {
	header := "<table><tr><td>Org name</td><td>Prj Name</td><td>Issue Id</td><td>Severity</td><td>IsIgnored</td><td>Introduced Through</td></tr>"
	footer := "</table>"
	var body string
	for _, prj := range r.Prjs {
		for _, issue := range prj.Issues {
			ignoredStr := "false"
			if issue.IsIgnored {
				ignoredStr = "true"
			}
			body += fmt.Sprintf("<tr><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td></tr>", r.Name, prj.Name, issue.Id, issue.IssueData.Severity, ignoredStr, issue.IntroducedThrough)
		}
	}

	return header + body + footer
}
