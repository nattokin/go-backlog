package main

import (
	"fmt"
	"log"

	backlog "github.com/nattokin/go-backlog"
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
	projectIDOrKey := "12345"
	// projectIDOrKey := "ProjectKey"

	// Get count of how many Wiki in project.
	count, err := c.Wiki.Count(projectIDOrKey)
	if err != nil {
		log.Fatalln(err)
	}

	// Output the number of count.
	fmt.Printf("%d\n", count)
}
