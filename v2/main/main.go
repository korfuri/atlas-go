package main

import (
	"fmt"

	"github.com/hashicorp/atlas-go/v2"
)

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	c, err := terraformenterprise.NewClient(terraformenterprise.DefaultClientOptions())
	checkErr(err)

	ws, err := c.ListWorkspaces("Grab-TestAPI")
	checkErr(err)
	fmt.Print(ws)

	w := &terraformenterprise.Workspace{}
	w.Name = "test-uriel-destroy-me"
	w2, err := c.CreateWorkspace("Grab-TestAPI", w)
	checkErr(err)
	fmt.Println(w2)

	ws, err = c.ListWorkspaces("Grab-TestAPI")
	checkErr(err)
	fmt.Println(ws)

	err = c.DeleteWorkspace("Grab-TestAPI", "test-uriel-destroy-me")
	checkErr(err)

	ws, err = c.ListWorkspaces("Grab-TestAPI")
	checkErr(err)
	fmt.Println(ws)

	tokens, err := c.ListOAuthTokens("Grab-TestAPI")
	checkErr(err)
	fmt.Println(tokens)
	fmt.Println(*tokens[0])

	cw := &terraformenterprise.CompoundWorkspace{
		Name: "test-compound-wspace",
		WorkingDirectory: "/qa/base/network",
		LinkableRepoID: "Uriel-Test-Org/test-repo",
		OAuthTokenID: tokens[0].ID,
		IngressTriggerAttributes: &terraformenterprise.IngressTriggerAttributesT{
			Branch: "",
			DefaultBranch: true,
			VCSRootPath: "",
		},
	}
	fmt.Println(*cw)
	fmt.Println(*cw.IngressTriggerAttributes)
	w, err = c.CreateCompoundWorkspace("Grab-TestAPI", cw)

	ws, err = c.ListWorkspaces("Grab-TestAPI")
	checkErr(err)
	fmt.Println(ws)

	err = c.DeleteWorkspace("Grab-TestAPI", "test-compound-wspace")
	checkErr(err)

	ws, err = c.ListWorkspaces("Grab-TestAPI")
	checkErr(err)
	fmt.Println(ws)
}
