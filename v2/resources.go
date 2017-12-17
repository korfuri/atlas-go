package terraformenterprise

import (
	"time"
)

type Organization struct {
	ID string `jsonapi:"primary,organizations"`
}

type IngressTriggerAttributesT struct {
	Branch            string `jsonapi:"attr,branch"`
	DefaultBranch     bool   `jsonapi:"attr,default-branch"`
	VCSRootPath       string `jsonapi:"attr,vcs-root-path"`
	IngressSubmodules bool   `jsonapi:"attr,ingress-submodules"`
}

// Workspace represents a workspace in Terraform Enterprise.  To
// create a Workspace with a VCS connection, use CompoundWorkspace
// instead.
// https://www.terraform.io/docs/enterprise/api/workspaces.html
type Workspace struct {
	// ID is the ID of the workspace. Generated server-side.
	ID string `jsonapi:"primary,workspaces"`

	// Name is the human-friendly name of the workspace.
	Name string `jsonapi:"attr,name"`

	// TODO what's that
	Environment string `jsonapi:"attr,environment"`

	// AutoApply is whether changes get applied without human
	// approval
	AutoApply bool `jsonapi:"attr,auto-apply"`

	// TODO
	Locked bool `jsonapi:"attr,locked"`

	// CreatedAt is the timestamp of this workspace's creation
	CreatedAt time.Time `jsonapi:"attr,created-at,iso8601"`

	// WorkingDirectory is the working directory within the VCS
	// repository used to run Terraform for this workspace.
	WorkingDirectory string `jsonapi:"attr,working-directory"`

	// TerraformVersion is the version of Terraform in use in this
	// workspace.
	TerraformVersion string `jsonapi:"attr,terraform-version"`

	// TODO what's that
	CanQueueDestroyPlan bool `jsonapi:"attr,can-queue-destroy-plan"`

	// IngressTriggerAttributes is the settings struct for VCS
	// integration
	IngressTriggerAttributes *IngressTriggerAttributesT `jsonapi:"relation,ingress-trigger-attributes,omitempty"`
}

// CompoundWorkspace is a special type of Workspace resource used only
// when creating a VCS-integrated workspace.
type CompoundWorkspace struct {
	// Name is the name of the new workspace.
	Name string `jsonapi:"attr,name"`

	// LinkableRepiID is the name of the repository this workspace
	// is linked to. If you're using Github or Bitbucket this is
	// in the format "$user/$repo".
	LinkableRepoID string `jsonapi:"attr,linkable-repo-id"`

	// OAuthTokenID is the ID of a previously registered OAuth
	// token for Terraform to connect to your VCS system (Github,
	// Bitbucket, Gitlab...).
	OAuthTokenID string `jsonapi:"attr,oauth-token-id"`

	// WorkingDirectory is the path under the repo to which
	// Terraform enterprise cd's before running
	// Terraform. Optional.
	WorkingDirectory string `jsonapi:"attr,working-directory,omitempty"`

	// IngressTriggerAttributes is the settings struct for VCS
	// integration
	IngressTriggerAttributes *IngressTriggerAttributesT `jsonapi:"attr,ingress-trigger-attributes"`
}

type OAuthClientT struct {
	ID string `jsonapi:"primary,oauth-clients"`
}

type OAuthToken struct {
	ID                  string       `jsonapi:"primary,oauth-tokens"`
	CreatedAt           time.Time    `jsonapi:"attr,created-at,iso8601"`
	ServiceProviderUser string       `jsonapi:"attr,service-provider-user"`
	HasSSHKey           bool         `jsonapi:"attr,has-ssh-key"`
	OAuthClient         OAuthClientT `jsonapi:"relation,oauth-client`
}

type Variable struct {
	ID        string     `jsonapi:"primary,vars"`
	Key       string     `jsonapi:"attr,key"`
	Value     string     `jsonapi:"attr,value"`
	Sensitive bool       `jsonapi:"attr,sensitive"`
	Category  string     `jsonapi:"attr,category"`
	HCL       bool       `jsonapi:"attr,hcl"`
	//Workspace *Workspace `jsonapi:"relation,configurable"`
}
