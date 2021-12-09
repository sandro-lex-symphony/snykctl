package domain

import (
	"fmt"
	"snykctl/internal/tools"
	"strings"
)

func ParseAttributes(filterEnv, filterLifecycle, filterCriticality string) error {
	if filterEnv != "" {
		if !tools.Contains(validEnvironments[:], filterEnv) {
			return fmt.Errorf("invalid environment value: %s\nValid values: %v", filterEnv, validEnvironments[:])
		}
	}

	if filterLifecycle != "" {
		if !tools.Contains(validLifecycle[:], filterLifecycle) {
			return fmt.Errorf("invalid lifecycle value: %s\nValid values: %v", filterLifecycle, validLifecycle[:])
		}
	}

	if filterCriticality != "" {
		if !tools.Contains(validCriticality[:], filterCriticality) {
			return fmt.Errorf("invalid lifecycle value: %s\nValid values: %v", filterCriticality, validCriticality[:])
		}
	}
	return nil
}

func BuildAttributesBody(env, lifecycle, criticality string) string {
	if env == "" && lifecycle == "" && criticality == "" {
		return ""
	}
	var values []string
	if env != "" {
		value := fmt.Sprintf(`"environment": ["%s"]`, env)
		values = append(values, value)
	}
	if lifecycle != "" {
		value := fmt.Sprintf(`"lifecycle": ["%s"]`, lifecycle)
		values = append(values, value)
	}
	if criticality != "" {
		value := fmt.Sprintf(`"criticality": ["%s"]`, criticality)
		values = append(values, value)
	}
	c := strings.Join(values, ",")
	return fmt.Sprintf("{ %s }", c)
}

func BuildAttributesFilter(env, lifecycle, criticality string) string {
	var attributes string

	if lifecycle == "" && env == "" && criticality == "" {
		return ""
	}

	var attrs []string
	if lifecycle != "" {
		at := fmt.Sprintf(`"lifecycle": [ "%s" ]`, lifecycle)
		attrs = append(attrs, at)
	}
	if env != "" {
		at := fmt.Sprintf(`"environment": [ "%s" ]`, env)
		attrs = append(attrs, at)
	}

	if criticality != "" {
		at := fmt.Sprintf(`"criticality": [ "%s" ]`, criticality)
		attrs = append(attrs, at)
	}

	attributes = strings.Join(attrs, ",")
	return fmt.Sprintf(`"attributes": { %s }`, attributes)
}

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

func ParseTags(filterTag []string) (map[string]string, error) {
	var mTags map[string]string
	if len(filterTag) > 0 {
		mTags = make(map[string]string)
		for _, tag := range filterTag {
			k, v, err := ParseTag(tag)
			if err != nil {
				return mTags, err
			}
			mTags[k] = v
		}
	}

	return mTags, nil
}

func BuildTagsFilter(mTags map[string]string) string {
	var tags string
	if len(mTags) > 0 {
		tags += ` "tags": { "includes": [`
		var ii []string
		for key, value := range mTags {
			i := fmt.Sprintf(`{ "key": "%s", "value": "%s" } `, key, value)
			ii = append(ii, i)
		}
		tag := strings.Join(ii, ", ")
		tags += tag
		tags += "] }"
	}
	return tags
}

func BuildFilterBody(env string, lifecycle string, criticality string, mTags map[string]string) string {
	var filters []string
	attributesContent := BuildAttributesFilter(env, lifecycle, criticality)
	if attributesContent != "" {
		filters = append(filters, attributesContent)
	}
	tagsContent := BuildTagsFilter(mTags)
	if tagsContent != "" {
		filters = append(filters, tagsContent)
	}
	filterContent := strings.Join(filters, ",")

	return fmt.Sprintf(`{ "filters": { %s } }`, filterContent)
}
