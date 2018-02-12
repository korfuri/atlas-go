package tfe

import (
	"fmt"
	"reflect"

	"github.com/google/jsonapi"
)

func (c *Client) ListLinkableRepos(oauthTokenID string) ([]*LinkableRepo, error) {
	request, err := c.NewRequest("GET", fmt.Sprintf("/oauth-tokens/%s/linkable-repos", oauthTokenID), nil)
	if err != nil {
		return nil, err
	}
	response, err := CheckResp(c.HTTPClient.Do(request))
	if err != nil {
		return nil, err
	}

	objs, err := jsonapi.UnmarshalManyPayload(response.Body, reflect.TypeOf(new(LinkableRepo)))
	if err != nil {
		return nil, err
	}
	tokens := make([]*LinkableRepo, len(objs))
	for i, o := range objs {
		t, ok := o.(*LinkableRepo)
		if !ok {
			return nil, fmt.Errorf("Invalid type during unmarshaling, data was %v", o)
		}
		tokens[i] = t
	}
	return tokens, nil
}
