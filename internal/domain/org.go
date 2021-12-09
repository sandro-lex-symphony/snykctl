package domain

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"snykctl/internal/config"
	"snykctl/internal/tools"
)

const orgsPath = "/orgs"
const orgPath = "/org"
const deleteOrgPath = "/org/%s"
const orgSettingsPath = "/org/%s/settings"

type Org struct {
	Id   string
	Name string
}

type Orgs struct {
	Orgs        []*Org
	sync        bool
	client      tools.HttpClient
	rawResponse string
}

func NewOrgs(c tools.HttpClient) *Orgs {
	o := new(Orgs)
	o.SetClient(c)
	return o
}

func (o *Orgs) GetOrgName(id string) (string, error) {
	if !o.Sync() {
		if err := o.Get(); err != nil {
			return "", err
		}
	}
	for _, org := range o.Orgs {
		if org.Id == id {
			return org.Name, nil
		}
	}
	return "", fmt.Errorf("getOrgName: org not found %s", id)
}

func (o *Orgs) SetClient(c tools.HttpClient) {
	o.client = c
}

func (o *Orgs) String() string {
	return o.toString("")
}

func (o *Orgs) Quiet() string {
	return o.toString("id")
}

func (o *Orgs) Names() string {
	return o.toString("name")
}

func (o *Orgs) toString(filter string) string {
	var ret string
	for _, org := range o.Orgs {
		if filter == "id" {
			ret += fmt.Sprintf("%s\n", org.Id)
		} else if filter == "name" {
			ret += fmt.Sprintf("%s\n", org.Name)
		} else {
			ret += fmt.Sprintf("%-38s %s\n", org.Id, org.Name)
		}

	}
	return ret
}

func (o Orgs) Print(quiet, names bool) {
	if quiet {
		fmt.Print(o.Quiet())
	} else if names {
		fmt.Print(o.Names())
	} else {
		fmt.Print(o.String())
	}
}

func (o *Orgs) Get() error {
	return o.baseGet(false)
}

func (o *Orgs) GetRaw() (string, error) {
	if err := o.baseGet(true); err != nil {
		return "", err
	}
	return o.rawResponse, nil
}

func (o *Orgs) baseGet(raw bool) error {
	resp := o.client.RequestGet(orgsPath)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("GetOrgs failed: %s", resp.Status)
	}

	if raw {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("GetOrgs failed: %s", err)
		}
		o.rawResponse = string(bodyBytes)
	} else {
		if err := json.NewDecoder(resp.Body).Decode(o); err != nil {
			return fmt.Errorf("GetOrgs failed: %s", err)
		}
	}
	o.sync = true

	return nil
}

func (o Orgs) Sync() bool {
	return o.sync
}

func CreateOrg(client tools.HttpClient, org_name string) error {
	jsonValue, _ := json.Marshal(map[string]string{
		"name":    org_name,
		"groupId": config.Instance.Id(),
	})
	resp := client.RequestPost(orgPath, jsonValue)
	if resp.StatusCode != http.StatusCreated {
		resp.Body.Close()
		return fmt.Errorf("addOrg failed: %s", resp.Status)
	}
	return nil
}

func DeleteOrg(client tools.HttpClient, org_id string) error {
	path := fmt.Sprintf(deleteOrgPath, org_id)
	resp := client.RequestDelete(path)
	if resp.StatusCode != http.StatusNoContent {
		resp.Body.Close()
		return fmt.Errorf("deleteOrg failed: %s", resp.Status)
	}
	return nil
}
