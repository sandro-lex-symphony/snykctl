package domain

import (
	"fmt"
	"strings"
)

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
