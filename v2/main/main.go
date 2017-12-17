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
	c, err := tfe.NewClient(tfe.DefaultClientOptions())
	checkErr(err)

	ws, err := c.ListWorkspaces("Grab-TestAPI")
	checkErr(err)
	fmt.Print(ws)

	w := &tfe.Workspace{}
	w.Name = "test-uriel-destroy-me"
	w2, err := c.CreateWorkspace("Grab-TestAPI", w)
	checkErr(err)
	fmt.Println(w2)

	ws, err = c.ListWorkspaces("Grab-TestAPI")
	checkErr(err)
	fmt.Println(ws)

	vars, err := c.ListVariables("Grab-TestAPI", "test-uriel-destroy-me")
	checkErr(err)
	fmt.Println(vars)
	//fmt.Println(vars[0])

	v, err := c.CreateVariable("Grab-TestAPI", "test-uriel-destroy-me", &tfe.Variable{
		Key:       "foo",
		Value:     "bar",
		Sensitive: false,
		Category:  "terraform",
		HCL:       false,
	})
	checkErr(err)
	fmt.Println(*v)

	v2, err := c.GetVariableByKey("Grab-TestAPI", "test-uriel-destroy-me", "foo")
	checkErr(err)
	fmt.Println(*v2)

	v2.Value = "baz"
	v3, err := c.UpdateVariable(v2)
	checkErr(err)
	fmt.Println(*v3)

	v4, err := c.GetVariableByKey("Grab-TestAPI", "test-uriel-destroy-me", "foo")
	checkErr(err)
	fmt.Println(*v4)

	err = c.DeleteWorkspace("Grab-TestAPI", "test-uriel-destroy-me")
	checkErr(err)

	ws, err = c.ListWorkspaces("Grab-TestAPI")
	checkErr(err)
	fmt.Println(ws)

	tokens, err := c.ListOAuthTokens("Grab-TestAPI")
	checkErr(err)
	fmt.Println(tokens)
	fmt.Println(*tokens[0])

	cw := &tfe.CompoundWorkspace{
		Name:             "test-compound-wspace",
		WorkingDirectory: "/qa/base/network",
		LinkableRepoID:   "Uriel-Test-Org/test-repo",
		OAuthTokenID:     tokens[0].ID,
		IngressTriggerAttributes: &tfe.IngressTriggerAttributesT{
			Branch:        "",
			DefaultBranch: true,
			VCSRootPath:   "",
		},
	}
	fmt.Println(*cw)
	fmt.Println(*cw.IngressTriggerAttributes)

	w3, err := c.CreateCompoundWorkspace("Grab-TestAPI", cw)
	checkErr(err)
	fmt.Println(*w3)

	ws, err = c.ListWorkspaces("Grab-TestAPI")
	checkErr(err)
	fmt.Println(ws)

	err = c.DeleteWorkspace("Grab-TestAPI", "test-compound-wspace")
	checkErr(err)

	ws, err = c.ListWorkspaces("Grab-TestAPI")
	checkErr(err)
	fmt.Println(ws)
}
