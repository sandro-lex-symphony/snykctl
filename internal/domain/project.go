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
	Name       string
	Id         string
	Attributes Attributes
	Tags       []*Tag
}

type Attributes struct {
	Criticality []string
	Environment []string
	Lifecycle   []string
}

type Tag struct {
	Key   string
	Value string
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
	if err := p.baseGet(true, path); err != nil {
		return "", err
	}

	return p.rawResponse, nil
}

func (p *Projects) Get() error {
	path := fmt.Sprintf(projectsPath, p.Org.Id)
	return p.baseGet(false, path)
}

func (p *Projects) GetRawProject(prj_id string) (string, error) {
	path := fmt.Sprintf(projectPath, p.Org.Id, prj_id)
	if err := p.baseGet(true, path); err != nil {
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

func (p *Projects) GetFiltered(env string, lifecycle string, criticality string, mTags map[string]string) error {
	return p.baseGetFiltered(false, env, lifecycle, criticality, mTags)
}

func (p *Projects) GetRawFiltered(env string, lifecycle string, criticality string, mTags map[string]string) (string, error) {
	if err := p.baseGetFiltered(true, env, lifecycle, criticality, mTags); err != nil {
		return "", err
	}
	return p.rawResponse, nil
}

func (p *Projects) baseGetFiltered(raw bool, env string, lifecycle string, criticality string, mTags map[string]string) error {
	path := fmt.Sprintf(projectsPath, p.Org.Id)
	filterBody := BuildFilterBody(env, lifecycle, criticality, mTags)

	var jsonStr = []byte(filterBody)
	resp := p.client.RequestPost(path, jsonStr)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Get filtered projects list failed: %s ", resp.Status)
	}

	if raw {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("GetProjects failed: %s", err)
		}
		p.rawResponse = string(bodyBytes)
	} else {
		if err := json.NewDecoder(resp.Body).Decode(&p); err != nil {
			return fmt.Errorf("Get filtered projects list failed: %s ", err)
		}
	}

	p.sync = true

	return nil
}

func (p *Projects) String() string {
	return p.toString("")
}

func (p *Projects) Quiet() string {
	return p.toString("id")
}

func (p *Projects) Names() string {
	return p.toString("name")
}

func (p *Projects) Verbose() string {
	return p.toString("verbose")
}

func (p *Projects) toString(filter string) string {
	var ret string
	for _, prj := range p.Projects {
		if filter == "id" {
			ret += fmt.Sprintf("%s\n", prj.Id)
		} else if filter == "name" {
			ret += fmt.Sprintf("%s\n", prj.Name)
		} else if filter == "verbose" {
			var attrs []string
			if len(prj.Attributes.Criticality) > 0 {
				attrs = append(attrs, prj.Attributes.Criticality[0])
			}
			if len(prj.Attributes.Environment) > 0 {
				attrs = append(attrs, prj.Attributes.Environment[0])
			}
			if len(prj.Attributes.Lifecycle) > 0 {
				attrs = append(attrs, prj.Attributes.Lifecycle[0])
			}
			if len(prj.Tags) > 0 {
				var attrTags []string
				for _, tag := range prj.Tags {
					t := fmt.Sprintf("%s=%s", tag.Key, tag.Value)
					attrTags = append(attrTags, t)
				}
				tagsStr := strings.Join(attrTags, ",")
				tagsStr = "[" + tagsStr + "]"
				attrs = append(attrs, tagsStr)
			}
			attributes := strings.Join(attrs, ",")
			ret += fmt.Sprintf("%-38s %-50s%s\n", prj.Id, prj.Name, attributes)
		} else {
			ret += fmt.Sprintf("%-38s %s\n", prj.Id, prj.Name)
		}
	}
	return ret
}

func (p Projects) Print(quiet, names, verbose bool) {
	var out string
	if quiet {
		out = p.Quiet()
	} else if names {
		out = p.Names()
	} else if verbose {
		out = p.Verbose()
	} else {
		out = p.String()
	}

	fmt.Print(out)
}

func (p Projects) IsSync() bool {
	return p.sync
}

func (p *Projects) AddAttributes(prj_id string, env string, lifecycle string, criticality string) error {
	if err := ParseAttributes(env, lifecycle, criticality); err != nil {
		return err
	}

	attrBody := BuildAttributesBody(env, lifecycle, criticality)

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

func (p *Projects) DeleteProject(prj_id string) error {
	path := fmt.Sprintf(projectPath, p.Org.Id, prj_id)
	resp := p.client.RequestDelete(path)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("deleteProject failed: %s", resp.Status)
	}
	return nil
}
