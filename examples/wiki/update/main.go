// Example: Update an existing Wiki page.
// This example demonstrates how to update the name and content
// of a Wiki page using the Backlog API client.
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

	// ID of the Wiki page to update.
	wikiID := 12345

	wiki, err := c.Wiki.Update(
		context.Background(),
		wikiID,
		c.Wiki.Option.WithName("Updated Name"),
		c.Wiki.Option.WithContent("Updated content."),
	)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("ID: %d, Name: %s\n", wiki.ID, wiki.Name)
}
