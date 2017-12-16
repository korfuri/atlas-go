package main

import (
	"fmt"
	"os"

	"github.com/hashicorp/atlas-go/v2"
)

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	c, err := terraformenterprise.NewClient(os.Getenv("ATLAS_TOKEN"))
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
