// Example: Delete a Wiki page by ID.
// This example demonstrates how to delete a Wiki page
// using the Backlog API client. The deleted Wiki is returned as a result.
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

	// ID of the Wiki page to delete.
	wikiID := 12345

	wiki, err := c.Wiki.Delete(context.Background(), wikiID)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("ID: %d, Name: %s\n", wiki.ID, wiki.Name)
}
