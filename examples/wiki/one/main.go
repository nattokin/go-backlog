// Example: Get a single Wiki page by ID.
// This example demonstrates how to retrieve a specific Wiki page
// using its ID via the Backlog API client.
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

	wikiID := 12345

	wiki, err := c.Wiki.One(context.Background(), wikiID)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("ID: %d, Name: %s\n", wiki.ID, wiki.Name)
}
