package terraformenterprise

import (
	"time"
)

type IngressTriggerAttributesT struct {
	Branch string `jsonapi:"attr,branch"`
	//   "branch": "",
	DefaultBranch string `jsonapi:"attr,default-branch"`
	//   "default-branch": true,
	VCSRootPath string `jsonapi:"attr,vcs-root-path"`
	//   "vcs-root-path": "",
	IngressSubmodules bool `jsonapi:"attr,ingress-submodules"`
	//   "ingress-submodules": false
}

type Workspace struct {
	ID   string `jsonapi:"primary,workspaces"`
	Name string `jsonapi:"attr,name"`
	// "name": "my-workspace-2",
	Environment string `jsonapi:"attr,environment"`
	// "environment": "default",
	AutoApply bool `jsonapi:"attr,auto-apply"`
	// "auto-apply": false,
	Locked bool `jsonapi:"attr,locked"`
	// "locked": false,
	CreatedAt time.Time `jsonapi:"attr,created-at,iso8601"`
	// "created-at": "2017-11-02T23:24:05.997Z",
	WorkingDirectory string `jsonapi:"attr,working-directory"`
	// "working-directory": "",
	TerraformVersion string `jsonapi:"attr,terraform-version"`
	// "terraform-version": "0.10.8",
	CanQueueDestroyPlan bool `jsonapi:"attr,can-queue-destroy-plan"`
	// "can-queue-destroy-plan": false,
	IngressTriggerAttributes *IngressTriggerAttributesT `jsonapi:"relation,ingress-trigger-attributes,omitempty"`
	// "ingress-trigger-attributes": { ... }
}

type OAuthClientT struct {
	ID string `jsonapi:"primary,oauth-clients"`
}

type OAuthToken struct {
	ID string `jsonapi:"primary,oauth-tokens"`
	CreatedAt time.Time `jsonapi:"attr,created-at,iso8601"`
	ServiceProviderUser string `jsonapi:"attr,service-provider-user"`
	HasSSHKey bool `jsonapi:"attr,has-ssh-key"`
	OAuthClient OAuthClientT `jsonapi:"relation,oauth-client`
}
