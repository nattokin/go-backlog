package main

import (
	"fmt"
	"log"

	backlog "github.com/nattokin/go-backlog"
)

func main() {
	// The base URL of Backlog API.
	baseURL := "https://zucc-hicc.backlog.com"
	// The tokun for request to Backlog API.
	token := "ftNsagLfrnSu9CleBraKl4rTdEcJ1CbYjZTtejBIu9D3wkCIhkxgmVdnfIMtqKKw"

	// Create Backlog API client.
	c, err := backlog.NewClient(baseURL, token)
	if err != nil {
		log.Fatalln(err)
	}

	// ID or Key of the project.
	projectIDOrKey := "API"
	// projectIDOrKey := "ProjectKey"

	// Get all Wikis by ID or Key of the project.
	// You get slice of Wiki.
	r, err := c.Wiki.All(projectIDOrKey)
	if err != nil {
		log.Fatalln(err)
	}

	// Out put
	for _, w := range r {
		fmt.Printf("%#v\n", w)
	}
}
