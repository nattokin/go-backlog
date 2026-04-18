// Example: List all Wiki pages in a project.
// This example demonstrates how to retrieve all Wiki pages
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

	projectKey := "MYPROJECT"

	wikis, err := c.Wiki.All(context.Background(), projectKey)
	// With keyword filter:
	// wikis, err := c.Wiki.All(context.Background(), projectKey, c.Wiki.Option.WithKeyword("keyword"))
	if err != nil {
		log.Fatalln(err)
	}

	for _, w := range wikis {
		fmt.Printf("ID: %d, Name: %s\n", w.ID, w.Name)
	}
}
