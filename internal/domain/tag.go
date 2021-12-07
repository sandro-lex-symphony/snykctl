package domain

import (
	"fmt"
	"net/http"
	"snykctl/internal/tools"
	"strings"
)

const tagPath = "/org/%s/project/%s/tags"

func ParseTag(tag string) (string, string, error) {
	if !strings.Contains(tag, "=") {
		return "", "", fmt.Errorf("invalid tag. Not a key=value format")
	}
	parts := strings.Split(tag, "=")
	if len(parts[0]) < 1 || len(parts[1]) < 1 {
		return "", "", fmt.Errorf("invalid tag. Not a key=value format")
	}
	return parts[0], parts[1], nil
}

func AddTag(client tools.HttpClient, org_id string, prj_id string, tag string) error {
	k, v, err := ParseTag(tag)
	if err != nil {
		return err
	}

	path := fmt.Sprintf(tagPath, org_id, prj_id)

	tagBody := fmt.Sprintf(`{"key": "%s", "value": "%s"}`, k, v)
	var jsonStr = []byte(tagBody)

	resp := client.RequestPost(path, jsonStr)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Failed to add tag %s", resp.Status)
	}
	return nil
}
