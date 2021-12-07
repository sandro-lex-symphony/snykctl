package domain

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"snykctl/internal/tools"
)

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
		err := o.Get()
		if err != nil {
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

func (o *Orgs) String() (string, error) {
	return o.toString("")
}

func (o *Orgs) Quiet() (string, error) {
	return o.toString("id")
}

func (o *Orgs) Names() (string, error) {
	return o.toString("name")
}

func (o *Orgs) toString(filter string) (string, error) {
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

	return ret, nil
}

func (o Orgs) Print(quiet, names bool) {
	var out string
	if quiet {
		out, _ = o.Quiet()
	} else if names {
		out, _ = o.Names()
	} else {
		out, _ = o.String()
	}

	fmt.Print(out)

}

func (o *Orgs) Get() error {
	return o.baseGet(false)
}

func (o *Orgs) GetRaw() (string, error) {
	err := o.baseGet(true)
	if err != nil {
		return "", err
	}
	return o.rawResponse, nil
}

func (o *Orgs) baseGet(raw bool) error {
	resp := o.client.RequestGet("/orgs")
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
