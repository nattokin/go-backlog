package main

import (
	"fmt"
	"log"

	"github.com/nattokin/go-backlog"
)

func main() {
	// The base URL of Backlog API.
	baseURL := "BACKLOG_BASE_URL"
	// The tokun for request to Backlog API.
	token := "BACKLOG_TOKEN"
	// Create Backlog API client.
	c, err := backlog.NewClient(baseURL, token)
	if err != nil {
		log.Fatalln(err)
	}
	// ID or Key of the project.
	projectKey := "PROJECTKEY"
	r, err := c.Wiki.All(backlog.ProjectKey(projectKey))
	// projectID := 1234
	// r, err := c.Wiki.All(backlog.ProjectID(projectID))

	if err != nil {
		log.Fatalln(err)
	}
	for _, w := range r {
		fmt.Printf("%#v\n", w)
	}
}
