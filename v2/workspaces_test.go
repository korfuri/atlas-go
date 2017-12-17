package tfe

import (
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"testing"

	"github.com/google/jsonapi"
	"github.com/hashicorp/atlas-go/v2/testutils"
	"github.com/stretchr/testify/assert"
)

// randomID returns a random id for a new workspace
func randomID() string {
	return fmt.Sprintf("ws-%d", rand.Int())
}

// This creates a mock handler for workspace listing, creation and
// deletion
func makeWorkspaceHandler(t *testing.T, srv *testutils.Server) func(rw http.ResponseWriter, req *http.Request) {
	workspaces := make([]*Workspace, 0)
	return func(rw http.ResponseWriter, req *http.Request) {
		srv.CheckHeaders(req)
		switch req.Method {
		case "GET":
			rw.WriteHeader(http.StatusOK)
			err := jsonapi.MarshalPayload(rw, workspaces)
			assert.NoError(t, err)
		case "POST":
			w := new(Workspace)
			err := jsonapi.UnmarshalPayload(req.Body, w)
			assert.NoError(t, err)
			rw.WriteHeader(http.StatusOK)
			w.ID = randomID()
			workspaces = append(workspaces, w)
			err = jsonapi.MarshalOnePayloadEmbedded(rw, w)
		case "DELETE":
			path := strings.Split(req.URL.Path, "/")
			wname := path[len(path)-1]
			for i, w := range workspaces {
				if w.Name == wname {
					workspaces = workspaces[:i]
					rw.WriteHeader(http.StatusOK)
					return
				}
			}
			rw.WriteHeader(http.StatusNotFound)
		default:
			t.Fail()
		}
	}
}

func TestLifecycle(t *testing.T) {
	client, server := testingClientServer(t)
	defer server.Stop()

	handler := makeWorkspaceHandler(t, server)
	server.Mux.HandleFunc("/organizations/TestOrg/workspaces", handler)
	server.Mux.HandleFunc("/organizations/TestOrg/workspaces/", handler)

	workspaces, err := client.ListWorkspaces("TestOrg")
	assert.NoError(t, err)
	assert.Len(t, workspaces, 0)

	newWorkspace, err := client.CreateWorkspace("TestOrg", &Workspace{Name: "my-workspace"})
	assert.NoError(t, err)
	assert.Equal(t, "my-workspace", newWorkspace.Name)

	workspaces, err = client.ListWorkspaces("TestOrg")
	assert.NoError(t, err)
	assert.Len(t, workspaces, 1)
	assert.Equal(t, "my-workspace", workspaces[0].Name)

	err = client.DeleteWorkspace("TestOrg", "my-workspace")
	assert.NoError(t, err)

	workspaces, err = client.ListWorkspaces("TestOrg")
	assert.NoError(t, err)
	assert.Len(t, workspaces, 0)
}
