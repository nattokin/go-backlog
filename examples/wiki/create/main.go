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

	// Create a new Wiki by ID of the project.
	// You get struct where represented the Wiki created.
	r, err := c.Wiki.Create(12345, "name", "content")
	if err != nil {
		log.Fatalln(err)
	}

	// Output result.
	fmt.Printf("%#v\n", r)
}
