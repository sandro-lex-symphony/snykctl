package domain

import (
	"encoding/json"
	"fmt"
	"net/http"
	"snykctl/internal/tools"
)

const ignorePath = "/org/%s/project/%s/ignores"

type IgnoreStar struct {
	Star Ignore `json:"*"`
}

type Ignore struct {
	Reason     string
	Created    string
	Expires    string
	reasonType string
	IgnoredBy  User
}

type IgnoreResult struct {
	Id      string
	Content Ignore
}

func GetProjectIgnores(client tools.HttpClient, org_id, prj_id string) ([]IgnoreResult, error) {
	var result []IgnoreResult
	path := fmt.Sprintf(ignorePath, org_id, prj_id)
	resp := client.RequestGet(path)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return result, fmt.Errorf("GetProjects failed: %s", resp.Status)
	}

	var ignore_result map[string][]IgnoreStar

	if err := json.NewDecoder(resp.Body).Decode(&ignore_result); err != nil {
		return result, err
	}

	for key, value := range ignore_result {
		for i := 0; i < len(value); i++ {
			var ii IgnoreResult
			ii.Id = key
			ii.Content = value[i].Star
			result = append(result, ii)
		}
	}

	return result, nil
}

func FormatIgnore(res IgnoreResult, prj string) string {
	if prj != "" {
		return fmt.Sprintf("%-38s%-30s%-30s%-30s%s\n", prj, res.Id, res.Content.Created, res.Content.IgnoredBy.Email, res.Content.Reason)
	} else {
		return fmt.Sprintf("%-30s%-30s%-30s%s\n", res.Id, res.Content.Created, res.Content.IgnoredBy.Email, res.Content.Reason)
	}
}

func FormatIgnoreResult(res []IgnoreResult, prj string) string {
	var out string
	for i := 0; i < len(res); i++ {
		out += FormatIgnore(res[i], prj)
	}

	return out
}
