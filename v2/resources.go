package tfe

import (
	"time"
)

// Organization represents an organization in Terraform Enterprise.
// The API to manipulate organizations is undocumented.
type Organization struct {
	ID             string `jsonapi:"primary,organizations"`
	Name           string `jsonapi:"attr,name"`
	Email          string `jsonapi:"attr,email"`
	EnterprisePlan string `jsonapi:"attr,enterprise-plan"`
}

type VCSRepoT struct {
	// LinkableRepiID is the name of the repository this workspace
	// is linked to. If you're using Github or Bitbucket this is
	// in the format "$user/$repo".
	LinkableRepoID string `json:"identifier"`

	// OAuthTokenID is the ID of a previously registered OAuth
	// token for Terraform to connect to your VCS system (Github,
	// Bitbucket, Gitlab...).
	OAuthTokenID string `json:"oauth-token-id"`

	Branch            string `json:"branch"`
	DefaultBranch     bool   `json:"default-branch"`
	IngressSubmodules bool   `json:"ingress-submodules"`
}

// Workspace represents a workspace in Terraform Enterprise.  To
// create a Workspace with a VCS connection, use CompoundWorkspace
// instead.
// https://www.terraform.io/docs/enterprise/api/workspaces.html
type Workspace struct {
	// ID is the ID of the workspace. Generated server-side.
	ID string `jsonapi:"primary,workspaces"`

	// Name is the human-friendly name of the workspace.
	Name string `jsonapi:"attr,name,omitempty"`

	// TODO what's that
	Environment string `jsonapi:"attr,environment,omitempty"`

	// AutoApply is whether changes get applied without human
	// approval
	AutoApply bool `jsonapi:"attr,auto-apply,omitempty"`

	// Locked is a read-only attribute that indicates whether the
	// workspace is locked.
	Locked bool `jsonapi:"attr,locked,omitempty"`

	// CreatedAt is the timestamp of this workspace's creation
	CreatedAt time.Time `jsonapi:"attr,created-at,iso8601,omitempty"`

	// WorkingDirectory is the working directory within the VCS
	// repository used to run Terraform for this workspace.
	WorkingDirectory string `jsonapi:"attr,working-directory"`

	// TerraformVersion is the version of Terraform in use in this
	// workspace.
	TerraformVersion string `jsonapi:"attr,terraform-version,omitempty"`

	// CanQueueDestroyPlan indicates whether the workspace allows
	// creating a destroy plan.
	CanQueueDestroyPlan bool `jsonapi:"attr,can-queue-destroy-plan,omitempty"`

	// IngressTriggerAttributes is the settings struct for VCS
	// integration
	VCSRepo VCSRepoT `jsonapi:"attr,vcs-repo,omitempty"`
}

// OauthClientT represents an OAuth Client. This is not directly
// manipulatable in the API, but it is used as part of the OAuthToken
// type.
type OAuthClientT struct {
	ID string `jsonapi:"primary,oauth-clients"`
}

// OAuthToken represents the OAuth Token associated with an OAuth Client.
type OAuthToken struct {
	ID                  string       `jsonapi:"primary,oauth-tokens"`
	CreatedAt           time.Time    `jsonapi:"attr,created-at,iso8601"`
	ServiceProviderUser string       `jsonapi:"attr,service-provider-user"`
	HasSSHKey           bool         `jsonapi:"attr,has-ssh-key"`
	OAuthClient         OAuthClientT `jsonapi:"relation,oauth-client`
}

// Variable represents a workspace variable.
// https://www.terraform.io/docs/enterprise/api/variables.html
type Variable struct {
	ID                   string     `jsonapi:"primary,vars"`
	Key                  string     `jsonapi:"attr,key"`
	Value                string     `jsonapi:"attr,value"`
	Sensitive            bool       `jsonapi:"attr,sensitive"`
	Category             string     `jsonapi:"attr,category"`
	HCL                  bool       `jsonapi:"attr,hcl"`
	Workspace            *Workspace `jsonapi:"relation,configurable"`
	workspaceForCreation string
	orgForCreation       string
}

// LinkableRepo represents a linkable repository, i.e. a repository
// hosted on a VCS hosting system like GitHub, Bitbucket, Gitlab, etc.
type LinkableRepo struct {
	// ID is a string identifying the repo. It appears to be the
	// name of the repo according to the repo hosting provider,
	// i.e. Grab/SecretProject for github.com/Grab/SecretProject.
	ID string `jsonapi:"primary,authorized-repos"`
}

// Team is a collection of users that may be granted permissions.
// https://www.terraform.io/docs/enterprise/api/teams.html
type Team struct {
	ID         string `jsonapi:"primary,teams"`
	Name       string `jsonapi:"attr,name"`
	UsersCount int    `jsonapi:"attr,users-count,omitempty"`
}

// TeamAccess represents a permission for a given Team to access a
// given Workspace.
// https://www.terraform.io/docs/enterprise/api/team-access.html
type TeamAccess struct {
	ID string `jsonapi:"primary,team-workspaces"`

	// Access should be "read", "write" or "admin"
	Access string `jsonapi:"attr,access"`

	Team      *Team      `jsonapi:"relation,team"`
	Workspace *Workspace `jsonapi:"relation,workspace"`
}
