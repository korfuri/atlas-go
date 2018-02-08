package tfe

import (
	"bytes"
	"fmt"
	"log"
	"reflect"

	"github.com/google/jsonapi"
)

func (c *Client) ListWorkspaces(organization string) ([]*Workspace, error) {
	request, err := c.NewRequest("GET", fmt.Sprintf("/organizations/%s/workspaces", organization), nil)
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

// CreateWorkspace creates a workspace without VCS integration. For
// VCS-integrated workspaces, see CreateCompoundWorkspace.
func (c *Client) CreateWorkspace(organization string, workspace *Workspace) (*Workspace, error) {
	buf := new(bytes.Buffer)
	if err := jsonapi.MarshalOnePayloadEmbedded(buf, workspace); err != nil {
		return nil, err
	}
	ro := &RequestOptions{
		Body: buf,
	}
	log.Printf("[DEBUG] Request body: %s", buf.String())
	request, err := c.NewRequest("POST", fmt.Sprintf("/organizations/%s/workspaces", organization), ro)
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

// CreateCompoundWorkspace creates a workspace with VCS integration.
func (c *Client) CreateCompoundWorkspace(organization string, workspace *CompoundWorkspace) (*Workspace, error) {
	buf := new(bytes.Buffer)
	if err := jsonapi.MarshalPayload(buf, workspace); err != nil {
		return nil, err
	}
	fmt.Println(buf)
	ro := &RequestOptions{
		Body: buf,
	}
	log.Printf("[DEBUG] Request body: %s", buf.String())
	request, err := c.NewRequest("POST", "/compound-workspaces", ro)
	if err != nil {
		return nil, err
	}
	response, err := CheckResp(c.HTTPClient.Do(request))
	if err != nil {
		return nil, err
	}

	out_workspace := new(Workspace)
	buf2 := new(bytes.Buffer)
	buf2.ReadFrom(response.Body)
	log.Printf("[DEBUG] Response body: `%s`", buf2.String())
	if err := jsonapi.UnmarshalPayload(buf2, out_workspace); err != nil {
		return nil, err
	}

	return out_workspace, nil
}

func (c *Client) DeleteWorkspace(organization string, workspaceName string) error {
	request, err := c.NewRequest("DELETE", fmt.Sprintf("/organizations/%s/workspaces/%s", organization, workspaceName), nil)
	if err != nil {
		return err
	}
	_, err = CheckResp(c.HTTPClient.Do(request))
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) UpdateCompoundWorkspace(organization string, w *CompoundWorkspace) (*Workspace, error) {
	buf := new(bytes.Buffer)
	if err := jsonapi.MarshalPayload(buf, w); err != nil {
		return nil, err
	}
	ro := &RequestOptions{
		Body: buf,
	}
	log.Printf("[DEBUG] Request body: %s", buf.String())
	request, err := c.NewRequest("PATCH", fmt.Sprintf("/compound-workspaces/%s", w.ID), ro)
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

func (c *Client) UpdateWorkspace(organization string, w *Workspace) (*Workspace, error) {
	buf := new(bytes.Buffer)
	if err := jsonapi.MarshalOnePayloadEmbedded(buf, w); err != nil {
		return nil, err
	}
	ro := &RequestOptions{
		Body: buf,
	}
	log.Printf("[DEBUG] Request body: %s", buf.String())
	name := w.Name
	request, err := c.NewRequest("PATCH", fmt.Sprintf("/organizations/%s/workspaces/%s", organization, name), ro)
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
