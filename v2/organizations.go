package tfe

import (
	"bytes"
	"fmt"

	"github.com/google/jsonapi"
)

const (
	TrialVersion = "trial"
)

func (c *Client) CreateOrganization(org *Organization) (*Organization, error) {
	buf := new(bytes.Buffer)
	if err := jsonapi.MarshalPayload(buf, org); err != nil {
		return nil, err
	}
	ro := &RequestOptions{
		Body: buf,
	}
	request, err := c.NewRequest("POST", "/organizations", ro)
	if err != nil {
		return nil, err
	}
	response, err := CheckResp(c.HTTPClient.Do(request))
	if err != nil {
		return nil, err
	}

	out_org := new(Organization)
	if err := jsonapi.UnmarshalPayload(response.Body, out_org); err != nil {
		return nil, err
	}

	return out_org, nil
}

func (c *Client) GetOrganizationByID(id string) (*Organization, error) {
	request, err := c.NewRequest("GET", fmt.Sprintf("/organizations/%s", id), nil)
	if err != nil {
		return nil, err
	}
	response, err := CheckResp(c.HTTPClient.Do(request))
	if err != nil {
		return nil, err
	}

	out_org := new(Organization)
	if err := jsonapi.UnmarshalPayload(response.Body, out_org); err != nil {
		return nil, err
	}

	return out_org, nil
}

func (c *Client) DeleteOrganization(id string) error {
	request, err := c.NewRequest("DELETE", fmt.Sprintf("/organizations/%s", id), nil)
	if err != nil {
		return err
	}
	_, err = CheckResp(c.HTTPClient.Do(request))
	if err != nil {
		return err
	}
	return nil
}
