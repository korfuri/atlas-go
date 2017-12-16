package terraformenterprise

import (
	"fmt"
	"github.com/google/jsonapi"
	"github.com/hashicorp/atlas-go/v2/testutils"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"net/http"
	"reflect"
	"strings"
	"testing"
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
			ws, err := jsonapi.UnmarshalManyPayload(req.Body, reflect.TypeOf(new(Workspace)))
			assert.NoError(t, err)
			rw.WriteHeader(http.StatusOK)
			ws[0].(*Workspace).ID = randomID()
			workspaces = append(workspaces, ws[0].(*Workspace))
			err = jsonapi.MarshalOnePayloadEmbedded(rw, ws[0].(*Workspace))
			assert.NoError(t, err)
		case "DELETE":
			wname := strings.Split(req.URL.Path, "/")[:]
			t.Log(wname)
		default:
			t.Fail()
		}
	}
}

func TestLifecycle(t *testing.T) {
	client, server := testingClientServer(t)
	defer server.Stop()

	server.Mux.HandleFunc("/organizations/TestOrg/workspaces/", makeWorkspaceHandler(t, server))

	workspaces, err := client.ListWorkspaces("TestOrg")
	assert.NoError(t, err)
	assert.Len(t, workspaces, 0)

	newWorkspace, err := client.CreateWorkspace("TestOrg", &Workspace{Name: "my-workspace"})
	assert.NoError(t, err)
	assert.Equal(t, "my-workspace", newWorkspace.Name)

	workspaces, err = client.ListWorkspaces("TestOrg")
	assert.NoError(t, err)
	assert.Len(t, workspaces, 1)
}
