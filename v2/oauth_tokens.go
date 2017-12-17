package terraformenterprise

import (
	"fmt"
	"github.com/google/jsonapi"
	"reflect"
)

func (c *Client) ListOAuthTokens(organization string) ([]*OAuthToken, error) {
	request, err := c.NewRequest("GET", fmt.Sprintf("/organizations/%s/oauth-tokens", organization), nil)
	if err != nil {
		return nil, err
	}
	response, err := CheckResp(c.HTTPClient.Do(request))
	if err != nil {
		return nil, err
	}

	objs, err := jsonapi.UnmarshalManyPayload(response.Body, reflect.TypeOf(new(OAuthToken)))
	if err != nil {
		return nil, err
	}
	tokens := make([]*OAuthToken, len(objs))
	for i, o := range objs {
		t, ok := o.(*OAuthToken)
		if !ok {
			return nil, fmt.Errorf("Invalid type during unmarshaling, data was %v", o)
		}
		tokens[i] = t
	}
	return tokens, nil
}
