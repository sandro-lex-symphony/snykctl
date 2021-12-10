package domain

import (
	"encoding/json"
	"fmt"
	"net/http"
	"snykctl/internal/tools"
)

const ignorePath = "/org/%s/project/%s/ignores"

type IgnoreResult map[string][]IgnoreStar

type IgnoreStar struct {
	Star IgnoreContent `json:"*"`
}

type IgnoreContent struct {
	Reason     string
	Created    string
	Expires    string
	ReasonType string
	IgnoredBy  User
}

type Ignore struct {
	Id      string
	Content IgnoreContent
}

func GetProjectIgnores(client tools.HttpClient, org_id, prj_id string) (IgnoreResult, error) {
	var res IgnoreResult
	path := fmt.Sprintf(ignorePath, org_id, prj_id)

	resp := client.RequestGet(path)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return res, fmt.Errorf("getProjectsIgnores failed: %s", resp.Status)
	}

	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return res, err
	}

	return res, nil
}

func FormatIgnore(res Ignore, prj string) string {
	if prj != "" {
		return fmt.Sprintf("%-38s%-30s%-30s%-30s%s\n", prj, res.Id, res.Content.Created, res.Content.IgnoredBy.Email, res.Content.Reason)
	} else {
		return fmt.Sprintf("%-30s%-30s%-30s%s\n", res.Id, res.Content.Created, res.Content.IgnoredBy.Email, res.Content.Reason)
	}
}

func FormatIgnoreResult(res IgnoreResult, prj string) string {
	var items []Ignore
	for key, value := range res {
		for i := 0; i < len(value); i++ {
			var ii Ignore
			ii.Id = key
			ii.Content = value[i].Star
			items = append(items, ii)
		}
	}

	var out string
	for i := 0; i < len(items); i++ {
		out += FormatIgnore(items[i], prj)
	}
	return out
}
