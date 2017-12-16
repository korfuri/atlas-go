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
	fmt.Print(w2)

	ws, err = c.ListWorkspaces("Grab-TestAPI")
	checkErr(err)
	fmt.Print(ws)

	err = c.DeleteWorkspace("Grab-TestAPI", "test-uriel-destroy-me")
	checkErr(err)

	ws, err = c.ListWorkspaces("Grab-TestAPI")
	checkErr(err)
	fmt.Print(ws)

	
}
