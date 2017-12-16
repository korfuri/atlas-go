package terraformenterprise_test

import (
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/google/jsonapi"
	"github.com/hashicorp/atlas-go/v2"
)

func TestUnserializeWorkspacesList(t *testing.T) {
	r, err := os.Open("./examples/api__v2__organizaations__GrabTerraform__workspaces")
	if err != nil {
		t.Fatal(err)
	}
	workspaces, err := jsonapi.UnmarshalManyPayload(r, reflect.TypeOf(new(terraformenterprise.Workspace)))
	if err != nil {
		t.Fatal(err)
	}
	assert.Len(t, workspaces, 6)
	t.Logf("Workspaces are: %v", workspaces)
	assert.Equal(t, workspaces[0].(*terraformenterprise.Workspace).Name, "qa__base__network")
	t.Logf("Workspace[0] is: %v", workspaces[0])
}
