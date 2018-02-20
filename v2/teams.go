package tfe

import (
	"bytes"
	"fmt"
	"reflect"

	"github.com/google/jsonapi"
)

func (c *Client) ListTeams(organization string) ([]*Team, error) {
	request, err := c.NewRequest("GET", fmt.Sprintf("/organizations/%s/teams", organization), nil)
	if err != nil {
		return nil, err
	}
	response, err := CheckResp(c.HTTPClient.Do(request))
	if err != nil {
		return nil, err
	}

	objs, err := jsonapi.UnmarshalManyPayload(response.Body, reflect.TypeOf(new(Team)))
	if err != nil {
		return nil, err
	}
	tokens := make([]*Team, len(objs))
	for i, o := range objs {
		t, ok := o.(*Team)
		if !ok {
			return nil, fmt.Errorf("Invalid type during unmarshaling, data was %v", o)
		}
		tokens[i] = t
	}
	return tokens, nil
}

func (c *Client) CreateTeam(org string, t *Team) (*Team, error) {
	buf := new(bytes.Buffer)
	if err := jsonapi.MarshalPayload(buf, t); err != nil {
		return nil, err
	}
	ro := &RequestOptions{
		Body: buf,
	}
	request, err := c.NewRequest("POST", fmt.Sprintf("/organizations/%s/teams", org), ro)
	if err != nil {
		return nil, err
	}
	response, err := CheckResp(c.HTTPClient.Do(request))
	if err != nil {
		return nil, err
	}

	out_t := new(Team)
	if err := jsonapi.UnmarshalPayload(response.Body, out_t); err != nil {
		return nil, err
	}

	return out_t, nil
}

func (c *Client) GetTeamByID(id string) (*Team, error) {
	request, err := c.NewRequest("GET", fmt.Sprintf("/teams/%s", id), nil)
	if err != nil {
		return nil, err
	}
	response, err := CheckResp(c.HTTPClient.Do(request))
	if err != nil {
		return nil, err
	}

	out := new(Team)
	if err := jsonapi.UnmarshalPayload(response.Body, out); err != nil {
		return nil, err
	}

	return out, nil
}

func (c *Client) DeleteTeam(id string) error {
	request, err := c.NewRequest("DELETE", fmt.Sprintf("/teams/%s", id), nil)
	if err != nil {
		return err
	}
	_, err = CheckResp(c.HTTPClient.Do(request))
	if err != nil {
		return err
	}
	return nil
}
