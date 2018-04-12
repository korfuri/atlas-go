package tfe

import (
	"bytes"
	"fmt"

	"github.com/korfuri/jsonapi"
)

const (
	AccessRead  = "read"
	AccessWrite = "write"
	AccessAdmin = "admin"
)

func (c *Client) CreateTeamAccess(ta *TeamAccess) (*TeamAccess, error) {
	buf := new(bytes.Buffer)
	if err := jsonapi.MarshalPayload(buf, ta); err != nil {
		return nil, err
	}
	ro := &RequestOptions{
		Body: buf,
	}
	request, err := c.NewRequest("POST", "/team-workspaces", ro)
	if err != nil {
		return nil, err
	}
	response, err := CheckResp(c.HTTPClient.Do(request))
	if err != nil {
		return nil, err
	}

	out_ta := new(TeamAccess)
	if err := jsonapi.UnmarshalPayload(response.Body, out_ta); err != nil {
		return nil, err
	}

	return out_ta, nil
}

func (c *Client) GetTeamAccessByID(id string) (*TeamAccess, error) {
	request, err := c.NewRequest("GET", fmt.Sprintf("/team-workspaces/%s", id), nil)
	if err != nil {
		return nil, err
	}
	response, err := CheckResp(c.HTTPClient.Do(request))
	if err != nil {
		return nil, err
	}

	out_ta := new(TeamAccess)
	if err := jsonapi.UnmarshalPayload(response.Body, out_ta); err != nil {
		return nil, err
	}

	return out_ta, nil
}

func (c *Client) DeleteTeamAccess(id string) error {
	request, err := c.NewRequest("DELETE", fmt.Sprintf("/team-workspaces/%s", id), nil)
	if err != nil {
		return err
	}
	_, err = CheckResp(c.HTTPClient.Do(request))
	if err != nil {
		return err
	}
	return nil
}
