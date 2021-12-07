package domain

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"snykctl/internal/tools"
	"strings"
)

const projectsPath = "/org/%s/projects"
const projectPath = "/org/%s/project/%s"
const tagPath = "/org/%s/project/%s/tags"
const attributesPath = "/org/%s/project/%s/attributes"

var validEnvironments = [9]string{"frontend", "backend", "internal", "external", "mobile", "saas", "on-prem", "hosted", "distributed"}
var validLifecycle = [3]string{"production", "development", "sandbox"}
var validCriticality = [4]string{"critical", "high", "medium", "low"}

type Projects struct {
	Org         Org
	Projects    []*Project
	sync        bool
	client      tools.HttpClient
	rawResponse string
}

type Project struct {
	Name string
	Id   string
}

func NewProjects(c tools.HttpClient, org_id string) *Projects {
	p := new(Projects)
	p.Org.Id = org_id
	p.SetClient(c)
	return p
}

func (p *Projects) SetClient(c tools.HttpClient) {
	p.client = c
}

func (p *Projects) GetRaw() (string, error) {
	path := fmt.Sprintf(projectsPath, p.Org.Id)
	err := p.baseGet(true, path)
	if err != nil {
		return "", err
	}

	return p.rawResponse, nil
}

func (p *Projects) Get() error {
	path := fmt.Sprintf(projectsPath, p.Org.Id)
	return p.baseGet(false, path)
}

func (p *Projects) GetProject(prj_id string) error {
	path := fmt.Sprintf(projectPath, p.Org.Id, prj_id)
	return p.baseGet(false, path)
}

func (p *Projects) GetRawProject(prj_id string) (string, error) {
	path := fmt.Sprintf(projectPath, p.Org.Id, prj_id)
	err := p.baseGet(true, path)
	if err != nil {
		return "", nil
	}
	return p.rawResponse, nil
}

func (p *Projects) baseGet(raw bool, path string) error {
	resp := p.client.RequestGet(path)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("GetProjects failed: %s", resp.Status)
	}

	if raw {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("GetProjects failed: %s", err)
		}
		p.rawResponse = string(bodyBytes)
	} else {
		if err := json.NewDecoder(resp.Body).Decode(&p); err != nil {
			return fmt.Errorf("GetProjects failed: %s", err)
		}
	}

	p.sync = true
	return nil
}

func (p *Projects) GetFiltered(env string, lifecycle string, mTags map[string]string) error {
	path := fmt.Sprintf(projectsPath, p.Org.Id)

	var attributes, filterContent, tags string

	if lifecycle != "" || env != "" {
		attributes += ` "attributes": { `

		if lifecycle != "" {
			attributes += fmt.Sprintf(`"lifecycle": [ "%s" ]`, lifecycle)
		}
		if lifecycle != "" && env != "" {
			attributes += ","
		}

		if env != "" {
			attributes += fmt.Sprintf(`"environment": [ "%s" ]`, env)
		}
		attributes += " }"
	}

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

	if attributes != "" && tags != "" {
		filterContent = attributes + ", " + tags
	} else if attributes != "" {
		filterContent = attributes
	} else {
		filterContent = tags
	}

	filters := fmt.Sprintf(`{ "filters": { %s } }`, filterContent)

	var jsonStr = []byte(filters)
	resp := p.client.RequestPost(path, jsonStr)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Get filtered projects list failed: %s ", resp.Status)
	}

	if err := json.NewDecoder(resp.Body).Decode(&p); err != nil {
		return fmt.Errorf("Get filtered projects list failed: %s ", err)
	}

	p.sync = true

	return nil
}

func (p *Projects) String() (string, error) {
	return p.toString("")
}

func (p *Projects) Quiet() (string, error) {
	return p.toString("id")
}

func (p *Projects) Names() (string, error) {
	return p.toString("name")
}

func (p *Projects) toString(filter string) (string, error) {

	var ret string
	for _, prj := range p.Projects {
		if filter == "id" {
			ret += fmt.Sprintf("%s\n", prj.Id)
		} else if filter == "name" {
			ret += fmt.Sprintf("%s\n", prj.Name)
		} else {
			ret += fmt.Sprintf("%-38s %s\n", prj.Id, prj.Name)
		}
	}
	return ret, nil
}

func (p Projects) Print(quiet, names bool) {
	var out string
	if quiet {
		out, _ = p.Quiet()
	} else if names {
		out, _ = p.Names()
	} else {
		out, _ = p.String()
	}

	fmt.Print(out)
}

func (p Projects) IsSync() bool {
	return p.sync
}

func (p *Projects) AddAttributes(prj_id string, env string, lifecycle string, criticality string) error {
	err := ParseAttributes(env, lifecycle, criticality)
	if err != nil {
		return err
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
	attrBody := fmt.Sprintf("{%s}", c)

	var jsonStr = []byte(attrBody)

	path := fmt.Sprintf(attributesPath, p.Org.Id, prj_id)
	resp := p.client.RequestPost(path, jsonStr)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to add attribute %s", resp.Status)
	}
	return nil
}

func (p *Projects) AddTag(prj_id string, tag string) error {
	k, v, err := ParseTag(tag)
	if err != nil {
		return err
	}

	path := fmt.Sprintf(tagPath, p.Org.Id, prj_id)

	tagBody := fmt.Sprintf(`{"key": "%s", "value": "%s"}`, k, v)
	var jsonStr = []byte(tagBody)

	resp := p.client.RequestPost(path, jsonStr)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to add tag %s", resp.Status)
	}
	return nil
}

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
