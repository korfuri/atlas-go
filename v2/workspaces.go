package terraformenterprise

import (
	"bytes"
	"fmt"
	"github.com/google/jsonapi"
	"reflect"
)

func (c *Client) ListWorkspaces(organization string) ([]*Workspace, error) {
	request, err := c.Request("GET", fmt.Sprintf("/organizations/%s/workspaces", organization), nil)
	if err != nil {
		return nil, err
	}
	response, err := CheckResp(c.HTTPClient.Do(request))
	if err != nil {
		return nil, err
	}

	objs, err := jsonapi.UnmarshalManyPayload(response.Body, reflect.TypeOf(new(Workspace)))
	if err != nil {
		return nil, err
	}
	workspaces := make([]*Workspace, len(objs))
	for i, o := range objs {
		w, ok := o.(*Workspace)
		if !ok {
			return nil, fmt.Errorf("Invalid type during unmarshaling, data was %v", o)
		}
		workspaces[i] = w
	}
	return workspaces, nil
}

func (c *Client) GetWorkspaceByID(organization string, workspaceId string) (*Workspace, error) {
	workspaces, err := c.ListWorkspaces(organization)
	if err != nil {
		return nil, err
	}

	for _, w := range workspaces {
		if w.ID == workspaceId {
			return w, nil
		}
	}

	// TODO: return a proper error type that can be tested for the NotFound case
	return nil, nil
}

// TODO: handle creating a workspace with VCS repo stuff
func (c *Client) CreateWorkspace(organization string, workspace *Workspace) (*Workspace, error) {
	buf := new(bytes.Buffer)
	if err := jsonapi.MarshalPayload(buf, workspace); err != nil {
		return nil, err
	}
	ro := &RequestOptions{
		Body: buf,
	}
	request, err := c.Request("POST", fmt.Sprintf("/organizations/%s/workspaces", organization), ro)
	if err != nil {
		return nil, err
	}
	response, err := CheckResp(c.HTTPClient.Do(request))
	if err != nil {
		return nil, err
	}

	out_workspace := new(Workspace)
	if err := jsonapi.UnmarshalPayload(response.Body, out_workspace); err != nil {
		return nil, err
	}

	return out_workspace, nil
}

func (c *Client) DeleteWorkspace(organization string, workspaceName string) error {
	request, err := c.Request("DELETE", fmt.Sprintf("/organizations/%s/workspaces/%s", organization, workspaceName), nil)
	if err != nil {
		return err
	}
	_, err = CheckResp(c.HTTPClient.Do(request))
	if err != nil {
		return err
	}
	return nil
}
