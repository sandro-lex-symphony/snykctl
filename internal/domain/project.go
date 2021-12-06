package domain

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"snykctl/internal/tools"
)

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
	err := p.baseGet(true)
	if err != nil {
		return "", err
	}

	return p.rawResponse, nil
}

func (p *Projects) Get() error {
	return p.baseGet(false)
}

func (p *Projects) baseGet(raw bool) error {
	// TODO: add filter support
	// if FilterLifecycle != "" || FilterEnvironment != "" || (Key != "" && Value != "") {
	// 	return GetFilteredProjects(org_id)
	// }

	path := fmt.Sprintf("/org/%s/projects", p.Org.Id)
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

func (p *Projects) GetFiltered(env string, lifecycle string) error {
	path := fmt.Sprintf("/org/%s/projects", p.Org.Id)

	var attributes, filterContent string

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

	// if Key != "" && Value != "" {
	// 	tags += fmt.Sprintf(` "tags": { "includes": [ { "key": "%s", "value": "%s"} ] }`, Key, Value)
	// }

	// if attributes != "" && tags != "" {
	// 	filterContent = attributes + ", " + tags
	// } else if attributes != "" {
	// 	filterContent = attributes
	// } else {
	// 	filterContent = tags
	// }

	filterContent = attributes

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

func (p Projects) IsSync() bool {
	return p.sync
}
