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
	projectID := 12345
	// projectKey := "ProjectKey"

	// Get count of how many Wiki in project.
	count, err := c.Wiki.Count(backlog.ProjectID(projectID))
	// count, err := c.Wiki.Count(backlog.ProjectKey(projectKey))

	if err != nil {
		log.Fatalln(err)
	}

	// Output the number of count.
	fmt.Printf("%d\n", count)
}
