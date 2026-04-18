// Example: Create a new Wiki page in a project.
// This example demonstrates how to create a Wiki page
// in a specified project using the Backlog API client.
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/nattokin/go-backlog"
)

func main() {
	baseURL := os.Getenv("BACKLOG_BASE_URL")
	token := os.Getenv("BACKLOG_TOKEN")
	if baseURL == "" || token == "" {
		log.Fatalln("BACKLOG_BASE_URL and BACKLOG_TOKEN must be set")
	}

	c, err := backlog.NewClient(baseURL, token)
	if err != nil {
		log.Fatalln(err)
	}

	projectID := 12345

	wiki, err := c.Wiki.Create(context.Background(), projectID, "My Wiki Page", "Page content here.")
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("ID: %d, Name: %s\n", wiki.ID, wiki.Name)
}
