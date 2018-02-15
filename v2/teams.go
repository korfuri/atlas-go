package tfe

import (
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
