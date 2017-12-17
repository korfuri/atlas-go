package terraformenterprise

import (
	"bytes"
	"fmt"
	"reflect"

	"github.com/google/jsonapi"
)

func filterString(organization, workspace string) string {
	return "filter%5Bworkspace%5D%5Bname%5D=test-uriel-destroy-me&filter%5Borganization%5D%5Busername%5D=Grab-TestAPI"
}

func (c *Client) ListVariables(organization string, workspace string) ([]*Variable, error) {
	request, err := c.NewRequest("GET", "/vars", nil)
	if err != nil {
		return nil, err
	}
	request.URL.RawQuery = filterString(organization, workspace)
	response, err := CheckResp(c.HTTPClient.Do(request))
	if err != nil {
		return nil, err
	}

	objs, err := jsonapi.UnmarshalManyPayload(response.Body, reflect.TypeOf(new(Variable)))
	if err != nil {
		return nil, err
	}
	vars := make([]*Variable, len(objs))
	for i, o := range objs {
		t, ok := o.(*Variable)
		if !ok {
			return nil, fmt.Errorf("Invalid type during unmarshaling, data was %v", o)
		}
		vars[i] = t
	}
	return vars, nil
}

func (c *Client) GetVariableByKey(organization string, workspace string, key string) (*Variable, error) {
	vars, err := c.ListVariables(organization, workspace)
	if err != nil {
		return nil, err
	}
	for _, v := range vars {
		if v.Key == key {
			return v, nil
		}
	}
	// TODO: return nil, Not Found
	return nil, nil
}

func (c *Client) CreateVariable(organization string, workspace string, variable *Variable) (*Variable, error) {
	buf := new(bytes.Buffer)
	if err := jsonapi.MarshalPayload(buf, variable); err != nil {
		return nil, err
	}
	ro := &RequestOptions{
		Body: buf,
	}
	request, err := c.NewRequest("POST", "/vars", ro)
	if err != nil {
		return nil, err
	}
	request.URL.RawQuery = filterString(organization, workspace)
	response, err := CheckResp(c.HTTPClient.Do(request))
	if err != nil {
		return nil, err
	}

	out_var := new(Variable)
	if err := jsonapi.UnmarshalPayload(response.Body, out_var); err != nil {
		return nil, err
	}
	
	return out_var, nil
}

func (c *Client) UpdateVariable(variable *Variable) (*Variable, error) {
	buf := new(bytes.Buffer)
	if err := jsonapi.MarshalPayload(buf, variable); err != nil {
		return nil, err
	}
	ro := &RequestOptions{
		Body: buf,
	}
	request, err := c.NewRequest("PATCH", fmt.Sprintf("/vars/%s", variable.ID), ro)
	if err != nil {
		return nil, err
	}
	response, err := CheckResp(c.HTTPClient.Do(request))
	if err != nil {
		return nil, err
	}

	out_var := new(Variable)
	if err := jsonapi.UnmarshalPayload(response.Body, out_var); err != nil {
		return nil, err
	}
	
	return out_var, nil	
}
